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
	_ datasource.DataSource              = (*PackagesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*PackagesDataSource)(nil)
)

type PackagesModel struct {
	Packages types.List `tfsdk:"packages"`
}

func NewPackagesDataSource() datasource.DataSource { //nolint:ireturn
	return &PackagesDataSource{}
}

type PackagesDataSource struct {
	client *pfsense.Client
}

func (m *PackagesModel) Set(ctx context.Context, pkgs pfsense.Packages) diag.Diagnostics {
	var diags diag.Diagnostics

	pkgModels := []PackageModel{}
	for _, p := range pkgs {
		var pkgModel PackageModel
		diags.Append(pkgModel.Set(ctx, p)...)
		pkgModels = append(pkgModels, pkgModel)
	}

	pkgsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: PackageModel{}.AttrTypes()}, pkgModels)
	diags.Append(newDiags...)
	m.Packages = pkgsValue

	return diags
}

func (d *PackagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_packages", req.ProviderTypeName)
}

func (d *PackagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all installed pfSense packages.",
		MarkdownDescription: "Retrieves all installed pfSense packages.",
		Attributes: map[string]schema.Attribute{
			"packages": schema.ListNestedAttribute{
				Description: "List of installed packages.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: PackageModel{}.descriptions()["name"].Description,
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *PackagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *PackagesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PackagesModel

	pkgs, err := d.client.GetPackages(ctx)
	if addError(&resp.Diagnostics, "Unable to get packages", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *pkgs)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
