package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*CronJobDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CronJobDataSource)(nil)
)

type CronJobDataSourceModel struct {
	CronJobModel
}

func NewCronJobDataSource() datasource.DataSource { //nolint:ireturn
	return &CronJobDataSource{}
}

type CronJobDataSource struct {
	client *pfsense.Client
}

func (d *CronJobDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cron_job", req.ProviderTypeName)
}

func (d *CronJobDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single cron job by its command. Cron jobs are scheduled tasks that run commands at specified intervals.",
		MarkdownDescription: "Retrieves a single [cron job](https://docs.netgate.com/pfsense/en/latest/packages/cron.html) by its command. Cron jobs are scheduled tasks that run commands at specified intervals.",
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				Description: CronJobModel{}.descriptions()["command"].Description,
				Required:    true,
			},
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
		},
	}
}

func (d *CronJobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CronJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CronJobDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	job, err := d.client.GetCronJob(ctx, data.Command.ValueString())
	if addError(&resp.Diagnostics, "Unable to get cron job", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *job)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
