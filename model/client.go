package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Ws   *websocket.Conn
	Send chan *Message
	Name string
	Msg  *Message
}

func (c *Client) Read() {
	defer func() {
		c.Ws.Close()
	}()
	m := &Message{}
	for {
		err := c.Ws.ReadJSON(&m)
		if err != nil {
			log.Println("Error ReadJSON():", err)
			return
		}
		c.Msg = m
		switch c.Name {
		case "web":
			fmt.Println("Message from web")
			c.WebEvent()

		case "dev":
			fmt.Println("Message from dev")
			c.DevEvent()
		default:
			fmt.Println("Case default: Message==>", m)
			m.Type = "response"
			m.Data = "Hello"
			c.Send <- m
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Ws.Close()
	}()
	for {
		select {
		case m, ok := <-c.Send:
			if !ok {
				c.Ws.WriteJSON(gin.H{"message": "Cannot send data"})
				return
			}
			fmt.Println("Client.Write():", m)
			c.Ws.WriteJSON(m)
		}
	}
}

func (c *Client) WebEvent() {
	switch c.Msg.Command {
	case "onhand":
		//H.GetEscrow <- c
		H.Onhand(c)
	case "cancel":
		//H.CancelOrder <- c
		H.Cancel(c)
	}
}

func (c *Client) DevEvent() {
	switch c.Msg.Device {
	case "coin_hopper":
		switch c.Msg.Command {
		case "machine_id": //ร้องขอหมายเลข Serial Number ของ อุปกรณ์ Coins Hopper
		case "status": // ร้องขอสถานะต่างๆของอุปกรณ์
		case "cash_amount": // ร้องขอจานวนเงินคงเหลือใน Coins Hopper
		case "coin_count": // ร้องขอจานวนเงินเหรียญคงเหลือใน Coins Hopper
		case "set_coin_count": // ตั้งค่าจำนวนเงินคงเหลือใน Coins Hopper
		case "payout_by_cash": // ร้องขอการจ่ายเหรียญออกทางด้านหน้าเครื่องโดยระบุจานวนเป็นยอดเงิน
		case "payout_by_coin": // ร้องขอการจ่ายเหรียญออกทางด้านหน้าเครื่องโดยระบุจานวนเป็นจานวนเหรียญ
		case "empty": // ร้องขอการปล่อยเหรียญทั้งหมดออกทางด้านล่าง
		case "reset": // ร้องขอการ Reset ตัวเครื่อง เพ่ือเคลียร์ค่า Error ต่างๆ
		case "status_change": // Event น้ีจะเกิดข้ึนเม่ือสถานะใดๆของ Coins Hopper มีการเปลี่ยนแปลง
		}
	case "coin_acc":
	case "bill_acc":
	case "printer":
	}
}
