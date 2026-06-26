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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallNATNPtRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallNATNPtRuleResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallNATNPtRuleResource)(nil)
)

type FirewallNATNPtRuleResourceModel struct {
	FirewallNATNPtRuleModel
}

func NewFirewallNATNPtRuleResource() resource.Resource { //nolint:ireturn
	return &FirewallNATNPtRuleResource{}
}

type FirewallNATNPtRuleResource struct {
	client *pfsense.Client
}

func (r *FirewallNATNPtRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_npt", req.ProviderTypeName)
}

func (r *FirewallNATNPtRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT NPt Rule. Maps an internal IPv6 prefix to an external IPv6 prefix (IPv6-to-IPv6 Network Prefix Translation).",
		MarkdownDescription: "[Firewall NAT NPt](https://docs.netgate.com/pfsense/en/latest/nat/npt.html). Maps an internal IPv6 prefix to an external IPv6 prefix (IPv6-to-IPv6 Network Prefix Translation).",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"source_prefix": schema.StringAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["source_prefix"].Description,
				Required:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["source_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_prefix": schema.StringAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["destination_prefix"].Description,
				Required:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["destination_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: FirewallNATNPtRuleModel{}.descriptions()["description"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
		},
	}
}

func (r *FirewallNATNPtRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallNATNPtRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallNATNPtRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATNPt
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateNATNPt(ctx, natReq)
	if addError(&resp.Diagnostics, "Error creating NAT NPt rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATNPtRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallNATNPtRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetNATNPt(ctx, data.Description.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading NAT NPt rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATNPtRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallNATNPtRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *FirewallNATNPtRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATNPt
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.UpdateNATNPt(ctx, state.Description.ValueString(), natReq)
	if addError(&resp.Diagnostics, "Error updating NAT NPt rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATNPtRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallNATNPtRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNATNPt(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Error deleting NAT NPt rule", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallNATNPtRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("description"), req, resp)
}
