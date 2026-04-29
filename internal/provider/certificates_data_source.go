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
	_ datasource.DataSource              = (*CertificatesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CertificatesDataSource)(nil)
)

type CertificatesModel struct {
	Certificates types.List `tfsdk:"certificates"`
}

func NewCertificatesDataSource() datasource.DataSource { //nolint:ireturn
	return &CertificatesDataSource{}
}

type CertificatesDataSource struct {
	client *pfsense.Client
}

func (m *CertificatesModel) Set(ctx context.Context, certs pfsense.Certificates) diag.Diagnostics {
	var diags diag.Diagnostics

	certModels := []CertificateModel{}
	for _, c := range certs {
		var certModel CertificateModel
		diags.Append(certModel.Set(ctx, c)...)
		certModels = append(certModels, certModel)
	}

	certsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: CertificateModel{}.AttrTypes()}, certModels)
	diags.Append(newDiags...)
	m.Certificates = certsValue

	return diags
}

func (d *CertificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificates", req.ProviderTypeName)
}

func (d *CertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all certificates from the pfSense certificate manager.",
		MarkdownDescription: "Retrieves all [certificates](https://docs.netgate.com/pfsense/en/latest/certificates/index.html) from the pfSense certificate manager.",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Description: "List of certificates.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"refid": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["refid"].Description,
							Computed:    true,
						},
						"descr": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["descr"].Description,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["type"].Description,
							Computed:    true,
						},
						"caref": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["caref"].Description,
							Computed:    true,
						},
						"certificate": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["certificate"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"private_key": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["private_key"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"subject": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["subject"].Description,
							Computed:    true,
						},
						"issuer": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["issuer"].Description,
							Computed:    true,
						},
						"serial": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["serial"].Description,
							Computed:    true,
						},
						"has_private_key": schema.BoolAttribute{
							Description: CertificateModel{}.descriptions()["has_private_key"].Description,
							Computed:    true,
						},
						"is_self_signed": schema.BoolAttribute{
							Description: CertificateModel{}.descriptions()["is_self_signed"].Description,
							Computed:    true,
						},
						"valid_from": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["valid_from"].Description,
							Computed:    true,
						},
						"valid_to": schema.StringAttribute{
							Description: CertificateModel{}.descriptions()["valid_to"].Description,
							Computed:    true,
						},
						"in_use": schema.BoolAttribute{
							Description: CertificateModel{}.descriptions()["in_use"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CertificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CertificatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificatesModel

	certs, err := d.client.GetCertificates(ctx)
	if addError(&resp.Diagnostics, "Unable to get certificates", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *certs)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
