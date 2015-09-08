package nm820

import (
	"errors"
	//"fmt"
	"github.com/huin/goserial" //引入串口库
)

//=========================================================
//---发送命令，得到nm820的状态，包括继电器，温湿度等
//---有些命令与参数是在nm820_main.go里面的
//====================================================

//============将要发送的协议命令，得到nm820的状态=============
var g_statepara = []byte{0x8a, 0x9b, 0x00, 0x01, 0x05, 0x02, 0x00, 0x09, 0x00}
var g_statepara_recBuf = make([]byte, 100) //发送g_statepara返回的byte，理论上是100个

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

//---------------------------------------------------
//--------发送协议命令--------------------------
//--------将返回的byte数组放到g_statepara_recBuf
//--------用串口发送-----
//--------------------------------------------------
func (nm *NM820_StatePara) sendCmd() {
	c := &goserial.Config{
		Name:     con_PORTNAME,
		Baud:     con_BAUD,
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
		Parity:   goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	defer s.Close()                //用完关闭
	checkerr(err)
	//发送协议命令
	_, err = s.Write(append(g_statepara, sumCheck(g_statepara))) //在原来的命令后面再加一个校验和比特再发送
	checkerr(err)
	//将接受到的数组赋值给g_statepara_recBuf
	_, err = s.Read(g_statepara_recBuf) //这里要注意只有完全读满buf才会完成这一步
	checkerr(err)

}

//--------------------------------------------------------------------------------------
//--------byte[100]更新为NM820_StatePara结构体
//---使用前先用sendCmd更新g_statepara_recBuf，再用g_statepara_recBuf作为参数给结构体赋值
//----------------------------------------------------------------------------------------
func (nm *NM820_StatePara) reflashValue() error { //默认输入的是100的byte[]
	//判断校验和是否一样
	//fmt.Printf("last:%x\n", g_statepara_recBuf[99])
	if sumCheck(g_statepara_recBuf[0:99]) != g_statepara_recBuf[99] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	//按小端的方式(具有就是高字节打头，如下)将byte赋值给结构体，重00 09 5a后面开始，就是byte[9]开始
	nm.GDay = twobyte_to_uint16(g_statepara_recBuf[10], g_statepara_recBuf[9])
	nm.Year = twobyte_to_uint16(g_statepara_recBuf[12], g_statepara_recBuf[11])
	nm.Month = twobyte_to_uint16(g_statepara_recBuf[14], g_statepara_recBuf[13])
	nm.Day = twobyte_to_uint16(g_statepara_recBuf[16], g_statepara_recBuf[15])
	nm.Hour = twobyte_to_uint16(g_statepara_recBuf[18], g_statepara_recBuf[17])
	nm.Min = twobyte_to_uint16(g_statepara_recBuf[20], g_statepara_recBuf[19])
	nm.Sec = twobyte_to_uint16(g_statepara_recBuf[22], g_statepara_recBuf[21])
	nm.TemAvg = twobyte_to_int16(g_statepara_recBuf[24], g_statepara_recBuf[23]) //
	for i := 0; i < 5; i++ {
		nm.Tem_1to5[i] = twobyte_to_int16(g_statepara_recBuf[26+i*2], g_statepara_recBuf[25+i*2])
	}
	nm.HumiAvg = twobyte_to_uint16(g_statepara_recBuf[36], g_statepara_recBuf[35])
	nm.Hmi_1to2[0] = twobyte_to_uint16(g_statepara_recBuf[38], g_statepara_recBuf[37])
	nm.Hmi_1to2[1] = twobyte_to_uint16(g_statepara_recBuf[40], g_statepara_recBuf[39])
	nm.NH3 = twobyte_to_uint16(g_statepara_recBuf[42], g_statepara_recBuf[41])
	nm.Light = twobyte_to_uint16(g_statepara_recBuf[44], g_statepara_recBuf[43])
	nm.FanLevel = twobyte_to_uint16(g_statepara_recBuf[46], g_statepara_recBuf[45])
	nm.Pos_SideWin = twobyte_to_uint16(g_statepara_recBuf[48], g_statepara_recBuf[47])
	nm.Pos_Curtain = twobyte_to_uint16(g_statepara_recBuf[50], g_statepara_recBuf[49])
	for i := 0; i < 4; i++ {
		nm.Pos_Roller[i] = twobyte_to_uint16(g_statepara_recBuf[52+i*2], g_statepara_recBuf[51+i*2])
	}
	for i := 0; i < 20; i++ {
		nm.RelayType[i] = g_statepara_recBuf[59+i]
	}
	for i := 0; i < 20; i++ {
		nm.RelayState[i] = g_statepara_recBuf[79+i]
	}
	return nil //什么错误都没就返回空
}
