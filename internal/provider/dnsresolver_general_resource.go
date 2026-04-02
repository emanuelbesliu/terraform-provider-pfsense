package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*DNSResolverGeneralResource)(nil)
	_ resource.ResourceWithConfigure   = (*DNSResolverGeneralResource)(nil)
	_ resource.ResourceWithImportState = (*DNSResolverGeneralResource)(nil)
)

type DNSResolverGeneralResourceModel struct {
	DNSResolverGeneralModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewDNSResolverGeneralResource() resource.Resource { //nolint:ireturn
	return &DNSResolverGeneralResource{}
}

type DNSResolverGeneralResource struct {
	client *pfsense.Client
}

func (r *DNSResolverGeneralResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_general", req.ProviderTypeName)
}

func (r *DNSResolverGeneralResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "DNS resolver (Unbound) general settings configuration. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[DNS resolver (Unbound)](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html) general settings configuration. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralEnable),
			},
			"port": schema.Int64Attribute{
				Description: DNSResolverGeneralModel{}.descriptions()["port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultDNSResolverGeneralPort)),
			},
			"enable_ssl": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["enable_ssl"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralEnableSSL),
			},
			"tls_port": schema.Int64Attribute{
				Description: DNSResolverGeneralModel{}.descriptions()["tls_port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultDNSResolverGeneralTLSPort)),
			},
			"ssl_cert_ref": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["ssl_cert_ref"].Description,
				Optional:    true,
			},
			"active_interfaces": schema.ListAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["active_interfaces"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["active_interfaces"].MarkdownDescription,
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"outgoing_interfaces": schema.ListAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["outgoing_interfaces"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["outgoing_interfaces"].MarkdownDescription,
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"system_domain_local_zone_type": schema.StringAttribute{
				Description:         DNSResolverGeneralModel{}.descriptions()["system_domain_local_zone_type"].Description,
				MarkdownDescription: DNSResolverGeneralModel{}.descriptions()["system_domain_local_zone_type"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultDNSResolverGeneralSystemDomainLocalZone),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.DNSResolverGeneral{}.SystemDomainLocalZoneOptions()...),
				},
			},
			"dnssec": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["dnssec"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralDNSSEC),
			},
			"forwarding": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["forwarding"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralForwarding),
			},
			"forward_tls_upstream": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["forward_tls_upstream"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralForwardTLSUpstream),
			},
			"register_dhcp_leases": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_dhcp_leases"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralRegisterDHCPLeases),
			},
			"register_dhcp_static_maps": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_dhcp_static_maps"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralRegisterDHCPStaticMaps),
			},
			"register_openvpn_clients": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["register_openvpn_clients"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralRegisterOpenVPNClients),
			},
			"strict_outgoing_interface": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["strict_outgoing_interface"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralStrictOutgoingInterface),
			},
			"python": schema.BoolAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultDNSResolverGeneralPython),
			},
			"python_order": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python_order"].Description,
				Optional:    true,
			},
			"python_script": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["python_script"].Description,
				Optional:    true,
			},
			"custom_options": schema.StringAttribute{
				Description: DNSResolverGeneralModel{}.descriptions()["custom_options"].Description,
				Optional:    true,
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

func (r *DNSResolverGeneralResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *DNSResolverGeneralResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DNSResolverGeneralResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var dgReq pfsense.DNSResolverGeneral
	resp.Diagnostics.Append(data.Value(ctx, &dgReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := r.client.UpdateDNSResolverGeneral(ctx, dgReq)
	if addError(&resp.Diagnostics, "Error creating DNS resolver general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *dg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverGeneralResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DNSResolverGeneralResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := r.client.GetDNSResolverGeneral(ctx)
	if addError(&resp.Diagnostics, "Error reading DNS resolver general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *dg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSResolverGeneralResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DNSResolverGeneralResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var dgReq pfsense.DNSResolverGeneral
	resp.Diagnostics.Append(data.Value(ctx, &dgReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := r.client.UpdateDNSResolverGeneral(ctx, dgReq)
	if addError(&resp.Diagnostics, "Error updating DNS resolver general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *dg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverGeneralResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DNSResolverGeneralResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultDG := pfsense.DNSResolverGeneral{
		Enable:                  pfsense.DefaultDNSResolverGeneralEnable,
		Port:                    pfsense.DefaultDNSResolverGeneralPort,
		EnableSSL:               pfsense.DefaultDNSResolverGeneralEnableSSL,
		TLSPort:                 pfsense.DefaultDNSResolverGeneralTLSPort,
		ActiveInterfaces:        []string{"all"},
		OutgoingInterfaces:      []string{"all"},
		SystemDomainLocalZone:   pfsense.DefaultDNSResolverGeneralSystemDomainLocalZone,
		DNSSEC:                  pfsense.DefaultDNSResolverGeneralDNSSEC,
		Forwarding:              pfsense.DefaultDNSResolverGeneralForwarding,
		ForwardTLSUpstream:      pfsense.DefaultDNSResolverGeneralForwardTLSUpstream,
		RegisterDHCPLeases:      pfsense.DefaultDNSResolverGeneralRegisterDHCPLeases,
		RegisterDHCPStaticMaps:  pfsense.DefaultDNSResolverGeneralRegisterDHCPStaticMaps,
		RegisterOpenVPNClients:  pfsense.DefaultDNSResolverGeneralRegisterOpenVPNClients,
		StrictOutgoingInterface: pfsense.DefaultDNSResolverGeneralStrictOutgoingInterface,
		Python:                  pfsense.DefaultDNSResolverGeneralPython,
	}

	_, err := r.client.UpdateDNSResolverGeneral(ctx, defaultDG)
	if addError(&resp.Diagnostics, "Error resetting DNS resolver general settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDNSResolverChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying DNS resolver changes", err)
	}
}

func (r *DNSResolverGeneralResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	dg, err := r.client.GetDNSResolverGeneral(ctx)
	if addError(&resp.Diagnostics, "Error importing DNS resolver general settings", err) {
		return
	}

	var data DNSResolverGeneralResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *dg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
