package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod    = "getUpdates"
	sendMessageMethod   = "sendMessage"
	deleteMessageMethod = "deleteMessage"
	apiEndpoint         = "api.telegram.org"
)

func New(token string) *Client {
	return &Client{
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (updates []Update, err error) {
	defer func() {
		if err != nil {
			err = errors.New("can't get updates: " + err.Error())
		}
	}()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	return c.formatResponse(data)
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text, keyboard string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	if keyboard != "" {
		q.Add("parse_mode", "HTML")
		q.Add("reply_markup", keyboard)
	}

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return errors.New("can't send message:" + err.Error())
	}

	return nil
}

func (c *Client) DeleteMessage(ctx context.Context, chatID int, messageId int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageId))

	_, err := c.doRequest(deleteMessageMethod, q)
	if err != nil {
		return errors.New("can't send message:" + err.Error())
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		if err != nil {
			err = errors.New("can't do request: " + err.Error())
		}
	}()

	u := url.URL{
		Scheme: "https",
		Host:   apiEndpoint,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) formatResponse(data []byte) ([]Update, error) {
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	var resRepl UpdatesReplyResponse
	if err := json.Unmarshal(data, &resRepl); err != nil {
		return nil, err
	}

	upds := make([]Update, 0)
	for _, v := range res.Result {
		if v.Message != nil {
			upd := Update{v.ID, v.Message}
			upds = append(upds, upd)
		}
	}
	for _, v := range resRepl.Result {
		if v.Message != nil {
			upd := Update{ID: v.ID, Message: &IncomingMessage{
				Chat:      v.Message.Message.Chat,
				Text:      v.Message.Data,
				From:      v.Message.From,
				MessageID: v.Message.Message.MessageID,
				Date:      v.Message.Message.Date,
			}}
			upds = append(upds, upd)
		}
	}

	return upds, nil
}
