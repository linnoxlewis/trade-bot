package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateTpSl struct {
	Id        int64  `json:"id"`
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	TpPercent string `json:"tpPercent"`
	SlPercent string `json:"slPercent"`
	TpPrice   string `json:"tpPrice"`
	SlPrice   string `json:"slPrice"`
}

func (u *UpdateTpSl) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Symbol, validation.Required,
			validation.Length(4, 20),
			validation.Match(symbolRegexp)),

		validation.Field(&u.Id, validation.Required),

		validation.Field(&u.Exchange, validation.Required),

		validation.Field(&u.TpPercent,
			validation.When(u.TpPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&u.SlPercent,
			validation.When(u.SlPrice == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&u.TpPrice,
			validation.When(u.TpPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
		validation.Field(&u.SlPrice,
			validation.When(u.SlPercent == "",
				validation.Required,
				validation.Match(intRegexp),
				validation.By(zeroString),
			),
		),
	)
}
