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
	_ resource.Resource                = (*FirewallNATOutboundRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallNATOutboundRuleResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallNATOutboundRuleResource)(nil)
)

type FirewallNATOutboundRuleResourceModel struct {
	FirewallNATOutboundRuleModel
}

func NewFirewallNATOutboundRuleResource() resource.Resource { //nolint:ireturn
	return &FirewallNATOutboundRuleResource{}
}

type FirewallNATOutboundRuleResource struct {
	client *pfsense.Client
}

func (r *FirewallNATOutboundRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_outbound", req.ProviderTypeName)
}

func (r *FirewallNATOutboundRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT Outbound Rule. Outbound NAT rules control how traffic leaving the firewall is translated.",
		MarkdownDescription: "[Firewall NAT Outbound](https://docs.netgate.com/pfsense/en/latest/nat/outbound.html). Outbound NAT rules control how traffic leaving the firewall is translated.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"protocol": schema.StringAttribute{
				Description:         FirewallNATOutboundRuleModel{}.descriptions()["protocol"].Description,
				MarkdownDescription: FirewallNATOutboundRuleModel{}.descriptions()["protocol"].MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOutboundRule{}.Protocols()...),
				},
			},
			"source_address": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"source_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_port"].Description,
				Optional:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"destination_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_port"].Description,
				Optional:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["destination_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"target": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"target_ip": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target_ip"].Description,
				Optional:    true,
			},
			"target_ip_subnet": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["target_ip_subnet"].Description,
				Optional:    true,
			},
			"nat_port": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["nat_port"].Description,
				Optional:    true,
			},
			"pool_options": schema.StringAttribute{
				Description:         FirewallNATOutboundRuleModel{}.descriptions()["pool_options"].Description,
				MarkdownDescription: FirewallNATOutboundRuleModel{}.descriptions()["pool_options"].MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOutboundRule{}.PoolOptions()...),
				},
			},
			"source_hash_key": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["source_hash_key"].Description,
				Optional:    true,
			},
			"static_nat_port": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["static_nat_port"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_sync": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["no_sync"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_nat": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["no_nat"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: FirewallNATOutboundRuleModel{}.descriptions()["description"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
		},
	}
}

func (r *FirewallNATOutboundRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallNATOutboundRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallNATOutboundRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATOutboundRule
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateNATOutboundRule(ctx, natReq)
	if addError(&resp.Diagnostics, "Error creating NAT outbound rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOutboundRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallNATOutboundRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetNATOutboundRule(ctx, data.Description.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading NAT outbound rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOutboundRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallNATOutboundRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATOutboundRule
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.UpdateNATOutboundRule(ctx, natReq)
	if addError(&resp.Diagnostics, "Error updating NAT outbound rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOutboundRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallNATOutboundRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNATOutboundRule(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Error deleting NAT outbound rule", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallNATOutboundRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("description"), req, resp)
}
