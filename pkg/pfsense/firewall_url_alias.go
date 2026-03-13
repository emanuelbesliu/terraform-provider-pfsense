package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type firewallURLAliasResponse struct {
	Name        string      `json:"name"`
	Description string      `json:"descr"`
	Type        string      `json:"type"`
	URL         string      `json:"url"`
	AliasURLs   interface{} `json:"aliasurl"`
	UpdateFreq  string      `json:"updatefreq"`
	Addresses   string      `json:"address"`
	Details     string      `json:"detail"`
	ControlID   int         `json:"controlID"` //nolint:tagliatelle
}

type FirewallURLAlias struct {
	Name            string
	Description     string
	Type            string
	Entries         []FirewallURLAliasEntry
	UpdateFrequency int
	controlID       int
}

type FirewallURLAliasEntry struct {
	URL         string
	Description string
}

func (FirewallURLAlias) Types() []string {
	return []string{"url", "url_ports", "urltable", "urltable_ports"}
}

func (urlAlias *FirewallURLAlias) SetName(name string) error {
	urlAlias.Name = name

	return nil
}

func (urlAlias *FirewallURLAlias) SetDescription(description string) error {
	urlAlias.Description = description

	return nil
}

func (urlAlias *FirewallURLAlias) SetType(t string) error {
	urlAlias.Type = t

	return nil
}

func (urlAlias *FirewallURLAlias) SetUpdateFrequency(freq int) error {
	urlAlias.UpdateFrequency = freq

	return nil
}

func (entry *FirewallURLAliasEntry) SetURL(u string) error {
	entry.URL = u

	return nil
}

func (entry *FirewallURLAliasEntry) SetDescription(description string) error {
	entry.Description = description

	return nil
}

type FirewallURLAliases []FirewallURLAlias

func (urlAliases FirewallURLAliases) GetByName(name string) (*FirewallURLAlias, error) {
	for _, urlAlias := range urlAliases {
		if urlAlias.Name == name {
			return &urlAlias, nil
		}
	}

	return nil, fmt.Errorf("url alias %w with name '%s'", ErrNotFound, name)
}

func (urlAliases FirewallURLAliases) GetControlIDByName(name string) (*int, error) {
	for _, urlAlias := range urlAliases {
		if urlAlias.Name == name {
			return &urlAlias.controlID, nil
		}
	}

	return nil, fmt.Errorf("url alias %w with name '%s'", ErrNotFound, name)
}

func parseURLAliasResponse(resp firewallURLAliasResponse) (FirewallURLAlias, error) {
	var urlAlias FirewallURLAlias

	if err := urlAlias.SetName(resp.Name); err != nil {
		return urlAlias, err
	}

	if err := urlAlias.SetDescription(resp.Description); err != nil {
		return urlAlias, err
	}

	if err := urlAlias.SetType(resp.Type); err != nil {
		return urlAlias, err
	}

	urlAlias.controlID = resp.ControlID

	// Parse update frequency for urltable types.
	if resp.UpdateFreq != "" {
		freq, err := strconv.Atoi(resp.UpdateFreq)
		if err != nil {
			return urlAlias, fmt.Errorf("unable to parse update frequency: %w", err)
		}

		urlAlias.UpdateFrequency = freq
	}

	// Parse URL entries depending on type.
	if isURLTableType(resp.Type) {
		// urltable/urltable_ports: single URL in 'url' field.
		if resp.URL != "" {
			entry := FirewallURLAliasEntry{URL: resp.URL}

			descriptions := safeSplit(resp.Details, aliasEntryDescriptionSep)
			if len(descriptions) > 0 && descriptions[0] != "" {
				entry.Description = descriptions[0]
			}

			urlAlias.Entries = append(urlAlias.Entries, entry)
		}
	} else {
		// url/url_ports: multiple URLs in 'aliasurl' field.
		aliasURLs := parseAliasURLs(resp.AliasURLs)
		descriptions := safeSplit(resp.Details, aliasEntryDescriptionSep)

		for i, u := range aliasURLs {
			entry := FirewallURLAliasEntry{URL: u}

			if i < len(descriptions) && descriptions[i] != "" {
				entry.Description = descriptions[i]
			}

			urlAlias.Entries = append(urlAlias.Entries, entry)
		}
	}

	return urlAlias, nil
}

// parseAliasURLs handles the aliasurl field which can be a string or an array.
func parseAliasURLs(raw interface{}) []string {
	if raw == nil {
		return nil
	}

	switch v := raw.(type) {
	case string:
		if v == "" {
			return nil
		}

		return safeSplit(v, aliasEntryAddressSep)
	case []interface{}:
		result := make([]string, 0, len(v))

		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				result = append(result, s)
			}
		}

		return result
	default:
		return nil
	}
}

func isURLTableType(t string) bool {
	return t == "urltable" || t == "urltable_ports"
}

func (pf *Client) getFirewallURLAliases(ctx context.Context) (*FirewallURLAliases, error) {
	unableToParseResErr := fmt.Errorf("%w url alias response", ErrUnableToParse)
	command := "$output = array();" +
		"foreach (config_get_path('aliases/alias', []) as $k => $v) {" +
		"if (in_array($v['type'], array('url', 'url_ports', 'urltable', 'urltable_ports'))) {" +
		"$v['controlID'] = $k; $output[] = $v;" +
		"}}" +
		"echo json_encode($output);"

	var urlAliasResp []firewallURLAliasResponse
	if err := pf.executePHPCommand(ctx, command, &urlAliasResp); err != nil {
		return nil, err
	}

	urlAliases := make(FirewallURLAliases, 0, len(urlAliasResp))
	for _, resp := range urlAliasResp {
		urlAlias, err := parseURLAliasResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w, %w", unableToParseResErr, err)
		}

		urlAliases = append(urlAliases, urlAlias)
	}

	return &urlAliases, nil
}

func (pf *Client) GetFirewallURLAliases(ctx context.Context) (*FirewallURLAliases, error) {
	defer pf.read(&pf.mutexes.FirewallAlias)()

	urlAliases, err := pf.getFirewallURLAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w url aliases, %w", ErrGetOperationFailed, err)
	}

	return urlAliases, nil
}

func (pf *Client) GetFirewallURLAlias(ctx context.Context, name string) (*FirewallURLAlias, error) {
	defer pf.read(&pf.mutexes.FirewallAlias)()

	urlAliases, err := pf.getFirewallURLAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w url aliases, %w", ErrGetOperationFailed, err)
	}

	urlAlias, err := urlAliases.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w url alias, %w", ErrGetOperationFailed, err)
	}

	return urlAlias, nil
}

func (pf *Client) createOrUpdateFirewallURLAlias(ctx context.Context, urlAliasReq FirewallURLAlias, controlID *int) error {
	relativeURL := url.URL{Path: "firewall_aliases_edit.php"}
	values := url.Values{
		"name":  {urlAliasReq.Name},
		"descr": {urlAliasReq.Description},
		"type":  {urlAliasReq.Type},
		"save":  {"Save"},
	}

	if isURLTableType(urlAliasReq.Type) {
		// urltable/urltable_ports: single URL with update frequency.
		if len(urlAliasReq.Entries) > 0 {
			values.Set("address0", urlAliasReq.Entries[0].URL)
			values.Set("address_subnet0", strconv.Itoa(urlAliasReq.UpdateFrequency))
			values.Set("detail0", urlAliasReq.Entries[0].Description)
		}
	} else {
		// url/url_ports: multiple URLs.
		for index, entry := range urlAliasReq.Entries {
			values.Set(fmt.Sprintf("address%d", index), entry.URL)
			values.Set(fmt.Sprintf("detail%d", index), entry.Description)
		}
	}

	if controlID != nil {
		q := relativeURL.Query()
		q.Set("id", strconv.Itoa(*controlID))
		relativeURL.RawQuery = q.Encode()

		values.Set("origname", urlAliasReq.Name)
	}

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateFirewallURLAlias(ctx context.Context, urlAliasReq FirewallURLAlias) (*FirewallURLAlias, error) {
	defer pf.write(&pf.mutexes.FirewallAlias)()

	if err := pf.createOrUpdateFirewallURLAlias(ctx, urlAliasReq, nil); err != nil {
		return nil, fmt.Errorf("%w url alias, %w", ErrCreateOperationFailed, err)
	}

	urlAliases, err := pf.getFirewallURLAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w url aliases after creating, %w", ErrGetOperationFailed, err)
	}

	urlAlias, err := urlAliases.GetByName(urlAliasReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w url alias after creating, %w", ErrGetOperationFailed, err)
	}

	return urlAlias, nil
}

func (pf *Client) UpdateFirewallURLAlias(ctx context.Context, urlAliasReq FirewallURLAlias) (*FirewallURLAlias, error) {
	defer pf.write(&pf.mutexes.FirewallAlias)()

	urlAliases, err := pf.getFirewallURLAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w url aliases, %w", ErrGetOperationFailed, err)
	}

	controlID, err := urlAliases.GetControlIDByName(urlAliasReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w url alias, %w", ErrGetOperationFailed, err)
	}

	if err := pf.createOrUpdateFirewallURLAlias(ctx, urlAliasReq, controlID); err != nil {
		return nil, fmt.Errorf("%w url alias, %w", ErrUpdateOperationFailed, err)
	}

	urlAliases, err = pf.getFirewallURLAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w url aliases after updating, %w", ErrGetOperationFailed, err)
	}

	urlAlias, err := urlAliases.GetByName(urlAliasReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w url alias after updating, %w", ErrGetOperationFailed, err)
	}

	return urlAlias, nil
}

func (pf *Client) DeleteFirewallURLAlias(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.FirewallAlias)()

	urlAliases, err := pf.getFirewallURLAliases(ctx)
	if err != nil {
		return fmt.Errorf("%w url aliases, %w", ErrGetOperationFailed, err)
	}

	controlID, err := urlAliases.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w url alias, %w", ErrGetOperationFailed, err)
	}

	if err := pf.deleteFirewallAlias(ctx, *controlID); err != nil {
		return fmt.Errorf("%w url alias, %w", ErrDeleteOperationFailed, err)
	}

	urlAliases, err = pf.getFirewallURLAliases(ctx)
	if err != nil {
		return fmt.Errorf("%w url aliases after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := urlAliases.GetByName(name); err == nil {
		return fmt.Errorf("%w url alias, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
