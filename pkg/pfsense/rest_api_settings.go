package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RESTAPISettings struct {
	AuthMethods           []string `json:"auth_methods"`
	ContentType           string   `json:"content_type"`
	EnableLoginProtection bool     `json:"enable_login_protection"`
	HATierSync            bool     `json:"ha_tier_sync"`
	HASync                bool     `json:"ha_sync"`
	HAPartnerIP           string   `json:"ha_partner_ip"`
	HAPartnerKey          string   `json:"ha_partner_key"`
}

type RESTAPISettingsUpdateRequest struct {
	AuthMethods []string `json:"auth_methods,omitempty"`
}

func (pf *Client) GetRESTAPISettings(ctx context.Context) (*RESTAPISettings, error) {
	defer pf.read(&pf.mutexes.RESTAPISettings)()

	respBody, statusCode, err := pf.callRESTAPI(ctx, http.MethodGet, "/api/v2/system/restapi/settings", nil)
	if err != nil {
		return nil, fmt.Errorf("%w REST API settings, %w", ErrGetOperationFailed, err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("%w REST API settings, HTTP %d: %s", ErrGetOperationFailed, statusCode, truncate(string(respBody), 500))
	}

	data, err := parseRESTAPIResponse(respBody)
	if err != nil {
		return nil, fmt.Errorf("%w REST API settings, %w", ErrGetOperationFailed, err)
	}

	var settings RESTAPISettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("%w REST API settings response, %w", ErrUnableToParse, err)
	}

	return &settings, nil
}

func (pf *Client) UpdateRESTAPISettings(ctx context.Context, opts RESTAPISettingsUpdateRequest) (*RESTAPISettings, error) {
	defer pf.write(&pf.mutexes.RESTAPISettings)()

	respBody, statusCode, err := pf.callRESTAPI(ctx, http.MethodPatch, "/api/v2/system/restapi/settings", opts)
	if err != nil {
		return nil, fmt.Errorf("%w REST API settings, %w", ErrUpdateOperationFailed, err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("%w REST API settings, HTTP %d: %s", ErrUpdateOperationFailed, statusCode, truncate(string(respBody), 500))
	}

	data, err := parseRESTAPIResponse(respBody)
	if err != nil {
		return nil, fmt.Errorf("%w REST API settings, %w", ErrUpdateOperationFailed, err)
	}

	var settings RESTAPISettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("%w REST API settings response, %w", ErrUnableToParse, err)
	}

	return &settings, nil
}
