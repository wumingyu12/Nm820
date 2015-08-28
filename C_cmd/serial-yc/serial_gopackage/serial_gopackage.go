package main

import (
	"fmt"
	"github.com/huin/goserial" //引入串口库
	"log"
)

//用到的常量
const (
	CON_PORTNAME = "/dev/ttySAC1" //要连接的串口名字，在window下可以用“COM1”
	CON_BAUD     = 115200         //要连接的串口波特率
)

//NM820状态结构体
type NM820_StatePara struct {
	GDay  []byte //日龄 0-1
	Year  []byte //当前年 2-3
	Month []byte //月	4-5
	Day   []byte //日	6-7
}

func main() {
	c := &goserial.Config{Name: CON_PORTNAME, Baud: CON_BAUD} //以波特率和串口名打开
	s, err := goserial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	cmd := []byte{0xff, 0x11, 0x44, 0xff}
	//n, err := s.Write([]byte("test,hello"))
	n, err := s.Write(cmd)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 4)
	for {
		n, err = s.Read(buf) //这里要注意只有完全读满buf才会完成这一步
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%x\n", buf[:n])

		//log.Println(string(buf[:n]))

		//nm820_StatePara := new(NM820_StatePara)
		//nm820_StatePara.GDay = buf[0:1]
		//nm820_StatePara.Year = buf[2:3]
		//nm820_StatePara.Month = buf[4:5]
		//nm820_StatePara.Day = buf[6:7]
		fmt.Printf("%x\n", buf[0:1])
		//fmt.Printf("%x\n", nm820_StatePara.Day)
		//fmt.Printf("%x\n", nm820_StatePara.Month)
	}
}
