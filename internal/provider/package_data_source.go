package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*PackageDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*PackageDataSource)(nil)
)

type PackageDataSourceModel struct {
	PackageModel
}

func NewPackageDataSource() datasource.DataSource { //nolint:ireturn
	return &PackageDataSource{}
}

type PackageDataSource struct {
	client *pfsense.Client
}

func (d *PackageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_package", req.ProviderTypeName)
}

func (d *PackageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single installed pfSense package by name.",
		MarkdownDescription: "Retrieves a single installed pfSense package by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"installed_version": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["installed_version"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *PackageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *PackageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PackageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pkg, err := d.client.GetPackage(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get package", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *pkg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
