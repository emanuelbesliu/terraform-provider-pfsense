package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallURLAliasResource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallURLAliasResource)(nil)
	_ resource.ResourceWithImportState = (*FirewallURLAliasResource)(nil)
)

type FirewallURLAliasResourceModel struct {
	FirewallURLAliasModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewFirewallURLAliasResource() resource.Resource { //nolint:ireturn
	return &FirewallURLAliasResource{}
}

type FirewallURLAliasResource struct {
	client *pfsense.Client
}

func (r *FirewallURLAliasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_url_alias", req.ProviderTypeName)
}

func (r *FirewallURLAliasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall URL alias, defines a group of hosts, networks, or ports imported from URLs. URL aliases are periodically updated and can handle large lists. Aliases can be referenced by firewall rules, port forwards, outbound NAT rules, and other places in the firewall.",
		MarkdownDescription: "Firewall URL [alias](https://docs.netgate.com/pfsense/en/latest/firewall/aliases.html), defines a group of hosts, networks, or ports imported from URLs. URL aliases are periodically updated and can handle large lists. Aliases can be referenced by firewall rules, port forwards, outbound NAT rules, and other places in the firewall.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringIsAlias(),
				},
			},
			"description": schema.StringAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"type": schema.StringAttribute{
				Description:         FirewallURLAliasModel{}.descriptions()["type"].Description,
				MarkdownDescription: FirewallURLAliasModel{}.descriptions()["type"].MarkdownDescription,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.FirewallURLAlias{}.Types()...),
				},
			},
			"update_frequency": schema.Int64Attribute{
				Description: FirewallURLAliasModel{}.descriptions()["update_frequency"].Description,
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(7),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
			"entries": schema.ListNestedAttribute{
				Description: FirewallURLAliasModel{}.descriptions()["entries"].Description,
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: FirewallURLAliasEntryModel{}.AttrTypes()}, []attr.Value{})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: FirewallURLAliasEntryModel{}.descriptions()["url"].Description,
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallURLAliasEntryModel{}.descriptions()["description"].Description,
							Computed:    true,
							Optional:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 200),
							},
						},
					},
				},
			},
		},
	}
}

func (r *FirewallURLAliasResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *FirewallURLAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FirewallURLAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var urlAliasReq pfsense.FirewallURLAlias
	resp.Diagnostics.Append(data.Value(ctx, &urlAliasReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlAlias, err := r.client.CreateFirewallURLAlias(ctx, urlAliasReq)
	if addError(&resp.Diagnostics, "Error creating URL alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *urlAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ReloadFirewallFilter(ctx)
		addWarning(&resp.Diagnostics, "Error applying URL alias", err)
	}
}

func (r *FirewallURLAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FirewallURLAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlAlias, err := r.client.GetFirewallURLAlias(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading URL alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *urlAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallURLAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FirewallURLAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var urlAliasReq pfsense.FirewallURLAlias
	resp.Diagnostics.Append(data.Value(ctx, &urlAliasReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	urlAlias, err := r.client.UpdateFirewallURLAlias(ctx, urlAliasReq)
	if addError(&resp.Diagnostics, "Error updating URL alias", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *urlAlias)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ReloadFirewallFilter(ctx)
		addWarning(&resp.Diagnostics, "Error applying URL alias", err)
	}
}

func (r *FirewallURLAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FirewallURLAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFirewallURLAlias(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting URL alias", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ReloadFirewallFilter(ctx)
		addWarning(&resp.Diagnostics, "Error applying URL alias", err)
	}
}

func (r *FirewallURLAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
