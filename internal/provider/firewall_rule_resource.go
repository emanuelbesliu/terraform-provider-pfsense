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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallRuleResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallRuleResource)(nil)
)

type FirewallRuleResourceModel struct {
	FirewallRuleModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewFirewallRuleResource() resource.Resource { //nolint:ireturn
	return &FirewallRuleResource{}
}

type FirewallRuleResource struct {
	client *pfsense.Client
}

func (r *FirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_rule", req.ProviderTypeName)
}

func (r *FirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall rule for controlling traffic flow. Rules are evaluated in order and the first matching rule wins.",
		MarkdownDescription: "[Firewall rule](https://docs.netgate.com/pfsense/en/latest/firewall/index.html) for controlling traffic flow. Rules are evaluated in order and the first matching rule wins.",
		Attributes: map[string]schema.Attribute{
			"tracker": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["tracker"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["type"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["type"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.FirewallRule{}.Types()...),
				},
			},
			"interface": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"ipprotocol": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("inet"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.FirewallRule{}.IPProtocols()...),
				},
			},
			"protocol": schema.StringAttribute{
				Description:         FirewallRuleModel{}.descriptions()["protocol"].Description,
				MarkdownDescription: FirewallRuleModel{}.descriptions()["protocol"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("any"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.FirewallRule{}.Protocols()...),
				},
			},
			"source_address": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_address"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"source_port": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_port"].Description,
				Optional:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["source_not"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_address"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"destination_port": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_port"].Description,
				Optional:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["destination_not"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["description"].Description,
				Optional:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["disabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"log": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["log"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
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

func (r *FirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ruleReq pfsense.FirewallRule
	resp.Diagnostics.Append(data.Value(ctx, &ruleReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.CreateFirewallRule(ctx, ruleReq)
	if addError(&resp.Diagnostics, "Error creating firewall rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyFirewallRuleChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying firewall rule changes", err)
	}
}

func (r *FirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetFirewallRule(ctx, data.Tracker.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading firewall rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to preserve the tracker ID
	var state *FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ruleReq pfsense.FirewallRule
	resp.Diagnostics.Append(data.Value(ctx, &ruleReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Use the tracker from state for the update lookup
	ruleReq.Tracker = state.Tracker.ValueString()

	rule, err := r.client.UpdateFirewallRule(ctx, ruleReq)
	if addError(&resp.Diagnostics, "Error updating firewall rule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *rule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyFirewallRuleChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying firewall rule changes", err)
	}
}

func (r *FirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFirewallRule(ctx, data.Tracker.ValueString())
	if addError(&resp.Diagnostics, "Error deleting firewall rule", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyFirewallRuleChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying firewall rule changes", err)
	}
}

func (r *FirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tracker"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
