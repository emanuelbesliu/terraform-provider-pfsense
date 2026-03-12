package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type GatewayModel struct {
	Name           types.String `tfsdk:"name"`
	Interface      types.String `tfsdk:"interface"`
	IPProtocol     types.String `tfsdk:"ipprotocol"`
	GatewayIP      types.String `tfsdk:"gateway"`
	Description    types.String `tfsdk:"description"`
	Disabled       types.Bool   `tfsdk:"disabled"`
	DefaultGW      types.Bool   `tfsdk:"default_gw"`
	Monitor        types.String `tfsdk:"monitor"`
	MonitorDisable types.Bool   `tfsdk:"monitor_disable"`
	ActionDisable  types.Bool   `tfsdk:"action_disable"`
	ForceDown      types.Bool   `tfsdk:"force_down"`
	Weight         types.Int64  `tfsdk:"weight"`
	NonLocalGW     types.Bool   `tfsdk:"non_local_gateway"`
	LatencyLow     types.Int64  `tfsdk:"latency_low"`
	LatencyHigh    types.Int64  `tfsdk:"latency_high"`
	LossLow        types.Int64  `tfsdk:"loss_low"`
	LossHigh       types.Int64  `tfsdk:"loss_high"`
	Interval       types.Int64  `tfsdk:"interval"`
	LossInterval   types.Int64  `tfsdk:"loss_interval"`
	TimePeriod     types.Int64  `tfsdk:"time_period"`
	AlertInterval  types.Int64  `tfsdk:"alert_interval"`
	DataPayload    types.Int64  `tfsdk:"data_payload"`
}

func (GatewayModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Name of the gateway. Must be unique.",
		},
		"interface": {
			Description: "Network interface for the gateway.",
		},
		"ipprotocol": {
			Description:         fmt.Sprintf("IP protocol. Options: %s.", wrapElementsJoin(pfsense.Gateway{}.IPProtocols(), "'")),
			MarkdownDescription: fmt.Sprintf("IP protocol. Options: %s.", wrapElementsJoin(pfsense.Gateway{}.IPProtocols(), "`")),
		},
		"gateway": {
			Description: "Gateway IP address, or 'dynamic' for DHCP/PPPoE interfaces.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"disabled": {
			Description: "Disable this gateway.",
		},
		"default_gw": {
			Description: "Mark this gateway as the default gateway.",
		},
		"monitor": {
			Description: "Alternative address to monitor link. Use this to monitor an IP address other than the gateway for link status.",
		},
		"monitor_disable": {
			Description: "Disable gateway monitoring. Consider this gateway as always being up.",
		},
		"action_disable": {
			Description: "Disable any action taken on gateway events. No action will be taken on gateway events but the gateway's status will be tracked.",
		},
		"force_down": {
			Description: "Mark gateway as down. Forces this gateway to be considered down.",
		},
		"weight": {
			Description: fmt.Sprintf("Weight for this gateway when used in a gateway group (range %d-%d). Defaults to '%d'.", pfsense.MinGatewayWeight, pfsense.MaxGatewayWeight, pfsense.DefaultGatewayWeight),
		},
		"non_local_gateway": {
			Description: "Allow the use of a gateway outside of the interface's subnet.",
		},
		"latency_low": {
			Description: fmt.Sprintf("Low latency threshold in milliseconds. Defaults to '%d'.", pfsense.DefaultGatewayLatencyLow),
		},
		"latency_high": {
			Description: fmt.Sprintf("High latency threshold in milliseconds. Defaults to '%d'.", pfsense.DefaultGatewayLatencyHigh),
		},
		"loss_low": {
			Description: fmt.Sprintf("Low packet loss threshold in percent. Defaults to '%d'.", pfsense.DefaultGatewayLossLow),
		},
		"loss_high": {
			Description: fmt.Sprintf("High packet loss threshold in percent. Defaults to '%d'.", pfsense.DefaultGatewayLossHigh),
		},
		"interval": {
			Description: fmt.Sprintf("Probe interval in milliseconds. Defaults to '%d'.", pfsense.DefaultGatewayInterval),
		},
		"loss_interval": {
			Description: fmt.Sprintf("Time in milliseconds before packets are treated as lost. Defaults to '%d'.", pfsense.DefaultGatewayLossInterval),
		},
		"time_period": {
			Description: fmt.Sprintf("Averaging time period in milliseconds for gateway quality calculation. Defaults to '%d'.", pfsense.DefaultGatewayTimePeriod),
		},
		"alert_interval": {
			Description: fmt.Sprintf("Alert interval in milliseconds. Defaults to '%d'.", pfsense.DefaultGatewayAlertInterval),
		},
		"data_payload": {
			Description: fmt.Sprintf("ICMP data payload size in bytes. Defaults to '%d'.", pfsense.DefaultGatewayDataPayload),
		},
	}
}

func (GatewayModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":              types.StringType,
		"interface":         types.StringType,
		"ipprotocol":        types.StringType,
		"gateway":           types.StringType,
		"description":       types.StringType,
		"disabled":          types.BoolType,
		"default_gw":        types.BoolType,
		"monitor":           types.StringType,
		"monitor_disable":   types.BoolType,
		"action_disable":    types.BoolType,
		"force_down":        types.BoolType,
		"weight":            types.Int64Type,
		"non_local_gateway": types.BoolType,
		"latency_low":       types.Int64Type,
		"latency_high":      types.Int64Type,
		"loss_low":          types.Int64Type,
		"loss_high":         types.Int64Type,
		"interval":          types.Int64Type,
		"loss_interval":     types.Int64Type,
		"time_period":       types.Int64Type,
		"alert_interval":    types.Int64Type,
		"data_payload":      types.Int64Type,
	}
}

func (m *GatewayModel) Set(_ context.Context, gw pfsense.Gateway) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(gw.Name)
	m.Interface = types.StringValue(gw.Interface)
	m.IPProtocol = types.StringValue(gw.IPProtocol)
	m.GatewayIP = types.StringValue(gw.GatewayIP)

	if gw.Description != "" {
		m.Description = types.StringValue(gw.Description)
	}

	m.Disabled = types.BoolValue(gw.Disabled)
	m.DefaultGW = types.BoolValue(gw.DefaultGW)

	if gw.Monitor != "" {
		m.Monitor = types.StringValue(gw.Monitor)
	}

	m.MonitorDisable = types.BoolValue(gw.MonitorDisable)
	m.ActionDisable = types.BoolValue(gw.ActionDisable)
	m.ForceDown = types.BoolValue(gw.ForceDown)
	m.Weight = types.Int64Value(int64(gw.Weight))
	m.NonLocalGW = types.BoolValue(gw.NonLocalGW)
	m.LatencyLow = types.Int64Value(int64(gw.LatencyLow))
	m.LatencyHigh = types.Int64Value(int64(gw.LatencyHigh))
	m.LossLow = types.Int64Value(int64(gw.LossLow))
	m.LossHigh = types.Int64Value(int64(gw.LossHigh))
	m.Interval = types.Int64Value(int64(gw.Interval))
	m.LossInterval = types.Int64Value(int64(gw.LossInterval))
	m.TimePeriod = types.Int64Value(int64(gw.TimePeriod))
	m.AlertInterval = types.Int64Value(int64(gw.AlertInterval))
	m.DataPayload = types.Int64Value(int64(gw.DataPayload))

	return diags
}

func (m GatewayModel) Value(_ context.Context, gw *pfsense.Gateway) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Name cannot be parsed",
		gw.SetName(m.Name.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		gw.SetInterface(m.Interface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("ipprotocol"),
		"IP protocol cannot be parsed",
		gw.SetIPProtocol(m.IPProtocol.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("gateway"),
		"Gateway IP cannot be parsed",
		gw.SetGatewayIP(m.GatewayIP.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			gw.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("disabled"),
		"Disabled cannot be parsed",
		gw.SetDisabled(m.Disabled.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("default_gw"),
		"Default gateway cannot be parsed",
		gw.SetDefaultGW(m.DefaultGW.ValueBool()),
	)

	if !m.Monitor.IsNull() {
		addPathError(
			&diags,
			path.Root("monitor"),
			"Monitor cannot be parsed",
			gw.SetMonitor(m.Monitor.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("monitor_disable"),
		"Monitor disable cannot be parsed",
		gw.SetMonitorDisable(m.MonitorDisable.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("action_disable"),
		"Action disable cannot be parsed",
		gw.SetActionDisable(m.ActionDisable.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("force_down"),
		"Force down cannot be parsed",
		gw.SetForceDown(m.ForceDown.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("weight"),
		"Weight cannot be parsed",
		gw.SetWeight(int(m.Weight.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("non_local_gateway"),
		"Non-local gateway cannot be parsed",
		gw.SetNonLocalGW(m.NonLocalGW.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("latency_low"),
		"Latency low cannot be parsed",
		gw.SetLatencyLow(int(m.LatencyLow.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("latency_high"),
		"Latency high cannot be parsed",
		gw.SetLatencyHigh(int(m.LatencyHigh.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("loss_low"),
		"Loss low cannot be parsed",
		gw.SetLossLow(int(m.LossLow.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("loss_high"),
		"Loss high cannot be parsed",
		gw.SetLossHigh(int(m.LossHigh.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("interval"),
		"Interval cannot be parsed",
		gw.SetInterval(int(m.Interval.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("loss_interval"),
		"Loss interval cannot be parsed",
		gw.SetLossInterval(int(m.LossInterval.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("time_period"),
		"Time period cannot be parsed",
		gw.SetTimePeriod(int(m.TimePeriod.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("alert_interval"),
		"Alert interval cannot be parsed",
		gw.SetAlertInterval(int(m.AlertInterval.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("data_payload"),
		"Data payload cannot be parsed",
		gw.SetDataPayload(int(m.DataPayload.ValueInt64())),
	)

	return diags
}
