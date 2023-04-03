package platform

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var tokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type Credentials struct {
	ClientId     string
	ClientSecret string
}

type ClientConfig struct {
	HostAuth    string
	Host        string
	Credentials *Credentials
	ProjectKey  string
	Scopes      []string
}

type Client struct {
	ApiUrl     string
	ProjectKey string
}

type MethodHead struct {
	EndPoint string
	Method   string
	Where    string
}

func NewClient(crd *ClientConfig) (*Client, error) {
	authStr := fmt.Sprintf("%s:%s", crd.Credentials.ClientId, crd.Credentials.ClientSecret)
	encodedAuthStr := base64.StdEncoding.EncodeToString([]byte(authStr))

	autToken := fmt.Sprintf("%s/oauth/token", crd.HostAuth)
	scopes := strings.Join(crd.Scopes, " ")
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", scopes)
	req, err := http.NewRequest("POST", autToken, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedAuthStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		fmt.Println("Error al decodificar la respuesta:", err)
		return nil, err
	}
	return &Client{
		ApiUrl:     crd.Host,
		ProjectKey: crd.ProjectKey,
	}, nil
}
