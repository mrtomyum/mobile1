package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
)

func main() {
	r := gin.Default()
	app := router(r)

}

func router(r *gin.Engine) *gin.Engine{
	r.LoadHTMLGlob("view/**/*.html")
	r.Static("/", "./view/html")
	r.Static("/js", "./view/public/js")
	r.Static("/css", "./view/public/css")
	r.Static("/img", "./view/public/img")
	r.Static("/json", "./view/public/json")
	r.Use(static.Serve("/", static.LocalFile("view", true)))
	return r
}