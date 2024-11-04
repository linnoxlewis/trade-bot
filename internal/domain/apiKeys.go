package domain

import (
	"fmt"
	"github.com/linnoxlewis/trade-bot/internal/helper"
)

type ApiKeys struct {
	UserId     int64
	PubKey     string
	PrivKey    string
	Passphrase string
	Exchange   string
	Date
}

func NewApiKeys(userId int64,
	exchange,
	pubKey,
	privKey,
	passphrase string,
) *ApiKeys {
	return &ApiKeys{
		UserId:     userId,
		Exchange:   exchange,
		PubKey:     pubKey,
		PrivKey:    privKey,
		Passphrase: passphrase,
	}
}

func (a *ApiKeys) EmptyPassphrase() bool {
	return a.Passphrase == ""
}

func (a *ApiKeys) DecryptPrivKey(secret string) error {
	decryptedPrivKey, err := helper.DecryptMessage(a.PrivKey, secret)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}
	a.PrivKey = decryptedPrivKey

	return nil
}

func (a *ApiKeys) DecryptPassphrase(secret string) error {
	decryptedPassphrase, err := helper.DecryptMessage(a.Passphrase, secret)
	if err != nil {
		return fmt.Errorf("failed to decrypt passphrase: %w", err)
	}
	a.Passphrase = decryptedPassphrase

	return nil
}

func (a *ApiKeys) DecryptField(field *string, secret string) error {
	decryptedValue, err := helper.DecryptMessage(*field, secret)
	if err != nil {
		return fmt.Errorf("failed to decrypt field: %w", err)
	}
	*field = decryptedValue
	
	return nil
}
