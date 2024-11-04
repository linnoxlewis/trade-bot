package telegram

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func intToString(value int) string {
	return strconv.Itoa(value)
}

func newBasePath(token string) string {
	return fmt.Sprintf("bot%s", token)
}

func parseUpdates(data []byte) ([]Update, error) {
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	updates := make([]Update, 0, len(res.Result))
	for _, v := range res.Result {
		if v.Message != nil {
			updates = append(updates, Update{
				ID:      v.ID,
				Message: v.Message,
			})
		}
	}

	return updates, nil
}
