package toledo

import (
	//"errors"
	//"encoding/json"
	"fmt"
	//"github.com/gorilla/mux" //路由库
	//"github.com/huin/goserial" //引入串口库
	"github.com/tarm/serial"
	//"io"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	"time"
)

//用到的常量
const (
	con_PORTNAME = "/dev/ttyO5" //周立功485串口2-----ttyO5
	con_BAUD     = 19200        //要连接的串口波特率
)

//串口发送接受的用到的通道
var chanWb = make(chan []byte, 1)      //发送的比特数组，缓冲0个
var chanRb = make(chan []byte, 1)      //接收的比特,缓冲0个
var chanRbNum = make(chan int, 1)      //无缓冲表面是互斥锁，只有这个有值才会让串口发送命令
var chanSerialBusy = make(chan int, 1) //有东西在里面代表busy，其他程序不要写上面的3个东西

/*=======================线程函数 init====================
	依赖：1.chekerr()函数
		  2.var chanSerialRecvice = make(chan []byte, 10) //接收的比特数组
		  3.
	作用：1.一直死循环向io.read发送命令并接受返回的n个字节，放到通道中
		  2.具有一个锁没有东西放进来会堵塞
	参数
		  2.wb.  chan []byte要发送到io里面的的[]byte
		  4.rb 将读到的数据放到rb中
		  5.rbnum 需要读取的字节长度，也是read函数要读取的字节数，注意也是阻塞的标志，只有取出后才会开始一次读取
	使用示范：可以参看/C_cmd/serial-yc/serial-package/goserial/serial_gopackage3.go
*====================================================*/
func goSendSerial(wb <-chan []byte, rb chan<- []byte, rbnum <-chan int) {
	//c := &goserial.Config{
	//	Name: con_PORTNAME,
	//	Baud: con_BAUD,
	//ReadTimeout: time.Second * 2, //读取超时
	//	Size:     goserial.Byte8,
	//	StopBits: goserial.StopBits1,
	//	Parity:   goserial.ParityNone,
	//} //以波特率和串口名打开
	c := &serial.Config{Name: con_PORTNAME, Baud: con_BAUD, ReadTimeout: time.Second * 4}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("线程goSendSerial启动,死循环发送命令")
	for {
		log.Println("堵塞中")
		rbnum := <-rbnum
		readbuf := make([]byte, rbnum) //堵塞，如果有值可以取出，开始一次发送命令

		log.Println("堵塞终止")
		send := <-wb
		log.Printf("重通道中得到发送命令：%x\n", send)

		//s, err := goserial.OpenPort(c) //打开串口
		s, err := serial.OpenPort(c) //打开串口
		checkerr(err)
		s.Write(send) //发送命令,wb中取出一个
		log.Printf("命令发送成功")

		b := make([]byte, 1)
		timer1 := time.NewTicker(2 * time.Second) //定时器到时间会退出下面的for
		defer timer1.Stop()
		istimeout := false //用来退出下面的for如果超时
		for i := 0; i < rbnum; i++ {
			select {
			case <-timer1.C:
				fmt.Println("接收数据超时")
				istimeout = true
				//如果2秒后还没完成读取数据的任务就退出

			default: //如果定时器还没动作就进行以下的默认操作
				s.Read(b)
				readbuf[i] = b[0]
			}
			if istimeout { //如果超时了就退出for循环
				break
			}
		}

		rb <- readbuf
		s.Close() //使用完后关闭

		log.Printf("已经发送：%x\n", send)
		fmt.Printf("接收到串口数据:%x\n", readbuf)
		//io.Close() //关闭io口
		//checkerr(err)
		//_, err = s.Write(append(g_cmd, sumCheck(g_cmd))) //在原来的命令后面再加一个校验和比特再发送
	}
}

func init() {
	go goSendSerial(chanWb, chanRb, chanRbNum)
}

//-----------------------------
//错误检查
//-----------------------------
func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//-----------------------------
//计算crc
//-----------------------------
func Crc(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	return crc16
}

//------------------------------------
//--uint16转换为2个byte
//----------------------------------
func uint16_to_twobyte(i uint16) (byte, byte) {
	bh := byte(i >> 8)   //高位
	bl := byte(i & 0xff) //低位
	return bh, bl
}

//==============================
//resetful 获取毛重，净重
//【00 地址】 【03 功能码】 【00 01 寄存器地址】【 00 02 读2个数据】 【94 1a CRC校验】只是读0001地址是毛重
//【00 地址】 【03 功能码】 【00 02 寄存器地址】【 00 02 读2个数据】 获取净重
//==========================================
func GetWeight(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始--resetful获取毛重净重。")

	cmd1 := []byte{0x00, 0x03, 0x00, 0x02, 0x00, 0x02}

	result_uint16 := Crc(cmd1) //uint16 格式,得到crc校验
	bh, bl := uint16_to_twobyte(result_uint16)

	//如果你得到的crc是1a94
	//但你要发送的命令应该是0x00, 0x03, 0x00, 0x01, 0x00, 0x02，0x94,0x1a
	//因为发送的命令crc应该是低位在前高位在后
	cmd2 := append(cmd1, bl)
	cmd2 = append(cmd2, bh) //在最后面加crc

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 9 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b := <-chanRb  //类型byte[20]
	<-chanSerialBusy

	fmt.Fprintf(w, "%x", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式

	log.Println("结束--resetful获取毛重净重。")
	//fmt.Printf("para:%s\n", b)
}
