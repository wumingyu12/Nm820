package main

import (
	"errors"
	"fmt"
	"github.com/huin/goserial" //引入串口库
	"log"
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
var g_cmd_recBuf = make([]byte, 100) //发送g_cmd返回的byte，理论上是100个

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

//---------------------------------------------
//--------byte[100]更新为NM820_StatePara结构体
//------------------------------------------------
func (nm *NM820_StatePara) reflashValue(b []byte) error { //默认输入的是100的byte[]
	//判断校验和是否一样
	fmt.Printf("last:%x\n", b[99])
	if sumCheck(b[0:99]) != b[99] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	//按小端的方式(具有就是高字节打头，如下)将byte赋值给结构体，重00 09 5a后面开始，就是byte[9]开始
	nm.GDay = twobyte_to_uint16(b[10], b[9])
	nm.Year = twobyte_to_uint16(b[12], b[11])
	nm.Month = twobyte_to_uint16(b[14], b[13])
	nm.Day = twobyte_to_uint16(b[16], b[15])
	nm.Hour = twobyte_to_uint16(b[18], b[17])
	nm.Min = twobyte_to_uint16(b[20], b[19])
	nm.Sec = twobyte_to_uint16(b[22], b[21])
	nm.TemAvg = twobyte_to_int16(b[24], b[23]) //
	for i := 0; i < 5; i++ {
		nm.Tem_1to5[i] = twobyte_to_int16(b[26+i*2], b[25+i*2])
	}
	nm.HumiAvg = twobyte_to_uint16(b[36], b[35])
	nm.Hmi_1to2[0] = twobyte_to_uint16(b[38], b[37])
	nm.Hmi_1to2[1] = twobyte_to_uint16(b[40], b[39])
	nm.NH3 = twobyte_to_uint16(b[42], b[41])
	nm.Light = twobyte_to_uint16(b[44], b[43])
	nm.FanLevel = twobyte_to_uint16(b[46], b[45])
	nm.Pos_SideWin = twobyte_to_uint16(b[48], b[47])
	nm.Pos_Curtain = twobyte_to_uint16(b[50], b[49])
	for i := 0; i < 4; i++ {
		nm.Pos_Roller[i] = twobyte_to_uint16(b[52+i*2], b[51+i*2])
	}
	for i := 0; i < 20; i++ {
		nm.RelayType[i] = b[59+i]
	}
	for i := 0; i < 20; i++ {
		nm.RelayState[i] = b[79+i]
	}
	return nil //什么错误都没就返回空
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
func twobyte_to_uint16(b1 byte, b2 byte) uint16 {
	return uint16(b1)<<8 + uint16(b2) //b1左移8位再加上低位的b2
}

//--------------------------------
//同上不过是int16
//-----------------------------------
func twobyte_to_int16(b1 byte, b2 byte) int16 {
	return int16(b1)<<8 + int16(b2) //b1左移8位再加上低位的b2
}

//============================================================
//============主函数=======================
//===================================================

func main() {
	c := &goserial.Config{
		Name:     CON_PORTNAME,
		Baud:     CON_BAUD,
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
		Parity:   goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	defer s.Close()                //用完关闭
	checkerr(err)
	//发送协议命令
	_, err = s.Write(append(g_cmd, sumCheck(g_cmd))) //在原来的命令后面再加一个校验和比特再发送
	checkerr(err)
	//接收协议命令
	for {
		_, err = s.Read(g_cmd_recBuf) //这里要注意只有完全读满buf才会完成这一步
		checkerr(err)

		fmt.Printf("get:%x\n", g_cmd_recBuf)

		nm820_StatePara := &NM820_StatePara{}
		//nm820_StatePara.Pos_Roller[1] = 0xfe

		fmt.Printf("Day:%x\n", nm820_StatePara.GDay)
		err = nm820_StatePara.reflashValue(g_cmd_recBuf)
		checkerr(err)

		fmt.Printf("GDay:%d\n", nm820_StatePara.GDay)
		fmt.Printf("Tem:%d\n", nm820_StatePara.TemAvg)
		fmt.Printf("Hour:%d\n", nm820_StatePara.Hour)
		fmt.Printf("FanLevel:%d\n", nm820_StatePara.FanLevel)
		fmt.Printf("relaystate:%x\n", nm820_StatePara.RelayState)
		fmt.Printf("RelayType:%x\n", nm820_StatePara.RelayType)
		//fmt.Printf("roller1:%x\n", nm820_StatePara.Pos_Roller[1])
		//fmt.Printf("%x\n B", nm820_StatePara.Month)
	}
}

//---------------------------------------------
//--------byte[100]转化为NM820_StatePara结构体
//------------------------------------------------
