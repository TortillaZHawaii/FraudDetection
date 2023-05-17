package main

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type CardOwner struct {
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
	Amount   int        `json:"amount"`
	Currency string     `json:"currency"`
	Card     *Card      `json:"card"`
	Owner    *CardOwner `json:"owner"`
	UTC      time.Time  `json:"utc"`
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
			ExpYear:  gofakeit.Number(2021, 2030),
			CVV:      gofakeit.CreditCardCvv(),
		}
	}

	for i := 0; i < int(cardOwnersCount); i++ {
		cardOwners[i] = CardOwner{
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

func (ts *TransactionSource) GetTransaction() *CardTransaction {
	i := gofakeit.Number(0, len(ts.Cards)-1)

	card := ts.Cards[i]
	cardOwner := ts.CardOwners[i%len(ts.CardOwners)]

	now := time.Now()
	maxLatency := time.Minute
	startRange := now.Add(-maxLatency)
	endRange := now

	return &CardTransaction{
		Amount:   gofakeit.Number(1, 100000),
		Currency: gofakeit.CurrencyShort(),
		Card:     &card,
		Owner:    &cardOwner,
		UTC:      gofakeit.DateRange(startRange, endRange).UTC(),
	}
}

func (tr *CardTransaction) ToJSON() ([]byte, error) {
	return json.Marshal(tr)
}