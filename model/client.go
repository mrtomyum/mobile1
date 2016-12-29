package model

import (
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
	"time"
)

type Client struct {
	Conn *websocket.Conn
	Send chan *Message
	Name string
}

func (c *Client) Read() {
	defer func() {
		c.Conn.Close()
	}()
	m := &Message{}
	for {
		err := c.Conn.ReadJSON(&m)
		if err != nil {
			log.Println("Error ReadJSON()")
		}
		switch c.Name {
		case "web":
			switch m.Payload.Command {
			case "onhand":
				H.Onhand(c)
			case "cancel":
				H.Cancel(c)
			}
		case "dev":
			switch m.Device {
			case "coin_hopper":
			case "coin_acc":
			case "bill_acc":
			case "printer":
			}
			return
		default:
			fmt.Println("Message:", m)
			time.Sleep(200 * time.Millisecond)
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

