package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*CronJobsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CronJobsDataSource)(nil)
)

type CronJobsModel struct {
	CronJobs types.List `tfsdk:"cron_jobs"`
}

func NewCronJobsDataSource() datasource.DataSource { //nolint:ireturn
	return &CronJobsDataSource{}
}

type CronJobsDataSource struct {
	client *pfsense.Client
}

func (m *CronJobsModel) Set(ctx context.Context, jobs pfsense.CronJobs) diag.Diagnostics {
	var diags diag.Diagnostics

	jobModels := []CronJobModel{}
	for _, j := range jobs {
		var jobModel CronJobModel
		diags.Append(jobModel.Set(ctx, j)...)
		jobModels = append(jobModels, jobModel)
	}

	jobsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CronJobModel{}.AttrTypes()}, jobModels)
	diags.Append(newDiags...)
	m.CronJobs = jobsValue

	return diags
}

func (d *CronJobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cron_jobs", req.ProviderTypeName)
}

func (d *CronJobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all cron jobs. Cron jobs are scheduled tasks that run commands at specified intervals.",
		MarkdownDescription: "Retrieves all [cron jobs](https://docs.netgate.com/pfsense/en/latest/packages/cron.html). Cron jobs are scheduled tasks that run commands at specified intervals.",
		Attributes: map[string]schema.Attribute{
			"cron_jobs": schema.ListNestedAttribute{
				Description: "List of cron jobs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"minute": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["minute"].Description,
							Computed:    true,
						},
						"hour": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["hour"].Description,
							Computed:    true,
						},
						"mday": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["mday"].Description,
							Computed:    true,
						},
						"month": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["month"].Description,
							Computed:    true,
						},
						"wday": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["wday"].Description,
							Computed:    true,
						},
						"who": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["who"].Description,
							Computed:    true,
						},
						"command": schema.StringAttribute{
							Description: CronJobModel{}.descriptions()["command"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CronJobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CronJobsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CronJobsModel

	jobs, err := d.client.GetCronJobs(ctx)
	if addError(&resp.Diagnostics, "Unable to get cron jobs", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *jobs)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
