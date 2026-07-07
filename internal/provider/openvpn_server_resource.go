package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = (*OpenVPNServerResource)(nil)
	_ resource.ResourceWithConfigure   = (*OpenVPNServerResource)(nil)
	_ resource.ResourceWithImportState = (*OpenVPNServerResource)(nil)
)

type OpenVPNServerResourceModel struct {
	OpenVPNServerModel
}

func NewOpenVPNServerResource() resource.Resource { //nolint:ireturn
	return &OpenVPNServerResource{}
}

type OpenVPNServerResource struct {
	client *pfsense.Client
}

func (r *OpenVPNServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_openvpn_server", req.ProviderTypeName)
}

func (r *OpenVPNServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := OpenVPNServerModel{}.descriptions()

	optionalString := func(key string) schema.StringAttribute {
		return schema.StringAttribute{
			Description: descriptions[key].Description,
			Optional:    true,
		}
	}

	optionalComputedBool := func(key string) schema.BoolAttribute {
		return schema.BoolAttribute{
			Description: descriptions[key].Description,
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		}
	}

	resp.Schema = schema.Schema{
		Description:         "OpenVPN server instance. Configures an OpenVPN server for site-to-site or remote access VPNs.",
		MarkdownDescription: "OpenVPN [server](https://docs.netgate.com/pfsense/en/latest/vpn/openvpn/index.html) instance. Configures an OpenVPN server for site-to-site or remote access VPNs.\n\nChanges are applied immediately by pfSense when the instance is saved.",
		Attributes: map[string]schema.Attribute{
			"vpn_id": schema.StringAttribute{
				Description: descriptions["vpn_id"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disable": optionalComputedBool("disable"),
			"mode": schema.StringAttribute{
				Description: descriptions["mode"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("p2p_tls", "p2p_shared_key", "server_tls", "server_user", "server_tls_user"),
				},
			},
			"auth_mode": schema.ListAttribute{
				Description: descriptions["auth_mode"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"dev_mode": schema.StringAttribute{
				Description: descriptions["dev_mode"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("tun", "tap"),
				},
			},
			"protocol": schema.StringAttribute{
				Description: descriptions["protocol"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("UDP4", "UDP6", "TCP4", "TCP6"),
				},
			},
			"interface": schema.StringAttribute{
				Description: descriptions["interface"].Description,
				Required:    true,
			},
			"ip_address": optionalString("ip_address"),
			"local_port": optionalString("local_port"),
			"description": schema.StringAttribute{
				Description: descriptions["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"custom_options": optionalString("custom_options"),
			"tls":            optionalString("tls"),
			"tls_type": schema.StringAttribute{
				Description: descriptions["tls_type"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("auth", "crypt"),
				},
			},
			"tls_auth_keydir": optionalString("tls_auth_keydir"),
			"ca_ref":          optionalString("ca_ref"),
			"crl_ref":         optionalString("crl_ref"),
			"ocsp_check":      optionalComputedBool("ocsp_check"),
			"ocsp_url":        optionalString("ocsp_url"),
			"cert_ref":        optionalString("cert_ref"),
			"dh_length":       optionalString("dh_length"),
			"ecdh_curve":      optionalString("ecdh_curve"),
			"cert_depth":      optionalString("cert_depth"),
			"strict_user_cn":  optionalComputedBool("strict_user_cn"),
			"remote_cert_tls": optionalComputedBool("remote_cert_tls"),
			"shared_key": schema.StringAttribute{
				Description: descriptions["shared_key"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"data_ciphers": schema.ListAttribute{
				Description: descriptions["data_ciphers"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"data_ciphers_fallback": optionalString("data_ciphers_fallback"),
			"digest":                optionalString("digest"),
			"tunnel_network":        optionalString("tunnel_network"),
			"tunnel_network_v6":     optionalString("tunnel_network_v6"),
			"local_network":         optionalString("local_network"),
			"local_network_v6":      optionalString("local_network_v6"),
			"remote_network":        optionalString("remote_network"),
			"remote_network_v6":     optionalString("remote_network_v6"),
			"gw_redir":              optionalComputedBool("gw_redir"),
			"gw_redir_v6":           optionalComputedBool("gw_redir_v6"),
			"topology": schema.StringAttribute{
				Description: descriptions["topology"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("subnet", "net30"),
				},
			},
			"max_clients":             optionalString("max_clients"),
			"connection_limit":        optionalString("connection_limit"),
			"client_to_client":        optionalComputedBool("client_to_client"),
			"duplicate_cn":            optionalComputedBool("duplicate_cn"),
			"dynamic_ip":              optionalComputedBool("dynamic_ip"),
			"compression":             optionalString("compression"),
			"compression_push":        optionalComputedBool("compression_push"),
			"allow_compression":       optionalString("allow_compression"),
			"pass_tos":                optionalComputedBool("pass_tos"),
			"dns_domain_enable":       optionalComputedBool("dns_domain_enable"),
			"dns_domain":              optionalString("dns_domain"),
			"dns_server_enable":       optionalComputedBool("dns_server_enable"),
			"dns_server1":             optionalString("dns_server1"),
			"dns_server2":             optionalString("dns_server2"),
			"dns_server3":             optionalString("dns_server3"),
			"dns_server4":             optionalString("dns_server4"),
			"ntp_server_enable":       optionalComputedBool("ntp_server_enable"),
			"ntp_server1":             optionalString("ntp_server1"),
			"ntp_server2":             optionalString("ntp_server2"),
			"push_register_dns":       optionalComputedBool("push_register_dns"),
			"push_block_outside_dns":  optionalComputedBool("push_block_outside_dns"),
			"username_as_common_name": optionalComputedBool("username_as_common_name"),
			"create_gw":               optionalString("create_gw"),
			"verbosity_level":         optionalString("verbosity_level"),
		},
	}
}

func (r *OpenVPNServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *OpenVPNServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OpenVPNServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var serverReq pfsense.OpenVPNServer
	resp.Diagnostics.Append(data.Value(ctx, &serverReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.CreateOpenVPNServer(ctx, serverReq)
	if addError(&resp.Diagnostics, "Error creating OpenVPN server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenVPNServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OpenVPNServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.GetOpenVPNServer(ctx, data.VPNID.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading OpenVPN server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenVPNServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OpenVPNServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var serverReq pfsense.OpenVPNServer
	resp.Diagnostics.Append(data.Value(ctx, &serverReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.UpdateOpenVPNServer(ctx, serverReq)
	if addError(&resp.Diagnostics, "Error updating OpenVPN server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpenVPNServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OpenVPNServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOpenVPNServer(ctx, data.VPNID.ValueString())
	if addError(&resp.Diagnostics, "Error deleting OpenVPN server", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *OpenVPNServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vpn_id"), types.StringValue(req.ID))...)
}
