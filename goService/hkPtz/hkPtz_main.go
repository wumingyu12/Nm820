package hkPtz

import (
	"encoding/xml"
	"github.com/gorilla/mux" //路由库
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//对应于hkcamConfig.xml的结构体，用于实现根据xml来初始化一些参数如ip，用户密码等等
type XmlHkcam struct { //顶级的
	Hkcameras []Camera `xml:"camera"`
}
type Camera struct {
	Ip       string `xml:"ip,attr"`
	User     string `xml:"user,attr"`
	Password string `xml:"password,attr"` //xml的属性
	//Interests Interests `xml:"interests"`
}

var xmlimfor XmlHkcam //用于存储得到的xml信息

/*========================================================================================
	init函数用于初始化
	作用：
		1.通过读取./myconfig/hkcamConfig.xml配置文件来初始化resetful命令要用到的一些参数
=============================================================================================*/
func init() {
	//打开配置文件
	content, err := ioutil.ReadFile("./myconfig/hkcamConfig.xml")
	if err != nil {
		log.Fatal(err)
	}
	//将配置文件的信息放到xmlimfor里面
	err = xml.Unmarshal(content, &xmlimfor)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(xmlimfor.Hkcameras)
}

/*=======================================resetful===============================================
	请求：GET ，前端将控制指令用get请求过去，后端再用put请求发送
	camNum 代表第几个摄像头，不同摄像头对应不同ip
	"/resetful/hkPtz/Continuous/{camNum}/{mode}/{speed}"
		camNum:代表控制第几个摄像头
			1 --摄像头1 对应
			2 --摄像头2 对应
		mode：
			panleft
			panright(左右)，
			tiltup
			tiltdown（上下），
			zoomfar
			zoomnear（远近），
			stop（直接停止）
			stopZoom (停止变焦)
			默认速度都为60，连续运动
		speed:
			运动的速度
			如果是用stop的话，这个speed就随便一个数
	作用：控制海康威视摄像头用isapi来实现上下左右远近连续运动
	返回：
	依赖：
		1.xmlimfor的初始化，里面有从xml读到的一些参数

=========================================================================================*/
func Continuous(w http.ResponseWriter, r *http.Request) {
	//从resetful中获取一些参数
	vars := mux.Vars(r)      //r为*http.Request
	camNum := vars["camNum"] //要移动的摄像头号数
	mode := vars["mode"]     //要移动的方式
	speed := vars["speed"]   //运动的速度

	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求控制摄像头--%s，运动模式为：%s\n", camNum, mode)

	//根据参数不同获取对应摄像头的ip地址,登录账户与密码
	var camip string       //摄像头的ip地址
	var camuser string     //摄像头的登录用户名 如admin
	var campassword string //摄像头的登录密码 wumingyu12
	switch camNum {
	case "1":
		camip = xmlimfor.Hkcameras[0].Ip //根据xml配置文件的信息获取
		camuser = xmlimfor.Hkcameras[0].User
		campassword = xmlimfor.Hkcameras[0].Password
	case "2":
		camip = xmlimfor.Hkcameras[1].Ip
		camuser = xmlimfor.Hkcameras[1].User
		campassword = xmlimfor.Hkcameras[1].Password
	default:
		camip = xmlimfor.Hkcameras[0].Ip
		camuser = xmlimfor.Hkcameras[0].User
		campassword = xmlimfor.Hkcameras[0].Password
	}

	//定义要put过去的xml命令,pan左右，tilt上下，zoom远近
	var sendstring string
	switch mode { //里面的60为移动的速度
	case "panleft": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>" + speed + "</pan><tilt>0</tilt></PTZData>"
	case "panright": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>-" + speed + "</pan><tilt>0</tilt></PTZData>"
	case "tiltup": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>" + speed + "</tilt></PTZData>"
	case "tiltdown": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>-" + speed + "</tilt></PTZData>"
	case "zoomfar":
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><zoom>" + speed + "</zoom></PTZData>"
	case "zoomnear":
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><zoom>-" + speed + "</zoom></PTZData>"
	case "stop": //停止左右上下运动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>0</tilt></PTZData>"
	case "stopzoom":
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><zoom>0</zoom></PTZData>"
	}
	sendbody := strings.NewReader(sendstring) //put要发送的内容根据mode而定

	client := &http.Client{}
	//要put的地址与命令前面的账号密码是为了防止401错误，验证身份
	//http://admin:wumingyu12@10.33.51.188/ISAPI/PTZCtrl/channels/1/continuous
	req, err := http.NewRequest("PUT", "http://"+camuser+":"+campassword+"@"+camip+"/ISAPI/PTZCtrl/channels/1/continuous", sendbody)
	if err != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	//将结果发送到前端
	w.Write(body)
	log.Printf("结束--resetful请求控制摄像头--%s，运动模式为：%s\n", camNum, mode)
}
