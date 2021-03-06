package nm820

import (
	//"errors"
	//"fmt"
	//"github.com/huin/goserial" //引入串口库
	"../mylib/Int2byte" //整数与byte的转换
	"log"
	//"time"
)

//=========================================================
//---发送命令，得到nm820的状态，包括继电器，温湿度等
//---有些命令与参数是在nm820_main.go里面的
//====================================================

//============将要发送的协议命令，得到nm820的状态=============
var g_statepara = []byte{0x8a, 0x9b, 0x00, 0x01, 0x05, 0x02, 0x00, 0x09, 0x00}

//var b = make([]byte, 100) //发送g_statepara返回的byte，理论上是100个

//NM820状态结构体,根据协议的txt修改，请注意其数值都是小端存储的在NM820中
type NM820_StatePara struct {
	GDay        uint16    //日龄 0-1
	Year        uint16    //当前年 2-3
	Month       uint16    //月	4-5
	Day         uint16    //日	6-7
	Hour        uint16    //小时   8-9
	Min         uint16    //分钟 10-11
	Sec         uint16    //秒 12-13
	TemAvg      int16     //各温度值平均值，也是面板显示的值
	Tem_1to5    [5]int16  //1-5个探头的温度
	HumiAvg     uint16    //平均湿度
	Hmi_1to2    [2]uint16 //湿度1-2
	NH3         uint16    //氨气值，错误为0xffff
	Light       uint16    //光照值
	FanLevel    uint16    //通风级别
	Pos_SideWin uint16    //侧风窗位置
	Pos_Curtain uint16    //幕帘位置
	Pos_Roller  [4]uint16 //卷帘位置
	RelayType   [20]byte  //继电器类型
	/*
			/继电器类型
		           0 = 手动 ,1 = 风机1  ,2 = 风机2  ,3 = 风机3    ,4 = 风机4
		           5 = 风机5 ,6 = 风机6 ,7 = 风机7 ,8 = 风机8 ,9 = 加热
		           10 = 冷却水泵 ,11 = 喷雾 ,12 = 回流阀  ,13 = 照明1  ,14 = 照明2
		           15 = 侧风窗开 ,16 = 侧风窗关 ,17 = 幕帘开  ,18 = 幕帘关  ,19 = 卷帘1开
		           20 = 卷帘1关 ,21 = 卷帘2开 ,22 = 卷帘2关  ,23 = 卷帘3开  ,24 = 卷帘3关
		           25 = 卷帘4开 ,26 = 卷帘4关 ,27 = 喂料 ,28 = 额外系统1 ,29 = 额外系统2
		           30 = 额外系统3 ,31 = 告警
	*/
	RelayState [20]byte //继电器状态 0=断开 1=闭合
}

//--------------------------------------------------------------------------------------
//--------byte[100]更新为NM820_StatePara结构体
//---使用前先用sendCmd更新b，再用b作为参数给结构体赋值
//----------------------------------------------------------------------------------------
func (nm *NM820_StatePara) reflashValue(b []byte) error { //默认输入的是100的byte[]
	//判断校验和是否一样
	//fmt.Printf("last:%x\n", b[99])
	if sumCheck(b[0:99]) != b[99] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		//return errors.New("sum check is wrong!!")
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
	}
	//按小端的方式(具有就是高字节打头，如下)将byte赋值给结构体，重00 09 5a后面开始，就是byte[9]开始
	nm.GDay = Int2byte.Twobyte_to_uint16(b[10], b[9])
	nm.Year = Int2byte.Twobyte_to_uint16(b[12], b[11])
	nm.Month = Int2byte.Twobyte_to_uint16(b[14], b[13])
	nm.Day = Int2byte.Twobyte_to_uint16(b[16], b[15])
	nm.Hour = Int2byte.Twobyte_to_uint16(b[18], b[17])
	nm.Min = Int2byte.Twobyte_to_uint16(b[20], b[19])
	nm.Sec = Int2byte.Twobyte_to_uint16(b[22], b[21])
	nm.TemAvg = Int2byte.Twobyte_to_int16(b[24], b[23]) //
	for i := 0; i < 5; i++ {
		nm.Tem_1to5[i] = Int2byte.Twobyte_to_int16(b[26+i*2], b[25+i*2])
	}
	nm.HumiAvg = Int2byte.Twobyte_to_uint16(b[36], b[35])
	nm.Hmi_1to2[0] = Int2byte.Twobyte_to_uint16(b[38], b[37])
	nm.Hmi_1to2[1] = Int2byte.Twobyte_to_uint16(b[40], b[39])
	nm.NH3 = Int2byte.Twobyte_to_uint16(b[42], b[41])
	nm.Light = Int2byte.Twobyte_to_uint16(b[44], b[43])
	nm.FanLevel = Int2byte.Twobyte_to_uint16(b[46], b[45])
	nm.Pos_SideWin = Int2byte.Twobyte_to_uint16(b[48], b[47])
	nm.Pos_Curtain = Int2byte.Twobyte_to_uint16(b[50], b[49])
	for i := 0; i < 4; i++ {
		nm.Pos_Roller[i] = Int2byte.Twobyte_to_uint16(b[52+i*2], b[51+i*2])
	}
	for i := 0; i < 20; i++ {
		nm.RelayType[i] = b[59+i]
	}
	for i := 0; i < 20; i++ {
		nm.RelayState[i] = b[79+i]
	}
	return nil //什么错误都没就返回空
}
