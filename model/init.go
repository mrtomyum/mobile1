package model

var (
	H  *Host
	B  *BillAcceptor
	C  *CoinAcceptor
	CH *CoinHopper
	P  *Printer
)

func init() {
	H = &Host{
		Id:            "001",
		Online:        true,
		TotalEscrow:   0,
		BillEscrow:    0,
		BillBox:       0,
		CoinHopperBox: 0,
		CoinBox:       0,
		TotalCash:     0,
		SetWebClient:  make(chan *Client),
		SetDevClient:  make(chan *Client),
		GetEscrow:     make(chan *Client),
		CancelOrder:   make(chan *Client),
	}
	B = &BillAcceptor{
		Status: "ok",
	}
	C = &CoinAcceptor{
		Status: "ok",
	}
	CH = &CoinHopper{
		Status: "ok",
	}
	P = &Printer{
		Status: "ok",
	}
}
