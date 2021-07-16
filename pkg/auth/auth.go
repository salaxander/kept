package auth

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	grantType  = "urn:ietf:params:oauth:grant-type:device_code"
	pollURI    = "https://github.com/login/oauth/access_token"
	requestURI = "https://github.com/login/device/code"
	clientID   = "Iv1.f314ed0472e20e82"
)

var client *resty.Client

func init() {
	client = resty.New()
}

type TokenSource struct{}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := viper.GetString("authToken")
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}, nil
}

type VerificationResponse struct {
	DeviceCode      string `json:"device_code"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
}

type AuthorizedResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

func Login() (*VerificationResponse, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&VerificationResponse{}).
		SetQueryParam("client_id", clientID).
		SetQueryParam("scope", "public_repo").
		Post(requestURI)
	if err != nil {
		return nil, err
	}

	return resp.Result().(*VerificationResponse), nil
}

func PollAuthRequest(verification *VerificationResponse, ch chan string) {
	token := ""

	for token == "" {
		time.Sleep(time.Second * time.Duration(verification.Interval))
		resp, err := sendAuthRequest(verification)
		if err != nil {
		}
		token = resp.Result().(*AuthorizedResponse).AccessToken
	}

	ch <- token
}

func sendAuthRequest(verification *VerificationResponse) (*resty.Response, error) {
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&AuthorizedResponse{}).
		SetQueryParam("client_id", clientID).
		SetQueryParam("device_code", verification.DeviceCode).
		SetQueryParam("grant_type", grantType).
		Post(pollURI)
}
