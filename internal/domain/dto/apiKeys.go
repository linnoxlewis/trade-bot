package dto

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

func (a *ApiKeys) EmptyPassPhrase() bool {
	return a.PassPhrase == ""
}
