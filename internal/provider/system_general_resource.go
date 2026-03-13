package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemGeneralResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemGeneralResource)(nil)
	_ resource.ResourceWithImportState = (*SystemGeneralResource)(nil)
)

type SystemGeneralResourceModel struct {
	SystemGeneralModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemGeneralResource() resource.Resource { //nolint:ireturn
	return &SystemGeneralResource{}
}

type SystemGeneralResource struct {
	client *pfsense.Client
}

func (r *SystemGeneralResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_general", req.ProviderTypeName)
}

func (r *SystemGeneralResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "System general setup configuration including hostname, domain, DNS servers, localization, and webConfigurator settings. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System general setup](https://docs.netgate.com/pfsense/en/latest/config/general.html) configuration including hostname, domain, DNS servers, localization, and webConfigurator settings. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["hostname"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemHostname),
				Validators: []validator.String{
					stringIsDNSLabel(),
				},
			},
			"domain": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["domain"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemDomain),
				Validators: []validator.String{
					stringIsDomain(),
				},
			},
			"dns_servers": schema.ListNestedAttribute{
				Description: SystemGeneralModel{}.descriptions()["dns_servers"].Description,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(pfsense.MaxDNSServers),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["address"].Description,
							Required:    true,
							Validators: []validator.String{
								stringIsIPAddress("Any"),
							},
						},
						"hostname": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["hostname"].Description,
							Optional:    true,
						},
						"gateway": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["gateway"].Description,
							Optional:    true,
						},
					},
				},
			},
			"dns_override": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["dns_override"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(pfsense.DefaultSystemDNSAllowOverride),
			},
			"dns_localhost": schema.StringAttribute{
				Description:         SystemGeneralModel{}.descriptions()["dns_localhost"].Description,
				MarkdownDescription: SystemGeneralModel{}.descriptions()["dns_localhost"].MarkdownDescription,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.SystemGeneral{}.DNSLocalhostOptions()...),
				},
			},
			"timezone": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["timezone"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemTimezone),
			},
			"timeservers": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["timeservers"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemTimeservers),
			},
			"language": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["language"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemLanguage),
			},
			"webgui_theme": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_theme"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemWebGUICSS),
			},
			"login_color": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["login_color"].Description,
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(pfsense.DefaultSystemLoginCSS),
			},
			"login_show_host": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["login_show_host"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"webgui_fixed_menu": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_fixed_menu"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dashboard_columns": schema.Int64Attribute{
				Description: SystemGeneralModel{}.descriptions()["dashboard_columns"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultSystemDashboardColumns)),
				Validators: []validator.Int64{
					int64validator.Between(1, 6),
				},
			},
			"webgui_left_column_hyper": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_left_column_hyper"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_alias_popup_detail": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["disable_alias_popup_detail"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dashboard_available_widgets_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["dashboard_available_widgets_panel"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"system_logs_filter_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["system_logs_filter_panel"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"system_logs_manage_log_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["system_logs_manage_log_panel"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"status_monitoring_settings_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["status_monitoring_settings_panel"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"row_order_dragging": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["row_order_dragging"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"interfaces_sort": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["interfaces_sort"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_state_filter": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["require_state_filter"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"hostname_in_menu": schema.StringAttribute{
				Description:         SystemGeneralModel{}.descriptions()["hostname_in_menu"].Description,
				MarkdownDescription: SystemGeneralModel{}.descriptions()["hostname_in_menu"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultSystemHostnameInMenu),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.SystemGeneral{}.HostnameInMenuOptions()...),
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

func (r *SystemGeneralResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemGeneralResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemGeneralResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var sgReq pfsense.SystemGeneral
	resp.Diagnostics.Append(data.Value(ctx, &sgReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sg, err := r.client.UpdateSystemGeneral(ctx, sgReq)
	if addError(&resp.Diagnostics, "Error creating system general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *sg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemGeneralChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system general changes", err)
	}
}

func (r *SystemGeneralResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemGeneralResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sg, err := r.client.GetSystemGeneral(ctx)
	if addError(&resp.Diagnostics, "Error reading system general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *sg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemGeneralResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemGeneralResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var sgReq pfsense.SystemGeneral
	resp.Diagnostics.Append(data.Value(ctx, &sgReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sg, err := r.client.UpdateSystemGeneral(ctx, sgReq)
	if addError(&resp.Diagnostics, "Error updating system general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *sg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemGeneralChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system general changes", err)
	}
}

func (r *SystemGeneralResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemGeneralResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultSG := pfsense.SystemGeneral{
		Hostname:         pfsense.DefaultSystemHostname,
		Domain:           pfsense.DefaultSystemDomain,
		DNSOverride:      pfsense.DefaultSystemDNSAllowOverride,
		Timezone:         pfsense.DefaultSystemTimezone,
		Timeservers:      pfsense.DefaultSystemTimeservers,
		Language:         pfsense.DefaultSystemLanguage,
		WebGUICSS:        pfsense.DefaultSystemWebGUICSS,
		LoginCSS:         pfsense.DefaultSystemLoginCSS,
		DashboardColumns: pfsense.DefaultSystemDashboardColumns,
		HostnameInMenu:   pfsense.DefaultSystemHostnameInMenu,
	}

	_, err := r.client.UpdateSystemGeneral(ctx, defaultSG)
	if addError(&resp.Diagnostics, "Error resetting system general settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemGeneralChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system general changes", err)
	}
}

func (r *SystemGeneralResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	sg, err := r.client.GetSystemGeneral(ctx)
	if addError(&resp.Diagnostics, "Error importing system general settings", err) {
		return
	}

	var data SystemGeneralResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *sg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
