package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*IPsecPhase1sDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*IPsecPhase1sDataSource)(nil)
)

func NewIPsecPhase1sDataSource() datasource.DataSource { //nolint:ireturn
	return &IPsecPhase1sDataSource{}
}

type IPsecPhase1sDataSource struct {
	client *pfsense.Client
}

func (d *IPsecPhase1sDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_ipsec_phase1s", req.ProviderTypeName)
}

func (d *IPsecPhase1sDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all IPsec Phase 1 tunnel configurations.",
		MarkdownDescription: "Retrieves all IPsec [Phase 1](https://docs.netgate.com/pfsense/en/latest/vpn/ipsec/configure.html) tunnel configurations.",
		Attributes: map[string]schema.Attribute{
			"all": schema.ListNestedAttribute{
				Description: "All IPsec Phase 1 entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ike_id": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["ike_id"].Description,
							Computed:    true,
						},
						"ike_type": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["ike_type"].Description,
							Computed:    true,
						},
						"interface": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["interface"].Description,
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["protocol"].Description,
							Computed:    true,
						},
						"remote_gateway": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["remote_gateway"].Description,
							Computed:    true,
						},
						"authentication_method": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["authentication_method"].Description,
							Computed:    true,
						},
						"pre_shared_key": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["pre_shared_key"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"my_id_type": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["my_id_type"].Description,
							Computed:    true,
						},
						"my_id_data": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["my_id_data"].Description,
							Computed:    true,
						},
						"peer_id_type": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["peer_id_type"].Description,
							Computed:    true,
						},
						"peer_id_data": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["peer_id_data"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"nat_traversal": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["nat_traversal"].Description,
							Computed:    true,
						},
						"mobike": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["mobike"].Description,
							Computed:    true,
						},
						"dpd_delay": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["dpd_delay"].Description,
							Computed:    true,
						},
						"dpd_max_fail": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["dpd_max_fail"].Description,
							Computed:    true,
						},
						"lifetime": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["lifetime"].Description,
							Computed:    true,
						},
						"rekey_time": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["rekey_time"].Description,
							Computed:    true,
						},
						"reauth_time": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["reauth_time"].Description,
							Computed:    true,
						},
						"rand_time": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["rand_time"].Description,
							Computed:    true,
						},
						"start_action": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["start_action"].Description,
							Computed:    true,
						},
						"close_action": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["close_action"].Description,
							Computed:    true,
						},
						"encryption": schema.ListNestedAttribute{
							Description: IPsecPhase1Model{}.descriptions()["encryption"].Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"algorithm": schema.StringAttribute{
										Description: IPsecPhase1EncryptionModel{}.descriptions()["algorithm"].Description,
										Computed:    true,
									},
									"key_length": schema.StringAttribute{
										Description: IPsecPhase1EncryptionModel{}.descriptions()["key_length"].Description,
										Computed:    true,
									},
									"hash_algorithm": schema.StringAttribute{
										Description: IPsecPhase1EncryptionModel{}.descriptions()["hash_algorithm"].Description,
										Computed:    true,
									},
									"prf_algorithm": schema.StringAttribute{
										Description: IPsecPhase1EncryptionModel{}.descriptions()["prf_algorithm"].Description,
										Computed:    true,
									},
									"dh_group": schema.StringAttribute{
										Description: IPsecPhase1EncryptionModel{}.descriptions()["dh_group"].Description,
										Computed:    true,
									},
								},
							},
						},
						"disabled": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["disabled"].Description,
							Computed:    true,
						},
						"cert_ref": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["cert_ref"].Description,
							Computed:    true,
						},
						"ca_ref": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["ca_ref"].Description,
							Computed:    true,
						},
						"pkcs11_cert_ref": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["pkcs11_cert_ref"].Description,
							Computed:    true,
						},
						"pkcs11_pin": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["pkcs11_pin"].Description,
							Computed:    true,
							Sensitive:   true,
						},
						"mobile": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["mobile"].Description,
							Computed:    true,
						},
						"ike_port": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["ike_port"].Description,
							Computed:    true,
						},
						"natt_port": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["natt_port"].Description,
							Computed:    true,
						},
						"gw_duplicates": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["gw_duplicates"].Description,
							Computed:    true,
						},
						"prf_select_enable": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["prf_select_enable"].Description,
							Computed:    true,
						},
						"split_conn": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["split_conn"].Description,
							Computed:    true,
						},
						"tfc_enable": schema.BoolAttribute{
							Description: IPsecPhase1Model{}.descriptions()["tfc_enable"].Description,
							Computed:    true,
						},
						"tfc_bytes": schema.StringAttribute{
							Description: IPsecPhase1Model{}.descriptions()["tfc_bytes"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *IPsecPhase1sDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *IPsecPhase1sDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPsecPhase1sModel

	phase1s, err := d.client.GetIPsecPhase1s(ctx)
	if addError(&resp.Diagnostics, "Unable to get IPsec phase 1 entries", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *phase1s)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
