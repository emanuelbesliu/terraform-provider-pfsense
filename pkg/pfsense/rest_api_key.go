package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type restAPIResponse struct {
	Code       int             `json:"code"`
	Status     string          `json:"status"`
	ReturnCode int             `json:"return_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

type RESTAPIKeyCreateRequest struct {
	Description string `json:"descr,omitempty"`
	HashAlgo    string `json:"hash_algo,omitempty"`
	LengthBytes int    `json:"length_bytes,omitempty"`
}

type RESTAPIKeyResponse struct {
	ID          int    `json:"id"`
	Description string `json:"descr"`
	Username    string `json:"username"`
	HashAlgo    string `json:"hash_algo"`
	LengthBytes int    `json:"length_bytes"`
	Hash        string `json:"hash"`
	Key         string `json:"key"`
}

// callRESTAPI performs a JSON REST API call with BasicAuth authentication.
// This is separate from the HTML-scraping call() method which uses CSRF tokens.
func (pf *Client) callRESTAPI(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	var reqBody *[]byte

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("unable to marshal request body, %w", err)
		}

		reqBody = &b
	}

	u := pf.Options.URL.ResolveReference(&url.URL{Path: path}).String()

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to create request, %s %s %w", method, path, err)
	}

	req.SetBasicAuth(pf.Options.Username, pf.Options.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-pfsense")

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = int64(len(*reqBody))
	}

	resp, err := pf.retryableDo(req, reqBody)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("unable to read response body, %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// parseRESTAPIResponse extracts the data field from the REST API v2 response envelope.
func parseRESTAPIResponse(respBody []byte) (json.RawMessage, error) {
	var envelope restAPIResponse
	if err := json.Unmarshal(respBody, &envelope); err == nil && envelope.Data != nil {
		if envelope.Code >= 400 {
			return nil, fmt.Errorf("REST API error %d: %s", envelope.Code, envelope.Message)
		}

		return envelope.Data, nil
	}

	return respBody, nil
}

func (pf *Client) CreateRESTAPIKey(ctx context.Context, opts RESTAPIKeyCreateRequest) (*RESTAPIKeyResponse, error) {
	defer pf.write(&pf.mutexes.RESTAPIKey)()

	respBody, statusCode, err := pf.callRESTAPI(ctx, http.MethodPost, "/api/v2/auth/key", opts)
	if err != nil {
		return nil, fmt.Errorf("%w REST API key, %w", ErrCreateOperationFailed, err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("%w REST API key, HTTP %d: %s", ErrCreateOperationFailed, statusCode, truncate(string(respBody), 500))
	}

	data, err := parseRESTAPIResponse(respBody)
	if err != nil {
		return nil, fmt.Errorf("%w REST API key, %w", ErrCreateOperationFailed, err)
	}

	var key RESTAPIKeyResponse
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("%w REST API key response, %w", ErrUnableToParse, err)
	}

	return &key, nil
}

func (pf *Client) GetRESTAPIKeys(ctx context.Context) ([]RESTAPIKeyResponse, error) {
	defer pf.read(&pf.mutexes.RESTAPIKey)()

	respBody, statusCode, err := pf.callRESTAPI(ctx, http.MethodGet, "/api/v2/auth/keys", nil)
	if err != nil {
		return nil, fmt.Errorf("%w REST API keys, %w", ErrGetOperationFailed, err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("%w REST API keys, HTTP %d: %s", ErrGetOperationFailed, statusCode, truncate(string(respBody), 500))
	}

	data, err := parseRESTAPIResponse(respBody)
	if err != nil {
		return nil, fmt.Errorf("%w REST API keys, %w", ErrGetOperationFailed, err)
	}

	var keys []RESTAPIKeyResponse
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("%w REST API keys response, %w", ErrUnableToParse, err)
	}

	return keys, nil
}

func (pf *Client) GetRESTAPIKeyByID(ctx context.Context, id int) (*RESTAPIKeyResponse, error) {
	keys, err := pf.GetRESTAPIKeys(ctx)
	if err != nil {
		return nil, err
	}

	for i := range keys {
		if keys[i].ID == id {
			return &keys[i], nil
		}
	}

	return nil, fmt.Errorf("%w, REST API key with id %d", ErrNotFound, id)
}

func (pf *Client) DeleteRESTAPIKey(ctx context.Context, id int) error {
	defer pf.write(&pf.mutexes.RESTAPIKey)()

	path := "/api/v2/auth/key?id=" + strconv.Itoa(id)

	respBody, statusCode, err := pf.callRESTAPI(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("%w REST API key %d, %w", ErrDeleteOperationFailed, id, err)
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("%w REST API key %d, HTTP %d: %s", ErrDeleteOperationFailed, id, statusCode, truncate(string(respBody), 500))
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}

	return s
}
