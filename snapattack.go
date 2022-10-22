package snapattack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const apiRoot = "https://app.snapattack.com/api/"

type Client struct {
	apiKey string

	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout:   1*time.Minute,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
			},
		},
	}
}

func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

func (c *Client) makeAPIRequest(verb string, path string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(verb, fmt.Sprintf("https://%s/%s", apiRoot, path), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-API-Key", c.apiKey)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 && resp.StatusCode > 202 {
		return dat, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	return dat, nil
}

func (c *Client) makeJSONAPIRequest(verb string, path string, data interface{}, out interface{}) error {
	dat, err := c.makeAPIRequest(verb, path, data)
	if err != nil {
		return err
	}
	return json.Unmarshal(dat, out)
}

func (c *Client) Export(ctx context.Context, filter Filter, target Target, format Format) ([]byte, error) {
	// Issue the export
	taskResp := struct {
		TaskID string `json:"task_id"`
	}{}
	err := c.makeJSONAPIRequest(http.MethodPost, "harbor/signatures/export/", map[string]interface{}{
		"analytic_compilation_target_id": target,
		"filter": filter,
		"format": []string{format},
	}, &taskResp)
	if err != nil {
		return nil, err
	}

	// Start polling for results
	for {
		task := Task{}
		
	}

}