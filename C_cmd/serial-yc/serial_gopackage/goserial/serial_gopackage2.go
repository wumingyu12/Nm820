package main

/*============================================================
	用chanenl的方式让串口一直存在，并返回状态值
=============================================================*/
import (
	"errors"
	"fmt"
	"github.com/huin/goserial" //引入串口库
	"io"
	"log"
	"time"
)

//用到的常量
const (
	CON_PORTNAME = "/dev/ttySAC3" //要连接的串口名字，在window下可以用“COM1”
	CON_BAUD     = 4800           //要连接的串口波特率
	//CON_ADDR_H   = 0x00           //通信地址高位
	//CON_ADDR_L   = 0x01           //通信地址低位，可以通过nm820的界面参看就是栏舍号
	//CMD []byte={}
)

//============将要发送的协议命令=============
var g_cmd = []byte{0x8a, 0x9b, 0x00, 0x01, 0x05, 0x02, 0x00, 0x09, 0x00}

//var g_cmd_recBuf = make([]byte, 100) //发送g_cmd返回的byte，理论上是100个

var chanSerialSend = make(chan []byte, 10)   //发送的比特数组，缓冲10个
var chanSerialRecvice = make(chan byte, 200) //接收的比特,缓冲100个
var chanSerialSendBegin = make(chan int)     //无缓冲表面是互斥锁，只有这个有值才会让串口发送命令

//-------------------------------------------------
//在计算byte的累加验证位，这是nm820采用的验证方式
//--------------------------------------------------
func sumCheck(date []byte) byte {
	var sum byte = 0x00
	for i := 0; i < len(date); i++ {
		sum = sum + date[i]
	}
	return sum
}

//-----------------------------
//错误检查
//-----------------------------
func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*==================消费者=====线程函数====================
	依赖：1.chekerr()函数
		  2.var chanSerialRecvice = make(chan []byte, 10) //接收的比特数组
		  3.
	作用：1.一直死循环向io.read发送命令
	参数：1.io。已经初始化的io.readWriteCloser
		  2.wb.  chan []byte要发送到io里面的的[]byte
		  3.sendbegin 通道只有改sendbegin有值的时候才发送命令，避免了在接收时还在发送命令
*====================================================*/
func goSendSerial(wb <-chan []byte, sendbegin <-chan int, io io.ReadWriteCloser) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("线程goSendSerial启动,死循环发送命令")
	for {
		log.Println("堵塞中")
		<-sendbegin //堵塞，如果有值可以取出，开始一次发送命令
		log.Println("堵塞终止")
		send := <-wb
		_, err := io.Write(send)
		//io.Close() //关闭io口
		checkerr(err)
		log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
		log.Printf("已经发送：%x\n", send)
		//_, err = s.Write(append(g_cmd, sumCheck(g_cmd))) //在原来的命令后面再加一个校验和比特再发送
	}
}

/*======================生产者===线程函数===================================
  依赖：1.checkerr
		2.需要 goSendSerial发送命令后才会有串口数据返回
		3.var chanSerialRecvice = make(chan []byte, 1024) //接收的比特数组

===========================================================================*/
func goRecviceSerial(rb chan<- byte, io io.ReadWriteCloser) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("线程goSendSerial启动,死循环接收数据")
	var te = make([]byte, 1) //一个个byte那样接收
	for {
		io.Read(te) //一个个byte那样接收
		rb <- te[0] //具有1024个byte的缓冲
		fmt.Printf("接收到串口数据:%x\n", te[0])
	}
}

/*======================从接收的cchanel中读取指定长度的字节数==================
	参数：1.num int 连续读取的channel中的字节数，int值不能小于head的长度
		  2.ch  make(chan []byte, 10) 带buf的chan  chanSerialRecvice
		  3.head []byte 为了屏蔽掉错误的数据,指定只有读到head头的时候才正式作为有效数
	内置参数：1.超时时间5秒
	过程：
		1.在Chanel中不断读取字节，如果读取到head开头的字节就开始存储
	    2.在ch中是串口不断死循环读回来的数，最多1024个
	返回：
		1.return nil, errors.New("num值要大于匹配头head的值")
		2.return nil, errors.New("超时5s不能匹配到头字节")
		3.return data, nil //正常
============================================================================*/
func readNumFromChanel(num int, head []byte, ch <-chan byte) ([]byte, error) {
	data := make([]byte, num)
	headnum := len(head)

	if headnum > num {
		log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
		log.Println("num值要大于匹配头head的值")
		return nil, errors.New("num值要大于匹配头head的值")
	}
	var headMatchNum int = 0 //匹配了多少个头字节
	//-------------------------------------------
	//死循环进行头匹配，避免前头的错误字节干扰
	//----------------------------------------
	for {
		if headMatchNum == headnum { //如果可以匹配的数和头匹配一样，就退出死循环
			break //退出到下面的data赋值
		}
		select {
		case byt := <-ch: //如果可以从通道拿到数据
			{
				log.Printf("接收:%x\n", byt)
				if byt == head[headMatchNum] { //如果可以匹配上就匹配数加1
					data[headMatchNum] = byt //将head字节头复制到将要发送的data里面
					headMatchNum++
				}
			}
		case <-time.After(5 * time.Second): //超时5s，读完ch里面的所以byte或者没有数据超时都没匹配成功
			{
				log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
				s := "超时5s匹配到的头字节个数：" + fmt.Sprintf("%d", headMatchNum) + "个"
				log.Println(s)
				return nil, errors.New(s)
			}
		}

	}
	s := "匹配到的头字节个数：" + fmt.Sprintf("%d", headMatchNum) + "个"
	log.Println(s)
	//------------------------------------------------------
	//成功匹配头字节后，开始将数据从通道取出来，并发送出去
	//---------------------------------------------------
	for i := headnum; i < num; i++ {
		select {
		case byt := <-ch: //如果可以从通道拿到数据
			data[i] = byt
			//log.Printf("接收:%d--%x\n", i, data[i])
		case <-time.After(5 * time.Second): //超时5s，读完ch里面的所以byte或者没有数据超时都没匹配成功
			log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
			log.Printf("超时5s不能读出%d数量的字节\n", num)
			return nil, errors.New("超时5s不能读出指定数量的字节")
		}
	}
	//可以完全读取后
	return data, nil
}

func test(ch <-chan byte) {
	for i := 0; i < 100; i++ {
		byt := <-ch
		log.Printf("取出数据：%d-%x\n", i, byt)
	}

}

//============================================================
//============主函数=======================
//===================================================

func main() {
	c := &goserial.Config{
		Name:        CON_PORTNAME,
		Baud:        CON_BAUD,
		ReadTimeout: time.Second * 5,
		Size:        goserial.Byte8,
		StopBits:    goserial.StopBits1,
		Parity:      goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	checkerr(err)

	go goSendSerial(chanSerialSend, chanSerialSendBegin, s)
	go goRecviceSerial(chanSerialRecvice, s)
	//for i := 0; i < 10; i++ {
	//head := []byte{0x8a, 0x9b}
	chanSerialSend <- append(g_cmd, sumCheck(g_cmd))
	chanSerialSendBegin <- 1 //开启一次锁让进程发送一次命令
	//test(chanSerialRecvice)
	chanSerialSend <- append(g_cmd, sumCheck(g_cmd))
	chanSerialSendBegin <- 1 //开启一次锁让进程发送一次命令
	//d, _ := readNumFromChanel(100, head, chanSerialRecvice) //从chanSerialRecvice中读出6个字节，以head字节头为准
	//log.Printf("data:%x\n", d)
	//}

	//for i := 0; i < 3; i++ {
	//	log.Println("chanserialwrite===")
	//	chanSerialWrite <- g_cmd
	//}
	//defer s.Close()                //用完关闭
	//将串口的打开一直用线程存在，用channel来发送数据与接收数据。
	time.Sleep(100 * time.Second)
}
