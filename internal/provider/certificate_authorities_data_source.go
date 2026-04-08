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
	_ datasource.DataSource              = (*CertificateAuthoritiesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CertificateAuthoritiesDataSource)(nil)
)

type CertificateAuthoritiesModel struct {
	CertificateAuthorities types.List `tfsdk:"certificate_authorities"`
}

func NewCertificateAuthoritiesDataSource() datasource.DataSource { //nolint:ireturn
	return &CertificateAuthoritiesDataSource{}
}

type CertificateAuthoritiesDataSource struct {
	client *pfsense.Client
}

func (m *CertificateAuthoritiesModel) Set(ctx context.Context, cas pfsense.CertificateAuthorities) diag.Diagnostics {
	var diags diag.Diagnostics

	caModels := []CertificateAuthorityModel{}
	for _, ca := range cas {
		var caModel CertificateAuthorityModel
		diags.Append(caModel.Set(ctx, ca)...)
		caModels = append(caModels, caModel)
	}

	casValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CertificateAuthorityModel{}.AttrTypes()}, caModels)
	diags.Append(newDiags...)
	m.CertificateAuthorities = casValue

	return diags
}

func (d *CertificateAuthoritiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificate_authorities", req.ProviderTypeName)
}

func (d *CertificateAuthoritiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all certificate authorities (CAs) from the pfSense trust store.",
		MarkdownDescription: "Retrieves all [certificate authorities](https://docs.netgate.com/pfsense/en/latest/certificates/cas.html) (CAs) from the pfSense trust store.",
		Attributes: map[string]schema.Attribute{
			"certificate_authorities": schema.ListNestedAttribute{
				Description: "List of certificate authorities.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"refid": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["refid"].Description,
							Computed:    true,
						},
						"descr": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["descr"].Description,
							Computed:    true,
						},
						"certificate": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["certificate"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"private_key": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["private_key"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"trust": schema.BoolAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["trust"].Description,
							Computed:    true,
						},
						"random_serial": schema.BoolAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["random_serial"].Description,
							Computed:    true,
						},
						"next_serial": schema.Int64Attribute{
							Description: CertificateAuthorityModel{}.descriptions()["next_serial"].Description,
							Computed:    true,
						},
						"subject": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["subject"].Description,
							Computed:    true,
						},
						"issuer": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["issuer"].Description,
							Computed:    true,
						},
						"serial": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["serial"].Description,
							Computed:    true,
						},
						"has_private_key": schema.BoolAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["has_private_key"].Description,
							Computed:    true,
						},
						"is_self_signed": schema.BoolAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["is_self_signed"].Description,
							Computed:    true,
						},
						"valid_from": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["valid_from"].Description,
							Computed:    true,
						},
						"valid_to": schema.StringAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["valid_to"].Description,
							Computed:    true,
						},
						"in_use": schema.BoolAttribute{
							Description: CertificateAuthorityModel{}.descriptions()["in_use"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CertificateAuthoritiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CertificateAuthoritiesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificateAuthoritiesModel

	cas, err := d.client.GetCertificateAuthorities(ctx)
	if addError(&resp.Diagnostics, "Unable to get certificate authorities", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *cas)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
