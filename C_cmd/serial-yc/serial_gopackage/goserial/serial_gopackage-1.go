package main

import (
	//"fmt"
	"github.com/huin/goserial" //引入串口库
	"log"
	//"time"
)

//用到的常量
const (
	CON_PORTNAME = "/dev/ttySAC3" //要连接的串口名字，在window下可以用“COM1”
	CON_BAUD     = 4800           //要连接的串口波特率

)

func main() {
	c := &goserial.Config{
		Name: CON_PORTNAME,
		Baud: CON_BAUD,
		//ReadTimeout: time.Second * 5,
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
		Parity:   goserial.ParityNone,
	} //以波特率和串口名打开
	s, err := goserial.OpenPort(c) //打开串口
	defer s.Close()                //用完关闭
	checkerr(err)
	//发送的命令，定义

	//n, err := s.Write([]byte("test,hello"))
	//============将要发送的协议命令=============
	cmd := []byte{0x8a, 0x9b, 0x00, 0x01, 0x05, 0x02, 0x00, 0x09, 0x00, 0x36}
	//发送协议命令
	_, err = s.Write(cmd)
	checkerr(err)

	//接收协议命令
	buf := make([]byte, 100)
	s.Read(buf)
	log.Printf("接收：%x\n", buf)
}

//错误检查
func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
