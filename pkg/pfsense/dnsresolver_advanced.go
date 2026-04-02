package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DefaultDNSResolverAdvancedHideIdentity             = false
	DefaultDNSResolverAdvancedHideVersion              = false
	DefaultDNSResolverAdvancedPrefetch                 = false
	DefaultDNSResolverAdvancedPrefetchKey              = false
	DefaultDNSResolverAdvancedDNSSECStripped           = false
	DefaultDNSResolverAdvancedAggressiveNSEC           = false
	DefaultDNSResolverAdvancedQNameMinimisation        = false
	DefaultDNSResolverAdvancedQNameMinimisationStrict  = false
	DefaultDNSResolverAdvancedUseCaps                  = false
	DefaultDNSResolverAdvancedDNSRecordCache           = false
	DefaultDNSResolverAdvancedDisableAutoAccessControl = false
	DefaultDNSResolverAdvancedDisableAutoHostEntries   = false
	DefaultDNSResolverAdvancedDNS64                    = false
	DefaultDNSResolverAdvancedDNS64Prefix              = ""
	DefaultDNSResolverAdvancedDNS64Netbits             = ""
	DefaultDNSResolverAdvancedMsgCacheSize             = "4"
	DefaultDNSResolverAdvancedOutgoingNumTCP           = "10"
	DefaultDNSResolverAdvancedIncomingNumTCP           = "10"
	DefaultDNSResolverAdvancedEDNSBufferSize           = "auto"
	DefaultDNSResolverAdvancedNumQueriesPerThread      = "512"
	DefaultDNSResolverAdvancedJostleTimeout            = "200"
	DefaultDNSResolverAdvancedCacheMaxTTL              = "86400"
	DefaultDNSResolverAdvancedCacheMinTTL              = "0"
	DefaultDNSResolverAdvancedInfraKeepProbing         = false
	DefaultDNSResolverAdvancedInfraHostTTL             = "900"
	DefaultDNSResolverAdvancedInfraCacheNumHosts       = "10000"
	DefaultDNSResolverAdvancedUnwantedReplyThreshold   = "disabled"
	DefaultDNSResolverAdvancedLogVerbosity             = "1"
	DefaultDNSResolverAdvancedSockQueueTimeout         = "0"
)

// DNSResolverAdvanced represents the Services > DNS Resolver > Advanced Settings.
type DNSResolverAdvanced struct {
	HideIdentity             bool
	HideVersion              bool
	Prefetch                 bool
	PrefetchKey              bool
	DNSSECStripped           bool
	AggressiveNSEC           bool
	QNameMinimisation        bool
	QNameMinimisationStrict  bool
	UseCaps                  bool
	DNSRecordCache           bool
	DisableAutoAccessControl bool
	DisableAutoHostEntries   bool
	DNS64                    bool
	DNS64Prefix              string
	DNS64Netbits             string
	MsgCacheSize             string
	OutgoingNumTCP           string
	IncomingNumTCP           string
	EDNSBufferSize           string
	NumQueriesPerThread      string
	JostleTimeout            string
	CacheMaxTTL              string
	CacheMinTTL              string
	InfraKeepProbing         bool
	InfraHostTTL             string
	InfraCacheNumHosts       string
	UnwantedReplyThreshold   string
	LogVerbosity             string
	SockQueueTimeout         string
}

// dnsResolverAdvancedResponse is the JSON shape returned by the PHP read command.
type dnsResolverAdvancedResponse struct {
	HideIdentity             *string `json:"hideidentity"`
	HideVersion              *string `json:"hideversion"`
	Prefetch                 *string `json:"prefetch"`
	PrefetchKey              *string `json:"prefetchkey"`
	DNSSECStripped           *string `json:"dnssecstripped"`
	AggressiveNSEC           *string `json:"aggressivensec"`
	QNameMinimisation        *string `json:"qname-minimisation"`
	QNameMinimisationStrict  *string `json:"qname-minimisation-strict"`
	UseCaps                  *string `json:"use_caps"`
	DNSRecordCache           *string `json:"dnsrecordcache"`
	DisableAutoAccessControl *string `json:"disable_auto_added_access_control"`
	DisableAutoHostEntries   *string `json:"disable_auto_added_host_entries"`
	DNS64                    *string `json:"dns64"`
	DNS64Prefix              string  `json:"dns64prefix"`
	DNS64Netbits             string  `json:"dns64netbits"`
	MsgCacheSize             string  `json:"msgcachesize"`
	OutgoingNumTCP           string  `json:"outgoing_num_tcp"`
	IncomingNumTCP           string  `json:"incoming_num_tcp"`
	EDNSBufferSize           string  `json:"edns_buffer_size"`
	NumQueriesPerThread      string  `json:"num_queries_per_thread"`
	JostleTimeout            string  `json:"jostle_timeout"`
	CacheMaxTTL              string  `json:"cache_max_ttl"`
	CacheMinTTL              string  `json:"cache_min_ttl"`
	InfraKeepProbing         string  `json:"infra_keep_probing"`
	InfraHostTTL             string  `json:"infra_host_ttl"`
	InfraCacheNumHosts       string  `json:"infra_cache_numhosts"`
	UnwantedReplyThreshold   string  `json:"unwanted_reply_threshold"`
	LogVerbosity             string  `json:"log_verbosity"`
	SockQueueTimeout         string  `json:"sock_queue_timeout"`
}

func parseDNSResolverAdvancedResponse(resp dnsResolverAdvancedResponse) DNSResolverAdvanced {
	var da DNSResolverAdvanced

	// Presence-based booleans
	da.HideIdentity = resp.HideIdentity != nil
	da.HideVersion = resp.HideVersion != nil
	da.Prefetch = resp.Prefetch != nil
	da.PrefetchKey = resp.PrefetchKey != nil
	da.DNSSECStripped = resp.DNSSECStripped != nil
	da.AggressiveNSEC = resp.AggressiveNSEC != nil
	da.QNameMinimisation = resp.QNameMinimisation != nil
	da.QNameMinimisationStrict = resp.QNameMinimisationStrict != nil
	da.UseCaps = resp.UseCaps != nil
	da.DNSRecordCache = resp.DNSRecordCache != nil
	da.DisableAutoAccessControl = resp.DisableAutoAccessControl != nil
	da.DisableAutoHostEntries = resp.DisableAutoHostEntries != nil
	da.DNS64 = resp.DNS64 != nil

	// String-based boolean — "enabled" means true, anything else means false
	da.InfraKeepProbing = resp.InfraKeepProbing == "enabled"

	// String fields with defaults
	da.DNS64Prefix = resp.DNS64Prefix
	da.DNS64Netbits = resp.DNS64Netbits

	da.MsgCacheSize = resp.MsgCacheSize
	if da.MsgCacheSize == "" {
		da.MsgCacheSize = DefaultDNSResolverAdvancedMsgCacheSize
	}

	da.OutgoingNumTCP = resp.OutgoingNumTCP
	if da.OutgoingNumTCP == "" {
		da.OutgoingNumTCP = DefaultDNSResolverAdvancedOutgoingNumTCP
	}

	da.IncomingNumTCP = resp.IncomingNumTCP
	if da.IncomingNumTCP == "" {
		da.IncomingNumTCP = DefaultDNSResolverAdvancedIncomingNumTCP
	}

	da.EDNSBufferSize = resp.EDNSBufferSize
	if da.EDNSBufferSize == "" {
		da.EDNSBufferSize = DefaultDNSResolverAdvancedEDNSBufferSize
	}

	da.NumQueriesPerThread = resp.NumQueriesPerThread
	if da.NumQueriesPerThread == "" {
		da.NumQueriesPerThread = DefaultDNSResolverAdvancedNumQueriesPerThread
	}

	da.JostleTimeout = resp.JostleTimeout
	if da.JostleTimeout == "" {
		da.JostleTimeout = DefaultDNSResolverAdvancedJostleTimeout
	}

	da.CacheMaxTTL = resp.CacheMaxTTL
	if da.CacheMaxTTL == "" {
		da.CacheMaxTTL = DefaultDNSResolverAdvancedCacheMaxTTL
	}

	da.CacheMinTTL = resp.CacheMinTTL
	if da.CacheMinTTL == "" {
		da.CacheMinTTL = DefaultDNSResolverAdvancedCacheMinTTL
	}

	da.InfraHostTTL = resp.InfraHostTTL
	if da.InfraHostTTL == "" {
		da.InfraHostTTL = DefaultDNSResolverAdvancedInfraHostTTL
	}

	da.InfraCacheNumHosts = resp.InfraCacheNumHosts
	if da.InfraCacheNumHosts == "" {
		da.InfraCacheNumHosts = DefaultDNSResolverAdvancedInfraCacheNumHosts
	}

	da.UnwantedReplyThreshold = resp.UnwantedReplyThreshold
	if da.UnwantedReplyThreshold == "" {
		da.UnwantedReplyThreshold = DefaultDNSResolverAdvancedUnwantedReplyThreshold
	}

	da.LogVerbosity = resp.LogVerbosity
	if da.LogVerbosity == "" {
		da.LogVerbosity = DefaultDNSResolverAdvancedLogVerbosity
	}

	da.SockQueueTimeout = resp.SockQueueTimeout
	if da.SockQueueTimeout == "" {
		da.SockQueueTimeout = DefaultDNSResolverAdvancedSockQueueTimeout
	}

	return da
}

func (pf *Client) getDNSResolverAdvanced(ctx context.Context) (*DNSResolverAdvanced, error) {
	command := "$ub = config_get_path('unbound', array());" +
		"$out = array(" +
		"'hideidentity' => array_key_exists('hideidentity', $ub) ? $ub['hideidentity'] : null," +
		"'hideversion' => array_key_exists('hideversion', $ub) ? $ub['hideversion'] : null," +
		"'prefetch' => array_key_exists('prefetch', $ub) ? $ub['prefetch'] : null," +
		"'prefetchkey' => array_key_exists('prefetchkey', $ub) ? $ub['prefetchkey'] : null," +
		"'dnssecstripped' => array_key_exists('dnssecstripped', $ub) ? $ub['dnssecstripped'] : null," +
		"'aggressivensec' => array_key_exists('aggressivensec', $ub) ? $ub['aggressivensec'] : null," +
		"'qname-minimisation' => array_key_exists('qname-minimisation', $ub) ? $ub['qname-minimisation'] : null," +
		"'qname-minimisation-strict' => array_key_exists('qname-minimisation-strict', $ub) ? $ub['qname-minimisation-strict'] : null," +
		"'use_caps' => array_key_exists('use_caps', $ub) ? $ub['use_caps'] : null," +
		"'dnsrecordcache' => array_key_exists('dnsrecordcache', $ub) ? $ub['dnsrecordcache'] : null," +
		"'disable_auto_added_access_control' => array_key_exists('disable_auto_added_access_control', $ub) ? $ub['disable_auto_added_access_control'] : null," +
		"'disable_auto_added_host_entries' => array_key_exists('disable_auto_added_host_entries', $ub) ? $ub['disable_auto_added_host_entries'] : null," +
		"'dns64' => array_key_exists('dns64', $ub) ? $ub['dns64'] : null," +
		"'dns64prefix' => isset($ub['dns64prefix']) ? $ub['dns64prefix'] : ''," +
		"'dns64netbits' => isset($ub['dns64netbits']) ? $ub['dns64netbits'] : ''," +
		"'msgcachesize' => isset($ub['msgcachesize']) ? $ub['msgcachesize'] : ''," +
		"'outgoing_num_tcp' => isset($ub['outgoing_num_tcp']) ? $ub['outgoing_num_tcp'] : ''," +
		"'incoming_num_tcp' => isset($ub['incoming_num_tcp']) ? $ub['incoming_num_tcp'] : ''," +
		"'edns_buffer_size' => isset($ub['edns_buffer_size']) ? $ub['edns_buffer_size'] : ''," +
		"'num_queries_per_thread' => isset($ub['num_queries_per_thread']) ? $ub['num_queries_per_thread'] : ''," +
		"'jostle_timeout' => isset($ub['jostle_timeout']) ? $ub['jostle_timeout'] : ''," +
		"'cache_max_ttl' => isset($ub['cache_max_ttl']) ? $ub['cache_max_ttl'] : ''," +
		"'cache_min_ttl' => isset($ub['cache_min_ttl']) ? $ub['cache_min_ttl'] : ''," +
		"'infra_keep_probing' => isset($ub['infra_keep_probing']) ? $ub['infra_keep_probing'] : ''," +
		"'infra_host_ttl' => isset($ub['infra_host_ttl']) ? $ub['infra_host_ttl'] : ''," +
		"'infra_cache_numhosts' => isset($ub['infra_cache_numhosts']) ? $ub['infra_cache_numhosts'] : ''," +
		"'unwanted_reply_threshold' => isset($ub['unwanted_reply_threshold']) ? $ub['unwanted_reply_threshold'] : ''," +
		"'log_verbosity' => isset($ub['log_verbosity']) ? $ub['log_verbosity'] : ''," +
		"'sock_queue_timeout' => isset($ub['sock_queue_timeout']) ? $ub['sock_queue_timeout'] : ''" +
		");" +
		"print(json_encode($out));"

	var resp dnsResolverAdvancedResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	da := parseDNSResolverAdvancedResponse(resp)

	return &da, nil
}

func (pf *Client) GetDNSResolverAdvanced(ctx context.Context) (*DNSResolverAdvanced, error) {
	defer pf.read(&pf.mutexes.DNSResolverAdvanced)()

	da, err := pf.getDNSResolverAdvanced(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver advanced, %w", ErrGetOperationFailed, err)
	}

	return da, nil
}

func dnsResolverAdvancedFormValues(da DNSResolverAdvanced) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// Checkbox fields — only set when enabled
	if da.HideIdentity {
		values.Set("hideidentity", "yes")
	}

	if da.HideVersion {
		values.Set("hideversion", "yes")
	}

	if da.Prefetch {
		values.Set("prefetch", "yes")
	}

	if da.PrefetchKey {
		values.Set("prefetchkey", "yes")
	}

	if da.DNSSECStripped {
		values.Set("dnssecstripped", "yes")
	}

	if da.AggressiveNSEC {
		values.Set("aggressivensec", "yes")
	}

	if da.QNameMinimisation {
		values.Set("qname-minimisation", "yes")
	}

	if da.QNameMinimisationStrict {
		values.Set("qname-minimisation-strict", "yes")
	}

	if da.UseCaps {
		values.Set("use_caps", "yes")
	}

	if da.DNSRecordCache {
		values.Set("dnsrecordcache", "yes")
	}

	if da.DisableAutoAccessControl {
		values.Set("disable_auto_added_access_control", "yes")
	}

	if da.DisableAutoHostEntries {
		values.Set("disable_auto_added_host_entries", "yes")
	}

	if da.DNS64 {
		values.Set("dns64", "yes")
	}

	if da.DNS64Prefix != "" {
		values.Set("dns64prefix", da.DNS64Prefix)
	}

	if da.DNS64Netbits != "" {
		values.Set("dns64netbits", da.DNS64Netbits)
	}

	// Numeric/string fields
	values.Set("msgcachesize", da.MsgCacheSize)
	values.Set("outgoing_num_tcp", da.OutgoingNumTCP)
	values.Set("incoming_num_tcp", da.IncomingNumTCP)
	values.Set("edns_buffer_size", da.EDNSBufferSize)
	values.Set("num_queries_per_thread", da.NumQueriesPerThread)
	values.Set("jostle_timeout", da.JostleTimeout)
	values.Set("cache_max_ttl", da.CacheMaxTTL)
	values.Set("cache_min_ttl", da.CacheMinTTL)

	if da.InfraKeepProbing {
		values.Set("infra_keep_probing", "yes")
	}

	values.Set("infra_host_ttl", da.InfraHostTTL)
	values.Set("infra_cache_numhosts", da.InfraCacheNumHosts)
	values.Set("unwanted_reply_threshold", da.UnwantedReplyThreshold)
	values.Set("log_verbosity", da.LogVerbosity)
	values.Set("sock_queue_timeout", da.SockQueueTimeout)

	return values
}

func (pf *Client) UpdateDNSResolverAdvanced(ctx context.Context, da DNSResolverAdvanced) (*DNSResolverAdvanced, error) {
	defer pf.write(&pf.mutexes.DNSResolverAdvanced)()

	relativeURL := url.URL{Path: "services_unbound_advanced.php"}
	values := dnsResolverAdvancedFormValues(da)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver advanced, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w dns resolver advanced, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getDNSResolverAdvanced(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver advanced after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}
