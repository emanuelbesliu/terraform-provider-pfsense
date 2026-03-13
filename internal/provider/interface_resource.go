package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                   = (*InterfaceResource)(nil)
	_ resource.ResourceWithConfigure      = (*InterfaceResource)(nil)
	_ resource.ResourceWithImportState    = (*InterfaceResource)(nil)
	_ resource.ResourceWithValidateConfig = (*InterfaceResource)(nil)
)

type InterfaceResourceModel struct {
	InterfaceModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewInterfaceResource() resource.Resource { //nolint:ireturn
	return &InterfaceResource{}
}

type InterfaceResource struct {
	client *pfsense.Client
}

func (r *InterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface", req.ProviderTypeName)
}

func (r *InterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Interface assignment. Assigns a physical or VLAN interface and configures its network settings.",
		MarkdownDescription: "[Interface assignment](https://docs.netgate.com/pfsense/en/latest/interfaces/index.html). Assigns a physical or VLAN interface and configures its network settings.",
		Attributes: map[string]schema.Attribute{
			"logical_name": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["logical_name"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["port"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["enabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ipv4_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv4_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv4_type"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.InterfaceIPv4Types...),
				},
			},
			"ipv4_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_address"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("IPv4"),
				},
			},
			"ipv4_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_subnet"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(subnetV4Regex, "must be a number between 1 and 32"),
				},
			},
			"ipv4_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_gateway"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv6_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv6_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv6_type"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.InterfaceIPv6Types...),
				},
			},
			"ipv6_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_address"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsIPAddress("IPv6"),
				},
			},
			"ipv6_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_subnet"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(subnetV6Regex, "must be a number between 1 and 128"),
				},
			},
			"ipv6_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_gateway"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"spoof_mac": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["spoof_mac"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringIsMACAddress(),
				},
			},
			"mtu": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mtu"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(576, 9000),
				},
			},
			"mss": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mss"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(536, 65535),
				},
			},
			"block_private": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_private"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"block_bogons": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_bogons"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *InterfaceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data InterfaceResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ipv4Type := "none"
	if !data.IPv4Type.IsNull() && !data.IPv4Type.IsUnknown() {
		ipv4Type = data.IPv4Type.ValueString()
	}

	ipv6Type := "none"
	if !data.IPv6Type.IsNull() && !data.IPv6Type.IsUnknown() {
		ipv6Type = data.IPv6Type.ValueString()
	}

	// When ipv4_type is staticv4, ipv4_address and ipv4_subnet are required.
	if ipv4Type == "staticv4" {
		if data.IPAddr.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv4_address"),
				"Missing required attribute",
				"ipv4_address is required when ipv4_type is 'staticv4'.",
			)
		}

		if data.Subnet.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv4_subnet"),
				"Missing required attribute",
				"ipv4_subnet is required when ipv4_type is 'staticv4'.",
			)
		}
	}

	// When ipv4_type is NOT staticv4, reject address/subnet/gateway fields.
	if ipv4Type != "staticv4" {
		if !data.IPAddr.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv4_address"),
				"Attribute not applicable",
				fmt.Sprintf("ipv4_address can only be set when ipv4_type is 'staticv4', current ipv4_type is '%s'.", ipv4Type),
			)
		}

		if !data.Subnet.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv4_subnet"),
				"Attribute not applicable",
				fmt.Sprintf("ipv4_subnet can only be set when ipv4_type is 'staticv4', current ipv4_type is '%s'.", ipv4Type),
			)
		}

		if !data.Gateway.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv4_gateway"),
				"Attribute not applicable",
				fmt.Sprintf("ipv4_gateway can only be set when ipv4_type is 'staticv4', current ipv4_type is '%s'.", ipv4Type),
			)
		}
	}

	// When ipv6_type is staticv6, ipv6_address and ipv6_subnet are required.
	if ipv6Type == "staticv6" {
		if data.IPAddrV6.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv6_address"),
				"Missing required attribute",
				"ipv6_address is required when ipv6_type is 'staticv6'.",
			)
		}

		if data.SubnetV6.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv6_subnet"),
				"Missing required attribute",
				"ipv6_subnet is required when ipv6_type is 'staticv6'.",
			)
		}
	}

	// When ipv6_type is NOT staticv6, reject address/subnet/gateway fields.
	if ipv6Type != "staticv6" {
		if !data.IPAddrV6.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv6_address"),
				"Attribute not applicable",
				fmt.Sprintf("ipv6_address can only be set when ipv6_type is 'staticv6', current ipv6_type is '%s'.", ipv6Type),
			)
		}

		if !data.SubnetV6.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv6_subnet"),
				"Attribute not applicable",
				fmt.Sprintf("ipv6_subnet can only be set when ipv6_type is 'staticv6', current ipv6_type is '%s'.", ipv6Type),
			)
		}

		if !data.GatewayV6.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ipv6_gateway"),
				"Attribute not applicable",
				fmt.Sprintf("ipv6_gateway can only be set when ipv6_type is 'staticv6', current ipv6_type is '%s'.", ipv6Type),
			)
		}
	}
}

func (r *InterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *InterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ifaceReq pfsense.Interface
	resp.Diagnostics.Append(data.Value(ctx, &ifaceReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.CreateInterface(ctx, ifaceReq)
	if addError(&resp.Diagnostics, "Error creating interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, iface.LogicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.GetInterface(ctx, data.LogicalName.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to preserve logical_name.
	var state *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure logical_name from state is used.
	data.LogicalName = state.LogicalName

	var ifaceReq pfsense.Interface
	resp.Diagnostics.Append(data.Value(ctx, &ifaceReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.UpdateInterface(ctx, ifaceReq)
	if addError(&resp.Diagnostics, "Error updating interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, iface.LogicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalName := data.LogicalName.ValueString()

	err := r.client.DeleteInterface(ctx, logicalName)
	if addError(&resp.Diagnostics, "Error deleting interface", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, logicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("logical_name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
