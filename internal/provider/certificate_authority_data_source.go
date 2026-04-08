package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*CertificateAuthorityDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CertificateAuthorityDataSource)(nil)
)

type CertificateAuthorityDataSourceModel struct {
	CertificateAuthorityModel
}

func NewCertificateAuthorityDataSource() datasource.DataSource { //nolint:ireturn
	return &CertificateAuthorityDataSource{}
}

type CertificateAuthorityDataSource struct {
	client *pfsense.Client
}

func (d *CertificateAuthorityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificate_authority", req.ProviderTypeName)
}

func (d *CertificateAuthorityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single certificate authority (CA) by descriptive name from the pfSense trust store.",
		MarkdownDescription: "Retrieves a single [certificate authority](https://docs.netgate.com/pfsense/en/latest/certificates/cas.html) (CA) by descriptive name from the pfSense trust store.",
		Attributes: map[string]schema.Attribute{
			"refid": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["refid"].Description,
				Computed:    true,
			},
			"descr": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["descr"].Description,
				Required:    true,
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
	}
}

func (d *CertificateAuthorityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CertificateAuthorityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificateAuthorityDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := d.client.GetCertificateAuthorityByDescr(ctx, data.Descr.ValueString())
	if addError(&resp.Diagnostics, "Unable to get certificate authority", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ca)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
