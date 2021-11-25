package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderFromJSON(t *testing.T) {
	var order Order
	json := `
		{
			"id"								:	"5a374064ec7740f287def7a66dff2dc4",
			"price"							:	7400,
			"amountTemp"				:	0,
			"amount"						:	1,
			"side" 							:	1,
			"createdAt" 				:	"2021-11-25 05:04:15.291438218 +0000 UTC",
			"fillOrKill"				:	false,
			"fillIndex" 				:	[],
			"reverseCalculate"	: 0,
			"idCalculate"				: ""
		}
	`
	result := Order{
		1,
		7400,
		"5a374064ec7740f287def7a66dff2dc4",
		1,
		"2021-11-25 05:04:15.291438218 +0000 UTC",
		false,
		0,
		[]int{},
		0,
		"",
	}

	if err := order.FromJSON([]byte(json)); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, result, order, "this should be equal")
}

func TestOrderToJSON(t *testing.T) {
	json := `{"amount":1,"price":7400,"id":"5a374064ec7740f287def7a66dff2dc4","side":1,"createdAt":"2021-11-25 05:04:15.291438218 +0000 UTC","fillOrKill":false,"amountTemp":0,"fillIndex":[],"reverseCalculate":0,"idCalculate":""}`
	order := Order{
		1,
		7400,
		"5a374064ec7740f287def7a66dff2dc4",
		1,
		"2021-11-25 05:04:15.291438218 +0000 UTC",
		false,
		0,
		[]int{},
		0,
		"",
	}

	assert.Equal(t, string(order.ToJSON()), json, "this should be equal")
}
