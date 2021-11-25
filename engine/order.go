package engine

import (
	"encoding/json"
)

type Order struct {
	Amount           uint64 `json:"amount"`
	Price            uint64 `json:"price"`
	ID               string `json:"id"`
	Side             int8   `json:"side"`
	CreatedAt        string `json:"createdAt"`
	FillOrKill       bool   `json:"fillOrKill"`
	AmountTemp       uint64 `json:"amountTemp"`
	FillIndex        []int  `json:"fillIndex"`
	ReverseCalculate int64  `json:"reverseCalculate"`
	IDCalculate      string `json:"idCalculate"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
