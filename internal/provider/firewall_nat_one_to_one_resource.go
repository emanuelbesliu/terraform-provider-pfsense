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
	_ resource.Resource                = (*FirewallNATOneToOneRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallNATOneToOneRuleResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallNATOneToOneRuleResource)(nil)
)

type FirewallNATOneToOneRuleResourceModel struct {
	FirewallNATOneToOneRuleModel
}

func NewFirewallNATOneToOneRuleResource() resource.Resource { //nolint:ireturn
	return &FirewallNATOneToOneRuleResource{}
}

type FirewallNATOneToOneRuleResource struct {
	client *pfsense.Client
}

func (r *FirewallNATOneToOneRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_one_to_one", req.ProviderTypeName)
}

func (r *FirewallNATOneToOneRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT 1:1 Rule. Maps an external IP address to an internal IP address (bidirectional NAT).",
		MarkdownDescription: "[Firewall NAT 1:1](https://docs.netgate.com/pfsense/en/latest/nat/1-1.html). Maps an external IP address to an internal IP address (bidirectional NAT).",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"external": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["external"].Description,
				Required:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description:         FirewallNATOneToOneRuleModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: FirewallNATOneToOneRuleModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("inet"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOneToOne{}.IPProtocols()...),
				},
			},
			"source_address": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["source_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["source_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["destination_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["destination_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_binat": schema.BoolAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["no_binat"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: FirewallNATOneToOneRuleModel{}.descriptions()["description"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"nat_reflection": schema.StringAttribute{
				Description:         FirewallNATOneToOneRuleModel{}.descriptions()["nat_reflection"].Description,
				MarkdownDescription: FirewallNATOneToOneRuleModel{}.descriptions()["nat_reflection"].MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOneToOne{}.NATReflectionModes()...),
				},
			},
		},
	}
}

func (r *FirewallNATOneToOneRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallNATOneToOneRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallNATOneToOneRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATOneToOne
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateNATOneToOne(ctx, natReq)
	if addError(&resp.Diagnostics, "Error creating NAT 1:1 rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOneToOneRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallNATOneToOneRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetNATOneToOne(ctx, data.Description.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading NAT 1:1 rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOneToOneRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallNATOneToOneRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *FirewallNATOneToOneRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var natReq pfsense.NATOneToOne
	resp.Diagnostics.Append(data.Value(ctx, &natReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.UpdateNATOneToOne(ctx, state.Description.ValueString(), natReq)
	if addError(&resp.Diagnostics, "Error updating NAT 1:1 rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNATOneToOneRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallNATOneToOneRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNATOneToOne(ctx, data.Description.ValueString())
	if addError(&resp.Diagnostics, "Error deleting NAT 1:1 rule", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallNATOneToOneRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("description"), req, resp)
}
