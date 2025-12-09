package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HasuraClient struct {
	endpoint    string
	adminSecret string
	httpClient  *http.Client
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphQLError  `json:"errors"`
}

func NewHasuraClient(endpoint, adminSecret string) *HasuraClient {
	return &HasuraClient{
		endpoint:    endpoint,
		adminSecret: adminSecret,
		httpClient:  &http.Client{},
	}
}

func (c *HasuraClient) execute(query string, variables map[string]interface{}, out interface{}) error {
	payload, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.adminSecret != "" {
		req.Header.Set("x-hasura-admin-secret", c.adminSecret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("graphql request failed with status %d", resp.StatusCode)
	}

	var gqlResp graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return err
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	if out != nil {
		if err := json.Unmarshal(gqlResp.Data, out); err != nil {
			return err
		}
	}

	return nil
}
