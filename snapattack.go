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

const apiRoot = "app.snapattack.com/api"

type Client struct {
	apiKey string

	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 1 * time.Minute,
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
		return nil, fmt.Errorf("json.Marshal(): %v", err)
	}
	url := fmt.Sprintf("https://%s/%s", apiRoot, path)
	r, err := http.NewRequest(verb, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest(): %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("htto.Do(): %v", err)
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll(): %v", err)
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
	if err := json.Unmarshal(dat, out); err != nil {
		return fmt.Errorf("json.Unmarshal(): %v: %s", err, string(dat))
	}
	return nil
}

func (c *Client) ExportSignatures(ctx context.Context, filter Filter, target Target) ([]map[string]interface{}, error) {
	// Issue the export
	taskResp := struct {
		TaskID string `json:"task_id"`
	}{}
	err := c.makeJSONAPIRequest(http.MethodPost, "harbor/signatures/export/", map[string]interface{}{
		"analytic_compilation_target_id": target,
		"filter":                         filter,
		"format":                         []string{"json"},
	}, &taskResp)
	if err != nil {
		return nil, fmt.Errorf("error starting export: %v", err)
	}

	// Start polling for results
	task := Task{}
	for {
		task = Task{}
		if err := c.makeJSONAPIRequest(http.MethodGet, fmt.Sprintf("harbor/signatures/export/%s", taskResp.TaskID), map[string]interface{}{}, &task); err != nil {
			return nil, fmt.Errorf("error getting task status: %v", err)
		}
		if task.Status == "PENDING" {
			time.Sleep(1 * time.Second)
			continue
		}
		if task.Status == "SUCCESS" {
			break
		}
		return nil, fmt.Errorf("unexpected export status: %+v", task)
	}

	// Go get the results.
	results := []map[string]interface{}{}
	if err := c.makeJSONAPIRequest(http.MethodGet, fmt.Sprintf("harbor/signatures/export/%s/result/", taskResp.TaskID), map[string]interface{}{}, &results); err != nil {
		return nil, fmt.Errorf("error fetching results: %v", err)
	}
	return results, nil
}
