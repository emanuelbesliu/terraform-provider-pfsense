package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*DynamicDNSResource)(nil)
	_ resource.ResourceWithConfigure   = (*DynamicDNSResource)(nil)
	_ resource.ResourceWithImportState = (*DynamicDNSResource)(nil)
)

type DynamicDNSResourceModel struct {
	DynamicDNSModel
}

func NewDynamicDNSResource() resource.Resource { //nolint:ireturn
	return &DynamicDNSResource{}
}

type DynamicDNSResource struct {
	client *pfsense.Client
}

func (r *DynamicDNSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dynamic_dns", req.ProviderTypeName)
}

func (r *DynamicDNSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	desc := DynamicDNSModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Manages a dynamic DNS client entry on pfSense.",
		MarkdownDescription: "Manages a [dynamic DNS](https://docs.netgate.com/pfsense/en/latest/services/dyndns/index.html) client entry on pfSense.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: desc["id"].Description,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: desc["type"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"interface": schema.StringAttribute{
				Description: desc["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"host": schema.StringAttribute{
				Description: desc["host"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"domain_name": schema.StringAttribute{
				Description: desc["domain_name"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"username": schema.StringAttribute{
				Description: desc["username"].Description,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				Description: desc["password"].Description,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
			},
			"mx": schema.StringAttribute{
				Description: desc["mx"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"wildcard": schema.BoolAttribute{
				Description: desc["wildcard"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"proxied": schema.BoolAttribute{
				Description: desc["proxied"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"verbose_log": schema.BoolAttribute{
				Description: desc["verbose_log"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"curl_ipresolve_v4": schema.BoolAttribute{
				Description: desc["curl_ipresolve_v4"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"curl_ssl_verifypeer": schema.BoolAttribute{
				Description: desc["curl_ssl_verifypeer"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"zone_id": schema.StringAttribute{
				Description: desc["zone_id"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ttl": schema.StringAttribute{
				Description: desc["ttl"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"max_cache_age": schema.StringAttribute{
				Description: desc["max_cache_age"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"update_url": schema.StringAttribute{
				Description: desc["update_url"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"result_match": schema.StringAttribute{
				Description: desc["result_match"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"request_interface": schema.StringAttribute{
				Description: desc["request_interface"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"curl_proxy": schema.StringAttribute{
				Description: desc["curl_proxy"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: desc["description"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"disabled": schema.BoolAttribute{
				Description: desc["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *DynamicDNSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *DynamicDNSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DynamicDNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entryReq pfsense.DynamicDNS
	resp.Diagnostics.Append(data.Value(ctx, &entryReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.CreateDynamicDNS(ctx, entryReq)
	if addError(&resp.Diagnostics, "Error creating dynamic DNS entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DynamicDNSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DynamicDNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.GetDynamicDNS(ctx, int(data.ID.ValueInt64()))

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading dynamic DNS entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DynamicDNSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DynamicDNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *DynamicDNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entryReq pfsense.DynamicDNS
	resp.Diagnostics.Append(data.Value(ctx, &entryReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.UpdateDynamicDNS(ctx, int(state.ID.ValueInt64()), entryReq)
	if addError(&resp.Diagnostics, "Error updating dynamic DNS entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DynamicDNSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DynamicDNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDynamicDNS(ctx, int(data.ID.ValueInt64()))
	if addError(&resp.Diagnostics, "Error deleting dynamic DNS entry", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DynamicDNSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", "Dynamic DNS import ID must be a numeric index.")

		return
	}

	entry, err := r.client.GetDynamicDNS(ctx, id)
	if addError(&resp.Diagnostics, "Error importing dynamic DNS entry", err) {
		return
	}

	var data DynamicDNSResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
