package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gookit/slog"
)

func ReadHttpResponseBody[T any](resp *http.Response) (*T, error) {
	var respBody T
	var bodyBytes []byte
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debugf("resp: %v", resp)
		return nil, err
	}
	slog.Debug("rawbody:", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		slog.Debugf("body: %v", string(bodyBytes))
		return nil, err
	}

	return &respBody, nil
}
