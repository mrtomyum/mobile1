package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/mrtomyum/mobile1/model"
	"fmt"
	"log"
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
	h := model.Host{}
	go h.Run()
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
	r.GET("/ws", func(c *gin.Context){
		Server(c.Writer, c.Request)
	})
	return r
}

func Server(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start Server...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		//http.NotFound(w, r)
		return
	}
	defer conn.Close()

	fmt.Println("start New Client instance...")
	clientName := r.Header.Get("Name") // "web" or "dev"
	c := &model.Client{
		Conn: conn,
		Send: make(chan *model.Message),
		Name: clientName,
	}
	switch c.Name {
	case "web":
		log.Println("WebClient:", c.Name)
		model.H.WebClient <- c
	case "dev":
		log.Println("DevClient:", c.Name)
		model.H.DevClient <- c
	default:
		log.Println("Default: No Name Provide")

	}
	go c.Write()
	c.Read()
}