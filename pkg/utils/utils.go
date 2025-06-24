package utils

import (
	"encoding/json"
	"net/http"
)

func ReadHttpResponseBody[T any](resp *http.Response) (T, error) {
	var respBody T
	var bodyBytes []byte
	if _, err := resp.Body.Read(bodyBytes); err != nil {
		return respBody, err
	}
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return respBody, err
	}

	return respBody, nil
}
