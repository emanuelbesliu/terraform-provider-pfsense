package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallNAT1to1DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallNAT1to1DataSource)(nil)
)

func NewFirewallNAT1to1DataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallNAT1to1DataSource{}
}

type FirewallNAT1to1DataSource struct {
	client *pfsense.Client
}

type FirewallNAT1to1DataSourceModel struct {
	Description   types.String `tfsdk:"description"`
	External      types.String `tfsdk:"external"`
	Interface     types.String `tfsdk:"interface"`
	IPProtocol    types.String `tfsdk:"ipprotocol"`
	SourceAddress types.String `tfsdk:"source_address"`
	SourceNot     types.Bool   `tfsdk:"source_not"`
	DestAddress   types.String `tfsdk:"destination_address"`
	DestNot       types.Bool   `tfsdk:"destination_not"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	NoBinat       types.Bool   `tfsdk:"no_binat"`
	NATReflection types.String `tfsdk:"nat_reflection"`
}

func (d *FirewallNAT1to1DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_1to1", req.ProviderTypeName)
}

func (d *FirewallNAT1to1DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT 1:1 (BINAT) data source. Retrieve details of an existing 1:1 NAT rule by description.",
		MarkdownDescription: "[Firewall NAT 1:1](https://docs.netgate.com/pfsense/en/latest/nat/1-to-1-nat.html) data source. Retrieve details of an existing 1:1 NAT rule by description.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "Description of the 1:1 NAT rule to retrieve.",
				Required:    true,
			},
			"external": schema.StringAttribute{
				Description: "External IP address for the 1:1 NAT rule.",
				Computed:    true,
			},
			"interface": schema.StringAttribute{
				Description: "Network interface this NAT rule applies to.",
				Computed:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description: "IP address family (inet, inet6, or inet46).",
				Computed:    true,
			},
			"source_address": schema.StringAttribute{
				Description: "Source address for the rule.",
				Computed:    true,
			},
			"source_not": schema.BoolAttribute{
				Description: "Whether the source address match is inverted.",
				Computed:    true,
			},
			"destination_address": schema.StringAttribute{
				Description: "Destination address for the rule.",
				Computed:    true,
			},
			"destination_not": schema.BoolAttribute{
				Description: "Whether the destination address match is inverted.",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Whether this 1:1 NAT rule is disabled.",
				Computed:    true,
			},
			"no_binat": schema.BoolAttribute{
				Description: "Whether 1:1 NAT (BINAT) is disabled for this rule.",
				Computed:    true,
			},
			"nat_reflection": schema.StringAttribute{
				Description: "NAT reflection mode for this rule.",
				Computed:    true,
			},
		},
	}
}

func (d *FirewallNAT1to1DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pfsense.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *pfsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *FirewallNAT1to1DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallNAT1to1DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := d.client.GetNATOneToOne(ctx, data.Description.ValueString())
	if err != nil {
		if errors.Is(err, pfsense.ErrNotFound) {
			resp.Diagnostics.AddError(
				"1:1 NAT Rule Not Found",
				fmt.Sprintf("Unable to find 1:1 NAT rule with description '%s'", data.Description.ValueString()),
			)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read 1:1 NAT Rule",
			err.Error(),
		)

		return
	}

	data.External = types.StringValue(rule.External)
	data.Interface = types.StringValue(rule.Interface)
	data.IPProtocol = types.StringValue(rule.IPProtocol)
	data.SourceAddress = types.StringValue(rule.SourceAddress)
	data.SourceNot = types.BoolValue(rule.SourceNot)
	data.DestAddress = types.StringValue(rule.DestAddress)
	data.DestNot = types.BoolValue(rule.DestNot)
	data.Disabled = types.BoolValue(rule.Disabled)
	data.NoBinat = types.BoolValue(rule.NoBinat)
	data.NATReflection = types.StringValue(rule.NATReflection)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
