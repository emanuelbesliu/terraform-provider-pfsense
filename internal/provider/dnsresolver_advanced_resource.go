package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*DNSResolverAdvancedResource)(nil)
	_ resource.ResourceWithConfigure   = (*DNSResolverAdvancedResource)(nil)
	_ resource.ResourceWithImportState = (*DNSResolverAdvancedResource)(nil)
)

type DNSResolverAdvancedResourceModel struct {
	DNSResolverAdvancedModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewDNSResolverAdvancedResource() resource.Resource { //nolint:ireturn
	return &DNSResolverAdvancedResource{}
}

type DNSResolverAdvancedResource struct {
	client *pfsense.Client
}

func (r *DNSResolverAdvancedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_advanced", req.ProviderTypeName)
}

func (r *DNSResolverAdvancedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "DNS resolver (Unbound) advanced settings configuration. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[DNS resolver (Unbound)](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html) advanced settings configuration. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			"hide_identity": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["hide_identity"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedHideIdentity),
			},
			"hide_version": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["hide_version"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedHideVersion),
			},
			"prefetch": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["prefetch"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedPrefetch),
			},
			"prefetch_key": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["prefetch_key"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedPrefetchKey),
			},
			"dnssec_stripped": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dnssec_stripped"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedDNSSECStripped),
			},
			"aggressive_nsec": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["aggressive_nsec"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedAggressiveNSEC),
			},
			"qname_minimisation": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["qname_minimisation"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedQNameMinimisation),
			},
			"qname_minimisation_strict": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["qname_minimisation_strict"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedQNameMinimisationStrict),
			},
			"use_caps": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["use_caps"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedUseCaps),
			},
			"dns_record_cache": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns_record_cache"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedDNSRecordCache),
			},
			"disable_auto_access_control": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["disable_auto_access_control"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedDisableAutoAccessControl),
			},
			"disable_auto_host_entries": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["disable_auto_host_entries"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedDisableAutoHostEntries),
			},
			"dns64": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedDNS64),
			},
			"dns64_prefix": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64_prefix"].Description,
				Optional:    true,
			},
			"dns64_netbits": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["dns64_netbits"].Description,
				Optional:    true,
			},
			"msg_cache_size": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["msg_cache_size"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedMsgCacheSize),
			},
			"outgoing_num_tcp": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["outgoing_num_tcp"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedOutgoingNumTCP),
			},
			"incoming_num_tcp": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["incoming_num_tcp"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedIncomingNumTCP),
			},
			"edns_buffer_size": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["edns_buffer_size"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedEDNSBufferSize),
			},
			"num_queries_per_thread": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["num_queries_per_thread"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedNumQueriesPerThread),
			},
			"jostle_timeout": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["jostle_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedJostleTimeout),
			},
			"cache_max_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["cache_max_ttl"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedCacheMaxTTL),
			},
			"cache_min_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["cache_min_ttl"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedCacheMinTTL),
			},
			"infra_keep_probing": schema.BoolAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_keep_probing"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverAdvancedInfraKeepProbing),
			},
			"infra_host_ttl": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_host_ttl"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedInfraHostTTL),
			},
			"infra_cache_num_hosts": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["infra_cache_num_hosts"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedInfraCacheNumHosts),
			},
			"unwanted_reply_threshold": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["unwanted_reply_threshold"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedUnwantedReplyThreshold),
			},
			"log_verbosity": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["log_verbosity"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedLogVerbosity),
			},
			"sock_queue_timeout": schema.StringAttribute{
				Description: DNSResolverAdvancedModel{}.descriptions()["sock_queue_timeout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultDNSResolverAdvancedSockQueueTimeout),
			},
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *DNSResolverAdvancedResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *DNSResolverAdvancedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DNSResolverAdvancedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var daReq pfsense.DNSResolverAdvanced
	resp.Diagnostics.Append(data.Value(ctx, &daReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	da, err := r.client.UpdateDNSResolverAdvanced(ctx, daReq)
	if addError(&resp.Diagnostics, "Error creating DNS resolver advanced settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *da)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverAdvancedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DNSResolverAdvancedResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	da, err := r.client.GetDNSResolverAdvanced(ctx)
	if addError(&resp.Diagnostics, "Error reading DNS resolver advanced settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *da)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSResolverAdvancedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DNSResolverAdvancedResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var daReq pfsense.DNSResolverAdvanced
	resp.Diagnostics.Append(data.Value(ctx, &daReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	da, err := r.client.UpdateDNSResolverAdvanced(ctx, daReq)
	if addError(&resp.Diagnostics, "Error updating DNS resolver advanced settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *da)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverAdvancedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DNSResolverAdvancedResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultDA := pfsense.DNSResolverAdvanced{
		HideIdentity:             pfsense.DefaultDNSResolverAdvancedHideIdentity,
		HideVersion:              pfsense.DefaultDNSResolverAdvancedHideVersion,
		Prefetch:                 pfsense.DefaultDNSResolverAdvancedPrefetch,
		PrefetchKey:              pfsense.DefaultDNSResolverAdvancedPrefetchKey,
		DNSSECStripped:           pfsense.DefaultDNSResolverAdvancedDNSSECStripped,
		AggressiveNSEC:           pfsense.DefaultDNSResolverAdvancedAggressiveNSEC,
		QNameMinimisation:        pfsense.DefaultDNSResolverAdvancedQNameMinimisation,
		QNameMinimisationStrict:  pfsense.DefaultDNSResolverAdvancedQNameMinimisationStrict,
		UseCaps:                  pfsense.DefaultDNSResolverAdvancedUseCaps,
		DNSRecordCache:           pfsense.DefaultDNSResolverAdvancedDNSRecordCache,
		DisableAutoAccessControl: pfsense.DefaultDNSResolverAdvancedDisableAutoAccessControl,
		DisableAutoHostEntries:   pfsense.DefaultDNSResolverAdvancedDisableAutoHostEntries,
		DNS64:                    pfsense.DefaultDNSResolverAdvancedDNS64,
		MsgCacheSize:             pfsense.DefaultDNSResolverAdvancedMsgCacheSize,
		OutgoingNumTCP:           pfsense.DefaultDNSResolverAdvancedOutgoingNumTCP,
		IncomingNumTCP:           pfsense.DefaultDNSResolverAdvancedIncomingNumTCP,
		EDNSBufferSize:           pfsense.DefaultDNSResolverAdvancedEDNSBufferSize,
		NumQueriesPerThread:      pfsense.DefaultDNSResolverAdvancedNumQueriesPerThread,
		JostleTimeout:            pfsense.DefaultDNSResolverAdvancedJostleTimeout,
		CacheMaxTTL:              pfsense.DefaultDNSResolverAdvancedCacheMaxTTL,
		CacheMinTTL:              pfsense.DefaultDNSResolverAdvancedCacheMinTTL,
		InfraKeepProbing:         pfsense.DefaultDNSResolverAdvancedInfraKeepProbing,
		InfraHostTTL:             pfsense.DefaultDNSResolverAdvancedInfraHostTTL,
		InfraCacheNumHosts:       pfsense.DefaultDNSResolverAdvancedInfraCacheNumHosts,
		UnwantedReplyThreshold:   pfsense.DefaultDNSResolverAdvancedUnwantedReplyThreshold,
		LogVerbosity:             pfsense.DefaultDNSResolverAdvancedLogVerbosity,
		SockQueueTimeout:         pfsense.DefaultDNSResolverAdvancedSockQueueTimeout,
	}

	_, err := r.client.UpdateDNSResolverAdvanced(ctx, defaultDA)
	if addError(&resp.Diagnostics, "Error resetting DNS resolver advanced settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverAdvancedResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	da, err := r.client.GetDNSResolverAdvanced(ctx)
	if addError(&resp.Diagnostics, "Error importing DNS resolver advanced settings", err) {
		return
	}

	var data DNSResolverAdvancedResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *da)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
