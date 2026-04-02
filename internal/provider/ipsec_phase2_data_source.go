package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*IPsecPhase2DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*IPsecPhase2DataSource)(nil)
)

type IPsecPhase2DataSourceModel struct {
	IPsecPhase2Model
}

func NewIPsecPhase2DataSource() datasource.DataSource { //nolint:ireturn
	return &IPsecPhase2DataSource{}
}

type IPsecPhase2DataSource struct {
	client *pfsense.Client
}

func (d *IPsecPhase2DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_ipsec_phase2", req.ProviderTypeName)
}

func (d *IPsecPhase2DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single IPsec Phase 2 tunnel configuration by unique ID.",
		MarkdownDescription: "Retrieves a single IPsec [Phase 2](https://docs.netgate.com/pfsense/en/latest/vpn/ipsec/configure.html) tunnel configuration by unique ID.",
		Attributes: map[string]schema.Attribute{
			"uniq_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["uniq_id"].Description,
				Required:    true,
			},
			"ike_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["ike_id"].Description,
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["mode"].Description,
				Computed:    true,
			},
			"req_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["req_id"].Description,
				Computed:    true,
			},
			"local_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["local_id"].Description,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["type"].Description,
						Computed:    true,
					},
					"address": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["address"].Description,
						Computed:    true,
					},
					"net_bits": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["net_bits"].Description,
						Computed:    true,
					},
				},
			},
			"remote_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["remote_id"].Description,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["type"].Description,
						Computed:    true,
					},
					"address": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["address"].Description,
						Computed:    true,
					},
					"net_bits": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["net_bits"].Description,
						Computed:    true,
					},
				},
			},
			"nat_local_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["nat_local_id"].Description,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["type"].Description,
						Computed:    true,
					},
					"address": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["address"].Description,
						Computed:    true,
					},
					"net_bits": schema.StringAttribute{
						Description: IPsecPhase2IDModel{}.descriptions()["net_bits"].Description,
						Computed:    true,
					},
				},
			},
			"protocol": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["protocol"].Description,
				Computed:    true,
			},
			"encryption_algorithm_option": schema.ListNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["encryption_algorithm_option"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: IPsecPhase2EncryptionAlgorithmModel{}.descriptions()["name"].Description,
							Computed:    true,
						},
						"key_length": schema.StringAttribute{
							Description: IPsecPhase2EncryptionAlgorithmModel{}.descriptions()["key_length"].Description,
							Computed:    true,
						},
					},
				},
			},
			"hash_algorithm_option": schema.ListAttribute{
				Description: IPsecPhase2Model{}.descriptions()["hash_algorithm_option"].Description,
				Computed:    true,
				ElementType: types.StringType,
			},
			"pfs_group": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["pfs_group"].Description,
				Computed:    true,
			},
			"lifetime": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["lifetime"].Description,
				Computed:    true,
			},
			"rekey_time": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["rekey_time"].Description,
				Computed:    true,
			},
			"rand_time": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["rand_time"].Description,
				Computed:    true,
			},
			"ping_host": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["ping_host"].Description,
				Computed:    true,
			},
			"keepalive": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["keepalive"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: IPsecPhase2Model{}.descriptions()["disabled"].Description,
				Computed:    true,
			},
			"mobile": schema.BoolAttribute{
				Description: IPsecPhase2Model{}.descriptions()["mobile"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *IPsecPhase2DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *IPsecPhase2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPsecPhase2DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p2, err := d.client.GetIPsecPhase2(ctx, data.UniqID.ValueString())
	if addError(&resp.Diagnostics, "Unable to get IPsec phase 2", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p2)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
