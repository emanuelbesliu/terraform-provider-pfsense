package pfsense

import (
	"context"
	"fmt"
)

type dynamicDNSResponse struct {
	Type              string `json:"type"`
	Interface         string `json:"interface"`
	Host              string `json:"host"`
	DomainName        string `json:"domainname"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	MX                string `json:"mx"`
	Wildcard          string `json:"wildcard"`
	Proxied           string `json:"proxied"`
	VerboseLog        string `json:"verboselog"`
	CurlIPResolveV4   string `json:"curl_ipresolve_v4"`
	CurlSSLVerifyPeer string `json:"curl_ssl_verifypeer"`
	ZoneID            string `json:"zoneid"`
	TTL               string `json:"ttl"`
	MaxCacheAge       string `json:"maxcacheage"`
	UpdateURL         string `json:"updateurl"`
	ResultMatch       string `json:"resultmatch"`
	RequestIf         string `json:"requestif"`
	CurlProxy         string `json:"curl_proxy"`
	Description       string `json:"descr"`
	Enable            string `json:"enable"`
	ControlID         int    `json:"controlID"` //nolint:tagliatelle
}

type DynamicDNS struct {
	Type              string
	Interface         string
	Host              string
	DomainName        string
	Username          string
	Password          string
	MX                string
	Wildcard          bool
	Proxied           bool
	VerboseLog        bool
	CurlIPResolveV4   bool
	CurlSSLVerifyPeer bool
	ZoneID            string
	TTL               string
	MaxCacheAge       string
	UpdateURL         string
	ResultMatch       string
	RequestIf         string
	CurlProxy         string
	Description       string
	Disabled          bool
	controlID         int
}

func (d *DynamicDNS) SetType(t string) error {
	if t == "" {
		return fmt.Errorf("%w, dynamic DNS type is required", ErrClientValidation)
	}

	d.Type = t

	return nil
}

func (d *DynamicDNS) SetInterface(iface string) error {
	if iface == "" {
		return fmt.Errorf("%w, dynamic DNS interface is required", ErrClientValidation)
	}

	d.Interface = iface

	return nil
}

func (d *DynamicDNS) SetHost(host string) error {
	d.Host = host

	return nil
}

func (d *DynamicDNS) SetDomainName(domain string) error {
	d.DomainName = domain

	return nil
}

func (d *DynamicDNS) SetUsername(username string) error {
	d.Username = username

	return nil
}

func (d *DynamicDNS) SetPassword(password string) error {
	d.Password = password

	return nil
}

func (d *DynamicDNS) SetDescription(desc string) error {
	d.Description = desc

	return nil
}

func (d *DynamicDNS) ControlID() int {
	return d.controlID
}

type DynamicDNSEntries []DynamicDNS

func (entries DynamicDNSEntries) GetByID(id int) (*DynamicDNS, error) {
	for _, e := range entries {
		if e.controlID == id {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("dynamic DNS entry %w with ID '%d'", ErrNotFound, id)
}

func parseDynamicDNSResponse(resp dynamicDNSResponse) DynamicDNS {
	return DynamicDNS{
		Type:              resp.Type,
		Interface:         resp.Interface,
		Host:              resp.Host,
		DomainName:        resp.DomainName,
		Username:          resp.Username,
		Password:          resp.Password,
		MX:                resp.MX,
		Wildcard:          resp.Wildcard != "",
		Proxied:           resp.Proxied != "",
		VerboseLog:        resp.VerboseLog != "",
		CurlIPResolveV4:   resp.CurlIPResolveV4 != "",
		CurlSSLVerifyPeer: resp.CurlSSLVerifyPeer != "",
		ZoneID:            resp.ZoneID,
		TTL:               resp.TTL,
		MaxCacheAge:       resp.MaxCacheAge,
		UpdateURL:         resp.UpdateURL,
		ResultMatch:       resp.ResultMatch,
		RequestIf:         resp.RequestIf,
		CurlProxy:         resp.CurlProxy,
		Description:       resp.Description,
		Disabled:          resp.Enable != "",
		controlID:         resp.ControlID,
	}
}

func (pf *Client) getDynamicDNSEntries(ctx context.Context) (*DynamicDNSEntries, error) {
	command := "$output = array();" +
		"$items = config_get_path('dyndnses/dyndns', array());" +
		"foreach ($items as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"$v['enable'] = array_key_exists('enable', $v) ? 'on' : '';" +
		"$v['wildcard'] = array_key_exists('wildcard', $v) ? 'on' : '';" +
		"$v['proxied'] = array_key_exists('proxied', $v) ? 'on' : '';" +
		"$v['verboselog'] = array_key_exists('verboselog', $v) ? 'on' : '';" +
		"$v['curl_ipresolve_v4'] = array_key_exists('curl_ipresolve_v4', $v) ? 'on' : '';" +
		"$v['curl_ssl_verifypeer'] = array_key_exists('curl_ssl_verifypeer', $v) ? 'on' : '';" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var resp []dynamicDNSResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	entries := make(DynamicDNSEntries, 0, len(resp))
	for _, r := range resp {
		entries = append(entries, parseDynamicDNSResponse(r))
	}

	return &entries, nil
}

func (pf *Client) GetDynamicDNSEntries(ctx context.Context) (*DynamicDNSEntries, error) {
	defer pf.read(&pf.mutexes.DynamicDNS)()

	entries, err := pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entries, %w", ErrGetOperationFailed, err)
	}

	return entries, nil
}

func (pf *Client) GetDynamicDNS(ctx context.Context, id int) (*DynamicDNS, error) {
	defer pf.read(&pf.mutexes.DynamicDNS)()

	entries, err := pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entries, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) CreateDynamicDNS(ctx context.Context, req DynamicDNS) (*DynamicDNS, error) {
	defer pf.write(&pf.mutexes.DynamicDNS)()

	command := buildDynamicDNSItem(req) +
		"$items = config_get_path('dyndnses/dyndns', array());" +
		"$items[] = $item;" +
		"config_set_path('dyndnses/dyndns', $items);" +
		"write_config('Terraform: created dynamic DNS entry');" +
		"print(json_encode(count($items) - 1));"

	var newID int
	if err := pf.executePHPCommand(ctx, command, &newID); err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrCreateOperationFailed, err)
	}

	if err := pf.applyDynamicDNSChanges(ctx); err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrCreateOperationFailed, err)
	}

	entries, err := pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entries after creating, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByID(newID)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry after creating, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) UpdateDynamicDNS(ctx context.Context, id int, req DynamicDNS) (*DynamicDNS, error) {
	defer pf.write(&pf.mutexes.DynamicDNS)()

	entries, err := pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entries, %w", ErrGetOperationFailed, err)
	}

	if _, err := entries.GetByID(id); err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrGetOperationFailed, err)
	}

	command := buildDynamicDNSItem(req) +
		fmt.Sprintf("config_set_path('dyndnses/dyndns/%d', $item);", id) +
		"write_config('Terraform: updated dynamic DNS entry');" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrUpdateOperationFailed, err)
	}

	if err := pf.applyDynamicDNSChanges(ctx); err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry, %w", ErrUpdateOperationFailed, err)
	}

	entries, err = pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entries after updating, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w dynamic DNS entry after updating, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) DeleteDynamicDNS(ctx context.Context, id int) error {
	defer pf.write(&pf.mutexes.DynamicDNS)()

	entries, err := pf.getDynamicDNSEntries(ctx)
	if err != nil {
		return fmt.Errorf("%w dynamic DNS entries, %w", ErrGetOperationFailed, err)
	}

	if _, err := entries.GetByID(id); err != nil {
		return fmt.Errorf("%w dynamic DNS entry, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"config_del_path('dyndnses/dyndns/%d');"+
			"$items = config_get_path('dyndnses/dyndns', array());"+
			"config_set_path('dyndnses/dyndns', array_values($items));"+
			"write_config('Terraform: deleted dynamic DNS entry');"+
			"print(json_encode(true));",
		id,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w dynamic DNS entry, %w", ErrDeleteOperationFailed, err)
	}

	if err := pf.applyDynamicDNSChanges(ctx); err != nil {
		return fmt.Errorf("%w dynamic DNS entry, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}

func (pf *Client) applyDynamicDNSChanges(ctx context.Context) error {
	command := "require_once('dyndns.class');" +
		"services_dyndns_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply dynamic DNS changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}

func buildDynamicDNSItem(req DynamicDNS) string {
	return "$item = array();" +
		fmt.Sprintf("$item['type'] = '%s';", phpEscape(req.Type)) +
		fmt.Sprintf("$item['interface'] = '%s';", phpEscape(req.Interface)) +
		fmt.Sprintf("$item['host'] = '%s';", phpEscape(req.Host)) +
		fmt.Sprintf("$item['domainname'] = '%s';", phpEscape(req.DomainName)) +
		fmt.Sprintf("$item['username'] = '%s';", phpEscape(req.Username)) +
		fmt.Sprintf("$item['password'] = '%s';", phpEscape(req.Password)) +
		fmt.Sprintf("$item['mx'] = '%s';", phpEscape(req.MX)) +
		boolToPHPField("$item", "wildcard", req.Wildcard) +
		boolToPHPField("$item", "proxied", req.Proxied) +
		boolToPHPField("$item", "verboselog", req.VerboseLog) +
		boolToPHPField("$item", "curl_ipresolve_v4", req.CurlIPResolveV4) +
		boolToPHPField("$item", "curl_ssl_verifypeer", req.CurlSSLVerifyPeer) +
		fmt.Sprintf("$item['zoneid'] = '%s';", phpEscape(req.ZoneID)) +
		fmt.Sprintf("$item['ttl'] = '%s';", phpEscape(req.TTL)) +
		fmt.Sprintf("$item['maxcacheage'] = '%s';", phpEscape(req.MaxCacheAge)) +
		fmt.Sprintf("$item['updateurl'] = '%s';", phpEscape(req.UpdateURL)) +
		fmt.Sprintf("$item['resultmatch'] = '%s';", phpEscape(req.ResultMatch)) +
		fmt.Sprintf("$item['requestif'] = '%s';", phpEscape(req.RequestIf)) +
		fmt.Sprintf("$item['curl_proxy'] = '%s';", phpEscape(req.CurlProxy)) +
		fmt.Sprintf("$item['descr'] = '%s';", phpEscape(req.Description)) +
		boolToPHPField("$item", "enable", req.Disabled)
}

func boolToPHPField(varName, field string, value bool) string {
	if value {
		return fmt.Sprintf("%s['%s'] = '';", varName, field)
	}

	return ""
}
