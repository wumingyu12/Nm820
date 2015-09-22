package nm820

//========================================
//==========通用函数部分==================
//本文件下的go文件，按照协议命令一个go文件组织
//=======================================

import (
	//"errors"
	"encoding/json"
	"fmt"
	"github.com/huin/goserial" //引入串口库
	"io"
	"log"
	"net/http"
)

//用到的常量
const (
	con_PORTNAME = "/dev/ttySAC3" //要连接的串口名字，在window下可以用“COM1”
	con_BAUD     = 4800           //要连接的串口波特率
	//CON_ADDR_H   = 0x00           //通信地址高位
	//CON_ADDR_L   = 0x01           //通信地址低位，可以通过nm820的界面参看就是栏舍号
	//CMD []byte={}
)

//串口发送接受的用到的通道
var chanWb = make(chan []byte, 1)      //发送的比特数组，缓冲0个
var chanRb = make(chan []byte, 1)      //接收的比特,缓冲0个
var chanRbNum = make(chan int, 1)      //无缓冲表面是互斥锁，只有这个有值才会让串口发送命令
var chanSerialBusy = make(chan int, 1) //有东西在里面代表busy，其他程序不要写上面的3个东西

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
	使用示范：可以参看/C_cmd/serial-yc/serial-package/goserial/serial_gopackage3.go
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
		log.Printf("重通道中得到发送命令：%x\n", send)
		io.Write(send) //发送命令,wb中取出一个
		b := make([]byte, 1)
		for i := 0; i < rbnum; i++ {
			io.Read(b)
			readbuf[i] = b[0]
			//log.Printf("接收到：%d--%x\n", i, b[0])
		}
		//io.Read(readbuf) //接收串口数据,一个个字节读取否则有bug
		rb <- readbuf
		log.Printf("已经发送：%x\n", send)
		fmt.Printf("接收到串口数据:%x\n", readbuf)
		//io.Close() //关闭io口
		//checkerr(err)
		//_, err = s.Write(append(g_cmd, sumCheck(g_cmd))) //在原来的命令后面再加一个校验和比特再发送
	}
}

/*==========================包初始化函数==================================
	在包引入时会被调用，隐式调用
	作用：
		1.开启一个go fuc 用来一直保持串口的打开
	使用示范：
			go goSendSerial(chanWb, chanRb, chanRbNum, s)
			chanWb <- append(g_cmd, sumCheck(g_cmd))
			chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
			rec := <-chanRb
=======================================================================*/
func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	c := &goserial.Config{
		Name: con_PORTNAME,
		Baud: con_BAUD,
		//ReadTimeout: time.Second * 5, //读取超时
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
		Parity:   goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	checkerr(err)

	go goSendSerial(chanWb, chanRb, chanRbNum, s)
}

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

//----------------------------------------------------------------------------
//将两个byte类型合并为一个uint16类型,组合后b1，b2排列，如果是小端请自行调换位置
//------------------------------------------------------------------------------
func twobyte_to_uint16(bh byte, bl byte) uint16 {
	return uint16(bh)<<8 + uint16(bl) //bh左移8位再加上低位的bl
}

//--------------------------------
//同上不过是int16
//-----------------------------------
func twobyte_to_int16(b1 byte, b2 byte) int16 {
	return int16(b1)<<8 + int16(b2) //b1左移8位再加上低位的b2
}

//------------------------------------
//--uint16转换为2个byte
//----------------------------------
func uint16_to_twobyte(i uint16) (byte, byte) {
	bh := byte(i >> 8)   //高位
	bl := byte(i & 0xff) //低位
	return bh, bl
}

/**************************************************************************
给resetful调用的函数，外部可访问,resetful
***************************************************************************/

/*=====================================================
	得到nm820的状态，主体函数在nm820_statePara.go里面
	发送的json根据结构图NM820_StatePara
	需要的外部参数：
		1.g_statepara 发送的命令
    请求url：/resetful/nm820/GetState
 	正常返回一个json：
				{"GDay":34,"Year":2000,"Month":1,"Day":13,"Hour":23,"Min":59,"Sec":9,"TemAvg":242,
				"Tem_1to5":[248,244,-1,-1,32767],"HumiAvg":863,"Hmi_1to2":[727,1000],"NH3":65535,"Light":65535,
				"FanLevel":7,"Pos_SideWin":60,"Pos_Curtain":0,"Pos_Roller":[0,0,0,0],
				"RelayType":[1,2,3,4,5,6,7,8,9,10,11,13,27,15,16,17,18,28,29,31],
				"RelayState":[1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1]}
  	错误的返回(字符串)：
				返回"error1:cant open serial port." 情况1不能打开串口，
				返回"error"
	依赖：1.func init()的运行
==========================================================================*/

func GetState(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始--resetful请求nm820获取状态。")
	para := &NM820_StatePara{}

	//发送数据并获取，前提func init()的运行,g_statepara在另一个go中
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- append(g_statepara, sumCheck(g_statepara))
	chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	rec := <-chanRb  //类型byte[100]
	<-chanSerialBusy

	//用返回的数据更新结构体
	err := para.reflashValue(rec)
	checkerr(err)

	//将para转换为json
	b, err := json.Marshal(para) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Println("结束--resetful请求nm820获取状态。")
	//fmt.Printf("para:%s\n", b)
}

/*======================================================================================
	请求：
	作用：以当前日龄为准，向前倒推30日的历史温度数据，并返回
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_StatePara.go
=========================================================================================*/
func GetTempHistory(w http.ResponseWriter, r *http.Request) {

	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始--resetful请求获取温度历史数据。")

	//先得到当前的日龄
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- append(g_statepara, sumCheck(g_statepara))
	chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	rec := <-chanRb  //类型byte[100]
	<-chanSerialBusy

	p := &NM820_StatePara{}
	err := p.reflashValue(rec) //用返回的数据更新结构体
	checkerr(err)
	day := p.GDay //得到当前日龄uint16格式

	hd := &NM820_History30{}
	//注意如果day为当前日的话是返回错误的,所以从day-1日开始
	hd.addData(day-1, "Tem") //可选类型"Tem"---温度"Humi"---湿度"NH3"---氨气"Light"--光照
	//log.Println(hd)

	//将hd转换为json
	b, err := json.Marshal(hd) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Println("结束--resetful请求获取温度历史数据。")
}
