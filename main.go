package main

import (
	"./goService/nm820" //nm820的引用库
	"fmt"
	"github.com/gorilla/mux" //路由库
	"html/template"
	"log"
	"net/http"
)

//===========================
//===404.html============
//=====================
func NotFoundHandler(w http.ResponseWriter, r *http.Request) { //如果路由规则不符合没有注册的如/2333,/22ww等
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/view/index.html", http.StatusFound) //地址重定向
	}

	t, err := template.ParseFiles("frontWeb/view/static/404/404.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

//================================================================

func main() {
	//nm820.GetState()
	http.Handle("/view/", http.FileServer(http.Dir("frontWeb")))
	//view/xxx/xxx的文件在frontweb里面找
	http.Handle("/frame/", http.FileServer(http.Dir("frontWeb")))
	http.Handle("/my_css/", http.FileServer(http.Dir("frontWeb")))
	http.Handle("/js/", http.FileServer(http.Dir("frontWeb")))
	http.Handle("/testjson/", http.FileServer(http.Dir("frontWeb"))) //放测试用的json文件模拟restful
	//这里的handle当一个连接过来的时候都会多开一个wshandler
	//http.Handle("/ws", websocket.Handler(wshandler.WsHandler)) //响应了ws://127.0.0.1/ws的websocket

	//http.HandleFunc("/index", IndexHandler) //不用这个带控制器的路由导致带angular的index无法正常加载
	//http.HandleFunc("/login", login)

	mux_router := mux.NewRouter()               //用mux库做路由
	mux_router.HandleFunc("/", NotFoundHandler) //初始化Session管理器
	//resetful得到nm820的状态
	//注意http://10.33.51.186:2234/resetful/nm820/GetState/是匹配不了的最后面不能有/
	mux_router.HandleFunc("/resetful/nm820/GetState", nm820.GetState).Methods("GET")
	//得到历史温度
	mux_router.HandleFunc("/resetful/nm820/GetTempHistory", nm820.GetTempHistory).Methods("GET")

	http.Handle("/", mux_router) //这一句别忘了 否则前面的mux_router是不作用的
	fmt.Println("正在监听2234端口,main.go")
	//http.HandleFunc("/", NotFoundHandler) //当没有找到路径名字时，后面改为用mux库了
	err1 := http.ListenAndServe(":2234", nil)
	if err1 != nil {
		log.Fatal("ListenAndServe:", err1)
	}
}
