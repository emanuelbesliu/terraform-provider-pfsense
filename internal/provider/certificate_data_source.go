package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*CertificateDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*CertificateDataSource)(nil)
)

type CertificateDataSourceModel struct {
	CertificateModel
}

func NewCertificateDataSource() datasource.DataSource { //nolint:ireturn
	return &CertificateDataSource{}
}

type CertificateDataSource struct {
	client *pfsense.Client
}

func (d *CertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificate", req.ProviderTypeName)
}

func (d *CertificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single certificate by descriptive name from the pfSense certificate manager.",
		MarkdownDescription: "Retrieves a single [certificate](https://docs.netgate.com/pfsense/en/latest/certificates/index.html) by descriptive name from the pfSense certificate manager.",
		Attributes: map[string]schema.Attribute{
			"refid": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["refid"].Description,
				Computed:    true,
			},
			"descr": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["descr"].Description,
				Required:    true,
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
	}
}

func (d *CertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *CertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := d.client.GetCertificateByDescr(ctx, data.Descr.ValueString())
	if addError(&resp.Diagnostics, "Unable to get certificate", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *cert)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
