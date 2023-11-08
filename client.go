package coleteonline

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// https://docs.api.colete-online.ro
type Client struct {
	authURL       string
	apiURL        string
	authBasic     string
	authBearer    string
	authBearerExp time.Time
	http          *http.Client
	mu            sync.Mutex
	timeNow       func() time.Time
}

type Config struct {
	ClientId      string
	ClientSecret  string
	UseProduction bool
	Timeout       time.Duration
}

func NewClient(config Config) *Client {
	client := &Client{
		authURL: "https://auth.colete-online.ro/token",
		authBasic: "Basic " + base64.StdEncoding.EncodeToString(
			[]byte(config.ClientId+":"+config.ClientSecret),
		),
		http: &http.Client{
			Timeout: config.Timeout,
		},
		timeNow: func() time.Time {
			return time.Now()
		},
		mu: sync.Mutex{},
	}
	if config.UseProduction {
		client.apiURL = "https://api.colete-online.ro/v1"
	} else {
		client.apiURL = "https://api.colete-online.ro/v1/staging"
	}
	return client
}

func (c *Client) GetAuthBearer() (*string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timeNow().After(c.authBearerExp) {
		req, err := http.NewRequest(
			"POST",
			c.authURL,
			strings.NewReader("grant_type=client_credentials"),
		)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.authBasic)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB
		if err != nil {
			return nil, err
		}
		if r.StatusCode != 200 {
			var rErr AuthResponseError
			err = json.Unmarshal(b, &rErr)
			if err != nil {
				return nil, err
			}
			return nil, &rErr
		}
		var res AuthToken
		err = json.Unmarshal(b, &res)
		if err != nil {
			return nil, err
		}
		c.authBearerExp, err = c.getExpiresAtFromJWT(res.AccessToken)
		if err != nil {
			return nil, err
		}
		c.authBearer = "Bearer " + res.AccessToken
	}
	return &c.authBearer, nil
}

func (c *Client) CreateOrder(order *Order) (*OrderResponse, error) {
	var res OrderResponse
	err := c.request("POST", "/order", order, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderPrice(order *Order) (*OrderPriceResponse, error) {
	var res OrderPriceResponse
	err := c.request("POST", "/order/price", order, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderStatus(uniqueIdOrAWB *string) (*OrderStatusResponse, error) {
	var res OrderStatusResponse
	err := c.request("GET", "/order/status/"+*uniqueIdOrAWB, nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) AddressList(page int64) (*AddressListResponse, error) {
	var res AddressListResponse
	err := c.request("GET", fmt.Sprintf("/address?page=%d", page), nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ServiceList() ([]ServiceResponse, error) {
	var res []ServiceResponse
	err := c.request("GET", "/service", nil, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UserBalance() (*UserBalance, error) {
	var res UserBalance
	err := c.request("GET", "/user/balance", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) request(
	method string,
	path string,
	body interface{},
	res interface{},
) error {
	token, err := c.GetAuthBearer()
	if err != nil {
		return err
	}
	var b []byte
	if body != nil {
		b, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest(
		method,
		c.apiURL+path,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", *token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	r, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode == 200 || r.StatusCode == 400 {
		b, err = io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB
		if err != nil {
			return err
		}
		if r.StatusCode == 400 {
			var rErr ResponseError
			err = json.Unmarshal(b, &rErr)
			if err != nil {
				return err
			}
			return &rErr
		}
		err = json.Unmarshal(b, res)
		if err != nil {
			return err
		}
		return nil
	}
	if r.StatusCode == 401 {
		c.authBearer = ""
		return c.request(method, path, body, res)
	}
	return &ResponseError{
		Message: "unexpected response status",
		Code:    r.StatusCode,
	}
}

// This does not guarantee that the token/payload is valid.
func (c *Client) getExpiresAtFromJWT(token string) (time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, errors.New("invalid JWT token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid JWT payload: %s", err)
	}
	return time.Unix(claims.Exp, 0), nil
}
