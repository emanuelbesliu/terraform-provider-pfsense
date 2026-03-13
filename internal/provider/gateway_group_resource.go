package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*GatewayGroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*GatewayGroupResource)(nil)
	_ resource.ResourceWithImportState = (*GatewayGroupResource)(nil)
)

type GatewayGroupResourceModel struct {
	GatewayGroupModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewGatewayGroupResource() resource.Resource { //nolint:ireturn
	return &GatewayGroupResource{}
}

type GatewayGroupResource struct {
	client *pfsense.Client
}

func (r *GatewayGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_gateway_group", req.ProviderTypeName)
}

func (r *GatewayGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Gateway group for failover and load balancing. Gateway groups organize multiple gateways into tiers for redundancy and traffic distribution.",
		MarkdownDescription: "[Gateway group](https://docs.netgate.com/pfsense/en/latest/routing/gateway-groups.html) for failover and load balancing. Gateway groups organize multiple gateways into tiers for redundancy and traffic distribution.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: GatewayGroupModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"description": schema.StringAttribute{
				Description: GatewayGroupModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"trigger": schema.StringAttribute{
				Description:         GatewayGroupModel{}.descriptions()["trigger"].Description,
				MarkdownDescription: GatewayGroupModel{}.descriptions()["trigger"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.GatewayGroup{}.Triggers()...),
				},
			},
			"keep_failover_states": schema.StringAttribute{
				Description: GatewayGroupModel{}.descriptions()["keep_failover_states"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.GatewayGroup{}.KeepFailoverStatesOptions()...),
				},
			},
			"members": schema.ListNestedAttribute{
				Description: GatewayGroupModel{}.descriptions()["members"].Description,
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: GatewayGroupMemberModel{}.AttrTypes()}, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gateway": schema.StringAttribute{
							Description: GatewayGroupMemberModel{}.descriptions()["gateway"].Description,
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"tier": schema.Int64Attribute{
							Description: GatewayGroupMemberModel{}.descriptions()["tier"].Description,
							Required:    true,
							Validators: []validator.Int64{
								int64validator.Between(int64(pfsense.MinGatewayGroupTier), int64(pfsense.MaxGatewayGroupTier)),
							},
						},
						"virtual_ip": schema.StringAttribute{
							Description: GatewayGroupMemberModel{}.descriptions()["virtual_ip"].Description,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
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

func (r *GatewayGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *GatewayGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GatewayGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var groupReq pfsense.GatewayGroup
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.CreateGatewayGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error creating gateway group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayGroupChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway group changes", err)
	}
}

func (r *GatewayGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GatewayGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetGatewayGroup(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading gateway group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GatewayGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GatewayGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var groupReq pfsense.GatewayGroup
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.UpdateGatewayGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error updating gateway group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayGroupChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway group changes", err)
	}
}

func (r *GatewayGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GatewayGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGatewayGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting gateway group", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGatewayGroupChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying gateway group changes", err)
	}
}

func (r *GatewayGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
