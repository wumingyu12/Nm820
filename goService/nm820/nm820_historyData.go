package nm820

import (
	"log"
)

//历史温度，每一条命令发送回来的结构体
/*
type HisTemp struct {
	Day uint16 //日龄
	Max uint16 //最大温度
	Min uint16 //最小值
	Avg uint16 //平均值
}
*/

//用于resetful返回的结构体,以当前日龄为准，向前存储30个数据
//这个结构体通用的，可以存储温度，湿度，光照，氨气等等都可以，下面的添加方法不同就会不同
type NM820_History30 struct {
	Days []uint16  //天龄列表
	Maxs []float32 //最大温度
	Mins []float32 //最小温度
	Avgs []float32
}

/*==============================================================
	依赖：
		1.nm820_main.go  中的4个串口发送通道，及其运行线程
		2.nm820_main.go  中的checkerr，sumCheck函数
		3.				 中的uint16_to_twobyte
		4.               中的twobyte_to_uint16
	参数：
		1.baseDay：基准日龄，以这个日龄为基准，将前30天的数据存入hs里面
		2.sensorType：要存储进去的数据类型有以下可选
			"Tem"---温度
			"Humi"---湿度
			"NH3"---氨气
			"Light"--光照
	注意：1.如要读日龄无记录，返回当日的记录 ，日龄为0
		  2.如果日龄恰好是当前日，也会返回日龄0
	遗留问题：
		  1.每次要请求历史数据都需要
	返回的包：
		8A 9B 00 01 0F 82 00 02 0A 00 20 【Days 20 00】【Maxs 45 01】【Mins FA 00】【Avgs FA 00】 3D
================================================================*/
func (hs *NM820_History30) addData(baseDay uint16, sensorType string) {
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	var cmd1 []byte

	//当要存储不同数据，要发送的头是不同的
	switch sensorType {
	case "Tem":
		//完整的命令后面还要加天龄00,20，sumcheck
		cmd1 = []byte{0x8A, 0x9B, 0x00, 0x01, 0x07, 0x02, 0x00, 0x02, 0x02}
	case "Humi":
		cmd1 = []byte{0x8A, 0x9B, 0x00, 0x01, 0x07, 0x02, 0x00, 0x03, 0x02}
	case "NH3":
		cmd1 = []byte{0x8A, 0x9B, 0x00, 0x01, 0x07, 0x02, 0x00, 0x05, 0x02}
	case "Light":
		cmd1 = []byte{0x8A, 0x9B, 0x00, 0x01, 0x07, 0x02, 0x00, 0x04, 0x02}
	}
	//循环读取以基准日期前30日的数据
	for day := baseDay - 29; day <= baseDay; day++ {
		if day < 0 { //如果baseday小于0，不动作
			continue
		}
		bh, bl := uint16_to_twobyte(day)
		cmd2 := append(cmd1, bh, bl) //将日龄添加到发送命令行中
		cmd3 := append(cmd2, sumCheck(cmd2))

		chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
		chanWb <- cmd3
		chanRbNum <- 20 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
		rec := <-chanRb //类型byte[20]
		<-chanSerialBusy

		//日龄,解析得到的包，并加到数组中
		hs.Days = append(hs.Days, twobyte_to_uint16(rec[12], rec[11]))
		hs.Maxs = append(hs.Maxs, float32(twobyte_to_uint16(rec[14], rec[13]))/10)
		hs.Mins = append(hs.Mins, float32(twobyte_to_uint16(rec[16], rec[15]))/10)
		hs.Avgs = append(hs.Avgs, float32(twobyte_to_uint16(rec[18], rec[17]))/10)
		//log.Printf("days:%d-%d-%d-%d", days_buf, maxs_buf, mins_buf, avgs_buf)
	}
}
