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
			"floating": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["floating"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"direction": schema.StringAttribute{
				Description: FirewallRuleModel{}.descriptions()["direction"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"quick": schema.BoolAttribute{
				Description: FirewallRuleModel{}.descriptions()["quick"].Description,
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

	if pfsense.BulkEnabled() && pfsense.RuleSupportedForBulk(ruleReq) && ruleReq.Description != "" {
		ph := r.client.EnqueueCreate(ruleReq)
		if r.client.PendingCount() >= 10 {
			if ferr := r.client.FlushPending(ctx); ferr != nil {
				addWarning(&resp.Diagnostics, "Bulk flush (create threshold) failed", ferr)
			}
		}
		if real, ok := r.client.ResolvePlaceholder(ph); ok {
			data.Tracker = types.StringValue(real)
		} else {
			data.Tracker = types.StringValue(ph)
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	if pfsense.BulkEnabled() && r.client.PendingCount() > 0 {
		if ferr := r.client.FlushPending(ctx); ferr != nil {
			addWarning(&resp.Diagnostics, "Bulk flush failed", ferr)
			return
		}
	}
	if pfsense.BulkEnabled() {
		tv := data.Tracker.ValueString()
		if len(tv) >= 8 && tv[:8] == "pending-" {
			if real, ok := r.client.ResolvePlaceholder(tv); ok {
				data.Tracker = types.StringValue(real)
			} else {
				if ferr := r.client.EnsureTrackerCache(ctx); ferr != nil {
					addWarning(&resp.Diagnostics, "Bulk tracker cache refresh failed", ferr)
				}
				if real := r.client.LookupTracker(data.Description.ValueString()); real != "" {
					data.Tracker = types.StringValue(real)
				} else if r.client.PendingCount() == 0 {
					// Dead placeholder: this provider process never enqueued it and no
					// live rule matches the description - removed out-of-band. Drop it
					// so plan self-heals with a clean recreate rather than keeping a
					// never-resolving id forever.
					resp.State.RemoveResource(ctx)
					return
				} else {
					addWarning(&resp.Diagnostics, "Bulk pending rule unresolved", fmt.Errorf("no live rule matches description %q after flush", data.Description.ValueString()))
					return
				}
			}
		}
	}

	if pfsense.BulkEnabled() {
		if errc := r.client.EnsureTrackerCache(ctx); errc == nil {
			if cached, found := r.client.BulkLookupRule(data.Tracker.ValueString()); found {
				resp.Diagnostics.Append(data.Set(ctx, *cached)...)
				if resp.Diagnostics.HasError() {
					return
				}
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				return
			}
		}
	}

	tv := data.Tracker.ValueString()
	if pfsense.BulkEnabled() && len(tv) >= 8 && tv[:8] == "pending-" {
		addWarning(&resp.Diagnostics, "Rule still queued in bulk apply mode",
			fmt.Errorf("tracker %q is a bulk placeholder; it will be resolved on the next refresh", tv))

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

	tvv := state.Tracker.ValueString()
	if pfsense.BulkEnabled() && pfsense.RuleSupportedForBulk(ruleReq) && ruleReq.Description != "" &&
		!(len(tvv) >= 8 && tvv[:8] == "pending-") {
		r.client.EnqueueUpdate(ruleReq)
		if r.client.PendingCount() >= 10 {
			if ferr := r.client.FlushPending(ctx); ferr != nil {
				addWarning(&resp.Diagnostics, "Bulk flush (update threshold) failed", ferr)
			}
		}
		data.Tracker = state.Tracker
		if cached, found := r.client.BulkLookupRule(state.Tracker.ValueString()); found {
			resp.Diagnostics.Append(data.Set(ctx, *cached)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

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

	// Floating rules must be deleted through the GUI "FloatingRules" tab;
	// their stored interface (e.g. "wan,lan") is not a tab name, and the
	// pfSense delete handler skips rules whose interface mismatches the tab.
	deleteIface := data.Interface.ValueString()
	if data.Floating.ValueBool() {
		deleteIface = "FloatingRules"
	}

	if pfsense.BulkEnabled() && !data.Floating.ValueBool() {
		tv := data.Tracker.ValueString()
		if len(tv) >= 8 && tv[:8] == "pending-" {
			if r.client.PendingCount() > 0 {
				if ferr := r.client.FlushPending(ctx); ferr != nil {
					addWarning(&resp.Diagnostics, "Bulk flush (pending-resolve) failed", ferr)
				}
			}
			if real, ok := r.client.ResolvePlaceholder(tv); ok {
				tv = real
			} else {
				if ferr := r.client.EnsureTrackerCache(ctx); ferr != nil {
					addWarning(&resp.Diagnostics, "Bulk tracker cache refresh failed", ferr)
					return
				}
				real := r.client.LookupTracker(data.Description.ValueString())
				if real == "" {
					addWarning(&resp.Diagnostics, "Bulk delete unresolved", fmt.Errorf("no live rule matches %q; resource kept in state", data.Description.ValueString()))
					return
				}
				tv = real
			}
		}
		r.client.EnqueueDelete(tv)
		if r.client.PendingCount() >= 2 {
			if ferr := r.client.FlushPending(ctx); ferr != nil {
				addWarning(&resp.Diagnostics, "Bulk flush (delete threshold) failed", ferr)
			}
		}
		resp.State.RemoveResource(ctx)
		return
	}

	err := r.client.DeleteFirewallRule(ctx, data.Tracker.ValueString(), deleteIface)
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
