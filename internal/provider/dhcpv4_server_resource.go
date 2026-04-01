package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*DHCPv4ServerResource)(nil)
	_ resource.ResourceWithConfigure   = (*DHCPv4ServerResource)(nil)
	_ resource.ResourceWithImportState = (*DHCPv4ServerResource)(nil)
)

type DHCPv4ServerResourceModel struct {
	DHCPv4ServerModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewDHCPv4ServerResource() resource.Resource { //nolint:ireturn
	return &DHCPv4ServerResource{}
}

type DHCPv4ServerResource struct {
	client *pfsense.Client
}

func (r *DHCPv4ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dhcpv4_server", req.ProviderTypeName)
}

func (r *DHCPv4ServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "DHCPv4 server configuration for a single network interface. This is a singleton-per-interface resource — only one instance can exist per pfSense interface. Removing the resource from Terraform will reset the DHCP server settings to defaults.",
		MarkdownDescription: "DHCPv4 [server configuration](https://docs.netgate.com/pfsense/en/latest/services/dhcp/ipv4.html) for a single network interface. This is a **singleton-per-interface** resource — only one instance can exist per pfSense interface. Removing the resource from Terraform will reset the DHCP server settings to defaults.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["interface"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"range_from": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["range_from"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("IPv4"),
				},
			},
			"range_to": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["range_to"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("IPv4"),
				},
			},
			"dns_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["dns_servers"].Description,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(pfsense.DHCPv4ServerMaxDNSServers),
				},
			},
			"gateway": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["gateway"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("IPv4"),
				},
			},
			"domain_name": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["domain_name"].Description,
				Optional:    true,
			},
			"domain_search_list": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["domain_search_list"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"default_lease_time": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["default_lease_time"].Description,
				Optional:    true,
				CustomType:  timetypes.GoDurationType{},
			},
			"maximum_lease_time": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["maximum_lease_time"].Description,
				Optional:    true,
				CustomType:  timetypes.GoDurationType{},
			},
			"wins_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["wins_servers"].Description,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(pfsense.DHCPv4ServerMaxWINSServers),
				},
			},
			"ntp_servers": schema.ListAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ntp_servers"].Description,
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtMost(pfsense.DHCPv4ServerMaxNTPServers),
				},
			},
			"tftp_server": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["tftp_server"].Description,
				Optional:    true,
			},
			"ldap_server": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ldap_server"].Description,
				Optional:    true,
			},
			"mac_allow": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["mac_allow"].Description,
				Optional:    true,
			},
			"mac_deny": schema.StringAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["mac_deny"].Description,
				Optional:    true,
			},
			"deny_unknown": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["deny_unknown"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ignore_client_uids": schema.BoolAttribute{
				Description: DHCPv4ServerModel{}.descriptions()["ignore_client_uids"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"static_arp": schema.BoolAttribute{
				Description:         DHCPv4ServerModel{}.descriptions()["static_arp"].Description,
				MarkdownDescription: DHCPv4ServerModel{}.descriptions()["static_arp"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
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

func (r *DHCPv4ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *DHCPv4ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DHCPv4ServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var srvReq pfsense.DHCPv4Server
	resp.Diagnostics.Append(data.Value(ctx, &srvReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	srv, err := r.client.UpdateDHCPv4Server(ctx, srvReq)
	if addError(&resp.Diagnostics, "Error creating DHCPv4 server configuration", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *srv)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDHCPv4Changes(ctx, srvReq.Interface)
		addWarning(&resp.Diagnostics, "Error applying DHCPv4 changes", err)
	}
}

func (r *DHCPv4ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DHCPv4ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	srv, err := r.client.GetDHCPv4Server(ctx, data.Interface.ValueString())
	if addError(&resp.Diagnostics, "Error reading DHCPv4 server configuration", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *srv)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DHCPv4ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DHCPv4ServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var srvReq pfsense.DHCPv4Server
	resp.Diagnostics.Append(data.Value(ctx, &srvReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	srv, err := r.client.UpdateDHCPv4Server(ctx, srvReq)
	if addError(&resp.Diagnostics, "Error updating DHCPv4 server configuration", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *srv)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDHCPv4Changes(ctx, srvReq.Interface)
		addWarning(&resp.Diagnostics, "Error applying DHCPv4 changes", err)
	}
}

func (r *DHCPv4ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DHCPv4ServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults — disable DHCP, clear all optional fields.
	defaultSrv := pfsense.DHCPv4Server{
		Interface: data.Interface.ValueString(),
		Enable:    false,
	}

	_, err := r.client.UpdateDHCPv4Server(ctx, defaultSrv)
	if addError(&resp.Diagnostics, "Error resetting DHCPv4 server configuration to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyDHCPv4Changes(ctx, data.Interface.ValueString())
		addWarning(&resp.Diagnostics, "Error applying DHCPv4 changes", err)
	}
}

func (r *DHCPv4ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	iface := req.ID

	srv, err := r.client.GetDHCPv4Server(ctx, iface)
	if addError(&resp.Diagnostics, "Error importing DHCPv4 server configuration", err) {
		return
	}

	var data DHCPv4ServerResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *srv)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
