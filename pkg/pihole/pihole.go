package pihole

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gookit/slog"
	"github.com/m50/traefik-pihole/pkg/utils"
	"github.com/spf13/viper"
)

var (
	ErrNoAuth = errors.New("auth not successful")

	authPath         = "auth"
	cnameRecordsPath = "config/dns/cnameRecords"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	c              HTTPClient
	baseAddr       string
	cname          string
	password       string
	authToken      string
	authTokenUntil time.Time
}

func NewClient(httpClient HTTPClient) *Client {
	addr := viper.GetString("pihole-address")
	addr = strings.TrimRight(addr, "/ ")
	return &Client{
		c:        httpClient,
		baseAddr: addr,
		cname:    viper.GetString("cname-address"),
		password: viper.GetString("pihole-password"),
	}
}

func (c *Client) Authenticate(ctx context.Context) error {
	b := struct{ Password string }{c.password}
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.baseAddr, authPath), bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	respBody, err := utils.ReadHttpResponseBody[AuthResponse](resp)
	if err != nil {
		return err
	}

	if !respBody.Session.Valid {
		return ErrNoAuth
	}

	c.authTokenUntil = time.Now().Add(time.Duration(respBody.Session.Validity) * time.Second)
	c.authToken = respBody.Session.SID
	return nil
}

func (c *Client) MergeHosts(ctx context.Context, hosts []string) error {
	if c.authToken == "" || c.authTokenUntil.Before(time.Now()) {
		if err := c.Authenticate(ctx); err != nil {
			return err
		}
	}
	cnames, err := c.GetCNames(ctx)
	if err != nil {
		return err
	}
	slog.WithContext(ctx).Debug("Found CNames:", cnames)
	newCNames := []string{}
	for _, h := range hosts {
		cname := fmt.Sprintf("%s,%s", h, c.cname)
		if !slices.Contains(cnames, cname) {
			slog.WithContext(ctx).Info("New host found", h)
			newCNames = append(newCNames, cname)
		}
	}

	return c.DeployCNames(ctx, newCNames)
}

func (c *Client) GetCNames(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseAddr, cnameRecordsPath), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	conf, err := utils.ReadHttpResponseBody[ConfigResponse](resp)
	if err != nil {
		return nil, err
	}
	return conf.Config.DNS.CNameRecords, nil
}

func (c *Client) DeployCNames(ctx context.Context, newCNames []string) error {
	for _, cname := range newCNames {
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s/%s", c.baseAddr, cnameRecordsPath, cname), nil)
		if err != nil {
			return fmt.Errorf("failed to add cname records [%s]: %s", cname, err)
		}
		_, err = c.c.Do(req)
		if err != nil {
			return fmt.Errorf("failed to add cname records [%s]: %s", cname, err)
		}
	}
	return nil
}
