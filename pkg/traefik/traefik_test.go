package traefik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type MockHTTPClient struct{}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	sPerPage := q.Get("per_page")
	sPage := q.Get("page")
	perPage, err := strconv.Atoi(sPerPage)
	if err != nil {
		return nil, err
	}
	page, err := strconv.Atoi(sPage)
	if err != nil {
		return nil, err
	}

	randomGen := make(RouterList, perPage)
	for i := range perPage {
		k := Router{
			Status: "enabled",
			Rule:   fmt.Sprintf("Host(`%s`)", gofakeit.DomainName()),
		}
		if i == 2 {
			k.Rule = "Path(`/path`)"
		}
		randomGen[i] = k
	}
	d, err := json.Marshal(randomGen)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       io.NopCloser(bytes.NewReader(d)),
		Header:     http.Header{},
	}
	nextPage := page + 1
	if nextPage > 3 {
		nextPage = 1
	}
	resp.Header.Add("X-Next-Page", strconv.Itoa(nextPage))
	return resp, nil
}

func TestListHosts(t *testing.T) {
	viper.Set("traefik-address", "http://traefik:8080/")

	c := NewClient(&MockHTTPClient{})
	h, err := c.ListHosts(t.Context())
	require.NoError(t, err)
	require.Len(t, h, 57)
}
