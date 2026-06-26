package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallScheduleResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallScheduleResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallScheduleResource)(nil)
)

var scheduleTimeRegex = regexp.MustCompile(`^([0-9]|[01][0-9]|2[0-3]):(00|15|30|45|59)$`)

type FirewallScheduleResourceModel struct {
	FirewallScheduleModel
}

func NewFirewallScheduleResource() resource.Resource { //nolint:ireturn
	return &FirewallScheduleResource{}
}

type FirewallScheduleResource struct {
	client *pfsense.Client
}

func (r *FirewallScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_schedule", req.ProviderTypeName)
}

func (r *FirewallScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall schedule defining time ranges during which firewall rules referencing it are active.",
		MarkdownDescription: "[Firewall schedule](https://docs.netgate.com/pfsense/en/latest/firewall/time-based-rules.html) defining time ranges during which firewall rules referencing it are active.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallScheduleModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
						"must consist only of the characters a-z, A-Z, 0-9 and _",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: FirewallScheduleModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"time_range": schema.ListNestedAttribute{
				Description: FirewallScheduleModel{}.descriptions()["time_range"].Description,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"position": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["position"].Description,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[1-7](,[1-7])*$`),
									"must be a comma separated list of weekday numbers 1-7 (Monday is 1, Sunday is 7)",
								),
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("month"),
									path.MatchRelative().AtParent().AtName("day"),
								),
							},
						},
						"month": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["month"].Description,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([1-9]|1[0-2])(,([1-9]|1[0-2]))*$`),
									"must be a comma separated list of month numbers 1-12",
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("day"),
								),
							},
						},
						"day": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["day"].Description,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([1-9]|[12][0-9]|3[01])(,([1-9]|[12][0-9]|3[01]))*$`),
									"must be a comma separated list of day-of-month numbers 1-31",
								),
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("month"),
								),
							},
						},
						"start_time": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["start_time"].Description,
							Required:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(scheduleTimeRegex, "must be in 'H:MM' format with minutes 00, 15, 30, 45 or 59"),
							},
						},
						"stop_time": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["stop_time"].Description,
							Required:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(scheduleTimeRegex, "must be in 'H:MM' format with minutes 00, 15, 30, 45 or 59"),
							},
						},
						"range_description": schema.StringAttribute{
							Description: FirewallScheduleTimeRangeModel{}.descriptions()["range_description"].Description,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
					},
				},
			},
			"label": schema.StringAttribute{
				Description: FirewallScheduleModel{}.descriptions()["label"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *FirewallScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scheduleReq pfsense.Schedule
	resp.Diagnostics.Append(data.Value(ctx, &scheduleReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.CreateSchedule(ctx, scheduleReq)
	if addError(&resp.Diagnostics, "Error creating firewall schedule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *schedule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.GetSchedule(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading firewall schedule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *schedule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallScheduleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *FirewallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scheduleReq pfsense.Schedule
	resp.Diagnostics.Append(data.Value(ctx, &scheduleReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.UpdateSchedule(ctx, state.Name.ValueString(), scheduleReq)
	if addError(&resp.Diagnostics, "Error updating firewall schedule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *schedule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallScheduleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSchedule(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting firewall schedule", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *FirewallScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
