package domain

import (
	"github.com/google/uuid"
	"github.com/linnoxlewis/trade-bot/internal/helper"
)

type ApiKeys struct {
	UserId     uuid.UUID
	PubKey     string
	PrivKey    string
	Passphrase string
	Exchange   string
	Date
}

func NewApiKeys(userId uuid.UUID,
	exchange,
	pubKey,
	privKey,
	passphrase string) *ApiKeys {
	return &ApiKeys{
		UserId:     userId,
		Exchange:   exchange,
		PubKey:     pubKey,
		PrivKey:    privKey,
		Passphrase: passphrase,
	}
}

func (a *ApiKeys) EmptyPassPhrase() bool {
	return a.Passphrase == ""
}

func (a *ApiKeys) DecodePrivKey(secret string) {
	a.PrivKey, _ = helper.DecryptMessage(a.PrivKey, secret)
}

func (a *ApiKeys) DecodePassKey(secret string) {
	a.PrivKey, _ = helper.DecryptMessage(a.Passphrase, secret)
}
