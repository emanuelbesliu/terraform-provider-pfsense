package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*GatewayResource)(nil)
	_ resource.ResourceWithConfigure   = (*GatewayResource)(nil)
	_ resource.ResourceWithImportState = (*GatewayResource)(nil)
)

type GatewayResourceModel struct {
	GatewayModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewGatewayResource() resource.Resource { //nolint:ireturn
	return &GatewayResource{}
}

type GatewayResource struct {
	client *pfsense.Client
}

func (r *GatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_gateway", req.ProviderTypeName)
}

func (r *GatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gateway for routing traffic. Gateways are used by static routes and can be organized into gateway groups for failover and load balancing.",
		MarkdownDescription: "[Gateway](https://docs.netgate.com/pfsense/en/latest/routing/gateways.html) for routing traffic. Gateways are used by static routes and can be organized into gateway groups for failover and load balancing.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"interface": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"ipprotocol": schema.StringAttribute{
				Description:         GatewayModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: GatewayModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.Gateway{}.IPProtocols()...),
				},
			},
			"gateway": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["gateway"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["disabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"monitor": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["monitor"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("Any"),
				},
			},
			"monitor_disable": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["monitor_disable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"action_disable": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["action_disable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"force_down": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["force_down"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"weight": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["weight"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayWeight)),
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinGatewayWeight), int64(pfsense.MaxGatewayWeight)),
				},
			},
			"non_local_gateway": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["non_local_gateway"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"default_gw": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["default_gw"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"latency_low": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["latency_low"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayLatencyLow)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"latency_high": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["latency_high"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayLatencyHigh)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"loss_low": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_low"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayLossLow)),
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"loss_high": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_high"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayLossHigh)),
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
			},
			"interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayInterval)),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"loss_interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayLossInterval)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"time_period": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["time_period"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayTimePeriod)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"alert_interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["alert_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayAlertInterval)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"data_payload": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["data_payload"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultGatewayDataPayload)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
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

func (r *GatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *GatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GatewayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var gwReq pfsense.Gateway
	resp.Diagnostics.Append(data.Value(ctx, &gwReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gw, err := r.client.CreateGateway(ctx, gwReq)
	if addError(&resp.Diagnostics, "Error creating gateway", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *gw)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway changes", err)
	}
}

func (r *GatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GatewayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gw, err := r.client.GetGateway(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading gateway", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *gw)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GatewayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var gwReq pfsense.Gateway
	resp.Diagnostics.Append(data.Value(ctx, &gwReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gw, err := r.client.UpdateGateway(ctx, gwReq)
	if addError(&resp.Diagnostics, "Error updating gateway", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *gw)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway changes", err)
	}
}

func (r *GatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GatewayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGateway(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting gateway", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway changes", err)
	}
}

func (r *GatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
