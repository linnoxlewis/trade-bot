package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	jsonType = "json"
	tmlType  = "toml"
)

type I18n struct {
	bundle *i18n.Bundle
}

func NewI18n(filepath string, languages []string) *I18n {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc(jsonType, json.Unmarshal)
	for _, lang := range languages {
		bundle.MustLoadMessageFile(fmt.Sprintf(filepath+".%v.json", lang))
	}

	return &I18n{
		bundle: bundle,
	}
}

func (i *I18n) T(identity string,
	params map[string]interface{},
	lang string) string {
	localizer := i18n.NewLocalizer(i.bundle, lang)

	result := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: identity,
		},
		TemplateData: params,
	})

	return result
}
