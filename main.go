package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/mrtomyum/mobile1/model"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	r := gin.Default()
	app := router(r)
	app.RunTLS(
		":8088",
		"api.nava.work.crt",
		"nava.work.key",
	)
}

func router(r *gin.Engine) *gin.Engine{
	//r.LoadHTMLGlob("view/**/*.html")
	//r.Static("/", "./view/html")
	//r.Static("/js", "./view/public/js")
	//r.Static("/css", "./view/public/css")
	//r.Static("/img", "./view/public/img")
	//r.Static("/json", "./view/public/json")
	r.Use(static.Serve("/", static.LocalFile("view", true)))
	r.GET("/ws", func(c *gin.Context){
		wsServ(c.Writer, c.Request)
	})
	return r
}

func wsServ(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer conn.Close()
	clientName := r.Header.Get("Name")
	c := model.Client{
		Conn: conn,
		Send: make(chan model.Message),
		Name: clientName,
	}
	go c.Read()
	go c.Write()
}