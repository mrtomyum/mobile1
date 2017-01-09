package model

import "fmt"

type Devicer interface {
	Status() string
}

type CoinHopperStatus int

const (
	DISABLE CoinHopperStatus = iota
	CALIBRATION_FAULT
	NO_KEY_SET
	COIN_JAMMED
	FRAUD
	HOPPER_EMPTY
	MEMORY_ERROR
	SENSORS_NOT_INITIALISED
	LID_REMOVED
)

type CoinHopper struct {
	Id     string
	Status string
}

func (ch *CoinHopper) Payout(v int) error {
	// command to send to devClient for "payout" value = v
	fmt.Println("CoinHopper Command=>Payout, Value:", v)
	return nil
}

type BillAcceptor struct {
	Id     string
	Status string
}

type CoinAcceptor struct {
	Id     string
	Status string
}

type Printer struct {
	Id     string
	Status string
}

type MainBoard struct {
	Id     string
	Status string
}
