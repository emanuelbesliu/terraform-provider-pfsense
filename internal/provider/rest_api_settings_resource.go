package provider

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource              = (*RESTAPISettingsResource)(nil)
	_ resource.ResourceWithConfigure = (*RESTAPISettingsResource)(nil)
)

type RESTAPISettingsResourceModel struct {
	AuthMethods types.List `tfsdk:"auth_methods"`
}

func NewRESTAPISettingsResource() resource.Resource { //nolint:ireturn
	return &RESTAPISettingsResource{}
}

type RESTAPISettingsResource struct {
	client *pfsense.Client
}

func (r *RESTAPISettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rest_api_settings"
}

func (r *RESTAPISettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages pfSense REST API package settings (authentication methods). Requires the REST API package to be installed.",
		MarkdownDescription: "Manages pfSense REST API package settings (authentication methods). Requires the REST API package to be installed.",
		Attributes: map[string]schema.Attribute{
			"auth_methods": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of enabled authentication methods. Valid values: BasicAuth, KeyAuth, JWTAuth.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("BasicAuth", "KeyAuth", "JWTAuth"),
					),
				},
			},
		},
	}
}

func (r *RESTAPISettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *RESTAPISettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RESTAPISettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	methods := expandStringList(ctx, data.AuthMethods, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.UpdateRESTAPISettings(ctx, pfsense.RESTAPISettingsUpdateRequest{
		AuthMethods: methods,
	})
	if addError(&resp.Diagnostics, "Error updating REST API settings", err) {
		return
	}

	data.AuthMethods = flattenStringList(ctx, settings.AuthMethods, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RESTAPISettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RESTAPISettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetRESTAPISettings(ctx)
	if addError(&resp.Diagnostics, "Error reading REST API settings", err) {
		return
	}

	data.AuthMethods = flattenStringList(ctx, settings.AuthMethods, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RESTAPISettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RESTAPISettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	methods := expandStringList(ctx, data.AuthMethods, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.UpdateRESTAPISettings(ctx, pfsense.RESTAPISettingsUpdateRequest{
		AuthMethods: methods,
	})
	if addError(&resp.Diagnostics, "Error updating REST API settings", err) {
		return
	}

	data.AuthMethods = flattenStringList(ctx, settings.AuthMethods, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RESTAPISettingsResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	_, err := r.client.UpdateRESTAPISettings(ctx, pfsense.RESTAPISettingsUpdateRequest{
		AuthMethods: []string{"BasicAuth"},
	})
	if addError(&resp.Diagnostics, "Error resetting REST API settings", err) {
		return
	}
}

func expandStringList(ctx context.Context, list types.List, diags *diag.Diagnostics) []string {
	var result []string

	d := list.ElementsAs(ctx, &result, false)
	diags.Append(d...)

	sort.Strings(result)

	return result
}

func flattenStringList(ctx context.Context, list []string, diags *diag.Diagnostics) types.List {
	sorted := make([]string, len(list))
	copy(sorted, list)
	sort.Strings(sorted)

	result, d := types.ListValueFrom(ctx, types.StringType, sorted)
	diags.Append(d...)

	return result
}
