package uaaclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    string
	Name       string
	Secret     string
	JSONClient jsonClient
}

//go:generate counterfeiter -o fakes/json_client.go --fake-name JSONClient . jsonClient
type jsonClient interface {
	MakeRequest(*http.Request, interface{}) error
}

func (c *Client) GetToken() (string, error) {
	reqURL := fmt.Sprintf("%s/oauth/token", c.BaseURL)
	bodyString := fmt.Sprintf("grant_type=client_credentials")

	request, err := http.NewRequest("POST", reqURL, strings.NewReader(bodyString))
	if err != nil {
		return "", err
	}

	request.SetBasicAuth(url.QueryEscape(c.Name), url.QueryEscape(c.Secret))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type getTokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	response := &getTokenResponse{}
	err = c.JSONClient.MakeRequest(request, response)
	if err != nil {
		return "", err
	}
	return response.AccessToken, nil
}
