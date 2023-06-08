package main

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type CardOwner struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Card struct {
	CardNumber string `json:"card_number"`
	CardType   string `json:"card_type"`
	ExpMonth   int    `json:"exp_month"`
	ExpYear    int    `json:"exp_year"`
	CVV        string `json:"cvv"`
}

type CardTransaction struct {
	Amount    int        `json:"amount"`
	LimitLeft int        `json:"limit_left"`
	Currency  string     `json:"currency"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Card      *Card      `json:"card"`
	Owner     *CardOwner `json:"owner"`
	UTC       time.Time  `json:"utc"`
}

type TransactionSource struct {
	CardOwners []CardOwner
	Cards      []Card
	Seed       int64
}

// new transaction source
func NewTransactionSource(seed, cardOwnersCount, cardCount int64) *TransactionSource {
	cards := make([]Card, cardCount)
	cardOwners := make([]CardOwner, cardOwnersCount)

	gofakeit.Seed(seed)

	for i := 0; i < int(cardCount); i++ {
		cards[i] = Card{
			CardNumber: gofakeit.CreditCardNumber(
				&gofakeit.CreditCardOptions{},
			),
			CardType: gofakeit.CreditCardType(),
			ExpMonth: gofakeit.Number(1, 12),
			ExpYear:  gofakeit.Number(2023, 2030),
			CVV:      gofakeit.CreditCardCvv(),
		}
	}

	for i := 0; i < int(cardOwnersCount); i++ {
		cardOwners[i] = CardOwner{
			Id:        i,
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
	}

	return &TransactionSource{
		Seed:       seed,
		CardOwners: cardOwners,
		Cards:      cards,
	}
}

func getRandomWarsawCoordinates() (float64, float64) {
	lowLat, highLat := 52.1, 52.4
	lowLong, highLong := 20.8, 21.2

	return gofakeit.Float64Range(lowLat, highLat), gofakeit.Float64Range(lowLong, highLong)
}

func (ts *TransactionSource) GetTransaction() *CardTransaction {
	i := gofakeit.Number(0, len(ts.Cards)-1)

	card := ts.Cards[i]
	cardOwner := ts.CardOwners[i%len(ts.CardOwners)]

	now := time.Now()
	amount := gofakeit.Number(1, 20000)
	limit := amount + gofakeit.Number(-200, 5000)
	if limit < 0 {
		limit = 0
	}
	lat, long := getRandomWarsawCoordinates()

	return &CardTransaction{
		Amount:    amount,
		LimitLeft: limit,
		Currency:  "PLN",
		Latitude:  lat,
		Longitude: long,
		Card:      &card,
		Owner:     &cardOwner,
		UTC:       now,
	}
}

func (tr *CardTransaction) ToJSON() ([]byte, error) {
	return json.Marshal(tr)
}
