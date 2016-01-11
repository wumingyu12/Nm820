package nm820

import (
	//"errors"
	"log"
	"reflect" //通过反射来初始化结构体
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
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
		//return errors.New("sum check is wrong!!")
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
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
		//return errors.New("sum check is wrong!!")
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
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
		//return errors.New("sum check is wrong!!")
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
		//return errors.New("sum check is wrong!!")
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
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

/*==================系统参数配置值表===================================================
		发送命令：8A 9B 00 01 05 00 【00 48】地址 88 长度  72 136十进制  结构体实际134个字节，因为内存是32个bit一起的就是8个字节
		返回：136+12=148个字节，从第10个数起来
		依赖：1.nm820_main.go中的串口发送

		更新数据命令：
		8a 9b 00 01 【8d 十进制136+5长度】【01 控制码】 【00 48 地址】 【88 数据长度】 【data】
		注意：
		 float32 ：实际值*10 原结构还是uint16为了可以除10改为float32。但还是用2个字节实例化
==================================================================================*/
type NM820_SysVal struct {
	DTemp_Des          float32 //*目标温度差 x.x 表7.2 *10
	Temp_Des           float32 //自设模式目标温度  表7.3 *10
	Temp_Heat          float32 //自设模式加热温度  表7.4 *10
	Temp_Cool          float32 //自设模式制冷温度  表7.5 *10
	Delay_FanUp        uint16  //*通风升级延迟xxx.x 分
	Delay_FanDown      uint16  //*通风降级延迟xxx.x 分
	Delay_HumiUp       uint16  //高湿润通风延迟
	Delay_HumiRst      uint16  //高湿润通风保持
	Time_LightChange   uint16  //光照变化时间xxx.x 分
	Time_Channel       uint16  //*隧道最短时间xxx.x 分
	Level_ChannelStart uint16  //*隧道开始级别 1-20
	Time_SideWin_ON    uint16  //*侧风窗开时间 xxx 秒
	Time_SideWin_OFF   uint16  //*侧风窗关时间 xxx 秒
	Time_Curtain_ON    uint16  //*幕帘开时间 xxx 秒
	Time_Curtain_OFF   uint16  //*幕帘关时间 xxx 秒
	//------------------------------------------------------
	Time_Roller_ON_1  uint16 //4         //*卷帘1-4开时间 xxx 秒
	Time_Roller_ON_2  uint16
	Time_Roller_ON_3  uint16
	Time_Roller_ON_4  uint16
	Time_Roller_OFF_1 uint16 //4         //*卷帘1-4关时间 xxx 秒
	Time_Roller_OFF_2 uint16
	Time_Roller_OFF_3 uint16
	Time_Roller_OFF_4 uint16
	//------------------------------------------------------------
	Delay_Return     uint16  //*回流阀延迟 xxx 秒
	DTemp_Curtain    uint16  //*湿帘水泵启动温差 x.x
	DTemp_Heater     uint16  //*加热器温差 x.x
	LevelUp_Max      uint16  //*最大增加的通风级别 xx 0-20
	Is_LevelUp_Humi  uint16  //*高湿是否增加通风级别 0-1 是否
	Is_LevelUp_NH3   uint16  //*高氨气是否增加通风级别 0-1 是否
	Is_MinWind_Auto  uint16  //*最小通风量自动 0-1 是否
	Is_AutoSpray     uint16  //*是否自动喷雾
	Alarm_Temp_Max   float32 //*报警最高温度 表6.1  实际值*10
	AlarmR_Temp_Max  float32 //*高温报警恢复回差 表6.2
	Alarm_Temp_Min   float32 //*报警最低温度 表6.3 实际值*10
	AlarmR_Temp_Min  float32 //*低温报警恢复回差 表6.4 实际值*10
	Alarm_dTemp_Max  float32 //*报警最高温差 表6.5 实际值*10
	AlarmR_dTemp_Max float32 //*高温差报警恢复回差 表6.6 实际值*10
	Alarm_dTemp_Min  float32 //*报警最低温差 表6.7 实际值*10
	AlarmR_dTemp_Min float32 //*低温差报警恢复回差 表6.8 实际值*10
	Alarm_Humi_Max   float32 //报警最高湿度 表6.9 实际值*10
	AlarmR_Humi_Max  float32 //报警最高湿度恢复回差 表6.10 实际值*10
	Alarm_Humi_Min   float32 //报警最低湿度  6.11 实际值*10
	AlarmR_Humi_Min  float32 //报警最低湿度恢复回差 6.12 实际值*10
	Alarm_Light_Min  uint16  //报警最低光照 6.13
	AlarmR_Light_Min uint16  //报警最低光照恢复回差 6.14
	Alarm_NH3_Max    uint16  //报警最高氨气 6.15
	AlarmR_NH3_Max   uint16  //报警最高氨气回复回差 6.16
	Mode             uint16  //*温度控制模式 2--曲线+通风 0--自设 1--曲线 表7.1
	Password         uint32  //*登录密码
	//-----------------------------------------------------
	Alarm_TempSF_1 uint16 //5    //温度探头故障是否报警
	Alarm_TempSF_2 uint16
	Alarm_TempSF_3 uint16
	Alarm_TempSF_4 uint16
	Alarm_TempSF_5 uint16
	//------------------------------------------------------------
	Alarm_HumiSF_1 uint16 //2    //湿度探头故障是否报警
	Alarm_HumiSF_2 uint16
	//-------------------------------------------------------------
	Relay_1  uint8 //20       //继电器编码
	Relay_2  uint8
	Relay_3  uint8
	Relay_4  uint8
	Relay_5  uint8
	Relay_6  uint8
	Relay_7  uint8
	Relay_8  uint8
	Relay_9  uint8
	Relay_10 uint8
	Relay_11 uint8
	Relay_12 uint8
	Relay_13 uint8
	Relay_14 uint8
	Relay_15 uint8
	Relay_16 uint8
	Relay_17 uint8
	Relay_18 uint8
	Relay_19 uint8
	Relay_20 uint8
}

func (nm *NM820_SysVal) addData() error {
	cmd1 := []byte{0x8A, 0x9B, 0x00, 0x01, 0x05, 0x00, 0x00, 0x48, 0x88}
	cmd2 := append(cmd1, sumCheck(cmd1))

	//串口发送
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd2
	chanRbNum <- 148 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	b := <-chanRb    //类型byte[20]
	<-chanSerialBusy

	if sumCheck(b[0:147]) != b[147] { //前面99个数的校验和是否等于最后一个校验位,b[0]--b[98]
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("sum check is wrong!!")
		//return errors.New("NM820_SysVal.addData() sum check is wrong!!")
	}

	//用反射的方式将byte数组赋值到结构体内
	var bpoint = 9                          //指代返回的byte数组的下标，随着每一次成功的赋值，其值也会增加,第一个数据位为b【9】
	refnm := reflect.ValueOf(nm).Elem()     //获取反射值的成员可设置
	for i := 0; i < refnm.NumField(); i++ { //i的数量代表该结构体的字段数
		f := refnm.Field(i)
		switch f.Kind() {
		case reflect.Uint16: //如果反射值为uint16值
			uint16value := twobyte_to_uint16(b[bpoint+1], b[bpoint]) //根据byte数组得到要设置的值
			bpoint = bpoint + 2                                      //因为是uint16的用来2个byte
			cansetvalue := reflect.ValueOf(uint16value)              //将得到的uint16变为可以用于reflect设置的类型
			f.Set(cansetvalue)

		case reflect.Uint32: //如果反射值为uint32值
			uint32value := twobyte_to_uint32(b[bpoint+3], b[bpoint+2], b[bpoint+1], b[bpoint]) //根据byte数组得到要设置的值
			bpoint = bpoint + 4                                                                //因为是uint16的用来2个byte
			cansetvalue := reflect.ValueOf(uint32value)                                        //将得到的uint16变为可以用于reflect设置的类型
			f.Set(cansetvalue)

		case reflect.Uint8: //如果反射值为uint8值
			uint8value := uint8(b[bpoint])             //根据byte数组得到要设置的值
			bpoint = bpoint + 1                        //因为是uint16的用来2个byte
			cansetvalue := reflect.ValueOf(uint8value) //将得到的uint16变为可以用于reflect设置的类型
			f.Set(cansetvalue)                         //如果有Setint就可以用

		case reflect.Float32: //请注意如果是Float32我们还是用2个字节来实例化，其本质为uint16我们为了可以显示为小数才这样赋值，比如温度*10
			float32value := float32(twobyte_to_uint16(b[bpoint+1], b[bpoint])) / 10 //根据byte数组得到要设置的值,比如都回来是温度233，我们要转为23.3
			bpoint = bpoint + 2                                                     //因为是uint16的用来2个byte
			cansetvalue := reflect.ValueOf(float32value)                            //将得到的uint16变为可以用于reflect设置的类型
			f.Set(cansetvalue)
		}
	}
	return nil
}

//更新参数表
func (nm *NM820_SysVal) uploadData() error {
	//用反射的方式将结构体的内容变为byte数组
	databyte1 := []byte{}                   //用来得到要发送的字节，注意是小端
	refnm := reflect.ValueOf(nm).Elem()     //获取反射值的成员可设置
	for i := 0; i < refnm.NumField(); i++ { //i的数量代表该结构体的字段数
		f := refnm.Field(i)
		switch f.Kind() {
		case reflect.Uint16: //如果反射值为uint16值
			bh, bl := uint16_to_twobyte(uint16(f.Uint()))
			databyte1 = append(databyte1, bl)
			databyte1 = append(databyte1, bh)

		case reflect.Uint32: //如果反射值为uint32值
			bhh, bh, bl, bll := uint32_to_fourbyte(uint32(f.Uint()))
			databyte1 = append(databyte1, bll) //先低位
			databyte1 = append(databyte1, bl)
			databyte1 = append(databyte1, bh) //先低位
			databyte1 = append(databyte1, bhh)

		case reflect.Uint8: //如果反射值为uint8值
			b8 := byte(f.Uint()) //根据byte数组得到要设置的值
			databyte1 = append(databyte1, b8)

		case reflect.Float32: //请注意如果是Float32我们还是用2个字节来实例化，其本质为uint16我们为了可以显示为小数才这样赋值，比如温度*10
			bh, bl := uint16_to_twobyte(uint16(f.Float() * 10))
			databyte1 = append(databyte1, bl)
			databyte1 = append(databyte1, bh)
		}
	}

	//组包添加包头和校验码
	cmd1 := []byte{0x8a, 0x9b, 0x00, 0x01, 0x8d, 0x01, 0x00, 0x48, 0x88}
	cmd2 := append(cmd1, databyte1...)
	cmd3 := append(cmd2, []byte{0xff, 0xff}...) //最后给补2个字节，这样发送的数据才是8的整数
	cmd4 := append(cmd3, sumCheck(cmd3))

	//发送更新命令
	chanSerialBusy <- 1 //为了其他地方使用串口时发送接受流程不被打断
	chanWb <- cmd4
	chanRbNum <- 11 //开启一次锁让进程发送一次命令,接收一次命令，接收字节数为100
	<-chanRb        //类型byte[11] 8a 9b 00 01 06 81 00 00 01 00 ae
	<-chanSerialBusy

	return nil
}
