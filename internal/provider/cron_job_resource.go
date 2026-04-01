package provider

import (
	"context"
	"errors"
	"fmt"

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
	_ resource.Resource                = (*CronJobResource)(nil)
	_ resource.ResourceWithConfigure   = (*CronJobResource)(nil)
	_ resource.ResourceWithImportState = (*CronJobResource)(nil)
)

type CronJobResourceModel struct {
	CronJobModel
}

func NewCronJobResource() resource.Resource { //nolint:ireturn
	return &CronJobResource{}
}

type CronJobResource struct {
	client *pfsense.Client
}

func (r *CronJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cron_job", req.ProviderTypeName)
}

func (r *CronJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a cron job on pfSense. Cron jobs are scheduled tasks that run commands at specified intervals.",
		MarkdownDescription: "Manages a [cron job](https://docs.netgate.com/pfsense/en/latest/packages/cron.html) on pfSense. Cron jobs are scheduled tasks that run commands at specified intervals.",
		Attributes: map[string]schema.Attribute{
			"minute": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["minute"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"hour": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["hour"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mday": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["mday"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"month": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["month"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"wday": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["wday"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"who": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["who"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"command": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["command"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *CronJobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *CronJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CronJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var jobReq pfsense.CronJob
	resp.Diagnostics.Append(data.Value(ctx, &jobReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	job, err := r.client.CreateCronJob(ctx, jobReq)
	if addError(&resp.Diagnostics, "Error creating cron job", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *job)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CronJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	job, err := r.client.GetCronJob(ctx, data.Command.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading cron job", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *job)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CronJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var jobReq pfsense.CronJob
	resp.Diagnostics.Append(data.Value(ctx, &jobReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	job, err := r.client.UpdateCronJob(ctx, jobReq)
	if addError(&resp.Diagnostics, "Error updating cron job", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *job)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CronJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CronJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCronJob(ctx, data.Command.ValueString())
	if addError(&resp.Diagnostics, "Error deleting cron job", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CronJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("command"), req, resp)
}
