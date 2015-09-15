package main

/*============================================================
	用chanenl的方式让串口一直存在，并返回状态值
	将读和写的线程合在一起
=============================================================*/
import (
	//"errors"
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

var chanWb = make(chan []byte, 1) //发送的比特数组，缓冲0个
var chanRb = make(chan []byte, 1) //接收的比特,缓冲0个
var chanRbNum = make(chan int, 1) //无缓冲表面是互斥锁，只有这个有值才会让串口发送命令

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

/*=======================线程函数====================
	依赖：1.chekerr()函数
		  2.var chanSerialRecvice = make(chan []byte, 10) //接收的比特数组
		  3.
	作用：1.一直死循环向io.read发送命令并接受返回的n个字节，放到通道中
		  2.具有一个锁没有东西放进来会堵塞
	参数：1.io。已经初始化的io.readWriteCloser
		  2.wb.  chan []byte要发送到io里面的的[]byte
		  4.rb 将读到的数据放到rb中
		  5.rbnum 需要读取的字节长度，也是read函数要读取的字节数，注意也是阻塞的标志，只有取出后才会开始一次读取
*====================================================*/
func goSendSerial(wb <-chan []byte, rb chan<- []byte, rbnum <-chan int, io io.ReadWriteCloser) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("线程goSendSerial启动,死循环发送命令")
	for {
		log.Println("堵塞中")
		rbnum := <-rbnum
		readbuf := make([]byte, rbnum) //堵塞，如果有值可以取出，开始一次发送命令
		log.Println("堵塞终止")
		send := <-wb
		io.Write(send) //发送命令,wb中取出一个
		b := make([]byte, 1)
		for i := 0; i < rbnum; i++ {
			io.Read(b)
			readbuf[i] = b[0]
		}
		//io.Read(readbuf) //接收串口数据,一个个字节读取否则有bug
		log.Println("3")
		rb <- readbuf
		log.Printf("已经发送：%x\n", send)
		fmt.Printf("接收到串口数据:%x\n", readbuf)
		//io.Close() //关闭io口
		//checkerr(err)
		//_, err = s.Write(append(g_cmd, sumCheck(g_cmd))) //在原来的命令后面再加一个校验和比特再发送
	}
}

//============================================================
//============主函数=======================
//===================================================

func main() {
	c := &goserial.Config{
		Name:        CON_PORTNAME,
		Baud:        CON_BAUD,
		ReadTimeout: time.Second * 5, //读取超时
		Size:        goserial.Byte8,
		StopBits:    goserial.StopBits1,
		Parity:      goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	checkerr(err)

	go goSendSerial(chanWb, chanRb, chanRbNum, s)
	chanWb <- append(g_cmd, sumCheck(g_cmd))
	chanRbNum <- 99 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	rec := <-chanRb
	log.Printf("接收通道：%x\n", rec)
	time.Sleep(100 * time.Second)
}
