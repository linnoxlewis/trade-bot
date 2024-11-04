package helper

import (
	"encoding/json"
	"io"
	"net/http"
)

// SendBadRequestResponse отправляет JSON-ответ с ошибкой
func SendBadRequestResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		panic(encodeErr)
	}
}

// SendSuccessResponse отправляет JSON-ответ с сообщением об успешном выполнении
func SendSuccessResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// ParseJSON чтение и анмаршалинг входящих данных
func ParseJSON(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyBytes, v)
}
