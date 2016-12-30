package model

import (
	"errors"
	"fmt"
	"log"
)

type Host struct {
	Id            string
	Online        bool
	TotalEscrow   int // มูลค่าเงินพักทั้งหมด
	BillEscrow    int // มูลค่าธนบัตรที่พักอยู่ในเครื่องรับธนบัตร
	BillBox       int // มูลค่าธนบัตรในกล่องเก็บธนบัตร
	CoinHopperBox int // มูลค่าเหรียญใน Coin Hopper
	CoinBox       int // มูลค่าเหรียญใน CoinBox
	TotalCash     int // รวมมูลค่าเงินในตู้นี้
	WebClient     *Client
	DevClient     *Client
	SetWebClient  chan *Client
	SetDevClient  chan *Client
	CheckOnhand   chan *Client
	CancelOrder   chan *Client
}

func (h *Host) Run() {
	for {
		fmt.Println("Host.Run()...wait for next Channel")
		select {
		case c := <-h.SetWebClient:
			log.Println("<-WebClient:", c)
			h.WebClient = c
		case c := <-h.SetDevClient:
			log.Println("<-devClient:", c)
			h.DevClient = c
		case c := <-h.CheckOnhand:
			log.Println("<-CheckOnhand:", c)
			h.Onhand(c)
		case c := <-h.CancelOrder:
			log.Println("<-CancelOrder:", c)
			err := h.Cancel(c)
			if err != nil {
				log.Println("error", err)
			}
		}
	}
}

// Onhand ส่งค่าเงินพัก Escrow ที่ Host เก็บไว้กลับไปให้ web
func (h *Host) Onhand(c *Client) {
	fmt.Println("Host.Onhand <-Message...")
	c.Msg.Payload.Result = true
	c.Msg.Payload.Type = "response"
	c.Msg.Payload.Data = h.TotalEscrow
	c.Send <- c.Msg
}

// Cancel คืนเงินจากทุก Device โดยตรวจสอบเงิน Escrow ใน Bill Acceptor ด้วยถ้ามีให้คืนเงิน
func (h *Host) Cancel(c *Client) error {
	fmt.Println("Host.Cancel()...")

	// Check Bill Acceptor
	if h.TotalEscrow == 0 { // ไม่มีเงินพัก
		log.Println("ไม่มีเงินพัก:")
		c.Msg.Payload.Type = "response"
		c.Msg.Payload.Result = false
		c.Msg.Payload.Data = "ไม่มีเงินพัก"
		c.Send <- c.Msg
		return errors.New("ไม่มีเงินพัก")
	}
	//  สั่งให้ BillAcceptor คืนเงินที่พักไว้
	m1 := &Message{
		Device: "bill_acc",
		Payload: Payload{
			Command: "escrow",
			Type:    "request",
			Result:  true,
			Data:    false,
		},
	}
	h.DevClient.Send <- m1

	// Todo: Check BillAcc response
	err := h.DevClient.Conn.ReadJSON(&m1)
	if err != nil {
		log.Println("Host.Cancel() error 1")
		return err
	}

	// Success
	coinHopperEscrow := h.TotalEscrow - h.BillEscrow
	h.BillEscrow = 0

	// CoinHopper
	// ให้จ่ายเหรียญที่คงค้างตามยอด Escrow ออกด้านหน้า
	m2 := &Message{
		Device: "coin_hopper",
		Payload: Payload{
			Command: "payout_by_cash",
			Type:    "request",
			Data:    coinHopperEscrow,
		},
	}
	h.DevClient.Send <- m2
	// Todo: Check if error from CoinHopper
	//if err != nil {
	//	m2.Payload.Result = false
	//	c.Send <- m2
	//}
	h.TotalEscrow = 0

	// Send message to Web Client
	c.Msg.Payload.Type = "response"
	c.Msg.Payload.Result = true
	c.Msg.Payload.Data = "sucess"
	c.Send <- c.Msg
	return nil
}
