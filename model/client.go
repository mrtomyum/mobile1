package model

import (
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
)

type Client struct {
	Conn *websocket.Conn
	Send chan Message
	Name string
}

func (c *Client) Read() {
	m := Message{}
	for {
		err := c.Conn.ReadJSON(&m)
		if err != nil {
			log.Println("Error ReadJSON()")
			break
		}
		switch m.Payload.Command{
		case "onhand":
			fmt.Println("Read Onhand Message...")
		case "cancel":
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case m, ok := <-c.Send:
			if !ok {
				c.Conn.WriteJSON(gin.H{"message": "Cannot send data"})
				return
			}
			fmt.Println("Client.Write():", m)
			c.Conn.WriteJSON(m)
		}
	}
}

