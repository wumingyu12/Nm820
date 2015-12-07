package nm820

//========================================
//==========通用函数部分==================
//本文件下的go文件，按照协议命令一个go文件组织
//=======================================

import (
	//"errors"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"   //路由库
	"github.com/huin/goserial" //引入串口库
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

var currentState = &NM820_StatePara{} //被一秒钟更新一次的状态变量，在init函数启动的线程里面被更新
var currentStatehasLink int = 1       //代表是否有请求得到状态变量，如果没有就停止,这里一开始不能为0因为要初始化一个数值

/*=======================线程函数 init====================
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

/*=======================线程函数 init====================
	作用：开启一个线程，每隔15分钟读一次数据作为24小时温度
	机制：先读控制器的小时数，如果当前小时数为现在正在记录的小时内就加到sum中，否则就求和，并更新新的小时数
	依赖：1串口发送通道
		  2.NM820_StatePara结构体
		  3.time包,iotil包，os包
	使用：需要在init中启动
	生成的东西：resetful/nm820Json/Get24TemHumi.json

	//修改说明：
		20151203 将其修改为按现在温度前的方式显示，不是绝对的时间显示
========================================================================*/
func goGet24TemHumi() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("线程goGet24Tem启动,每隔15分钟取一次温度。")

	type TemHumi24Hour struct {
		Nowhour uint16 //增加 20151203 代表当前时间
		Time    []uint16
		Tavg    []float32
		Havg    []float32
	}
	var mytime uint16 = 0 //一开始的时间点为0时
	var temSum int16 = 0
	var humiSum uint16 = 0
	var sum uint16 = 0 //计数每次时间点的数据个数,不赋值为0，是为了避免第一次sum=0作为分母
	para := NM820_StatePara{}
	data := TemHumi24Hour{}
	var tAvg, hAvg float32

	//第一次运行时获取当前的小时值
	log.Println("获取当前小时值。")
	//发送数据并获取，前提func init()的运行,g_statepara在另一个go中
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- append(g_statepara, sumCheck(g_statepara))
	chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	rec := <-chanRb  //类型byte[100]
	<-chanSerialBusy
	para.reflashValue(rec) //当前的状态
	mytime = para.Hour

	//判断之前的存放历史数据的json文件是否存在
	_, err := os.Stat("./resetful/nm820Json/Get24TemHumi.json")
	if err != nil { //如果不存在,初始化一个json
		os.MkdirAll("./resetful/nm820Json", 0777)
		os.Create("./resetful/nm820Json/Get24TemHumi.json")
		//初始化一个json
		p := &TemHumi24Hour{}
		p.Nowhour = mytime //增加 20151203 代表当前时间
		p.Time = []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
		p.Tavg = []float32{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20}
		p.Havg = []float32{80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80, 80}
		//将p转换为json
		b, err := json.Marshal(p) //用这个函数时一定要确保字段名首位大写
		checkerr(err)
		ioutil.WriteFile("./resetful/nm820Json/Get24TemHumi.json", b, 0777) //写入到指定位置
	}

	//主循环
	for {
		log.Println("每隔6分钟取一次温度。")
		//发送数据并获取，前提func init()的运行,g_statepara在另一个go中
		chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
		chanWb <- append(g_statepara, sumCheck(g_statepara))
		chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
		rec := <-chanRb  //类型byte[100]
		<-chanSerialBusy
		para.reflashValue(rec) //当前的状态

		if mytime == para.Hour { //如果当前时间点为要记录的时间点
			temSum = temSum + para.TemAvg    //温度的总和
			humiSum = humiSum + para.HumiAvg //湿度的总和
			sum++                            //这个时间点的数据量加1
		} else { //记录时间点已经过去了，放到json文件中
			//先读出原纪录的json

			js, _ := ioutil.ReadFile("./resetful/nm820Json/Get24TemHumi.json")
			json.Unmarshal([]byte(js), &data) //将json解码为struct

			tAvg = float32(temSum) / (float32(sum) + 0.00000000001) / 10  //平均温度，这里还是存在bug假设程序运行的下一条命令恰好是从前一个小时数到后一个小时数，sum就会以0进入这里
			hAvg = float32(humiSum) / (float32(sum) + 0.00000000001) / 10 //避免出现除0的情况

			tFl := float32(int32(tAvg*10)) / 10 //保留1位小数
			hFl := float32(int32(hAvg*10)) / 10

			//刚运行那会mytime为0并且sum也为0的的第一次就不修改json
			//if mytime != 0 && sum != 0 {
			data.Nowhour = mytime //增加 20151203 代表当前时间
			data.Tavg[mytime] = tFl
			data.Havg[mytime] = hFl
			b, err := json.Marshal(data) //用这个函数时一定要确保字段名首位大写
			checkerr(err)
			ioutil.WriteFile("./resetful/nm820Json/Get24TemHumi.json", b, 0777) //将结果写回到json中
			//}
			//累加清零
			temSum = 0
			humiSum = 0
			sum = 0

			mytime = para.Hour //等于新的时间
		}

		//<-time.After(5 * time.Second)
		//隔10秒后继续
		time.Sleep(6 * time.Minute) //6分钟采样一次一小时采样10次

	}
}

/*=======================状态心跳线程函数 init====================
	作用：开启一个线程，每隔一秒读取一次NM820的数据，这个数据保存到var currentState = &NM820_StatePara{}
	依赖：1串口发送通道
		  2.NM820_StatePara结构体
		  3.time包,iotil包，os包
	使用：需要在init中启动
	注意：这个线程要在goSendSerial后面启动，不用先于
		  更新的结构体在resetful请求中/resetful/nm820/GetState会用到
		  var currentStatehasLink int = 0 //代表是否有请求得到状态变量，如果没有就停止
		  currentStatehasLink在resetful请求中会更改
	生成的东西：currentState = &NM820_StatePara{}
========================================================================*/
func goFlashCurrentState() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("init--启动线程每秒更新状态结构体。")

	for {
		if currentStatehasLink == 1 { //只有有连接请求时才会更新结构体没有就不更新
			//发送数据并获取，前提func init()的运行,g_statepara在另一个go中
			chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
			chanWb <- append(g_statepara, sumCheck(g_statepara))
			chanRbNum <- 100 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
			rec := <-chanRb  //类型byte[100]
			<-chanSerialBusy

			//用返回的数据更新结构体
			err := currentState.reflashValue(rec)
			checkerr(err)

			time.Sleep(1 * time.Second) //延迟1秒再次更新
			currentStatehasLink = 0     //将请求变回0
		} else {
			time.Sleep(1 * time.Second) //延迟1秒再次更新，不加这个会导致很严重的死循环
		}
	}

}

/*==========================包初始化函数===init===============================
	在包引入时会被调用，隐式调用
	作用：
		1.开启一个go fuc 用来一直保持串口的打开
		2.开启一个线程，每隔15分钟读一次数据作为24小时温度
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

	go goSendSerial(chanWb, chanRb, chanRbNum, s) //串口收发线程
	go goGet24TemHumi()                           //24小时温湿度读取
	go goFlashCurrentState()                      //每一秒更新nm820的状态值
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

//----------------------------------------------------------------------------
//将4个byte类型合并为一个uint32类型,组合后b1，b2排列，如果是小端请自行调换位置
//------------------------------------------------------------------------------
func twobyte_to_uint32(b1 byte, b2 byte, b3 byte, b4 byte) uint32 {
	return uint32(b1)<<24 + uint32(b2)<<16 + uint32(b3)<<8 + uint32(b4) //bh左移8位再加上低位的bl
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

//------------------------------------
//--int16转换为2个byte
//----------------------------------
func int16_to_twobyte(i int16) (byte, byte) {
	bh := byte(i >> 8)   //高位
	bl := byte(i & 0xff) //低位
	return bh, bl
}

//------------------------------------
//--uint32转换为4个byte
//----------------------------------
func uint32_to_fourbyte(i uint32) (byte, byte, byte, byte) {
	bhh := byte(i >> 24) //高位
	bh := byte(i >> 16)  //低位
	bl := byte(i >> 8)
	bll := byte(i & 0xff) //低位
	return bhh, bh, bl, bll
}

/**************************************************************************
给resetful调用的函数，外部可访问,resetful
***************************************************************************/

/*===============================resetful======================
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
	注意：currentStatehasLink代表有请求，线程的循环读取数据就会开启，在线程中一个循环后会将其变回0
==========================================================================*/

func GetState(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始--resetful请求nm820获取状态。")
	//告诉线程有请求
	currentStatehasLink = 1

	//将para转换为json
	b, err := json.Marshal(currentState) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Println("结束--resetful请求nm820获取状态。")
	//fmt.Printf("para:%s\n", b)
}

/*=======================================resetful===============================================
	请求：/resetful/nm820/GetDataHistory/{type}
		 type: Tem 返回温度 Humi 湿度  NH3氨气 Light光照
	作用：以当前日龄为准，向前倒推30日的历史温度数据，并返回
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_StatePara.go
=========================================================================================*/
func GetTempHistory(w http.ResponseWriter, r *http.Request) {
	//先判断要请求的历史数据类型
	vars := mux.Vars(r) //r为*http.Request
	datatype := vars["type"]

	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求获取历史数据--%s\n", datatype)

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
	log.Printf("当前日龄：%d", day)

	day1 := int16(day) //将uint16转换为int16 避免出现负数的情况时无法处理
	hd := &NM820_History30{}
	//注意如果day为当前日的话是返回错误的,所以从day-1日开始
	hd.addData(day1-1, datatype) //可选类型"Tem"---温度"Humi"---湿度"NH3"---氨气"Light"--光照
	//log.Println(hd)

	//将hd转换为json
	b, err := json.Marshal(hd) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Println("结束--resetful请求获取历史数据--%s\n", datatype)
}

/*=====================================resetful=================================================
	请求：/resetful/nm820/sysPara/WenduCurve
	作用：获取温度曲线表
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
=========================================================================================*/
func WenduCurve(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求获取温度曲线表数据\n")

	p := &NM820_WenduCurve{}
	err := p.addData() //用返回的串口数据更新结构体
	checkerr(err)

	//将hd转换为json
	b, err := json.Marshal(p) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Printf("结束--resetful请求获取温度曲线表数据\n")
}

/*=====================================resetful 获取最大最小通风等级=================================================
	请求：/resetful/nm820/sysPara/WindLevel
	作用：获取最大最小通风等级
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
=========================================================================================*/
func WindLevel(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求获取最大最小通风等级数据\n")

	p := &NM820_WindLevel{}
	err := p.addData() //用返回的串口数据更新结构体
	checkerr(err)

	//将hd转换为json
	b, err := json.Marshal(p) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Printf("结束--resetful请求获取最大最小通风等级数据\n")
}

/*=====================================resetful 获取温度曲线表=================================================
	请求：/resetful/nm820/sysPara/WindTables
	作用： 获取温度曲线表
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
=========================================================================================*/
func WindTables(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求获取通风等级数据\n")

	p := &NM820_WindTables{}
	err := p.addData() //用返回的串口数据更新结构体
	checkerr(err)

	//将hd转换为json
	b, err := json.Marshal(p) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Printf("结束--resetful请求获取通风等级数据\n")
}

/*=====================================resetful=  POST================================================
	请求：POST：/resetful/nm820/sysPara/WindTables
	作用：更新获取温度曲线表,post的数据实例
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
=========================================================================================*/
func ReflashWindTables(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful:POST请求更新通风等级数据\n")

	postData, _ := ioutil.ReadAll(r.Body) //读出发送的的post数据
	log.Println(string(postData))
	r.Body.Close()
	//初始化化一个20个等级的结构体
	p := &NM820_WindTables{}
	//p := &NM820_WindTables{}
	//n := &NM820_WindTable{}
	//for i := 0; i < 20; i++ {
	//	p.WindTables = append(p.WindTables, n)
	//}
	json.Unmarshal([]byte(postData), p) //将post的数据解析为结构体
	p.uploadData()
	log.Println(p)
	log.Printf("结束--resetful请求获取通风等级数据\n")
}

/*=====================================resetful  POST=================================================
	请求：POST：/resetful/nm820/sysPara/WenduCurve
	作用：修改温度曲线表
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
	注意：1.前端经过xeditable的使用后是string类型的，用`json:",string"`是不适用于数组，所以要在前端进行将string转int
		  2.
=========================================================================================*/
func ReflashWenduCurve(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--Post请求修改温度曲线表数据\n")

	postData, _ := ioutil.ReadAll(r.Body) //读出发送的的post数据
	log.Println(string(postData))
	r.Body.Close()
	p := &NM820_WenduCurve{}
	json.Unmarshal([]byte(postData), p) //将post的数据解析为结构体
	p.uploadData()
	log.Println(p)
	log.Printf("结束--Post请求修改温度曲线表数据\n")
}

/*=====================================resetful POST修改最大最小通风等级表=================================================
	请求：POST：/resetful/nm820/sysPara/WindLevel
	作用：修改温度曲线表
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
	注意：1.前端经过xeditable的使用后是string类型的，用`json:",string"`是不适用于数组，所以要在前端进行将string转int
		  2.
=========================================================================================*/
func ReflashWindLevel(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--Post请求最大最小通风等级数据\n")

	postData, _ := ioutil.ReadAll(r.Body) //读出发送的的post数据
	log.Println(string(postData))
	r.Body.Close()
	p := &NM820_WindLevel{}
	json.Unmarshal([]byte(postData), p) //将post的数据解析为结构体
	p.uploadData()
	log.Println(p)
	log.Printf("结束--Post请求最大最小通风等级数据\n")
}

/*=====================================resetful 获取系统变量参数表=================================================
	请求：GET /resetful/nm820/sysPara/SysValTable
	作用： 获取系统变量参数表
	返回：
	依赖的函数：
		2.NM820_sysPara.go
=========================================================================================*/
func GetSysVal(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求获取系统参数表\n")

	p := &NM820_SysVal{}
	err := p.addData() //用返回的串口数据更新结构体
	checkerr(err)

	//将hd转换为json
	b, err := json.Marshal(p) //用这个函数时一定要确保字段名首位大写
	checkerr(err)
	//必须要string,确保没发送其他了否则解释不了为json在angular
	fmt.Fprintf(w, "%s", b) //注意在armlinux下面不能用fmt.Fprintf(w, string(b))的方式
	//fmt.Printf("para:%s\n", b)
	log.Printf("结束--resetful请求获取通风等级数据\n")
}

/*=====================================resetful POST修改系统变量参数表=================================================
	请求：POST：/resetful/nm820/sysPara/SysValTable
	作用：修改系统变量参数表
	返回：
	依赖的函数：
		1.uint16_to_twobyte
		2.NM820_sysPara.go
	注意：1.前端经过xeditable的使用后是string类型的，用`json:",string"`是不适用于数组，所以要在前端进行将string转int
		  2.
=========================================================================================*/
func ReflashSysVal(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--Post请求修改系统变量参数表\n")

	postData, _ := ioutil.ReadAll(r.Body) //读出发送的的post数据
	log.Println(string(postData))
	r.Body.Close()
	p := &NM820_SysVal{}
	json.Unmarshal([]byte(postData), p) //将post的数据解析为结构体
	p.uploadData()
	log.Println(p)
	log.Printf("结束--Post请求修改系统变量参数表\n")
}
