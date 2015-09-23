package nm820

import (
	"errors"
)

/*==================温度曲线===========================
		发送命令：8A 9B 00 01 05 00 【00 da】地址 50 长度
		返回：92个字节，从第10个数起来
		依赖：1.nm820_main.go中的串口发送
======================================================*/
//单日的
type NM820_WenduCurve struct {
	Day    []uint16
	Target []float32
	Heat   []float32
	Cool   []float32
}

//	发送命令：8A 9B 00 01 05 00 【00 da】地址 50 长度
//	返回：92个字节，从第10个数起来
func (nm *NM820_WenduCurve) addData() error {
	cmd1 := []byte{0x8A, 0x9B, 0x00, 0x01, 0x05, 0x00, 0x00, 0xda, 0x50}
	cmd2 := append(cmd1, sumCheck(cmd1))

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 92 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b := <-chanRb   //类型byte[20]
	<-chanSerialBusy

	//判断校验和是否一样
	if sumCheck(b[0:91]) != b[91] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	for i := 0; i < 10; i++ {
		nm.Day = append(nm.Day, twobyte_to_uint16(b[10+i*8], b[9+i*8]))
		nm.Target = append(nm.Target, float32(twobyte_to_uint16(b[12+i*8], b[11+i*8]))/10)
		nm.Heat = append(nm.Heat, float32(twobyte_to_uint16(b[14+i*8], b[13+i*8]))/10)
		nm.Cool = append(nm.Cool, float32(twobyte_to_uint16(b[16+i*8], b[15+i*8]))/10)
	}
	return nil
}
