package nm820

//========================================
//==========通用函数部分==================
//本文件下的go文件，按照协议命令一个go文件组织
//=======================================

import (
	//"errors"
	"fmt"
	//"github.com/huin/goserial" //引入串口库
	"log"
)

//用到的常量
const (
	con_PORTNAME = "/dev/ttySAC3" //要连接的串口名字，在window下可以用“COM1”
	con_BAUD     = 4800           //要连接的串口波特率
	//CON_ADDR_H   = 0x00           //通信地址高位
	//CON_ADDR_L   = 0x01           //通信地址低位，可以通过nm820的界面参看就是栏舍号
	//CMD []byte={}
)

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

/**************************************************************************
给resetful调用的函数，外部可访问
***************************************************************************/
//得到nm820的状态，主体函数在nm820_statePara.go里面

func GetState() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始nm820获取状态。")
	para := &NM820_StatePara{}
	para.sendCmd()
	err := para.reflashValue()
	checkerr(err)
	fmt.Printf("tem:%d\n", para.TemAvg)
}
