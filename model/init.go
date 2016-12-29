package model

var (
	H  Host
	B  BillAcceptor
	C  CoinAcceptor
	CH CoinHopper
)

func init() {
	H = Host{
		Id:          "001",
		Online:      true,
		TotalEscrow: 0,
		BillEscrow:  0,
		BillBox:     0,
		CoinBox:     0,
		SendToWeb:   make(chan *Message),
		SendToDev:   make(chan *Message),
	}
	B = BillAcceptor{
		status: "ok",
	}
	C = CoinAcceptor{
		status: "ok",
	}
	CH = CoinHopper{
		status: "ok",
	}
}
