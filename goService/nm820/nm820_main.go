package nm820

//========================================
//==========通用函数部分==================
//本文件下的go文件，按照协议命令一个go文件组织
//=======================================

import (
	//"errors"
	"fmt"
	//"github.com/huin/goserial" //引入串口库
	"encoding/json"
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
给resetful调用的函数，外部可访问,resetful
***************************************************************************/

//=====================================================
//得到nm820的状态，主体函数在nm820_statePara.go里面
//发送的json根据结构图NM820_StatePara
//  请求url：/resetful/nm820/GetState
//  正常返回一个json：
//				{"GDay":34,"Year":2000,"Month":1,"Day":13,"Hour":23,"Min":59,"Sec":9,"TemAvg":242,
//				"Tem_1to5":[248,244,-1,-1,32767],"HumiAvg":863,"Hmi_1to2":[727,1000],"NH3":65535,"Light":65535,
//				"FanLevel":7,"Pos_SideWin":60,"Pos_Curtain":0,"Pos_Roller":[0,0,0,0],
//				"RelayType":[1,2,3,4,5,6,7,8,9,10,11,13,27,15,16,17,18,28,29,31],
//				"RelayState":[1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1]}
//  错误的返回(字符串)：
//				返回"error1:cant open serial port." 情况1不能打开串口，
//				返回"error"
//==========================================================================

func GetState(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Println("开始--resetful请求nm820获取状态。")
	para := &NM820_StatePara{}
	para.sendCmd()             //发送串口协议命令
	err := para.reflashValue() //用返回的数据更新结构体
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
