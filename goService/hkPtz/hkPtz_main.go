package hkPtz

import (
	"github.com/gorilla/mux" //路由库
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*=======================================resetful===============================================
	请求：GET /mode--panleft panright(左右)，tiltup tiltdown（上下），zoomfar zoomnear（远近），stop（直接停止），默认速度都为60，连续运动
	camNum 代表第几个摄像头，不同摄像头对应不同ip
	"/resetful/hkPtz/Continuous/{camNum}/{mode}"
	作用：控制海康威视摄像头用isapi来实现上下左右远近连续运动
	返回：
	依赖的函数：

=========================================================================================*/
func Continuous(w http.ResponseWriter, r *http.Request) {
	//先判断要请求的历史数据类型
	vars := mux.Vars(r)      //r为*http.Request
	camNum := vars["camNum"] //要移动的摄像头号数
	mode := vars["mode"]     //要移动的方式

	log.SetFlags(log.Lshortfile | log.LstdFlags) //设置打印时添加上所在文件，行数
	log.Printf("开始--resetful请求控制摄像头--%s，运动模式为：%s\n", camNum, mode)

	var sendstring string
	switch mode { //里面的60为移动的速度
	case "panleft": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>60</pan><tilt>0</tilt></PTZData>"
	case "panright": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>-60</pan><tilt>0</tilt></PTZData>"
	case "tiltup": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>60</tilt></PTZData>"
	case "tiltdown": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>-60</tilt></PTZData>"
	case "stop": //如果是向左移动
		sendstring = "<?xml version='1.0' encoding='UTF-8'?><PTZData><pan>0</pan><tilt>0</tilt></PTZData>"
	}
	sendbody := strings.NewReader(sendstring) //put要发送的内容根据mode而定
	//先得到当前的日龄
	client := &http.Client{}

	req, err := http.NewRequest("PUT", "http://admin:wumingyu12@10.33.51.187/ISAPI/PTZCtrl/channels/1/continuous", sendbody)
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
	log.Println(string(body))

	log.Println("结束--resetful请求控制摄像头--%s，运动模式为：%s\n", camNum, mode)
}
