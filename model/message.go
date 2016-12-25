package model

type Message struct {
	Device  string
	Payload Payload
}

type Payload struct {
	Type    string
	Command string
	Result  bool
	Data    interface{}
}
