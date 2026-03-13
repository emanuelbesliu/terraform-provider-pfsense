package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*InterfaceGroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*InterfaceGroupResource)(nil)
	_ resource.ResourceWithImportState = (*InterfaceGroupResource)(nil)
)

type InterfaceGroupResourceModel struct {
	InterfaceGroupModel
}

func NewInterfaceGroupResource() resource.Resource { //nolint:ireturn
	return &InterfaceGroupResource{}
}

type InterfaceGroupResource struct {
	client *pfsense.Client
}

func (r *InterfaceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_group", req.ProviderTypeName)
}

func (r *InterfaceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Interface group. Interface groups allow applying firewall rules to multiple interfaces at once.",
		MarkdownDescription: "[Interface group](https://docs.netgate.com/pfsense/en/latest/interfaces/groups.html). Interface groups allow applying firewall rules to multiple interfaces at once.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: InterfaceGroupModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 15),
				},
			},
			"members": schema.ListAttribute{
				Description: InterfaceGroupModel{}.descriptions()["members"].Description,
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.LengthAtLeast(1),
					),
				},
			},
			"description": schema.StringAttribute{
				Description: InterfaceGroupModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *InterfaceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *InterfaceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InterfaceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var groupReq pfsense.InterfaceGroup
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.CreateInterfaceGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error creating interface group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InterfaceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetInterfaceGroup(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading interface group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InterfaceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var groupReq pfsense.InterfaceGroup
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.UpdateInterfaceGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error updating interface group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InterfaceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInterfaceGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting interface group", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *InterfaceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
