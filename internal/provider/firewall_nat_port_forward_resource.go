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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallNATPortForwardResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallNATPortForwardResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallNATPortForwardResource)(nil)
)

type FirewallNATPortForwardResourceModel struct {
	FirewallNATPortForwardModel
}

func NewFirewallNATPortForwardResource() resource.Resource { //nolint:ireturn
	return &FirewallNATPortForwardResource{}
}

type FirewallNATPortForwardResource struct {
	client *pfsense.Client
}

func (r *FirewallNATPortForwardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_port_forward", req.ProviderTypeName)
}

func (r *FirewallNATPortForwardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT Port Forward. Port forward rules redirect incoming traffic from external ports to internal hosts and ports.",
		MarkdownDescription: "[Firewall NAT Port Forward](https://docs.netgate.com/pfsense/en/latest/nat/port-forwards.html). Port forward rules redirect incoming traffic from external ports to internal hosts and ports.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description:         FirewallNATPortForwardModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: FirewallNATPortForwardModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATPortForward{}.IPProtocols()...),
				},
			},
			"protocol": schema.StringAttribute{
				Description:         FirewallNATPortForwardModel{}.descriptions()["protocol"].Description,
				MarkdownDescription: FirewallNATPortForwardModel{}.descriptions()["protocol"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATPortForward{}.Protocols()...),
				},
			},
			"source_address": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["source_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"source_port": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["source_port"].Description,
				Optional:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["source_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["destination_address"].Description,
				Required:    true,
			},
			"destination_port": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["destination_port"].Description,
				Optional:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["destination_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"target": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["target"].Description,
				Required:    true,
			},
			"local_port": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["local_port"].Description,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["description"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_rdr": schema.BoolAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["no_rdr"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"nat_reflection": schema.StringAttribute{
				Description:         FirewallNATPortForwardModel{}.descriptions()["nat_reflection"].Description,
				MarkdownDescription: FirewallNATPortForwardModel{}.descriptions()["nat_reflection"].MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enable", "disable", "purenat"),
				},
			},
			"associated_rule_id": schema.StringAttribute{
				Description: FirewallNATPortForwardModel{}.descriptions()["associated_rule_id"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATPortForward{}.AssociatedRuleIDOptions()...),
				},
			},
		},
	}
}

func (r *FirewallNATPortForwardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallNATPortForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallNATPortForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATPortForward
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateNATPortForward(ctx, natReq)
	if addError(&resp.Diagnostics, "Error creating NAT port forward", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATPortForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallNATPortForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetNATPortForward(ctx, data.Description.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading NAT port forward", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATPortForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallNATPortForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATPortForward
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.UpdateNATPortForward(ctx, natReq)
	if addError(&resp.Diagnostics, "Error updating NAT port forward", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATPortForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallNATPortForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNATPortForward(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Error deleting NAT port forward", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallNATPortForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("description"), req, resp)
}
