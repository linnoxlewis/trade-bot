package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApiKeys struct {
	PubKey     string `json:"pub_key"`
	PrivKey    string `json:"priv_key"`
	PassPhrase string `json:"pass_phrase"`
	Exchange   string `json:"exchange"`
	UserId     int64  `json:"user_id"`
}

func NewApiKeys(pub, priv, phrase, exchange string, userId int64) *ApiKeys {
	return &ApiKeys{
		PubKey:     pub,
		PrivKey:    priv,
		PassPhrase: phrase,
		Exchange:   exchange,
		UserId:     userId,
	}
}

func (a *ApiKeys) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.PubKey, validation.Required.Error("empty_field")),
		validation.Field(&a.PrivKey, validation.Required.Error("empty_field")),
		validation.Field(&a.UserId, validation.Required.Error("empty_field")),
		validation.Field(&a.Exchange, validation.Required.Error("empty_field")))

}

func (a *ApiKeys) EmptyPassPhrase() bool {
	return a.PassPhrase == ""
}
