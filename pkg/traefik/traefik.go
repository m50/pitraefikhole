package traefik

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gookit/slog"
	"github.com/m50/traefik-pihole/pkg/utils"
	"github.com/spf13/viper"
)

const (
	routersPath = "api/http/routers"
	perPage = 20
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

// ListHosts gets a list of all the hosts provided by Traefik.
func (c *Client) ListHosts(ctx context.Context) ([]string, error) {
	var routers RouterList
	curPage := 1
	for {
		resp, err := c.requestRouters(ctx, curPage)
		if err != nil {
			return nil, err
		}
		pageRouters, err := utils.ReadHttpResponseBody[RouterList](resp)
		if err != nil {
			return nil, err
		}
		routers = append(routers, *pageRouters...)

		nextPage := resp.Header.Get("X-Next-Page")
		slog.Debug("next page:", nextPage)
		if nextPage == "" || nextPage == strconv.Itoa(curPage) || nextPage == "1" {
			break
		}
		nPage, err := strconv.Atoi(nextPage)
		if err != nil {
			slog.WithContext(ctx).Error("failed to fetch next page:", err)
			break
		}
		curPage = nPage
	}
	return routers.ToHosts(), nil
}

func (c *Client) requestRouters(ctx context.Context, page int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s?page=%d&per_page=%d", c.baseAddr, routersPath, page, perPage), nil)
	if err != nil {
		return nil, err
	}

	return c.c.Do(req)
}
