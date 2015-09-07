package main

import (
	"fmt"
	"github.com/huin/goserial" //引入串口库
	"log"
)

//用到的常量
const (
	CON_PORTNAME = "/dev/ttySAC3" //要连接的串口名字，在window下可以用“COM1”
	CON_BAUD     = 9600           //要连接的串口波特率

)

//NM820状态结构体
type NM820_StatePara struct {
	GDay  byte //日龄 0-1
	Year  byte //当前年 2-3
	Month byte //月	4-5
	Day   byte //日	6-7
}

func main() {
	c := &goserial.Config{
		Name:     CON_PORTNAME,
		Baud:     CON_BAUD,
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
		Parity:   goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	checkerr(err)
	//发送的命令，定义

	//n, err := s.Write([]byte("test,hello"))
	//============将要发送的协议命令=============
	cmd := []byte{0xff, 0xff, 0xff, 0xff}
	//发送协议命令
	_, err = s.Write(cmd)
	checkerr(err)
	//接收协议命令
	buf := make([]byte, 4) //4个比特的bu
	for {
		_, err = s.Read(buf) //这里要注意只有完全读满buf才会完成这一步
		checkerr(err)
		fmt.Printf("%x\n", buf[0])
		fmt.Printf("%x\n", buf[1])
		fmt.Printf("%x\n", buf[2])
		fmt.Printf("%x\n", buf[3])

		_, err = s.Write(buf)
		checkerr(err)
		//log.Println(string(buf[:n]))

		//nm820_StatePara := &NM820_StatePara{}
		//nm820_StatePara.GDay = buf[0]
		//nm820_StatePara.Year = buf[1]
		//nm820_StatePara.Month = buf[2]
		//nm820_StatePara.Day = buf[3]
		//fmt.Printf("%x\n", buf[0:1])
		//fmt.Printf("%x\n A", nm820_StatePara.Day)
		//fmt.Printf("%x\n B", nm820_StatePara.Month)
	}
}

//错误检查
func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
