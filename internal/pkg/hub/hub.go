package hub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

type Client struct {
	resty *resty.Client
}

func NewClient() *Client {
	restyClient := resty.New()
	restyClient.SetBaseURL("https://hubapi.quakeworld.nu/v2")
	return &Client{resty: restyClient}
}

func (c *Client) MvdsvServers(queryParams map[string]string) ([]mvdsv.Mvdsv, error) {
	resp, err := c.resty.R().SetResult([]mvdsv.Mvdsv{}).SetQueryParams(queryParams).Get("servers/mvdsv")

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make([]mvdsv.Mvdsv, 0), err
	}

	if resp.IsError() {
		return make([]mvdsv.Mvdsv, 0), errors.New(resp.Status())
	}

	servers := resp.Result().(*[]mvdsv.Mvdsv)
	return *servers, nil
}

func (c *Client) FindPlayer(pattern string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(pattern) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	if !strings.Contains(pattern, "*") {
		pattern = fmt.Sprintf("*%s*", pattern)
	}

	servers, _ := c.MvdsvServers(map[string]string{"has_player": pattern})

	if 0 == len(servers) {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`player "%s" not found.`, pattern))
	}

	return servers[0], nil
}

type Stream struct {
	Channel       string `json:"channel"`
	Url           string `json:"url"`
	Title         string `json:"title"`
	ViewerCount   int    `json:"viewers"`
	Language      string `json:"language"`
	ClientName    string `json:"client_name"`
	ServerAddress string `json:"server_address"`
}

func (c *Client) Streams() ([]Stream, error) {
	resp, err := c.resty.R().SetResult([]Stream{}).Get("streams")

	if err != nil {
		return make([]Stream, 0), err
	}

	if resp.IsError() {
		return make([]Stream, 0), errors.New(resp.Status())
	}

	return *resp.Result().(*[]Stream), nil
}
