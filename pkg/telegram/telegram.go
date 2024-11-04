package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	apiEndpoint = "http.telegram.org"

	getUpdatesMethod    = "getUpdates"
	sendMessageMethod   = "sendMessage"
	deleteMessageMethod = "deleteMessage"
)

type ItgClient interface {
	Updates(ctx context.Context, offset, limit int) ([]Update, error)
	SendMessage(ctx context.Context, chatID int, text, keyboard string) error
	DeleteMessage(ctx context.Context, chatID, messageId int) error
}

type TgClient struct {
	basePath string
	client   *http.Client
}

func NewTgClient(token string) *TgClient {
	return &TgClient{
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

func (c *TgClient) Updates(ctx context.Context, offset, limit int) ([]Update, error) {
	queryParams := url.Values{}
	queryParams.Add("offset", intToString(offset))
	queryParams.Add("limit", intToString(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, queryParams)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	updates, err := parseUpdates(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updates: %w", err)
	}

	return updates, nil
}

func (c *TgClient) SendMessage(ctx context.Context, chatID int, text, keyboard string) error {
	queryParams := url.Values{}
	queryParams.Add("chat_id", intToString(chatID))
	queryParams.Add("text", text)
	if keyboard != "" {
		queryParams.Add("parse_mode", "HTML")
		queryParams.Add("reply_markup", keyboard)
	}

	_, err := c.doRequest(ctx, sendMessageMethod, queryParams)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *TgClient) DeleteMessage(ctx context.Context, chatID, messageId int) error {
	queryParams := url.Values{}
	queryParams.Add("chat_id", intToString(chatID))
	queryParams.Add("message_id", intToString(messageId))

	_, err := c.doRequest(ctx, deleteMessageMethod, queryParams)
	if err != nil {
		return fmt.Errorf("can't delete message: %w", err)
	}

	return nil
}

func (c *TgClient) doRequest(ctx context.Context, method string, queryParams url.Values) ([]byte, error) {
	requestURL := url.URL{
		Scheme:   "https",
		Host:     apiEndpoint,
		Path:     path.Join(c.basePath, method),
		RawQuery: queryParams.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
