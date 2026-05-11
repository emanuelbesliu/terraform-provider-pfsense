package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*FirewallNAT1to1Resource)(nil)
	_ resource.ResourceWithConfigure   = (*FirewallNAT1to1Resource)(nil)
	_ resource.ResourceWithImportState = (*FirewallNAT1to1Resource)(nil)
)

type FirewallNAT1to1ResourceModel struct {
	FirewallNAT1to1Model
}

func NewFirewallNAT1to1Resource() resource.Resource { //nolint:ireturn
	return &FirewallNAT1to1Resource{}
}

type FirewallNAT1to1Resource struct {
	client *pfsense.Client
}

func (r *FirewallNAT1to1Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_nat_1to1", req.ProviderTypeName)
}

func (r *FirewallNAT1to1Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Firewall NAT 1:1 (BINAT). 1:1 NAT rules provide a one-to-one mapping between external and internal IP addresses.",
		MarkdownDescription: "[Firewall NAT 1:1](https://docs.netgate.com/pfsense/en/latest/nat/1-to-1-nat.html). 1:1 NAT rules provide a one-to-one mapping between external and internal IP addresses.",
		Attributes: map[string]schema.Attribute{
			"external": schema.StringAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["external"].Description,
				Required:    true,
			},
			"interface": schema.StringAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description:         FirewallNAT1to1Model{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: FirewallNAT1to1Model{}.descriptions()["ipprotocol"].MarkdownDescription,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOneToOne{}.IPProtocols()...),
				},
			},
			"source_address": schema.StringAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["source_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"source_not": schema.BoolAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["source_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"destination_address": schema.StringAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["destination_address"].Description,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("any"),
			},
			"destination_not": schema.BoolAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["destination_not"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["description"].Description,
				Required:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"no_binat": schema.BoolAttribute{
				Description: FirewallNAT1to1Model{}.descriptions()["no_binat"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"nat_reflection": schema.StringAttribute{
				Description:         FirewallNAT1to1Model{}.descriptions()["nat_reflection"].Description,
				MarkdownDescription: FirewallNAT1to1Model{}.descriptions()["nat_reflection"].MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.NATOneToOne{}.NATReflectionModes()...),
				},
			},
		},
	}
}

func (r *FirewallNAT1to1Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*pfsense.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pfsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *FirewallNAT1to1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FirewallNAT1to1ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule := data.toClient(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateNATOneToOne(ctx, rule); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create 1:1 NAT Rule",
			err.Error(),
		)

		return
	}

	data.fromClient(ctx, rule, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNAT1to1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FirewallNAT1to1ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetNATOneToOne(ctx, data.Description.ValueString())
	if err != nil {
		if errors.Is(err, pfsense.ErrNotFound) {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"Unable to Read 1:1 NAT Rule",
			err.Error(),
		)

		return
	}

	data.fromClient(ctx, rule, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNAT1to1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FirewallNAT1to1ResourceModel
	var state FirewallNAT1to1ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rule := data.toClient(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateNATOneToOne(ctx, state.Description.ValueString(), rule); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update 1:1 NAT Rule",
			err.Error(),
		)

		return
	}

	data.fromClient(ctx, rule, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallNAT1to1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FirewallNAT1to1ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNATOneToOne(ctx, data.Description.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete 1:1 NAT Rule",
			err.Error(),
		)

		return
	}
}

func (r *FirewallNAT1to1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("description"), req, resp)
}
