package model_test

import (
	"github.com/mrtomyum/mobile1/model"
	"testing"
)

func Mock() (*model.Client, *model.Host) {
	mh := &model.Host{
		Id:            "001",
		Online:        true,
		TotalEscrow:   0,
		BillEscrow:    0,
		BillBox:       0,
		CoinHopperBox: 0,
		CoinBox:       0,
		TotalCash:     0,
		SetWebClient:  make(chan *model.Client),
		SetDevClient:  make(chan *model.Client),
		GetEscrow:     make(chan *model.Client),
		CancelOrder:   make(chan *model.Client),
	}
	mc := &model.Client{
		//Ws: conn,
		Send: make(chan *model.Message),
		Name: "web",
	}

	return mc, mh
}

func TestHost_Cancel_ZeroEscrow(t *testing.T) {
	//arrange
	c, h := Mock()
	h.TotalEscrow = 0

	//act

	//assert
	if err := h.Cancel(c); err != nil {
		t.Errorf("Error:", err)
	}

}
