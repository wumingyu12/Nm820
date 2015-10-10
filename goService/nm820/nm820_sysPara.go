package nm820

import (
	"errors"
	"log"
)

/*==================通风等级 ===================================================
		函数：addData
		发送命令：8A 9B 00 01 05 00 【01 2a】地址 3c 长度  1120 480十进制
		命令1：8A 9B 00 01 05 00 【04 60】地址 f0 长度  1120 240 地址04为高位
		命令2: 8A 9B 00 01 05 00 【05 50】地址 f0 长度  1360 240
		因为长度480超出了数据位的长度1byte所以我们分成两条命令来读取，每条读取10条数据
		返回：每条命令返回252个字节，从第10个数起来
		依赖：1.nm820_main.go中的串口发送

		函数：uploadData 将数据更新到nm820
		发送命令：8a 9b 00 01 【05+data的长度】【01 写控制码 】【04 60 起始地址】【len(data)】【data】【CS】校验和
		示范      8a 9b 00 01   07               01              04 60             02           07 01
		更新1-10号通风等级：    5+240=245
		          8a 9b 00 01   f5               01              04 60             f0           data
==================================================================================*/

type NM820_WindTable struct {
	//注意On到DTemp事实上得到的是2位byte，但为了转换所以才用了float32
	On         float32 `json:",string"` //开时间x10 `json:",string"`是因为在前端会用xeditable来修改值，修改后为string类型
	Off        float32 `json:",string"` //关时间x10 时间分钟 如十进制110 代表11分钟
	DTemp      float32 `json:",string"` //温差x10
	SideWindow uint16  `json:",string"` //侧风窗 0-100
	Curtain    uint16  `json:",string"` //幕帘 0-100
	VSFan      uint16  `json:",string"` //变速风机 0-100
	Roller1    uint16  `json:",string"` //卷帘1-4
	Roller2    uint16  `json:",string"`
	Roller3    uint16  `json:",string"`
	Roller4    uint16  `json:",string"`
	Fan        uint32  //如1234 表示风机1234

}

type NM820_WindTables struct {
	WindTables []NM820_WindTable
}

//从820得到数据
func (nm *NM820_WindTables) addData() error {

	//----------------------------------------------------------------
	//------------------获取1-10号数组-------------------------------
	//-----------------------------------------------------------------
	cmd1 := []byte{0x8A, 0x9B, 0x00, 0x01, 0x05, 0x00, 0x04, 0x60, 0xf0}
	cmd2 := append(cmd1, sumCheck(cmd1))

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 252 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b := <-chanRb    //类型byte[20]
	<-chanSerialBusy

	//判断校验和是否一样
	if sumCheck(b[0:251]) != b[251] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	//将有效数组保存
	var tem []byte
	tem = b[9:249]

	//----------------------------------------------------------------
	//------------------获取11-20号数组-------------------------------
	//-----------------------------------------------------------------
	cmd1 = []byte{0x8A, 0x9B, 0x00, 0x01, 0x05, 0x00, 0x05, 0x50, 0xf0}
	cmd2 = append(cmd1, sumCheck(cmd1))

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 252 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b1 := <-chanRb   //类型byte[20]
	<-chanSerialBusy

	//判断校验和是否一样
	if sumCheck(b1[0:251]) != b1[251] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	//将有效数组保存
	tem = append(tem, b1[9:249]...)

	for i := 0; i < 20; i++ {
		n := NM820_WindTable{}
		n.On = float32(twobyte_to_uint16(tem[1+i*24], tem[0+i*24])) / 10
		n.Off = float32(twobyte_to_uint16(tem[3+i*24], tem[2+i*24])) / 10
		n.DTemp = float32(twobyte_to_uint16(tem[5+i*24], tem[4+i*24])) / 10
		n.SideWindow = twobyte_to_uint16(tem[7+i*24], tem[6+i*24])
		n.Curtain = twobyte_to_uint16(tem[9+i*24], tem[8+i*24])
		n.VSFan = twobyte_to_uint16(tem[11+i*24], tem[10+i*24])
		n.Roller1 = twobyte_to_uint16(tem[13+i*24], tem[12+i*24])
		n.Roller2 = twobyte_to_uint16(tem[15+i*24], tem[14+i*24])
		n.Roller3 = twobyte_to_uint16(tem[17+i*24], tem[16+i*24])
		n.Roller4 = twobyte_to_uint16(tem[19+i*24], tem[18+i*24])
		//n.Fan = tem[20+i*24] //默认是下面的方式的用了4个字节，但我们其实就一个字节数据
		n.Fan = twobyte_to_uint32(tem[23+i*24], tem[22+i*24], tem[21+i*24], tem[20+i*24])
		nm.WindTables = append(nm.WindTables, n)
	}
	return nil
}

//将数据组包上传到nm820中
func (nm *NM820_WindTables) uploadData() error {

	databyte1 := []byte{} //用来得到要发送的字节，注意是小端
	//先将结构体序列化为byte，长度为480
	for i := 0; i < 20; i++ {
		bh, bl := uint16_to_twobyte(uint16(nm.WindTables[i].On * 10))
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(uint16(nm.WindTables[i].Off * 10))
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(uint16(nm.WindTables[i].DTemp * 10))
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].SideWindow)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].Curtain)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].VSFan)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].Roller1)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].Roller2)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].Roller3)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.WindTables[i].Roller4)
		databyte1 = append(databyte1, bl) //先低位
		databyte1 = append(databyte1, bh)
		bhh, bh, bl, bll := uint32_to_fourbyte(nm.WindTables[i].Fan)
		databyte1 = append(databyte1, bll) //先低位
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh) //先低位
		databyte1 = append(databyte1, bhh)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("数据：%x\n", databyte1)
	//组包发送更新1-10号
	cmd1 := []byte{0x8a, 0x9b, 0x00, 0x01, 0xf5, 0x01, 0x04, 0x60, 0xf0}
	cmd2 := append(cmd1, databyte1[0:240]...)
	cmd3 := append(cmd2, sumCheck(cmd2))

	//发送命令，接收到的回复为
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd3
	chanRbNum <- 11 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	<-chanRb        //类型byte[11] 8a 9b 00 01 06 81 00 00 01 00 ae
	<-chanSerialBusy

	//组包发送更新10-20号
	cmd1 = []byte{0x8a, 0x9b, 0x00, 0x01, 0xf5, 0x01, 0x05, 0x50, 0xf0}
	cmd2 = append(cmd1, databyte1[240:480]...)
	cmd3 = append(cmd2, sumCheck(cmd2))

	//发送命令，接收到的回复为
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd3
	chanRbNum <- 11 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	<-chanRb        //类型byte[11] 8a 9b 00 01 06 81 00 00 01 00 ae
	<-chanSerialBusy

	return nil
}

/*==================温度曲线===========================
		发送命令：8A 9B 00 01 05 00 【00 da】地址218 50 长度
		返回：92个字节，从第10个数起来 80个字节数据
		依赖：1.nm820_main.go中的串口发送

		更新：
		8a 9b 00 01 【55 十进制85长度】【01 控制码】 【00 da 地址】 【50 数据长度】 【data】
======================================================*/
//单日的
type NM820_WenduCurve struct {
	Day    []uint16
	Target []float32
	Heat   []float32
	Cool   []float32
}

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

//将数据组包上传到nm820中
func (nm *NM820_WenduCurve) uploadData() error {
	databyte1 := []byte{} //用来得到要发送的字节，注意是小端
	//组数据包
	for i := 0; i < 10; i++ {
		bh, bl := uint16_to_twobyte(nm.Day[i])
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(uint16(nm.Target[i] * 10))
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(uint16(nm.Heat[i] * 10))
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(uint16(nm.Cool[i] * 10))
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
	}

	//组包添加包头和校验码
	cmd1 := []byte{0x8a, 0x9b, 0x00, 0x01, 0x55, 0x01, 0x00, 0xda, 0x50}
	cmd2 := append(cmd1, databyte1...)
	cmd3 := append(cmd2, sumCheck(cmd2))

	//发送更新命令
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd3
	chanRbNum <- 11 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	<-chanRb        //类型byte[11] 8a 9b 00 01 06 81 00 00 01 00 ae
	<-chanSerialBusy

	return nil
}

/*==================最大最小通风等级曲线===================================================
		发送命令：8A 9B 00 01 05 00 【01 2a】地址 3c 长度  298 60十进制
		返回：72个字节，从第10个数起来
		依赖：1.nm820_main.go中的串口发送

		更新数据命令：
		8a 9b 00 01 【41 十进制65长度】【01 控制码】 【01 2a 地址】 【3c 数据长度】 【data】
==================================================================================*/

type NM820_WindLevel struct {
	Day []uint16
	Min []uint16
	Max []uint16
}

func (nm *NM820_WindLevel) addData() error {
	cmd1 := []byte{0x8A, 0x9B, 0x00, 0x01, 0x05, 0x00, 0x01, 0x2a, 0x3c}
	cmd2 := append(cmd1, sumCheck(cmd1))

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 72 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b := <-chanRb   //类型byte[20]
	<-chanSerialBusy

	//判断校验和是否一样
	if sumCheck(b[0:71]) != b[71] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		return errors.New("sum check is wrong!!")
	}
	for i := 0; i < 10; i++ {
		nm.Day = append(nm.Day, twobyte_to_uint16(b[10+i*6], b[9+i*6]))
		nm.Min = append(nm.Min, twobyte_to_uint16(b[12+i*6], b[11+i*6]))
		nm.Max = append(nm.Max, twobyte_to_uint16(b[14+i*6], b[13+i*6]))
	}
	return nil
}

//函数更新到nm820将数据
func (nm *NM820_WindLevel) uploadData() error {
	databyte1 := []byte{} //用来得到要发送的字节，注意是小端
	//组数据包
	for i := 0; i < 10; i++ {
		bh, bl := uint16_to_twobyte(nm.Day[i])
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.Min[i])
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
		bh, bl = uint16_to_twobyte(nm.Max[i])
		databyte1 = append(databyte1, bl)
		databyte1 = append(databyte1, bh)
	}

	//组包添加包头和校验码
	cmd1 := []byte{0x8a, 0x9b, 0x00, 0x01, 0x41, 0x01, 0x01, 0x2a, 0x3c}
	cmd2 := append(cmd1, databyte1...)
	cmd3 := append(cmd2, sumCheck(cmd2))

	//发送更新命令
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd3
	chanRbNum <- 11 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	<-chanRb        //类型byte[11] 8a 9b 00 01 06 81 00 00 01 00 ae
	<-chanSerialBusy

	return nil
}
