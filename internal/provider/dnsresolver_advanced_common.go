package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// DNSResolverAdvancedModel represents the Terraform model for DNS resolver advanced settings.
type DNSResolverAdvancedModel struct {
	HideIdentity             types.Bool   `tfsdk:"hide_identity"`
	HideVersion              types.Bool   `tfsdk:"hide_version"`
	Prefetch                 types.Bool   `tfsdk:"prefetch"`
	PrefetchKey              types.Bool   `tfsdk:"prefetch_key"`
	DNSSECStripped           types.Bool   `tfsdk:"dnssec_stripped"`
	AggressiveNSEC           types.Bool   `tfsdk:"aggressive_nsec"`
	QNameMinimisation        types.Bool   `tfsdk:"qname_minimisation"`
	QNameMinimisationStrict  types.Bool   `tfsdk:"qname_minimisation_strict"`
	UseCaps                  types.Bool   `tfsdk:"use_caps"`
	DNSRecordCache           types.Bool   `tfsdk:"dns_record_cache"`
	DisableAutoAccessControl types.Bool   `tfsdk:"disable_auto_access_control"`
	DisableAutoHostEntries   types.Bool   `tfsdk:"disable_auto_host_entries"`
	DNS64                    types.Bool   `tfsdk:"dns64"`
	DNS64Prefix              types.String `tfsdk:"dns64_prefix"`
	DNS64Netbits             types.String `tfsdk:"dns64_netbits"`
	MsgCacheSize             types.String `tfsdk:"msg_cache_size"`
	OutgoingNumTCP           types.String `tfsdk:"outgoing_num_tcp"`
	IncomingNumTCP           types.String `tfsdk:"incoming_num_tcp"`
	EDNSBufferSize           types.String `tfsdk:"edns_buffer_size"`
	NumQueriesPerThread      types.String `tfsdk:"num_queries_per_thread"`
	JostleTimeout            types.String `tfsdk:"jostle_timeout"`
	CacheMaxTTL              types.String `tfsdk:"cache_max_ttl"`
	CacheMinTTL              types.String `tfsdk:"cache_min_ttl"`
	InfraKeepProbing         types.Bool   `tfsdk:"infra_keep_probing"`
	InfraHostTTL             types.String `tfsdk:"infra_host_ttl"`
	InfraCacheNumHosts       types.String `tfsdk:"infra_cache_num_hosts"`
	UnwantedReplyThreshold   types.String `tfsdk:"unwanted_reply_threshold"`
	LogVerbosity             types.String `tfsdk:"log_verbosity"`
	SockQueueTimeout         types.String `tfsdk:"sock_queue_timeout"`
}

func (DNSResolverAdvancedModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"hide_identity": {
			Description: "Hide the identity (hostname) of the DNS resolver from queries. Defaults to 'false'.",
		},
		"hide_version": {
			Description: "Hide the version of the DNS resolver from queries. Defaults to 'false'.",
		},
		"prefetch": {
			Description: "Prefetch DNS cache entries before they expire. Defaults to 'false'.",
		},
		"prefetch_key": {
			Description: "Prefetch DNSSEC keys earlier in the validation process. Defaults to 'false'.",
		},
		"dnssec_stripped": {
			Description: "Harden DNSSEC data — require DNSSEC data for trust-anchored zones. Defaults to 'false'.",
		},
		"aggressive_nsec": {
			Description: "Use aggressive NSEC/NSEC3 to synthesize NXDOMAIN and NODATA responses. Defaults to 'false'.",
		},
		"qname_minimisation": {
			Description: "Send minimum amount of information to upstream servers to enhance privacy. Defaults to 'false'.",
		},
		"qname_minimisation_strict": {
			Description: "Do not fall back to sending full QNAME to upstream servers. Defaults to 'false'.",
		},
		"use_caps": {
			Description: "Use 0x20-encoded random bits in the DNS query to foil spoof attempts. Defaults to 'false'.",
		},
		"dns_record_cache": {
			Description: "Serve DNS records from the cache even if they have expired. Defaults to 'false'.",
		},
		"disable_auto_access_control": {
			Description: "Disable the automatic addition of access control entries based on listening interfaces. Defaults to 'false'.",
		},
		"disable_auto_host_entries": {
			Description: "Disable the automatic addition of host entries for the system hostname and configured interfaces. Defaults to 'false'.",
		},
		"dns64": {
			Description: "Enable DNS64 support. Defaults to 'false'.",
		},
		"dns64_prefix": {
			Description: "DNS64 prefix (e.g., '64:ff9b::/96').",
		},
		"dns64_netbits": {
			Description: "DNS64 prefix netbits.",
		},
		"msg_cache_size": {
			Description: "Message cache size in MB. Defaults to '4'.",
		},
		"outgoing_num_tcp": {
			Description: "Number of outgoing TCP buffers per thread. Defaults to '10'.",
		},
		"incoming_num_tcp": {
			Description: "Number of incoming TCP buffers per thread. Defaults to '10'.",
		},
		"edns_buffer_size": {
			Description: "EDNS reassembly buffer size in bytes. Defaults to 'auto'.",
		},
		"num_queries_per_thread": {
			Description: "Number of queries that every thread will service simultaneously. Defaults to '512'.",
		},
		"jostle_timeout": {
			Description: "Timeout in milliseconds used when the server is very busy. Defaults to '200'.",
		},
		"cache_max_ttl": {
			Description: "Maximum time to live (seconds) for RRsets and messages in the cache. Defaults to '86400'.",
		},
		"cache_min_ttl": {
			Description: "Minimum time to live (seconds) for RRsets and messages in the cache. Defaults to '0'.",
		},
		"infra_keep_probing": {
			Description: "Keep probing infrastructure hosts that are down. Defaults to 'false'.",
		},
		"infra_host_ttl": {
			Description: "Time to live (seconds) for entries in the infrastructure host cache. Defaults to '900'.",
		},
		"infra_cache_num_hosts": {
			Description: "Number of infrastructure hosts for which information is cached. Defaults to '10000'.",
		},
		"unwanted_reply_threshold": {
			Description: "Total number of unwanted replies to keep track of in every thread. Defaults to 'disabled'.",
		},
		"log_verbosity": {
			Description: "Log verbosity level (0-5). Defaults to '1'.",
		},
		"sock_queue_timeout": {
			Description: "Socket queue timeout in seconds. Defaults to '0'.",
		},
	}
}

func (DNSResolverAdvancedModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hide_identity":               types.BoolType,
		"hide_version":                types.BoolType,
		"prefetch":                    types.BoolType,
		"prefetch_key":                types.BoolType,
		"dnssec_stripped":             types.BoolType,
		"aggressive_nsec":             types.BoolType,
		"qname_minimisation":          types.BoolType,
		"qname_minimisation_strict":   types.BoolType,
		"use_caps":                    types.BoolType,
		"dns_record_cache":            types.BoolType,
		"disable_auto_access_control": types.BoolType,
		"disable_auto_host_entries":   types.BoolType,
		"dns64":                       types.BoolType,
		"dns64_prefix":                types.StringType,
		"dns64_netbits":               types.StringType,
		"msg_cache_size":              types.StringType,
		"outgoing_num_tcp":            types.StringType,
		"incoming_num_tcp":            types.StringType,
		"edns_buffer_size":            types.StringType,
		"num_queries_per_thread":      types.StringType,
		"jostle_timeout":              types.StringType,
		"cache_max_ttl":               types.StringType,
		"cache_min_ttl":               types.StringType,
		"infra_keep_probing":          types.BoolType,
		"infra_host_ttl":              types.StringType,
		"infra_cache_num_hosts":       types.StringType,
		"unwanted_reply_threshold":    types.StringType,
		"log_verbosity":               types.StringType,
		"sock_queue_timeout":          types.StringType,
	}
}

func (m *DNSResolverAdvancedModel) Set(_ context.Context, da pfsense.DNSResolverAdvanced) diag.Diagnostics {
	var diags diag.Diagnostics

	m.HideIdentity = types.BoolValue(da.HideIdentity)
	m.HideVersion = types.BoolValue(da.HideVersion)
	m.Prefetch = types.BoolValue(da.Prefetch)
	m.PrefetchKey = types.BoolValue(da.PrefetchKey)
	m.DNSSECStripped = types.BoolValue(da.DNSSECStripped)
	m.AggressiveNSEC = types.BoolValue(da.AggressiveNSEC)
	m.QNameMinimisation = types.BoolValue(da.QNameMinimisation)
	m.QNameMinimisationStrict = types.BoolValue(da.QNameMinimisationStrict)
	m.UseCaps = types.BoolValue(da.UseCaps)
	m.DNSRecordCache = types.BoolValue(da.DNSRecordCache)
	m.DisableAutoAccessControl = types.BoolValue(da.DisableAutoAccessControl)
	m.DisableAutoHostEntries = types.BoolValue(da.DisableAutoHostEntries)
	m.DNS64 = types.BoolValue(da.DNS64)

	if da.DNS64Prefix != "" {
		m.DNS64Prefix = types.StringValue(da.DNS64Prefix)
	} else {
		m.DNS64Prefix = types.StringNull()
	}

	if da.DNS64Netbits != "" {
		m.DNS64Netbits = types.StringValue(da.DNS64Netbits)
	} else {
		m.DNS64Netbits = types.StringNull()
	}

	m.MsgCacheSize = types.StringValue(da.MsgCacheSize)
	m.OutgoingNumTCP = types.StringValue(da.OutgoingNumTCP)
	m.IncomingNumTCP = types.StringValue(da.IncomingNumTCP)
	m.EDNSBufferSize = types.StringValue(da.EDNSBufferSize)
	m.NumQueriesPerThread = types.StringValue(da.NumQueriesPerThread)
	m.JostleTimeout = types.StringValue(da.JostleTimeout)
	m.CacheMaxTTL = types.StringValue(da.CacheMaxTTL)
	m.CacheMinTTL = types.StringValue(da.CacheMinTTL)
	m.InfraKeepProbing = types.BoolValue(da.InfraKeepProbing)
	m.InfraHostTTL = types.StringValue(da.InfraHostTTL)
	m.InfraCacheNumHosts = types.StringValue(da.InfraCacheNumHosts)
	m.UnwantedReplyThreshold = types.StringValue(da.UnwantedReplyThreshold)
	m.LogVerbosity = types.StringValue(da.LogVerbosity)
	m.SockQueueTimeout = types.StringValue(da.SockQueueTimeout)

	return diags
}

func (m DNSResolverAdvancedModel) Value(_ context.Context, da *pfsense.DNSResolverAdvanced) diag.Diagnostics {
	var diags diag.Diagnostics

	da.HideIdentity = m.HideIdentity.ValueBool()
	da.HideVersion = m.HideVersion.ValueBool()
	da.Prefetch = m.Prefetch.ValueBool()
	da.PrefetchKey = m.PrefetchKey.ValueBool()
	da.DNSSECStripped = m.DNSSECStripped.ValueBool()
	da.AggressiveNSEC = m.AggressiveNSEC.ValueBool()
	da.QNameMinimisation = m.QNameMinimisation.ValueBool()
	da.QNameMinimisationStrict = m.QNameMinimisationStrict.ValueBool()
	da.UseCaps = m.UseCaps.ValueBool()
	da.DNSRecordCache = m.DNSRecordCache.ValueBool()
	da.DisableAutoAccessControl = m.DisableAutoAccessControl.ValueBool()
	da.DisableAutoHostEntries = m.DisableAutoHostEntries.ValueBool()
	da.DNS64 = m.DNS64.ValueBool()

	if !m.DNS64Prefix.IsNull() {
		da.DNS64Prefix = m.DNS64Prefix.ValueString()
	}

	if !m.DNS64Netbits.IsNull() {
		da.DNS64Netbits = m.DNS64Netbits.ValueString()
	}

	da.MsgCacheSize = m.MsgCacheSize.ValueString()
	da.OutgoingNumTCP = m.OutgoingNumTCP.ValueString()
	da.IncomingNumTCP = m.IncomingNumTCP.ValueString()
	da.EDNSBufferSize = m.EDNSBufferSize.ValueString()
	da.NumQueriesPerThread = m.NumQueriesPerThread.ValueString()
	da.JostleTimeout = m.JostleTimeout.ValueString()
	da.CacheMaxTTL = m.CacheMaxTTL.ValueString()
	da.CacheMinTTL = m.CacheMinTTL.ValueString()
	da.InfraKeepProbing = m.InfraKeepProbing.ValueBool()
	da.InfraHostTTL = m.InfraHostTTL.ValueString()
	da.InfraCacheNumHosts = m.InfraCacheNumHosts.ValueString()
	da.UnwantedReplyThreshold = m.UnwantedReplyThreshold.ValueString()
	da.LogVerbosity = m.LogVerbosity.ValueString()
	da.SockQueueTimeout = m.SockQueueTimeout.ValueString()

	return diags
}
