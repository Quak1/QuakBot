package riotapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type baseClient struct {
	region     Region
	baseURL    string
	httpClient *http.Client
	apiKey     string
	limiterSec *rate.Limiter
	limiterMin *rate.Limiter
}

type errorResponse struct {
	Status struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	} `json:"status"`
}

const (
	urlFormat = "https://%s.%s/%s"
	baseURL   = "api.riotgames.com"
)

func newBaseClient(region Region, apiKey string, httpClient *http.Client) *baseClient {
	c := &baseClient{
		region:     region,
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
		limiterSec: rate.NewLimiter(rate.Every(time.Second/20), 20),     // limit 20 req/s
		limiterMin: rate.NewLimiter(rate.Every(2*time.Minute/100), 100), // limit 100 req/2min
	}

	return c
}

func (c *baseClient) do(req *http.Request, data any) error {
	req.Header.Set("X-Riot-Token", c.apiKey)
	// req.Header.Set("Content-Type", "application/json")

	if err := c.limiterSec.Wait(context.Background()); err != nil {
		return err
	}
	if err := c.limiterMin.Wait(context.Background()); err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return fmt.Errorf("error with status code: %d", res.StatusCode)
		}

		return fmt.Errorf("API error response %d: %s", res.StatusCode, errRes.Status.Message)
	}

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}

	return nil
}

func (c *baseClient) DoGet(endpoint string, data any) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlFormat, c.region, c.baseURL, endpoint), nil)
	if err != nil {
		return err
	}

	return c.do(req, data)
}

func (c *baseClient) DoAreaGet(endpoint string, data any) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlFormat, regionToArea[c.region], c.baseURL, endpoint), nil)
	if err != nil {
		return err
	}

	return c.do(req, data)
}
