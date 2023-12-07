package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ApiKeys struct {
	PubKey     string    `json:"pub_key"`
	PrivKey    string    `json:"priv_key"`
	PassPhrase string    `json:"pass_phrase"`
	Exchange   string    `json:"exchange"`
	UserId     uuid.UUID `json:"user_id"`
}

func NewApiKeys(pub, priv, phrase, exchange, userId string) *ApiKeys {
	uid := uuid.MustParse(userId)
	return &ApiKeys{
		PubKey:     pub,
		PrivKey:    priv,
		PassPhrase: phrase,
		Exchange:   exchange,
		UserId:     uid,
	}
}

func (a *ApiKeys) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.PubKey, validation.Required),
		validation.Field(&a.PrivKey, validation.Required),
		validation.Field(&a.Exchange, validation.Required))
}

func (a *ApiKeys) EmptyPassPhrase() bool {
	return a.PassPhrase == ""
}
