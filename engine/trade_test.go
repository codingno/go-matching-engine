package engine

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTradeFromJSON(t *testing.T) {
	var trade Trade
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()

	json := fmt.Sprintf(`{"taker_order_id":"%v","maker_order_id":"%v","amount":5,"price":7400,"createdAt":"%v"}`, uniqueTaker, uniqueMaker, createdAt)

	result := Trade{
		TakerOrderID: uniqueTaker,
		MakerOrderID: uniqueMaker,
		Amount:       5,
		Price:        7400,
		CreatedAt:    createdAt,
	}

	if err := trade.FromJSON([]byte(json)); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, result, trade, "this trade from JSON should be equal")
}

func TestTradeToJSON(t *testing.T) {
	uniqueTaker := strings.Replace(uuid.New().String(), "-", "", -1)
	uniqueMaker := strings.Replace(uuid.New().String(), "-", "", -1)
	createdAt := time.Now().UTC().String()

	json := fmt.Sprintf(`{"taker_order_id":"%v","maker_order_id":"%v","amount":5,"price":7400,"createdAt":"%v"}`, uniqueTaker, uniqueMaker, createdAt)
	result := Trade{
		TakerOrderID: uniqueTaker,
		MakerOrderID: uniqueMaker,
		Amount:       5,
		Price:        7400,
		CreatedAt:    createdAt,
	}

	assert.Equal(t, json, string(result.ToJSON()), "this trade to JSON should be equal")
}
