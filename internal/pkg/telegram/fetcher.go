package telegram

import (
	"context"
	"errors"
	telegramCli "github.com/linnoxlewis/trade-bot/pkg/telegram"
	"time"
)

type Meta struct {
	ChatID    int
	MessageId int
	Username  string
	Lang      string
	Date      int
}

type Fetcher struct {
	tg     *telegramCli.Client
	offset int
}

func NewFetcher(client *telegramCli.Client) *Fetcher {
	return &Fetcher{
		tg: client,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, limit int) ([]Event, error) {
	updates, err := f.tg.Updates(ctx, f.offset, limit)
	if err != nil {
		return nil, errors.New("can't get events:" + err.Error())
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, f.event(u))
	}

	f.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (f *Fetcher) event(upd telegramCli.Update) Event {
	updType := f.fetchType(upd)

	res := Event{
		Type: updType,
		Text: f.fetchText(upd),
	}

	if updType == Message {
		meta := Meta{
			ChatID:    upd.Message.Chat.ID,
			Username:  upd.Message.From.Username,
			Lang:      upd.Message.From.Language,
			MessageId: upd.Message.MessageID,
			Date:      upd.Message.Date,
		}
		if err := f.CheckTime(meta); err != nil {
			return Event{}
		}

		res.Meta = meta
	}

	return res
}

func (f *Fetcher) CheckTime(meta Meta) error {
	timestamp := time.Unix(int64(meta.Date), 0)
	maxMessageAge := time.Second * 10
	if time.Since(timestamp) > maxMessageAge {
		return errors.New("not actual")
	} else {
		return nil
	}
}

func (f *Fetcher) fetchText(upd telegramCli.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func (f *Fetcher) fetchType(upd telegramCli.Update) Type {
	if upd.Message == nil {
		return Unknown
	}

	return Message
}
