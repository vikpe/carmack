package hub

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func GetServerInfo(address string) (*mvdsv.Mvdsv, error) {
	url := fmt.Sprintf("https://hubapi.quakeworld.nu/v2/servers/%s", address)
	resp, err := resty.New().R().SetResult(&mvdsv.Mvdsv{}).Get(url)

	if err != nil {
		err = errors.New("unable to fetch server information")
	}

	if resp.StatusCode() != http.StatusOK {
		err = errors.New(resp.String())
	}

	if err != nil {
		return &mvdsv.Mvdsv{}, err
	}

	server := resp.Result().(*mvdsv.Mvdsv)

	return server, nil
}

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
