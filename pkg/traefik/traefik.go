package traefik

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/m50/traefik-pihole/pkg/utils"
	"github.com/spf13/viper"
)

const (
	routersPath = "api/http/routers"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	c        HTTPClient
	baseAddr string
}

func NewClient(httpClient HTTPClient) *Client {
	addr := viper.GetString("traefik-address")
	addr = strings.TrimRight(addr, "/ ")
	return &Client{
		c:        httpClient,
		baseAddr: addr,
	}
}

func (c *Client) ListHosts(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseAddr, routersPath), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	routers, err := utils.ReadHttpResponseBody[RouterList](resp)
	if err != nil {
		return nil, err
	}

	return routers.ToHosts(), nil
}
