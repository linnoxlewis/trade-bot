package telegram

import (
	"encoding/json"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
)

type ClickBoard struct{}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func (c ClickBoard) MakeActiveOrdersExchangeKeyboard() string {
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{
					Text:         consts.Binance,
					CallbackData: activeOrdersExchangeBinanceCmd,
				},
				InlineKeyboardButton{
					Text:         consts.Kucoin,
					CallbackData: activeOrdersExchangeKucoinCmd,
				},
				InlineKeyboardButton{
					Text:         consts.Okx,
					CallbackData: activeOrdersExchangeOkxCmd,
				},
			},
		},
	}

	result, _ := json.Marshal(keyboard)

	return string(result)
}

func (c ClickBoard) MakeBalanceExchangeKeyboard() string {
	keyboard := InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				InlineKeyboardButton{
					Text:         consts.Binance,
					CallbackData: balanceExchangeBinanceCmd,
				},
				InlineKeyboardButton{
					Text:         consts.Kucoin,
					CallbackData: balanceExchangeKucoinCmd,
				},
				InlineKeyboardButton{
					Text:         consts.Okx,
					CallbackData: balanceExchangeOkxCmd,
				},
			},
		},
	}

	result, _ := json.Marshal(keyboard)

	return string(result)
}
