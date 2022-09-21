package hub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func GetMvdsvServers(queryParams map[string]string) []mvdsv.Mvdsv {
	serversUrl := "https://hubapi.quakeworld.nu/v2/servers/mvdsv"
	resp, err := resty.New().R().SetResult([]mvdsv.Mvdsv{}).SetQueryParams(queryParams).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make([]mvdsv.Mvdsv, 0)
	}

	servers := resp.Result().(*[]mvdsv.Mvdsv)
	return *servers
}

func FindPlayerOnServer(pattern string) (mvdsv.Mvdsv, error) {
	const minFindLength = 2

	if len(pattern) < minFindLength {
		return mvdsv.Mvdsv{}, errors.New(fmt.Sprintf(`provide at least %d characters.`, minFindLength))
	}

	if !strings.Contains(pattern, "*") {
		pattern = fmt.Sprintf("*%s*", pattern)
	}

	servers := GetMvdsvServers(map[string]string{"has_player": pattern})

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

func Streams() ([]Stream, error) {
	serversUrl := "https://hubapi.quakeworld.nu/v2/streams"
	resp, err := resty.New().R().SetResult([]Stream{}).Get(serversUrl)

	if err != nil {
		fmt.Println("server fetch error", err.Error())
		return make([]Stream, 0), err
	}

	streams := resp.Result().(*[]Stream)
	return *streams, nil
}
