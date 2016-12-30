package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/mrtomyum/mobile1/model"
	"fmt"
	//"log"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	//h := model.Host{}
	go model.H.Run()
	r := gin.Default()
	app := router(r)
	app.Run(":8088")
	//app.RunTLS(
	//	":8088",
	//	"api.nava.work.crt",
	//	"nava.work.key",
	//)
}

func router(r *gin.Engine) *gin.Engine{
	//r.LoadHTMLGlob("view/**/*.html")
	//r.Static("/", "./view/html")
	r.Use(static.Serve("/", static.LocalFile("view", true)))
	r.GET("/web", func(c *gin.Context) {
		ServWeb(c.Writer, c.Request)
	})
	r.GET("/dev", func(c *gin.Context) {
		ServDev(c.Writer, c.Request)
	})
	return r
}

func ServWeb(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start ServWeb Websocket for Web...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("start New Web connection success...")
	c := &model.Client{
		Conn: conn,
		Send: make(chan *model.Message),
		Name: "web",
	}
	fmt.Println("WebClient:", c.Name)
	model.H.WebClient <- c
	fmt.Println("start go c.Write()")
	go c.Write()
	c.Read()
}

func ServDev(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start ServDev Websocket for Device...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Println("start New Device connection success...")
	//clientName := r.Header.Get("Name") // "web" or "dev"
	c := &model.Client{
		Conn: conn,
		Send: make(chan *model.Message),
		Name: "dev",
	}
	fmt.Println("DevClient:", c.Name)
	model.H.DevClient <- c
	go c.Write()
	c.Read()
}