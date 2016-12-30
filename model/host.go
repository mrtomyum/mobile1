package model

import (
	"log"
	"fmt"
)

type Host struct {
	Id            string
	Online        bool
	WebOnline     bool
	DevOnline     bool
	TotalEscrow   int // มูลค่าเงินพักทั้งหมด
	BillEscrow    int // มูลค่าธนบัตรที่พักอยู่ในเครื่องรับธนบัตร
	BillBox       int // มูลค่าธนบัตรในกล่องเก็บธนบัตร
	CoinHopperBox int // มูลค่าเหรียญใน Coin Hopper
	CoinBox       int // มูลค่าเหรียญใน CoinBox
	TotalCash     int // รวมมูลค่าเงินในตู้นี้
	WebClient     chan *Client
	DevClient     chan *Client
	Default       chan *Client
	SendToWeb     chan *Message
	RespFromWeb   chan *Message
	SendToDev     chan *Message
	RespFromDev   chan *Message
}

func (h *Host) Run() {
	var web, dev *Client
	for {
		select {
		case web = <-h.WebClient:
			h.WebOnline = true
			fmt.Println("Host.Run() case <-h.WebClient, WebOneline = ", h.WebOnline)
		//return
		case dev = <-h.DevClient:
			h.DevOnline = true
		//return
		//case def = <-h.Default:
		case m := <-h.SendToWeb:
			web.Send <- m
		//break
		case m := <-h.SendToDev:
			dev.Send <- m
		//break
		//dev.Conn.WriteJSON(m)
		}
	}
}

func (h *Host) Cancel(web *Client) {
	// คืนเงินจากทุก Device โดยตรวจสอบเงิน Escrow ใน Bill Acceptor ด้วยถ้ามีให้คืนเงิน
	// Check Bill Acceptor
	if h.TotalEscrow == 0 { // ไม่มีเงินพัก
		log.Println("ไม่มีเงินพัก:")
	}
	//  สั่งให้ BillAcceptor คืนเงินที่พักไว้
	m1 := &Message{
		Device:"bill_acc",
		Payload: Payload{
			Command:"escrow",
			Type:   "request",
			Result: true,
			Data:   false,
		},
	}
	fmt.Println("setup Message:", m1)
	//h.SendToDev <- m1
	go func() { // Todo: ไม่ควรใช้ go func? แล้วจะดัก Channel ยังไง
		select {
		case m := <-h.RespFromDev:
			if m.Payload.Result == false {
				// Error
				break
			}
			if m.Payload.Result == false {
				// Error
			}
			// Success
			h.TotalEscrow = h.TotalEscrow - h.BillEscrow
			h.BillEscrow = 0
			break
		}
	}()
	// Check CoinHopper
	// ให้จ่ายเหรียญที่คงค้างตามยอด Escrow ออกด้านหน้า
	m2 := &Message{
		Device:"host",
		Payload: Payload{
			Command:"cancel",
			Type:   "response",
		},
	}
	err := CH.Payout(h.TotalEscrow)
	if err != nil {
		m2.Payload.Result = false
		web.Send <- m2
	}
	h.TotalEscrow = 0
	m2.Payload.Result = true
	web.Send <- m2
}

// Onhand คืนค่า
func (h *Host) Onhand(web *Client) {
	fmt.Println("Onhand Message...")
	m := &Message{
		Device:"host",
		Payload: Payload{
			Command:"onhand",
			Type:   "response",
			Result: true,
			Data:   h.TotalEscrow,
		},
	}
	web.Send <- m
}
