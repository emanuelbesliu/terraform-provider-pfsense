package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = (*RouteResource)(nil)
	_ resource.ResourceWithConfigure   = (*RouteResource)(nil)
	_ resource.ResourceWithImportState = (*RouteResource)(nil)
)

type RouteResourceModel struct {
	RouteModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewRouteResource() resource.Resource { //nolint:ireturn
	return &RouteResource{}
}

type RouteResource struct {
	client *pfsense.Client
}

func (r *RouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_static_route", req.ProviderTypeName)
}

func (r *RouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Static route for directing traffic to a specific network through a gateway.",
		MarkdownDescription: "[Static route](https://docs.netgate.com/pfsense/en/latest/routing/static.html) for directing traffic to a specific network through a gateway.",
		Attributes: map[string]schema.Attribute{
			"network": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["network"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringIsNetwork(),
				},
			},
			"gateway": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["gateway"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: RouteModel{}.descriptions()["disabled"].Description,
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

func (r *RouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *RouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var routeReq pfsense.Route
	resp.Diagnostics.Append(data.Value(ctx, &routeReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.client.CreateRoute(ctx, routeReq)
	if addError(&resp.Diagnostics, "Error creating static route", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *route)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyRouteChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying route changes", err)
	}
}

func (r *RouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.client.GetRoute(ctx, data.Network.ValueString(), data.Gateway.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading static route", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *route)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *RouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *RouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var routeReq pfsense.Route
	resp.Diagnostics.Append(data.Value(ctx, &routeReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.client.UpdateRoute(ctx, routeReq, state.Network.ValueString(), state.Gateway.ValueString())
	if addError(&resp.Diagnostics, "Error updating static route", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *route)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyRouteChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying route changes", err)
	}
}

func (r *RouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRoute(ctx, data.Network.ValueString(), data.Gateway.ValueString())
	if addError(&resp.Diagnostics, "Error deleting static route", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyRouteChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying route changes", err)
	}
}

func (r *RouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: "network/gateway" (e.g., "10.0.0.0/24/WAN_GW")
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'network_cidr/gateway_name' (e.g., '10.0.0.0/24/WAN_GW'), got: %s", req.ID),
		)

		return
	}

	network := parts[0] + "/" + parts[1]
	gateway := parts[2]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network"), network)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("gateway"), gateway)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
