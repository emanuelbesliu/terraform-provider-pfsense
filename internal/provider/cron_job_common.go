package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type CronJobModel struct {
	Minute  types.String `tfsdk:"minute"`
	Hour    types.String `tfsdk:"hour"`
	MDay    types.String `tfsdk:"mday"`
	Month   types.String `tfsdk:"month"`
	WDay    types.String `tfsdk:"wday"`
	Who     types.String `tfsdk:"who"`
	Command types.String `tfsdk:"command"`
}

func (CronJobModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"minute": {
			Description: "Minute schedule field (0-59, or '*' for every minute).",
		},
		"hour": {
			Description: "Hour schedule field (0-23, or '*' for every hour).",
		},
		"mday": {
			Description: "Day of month schedule field (1-31, or '*' for every day).",
		},
		"month": {
			Description: "Month schedule field (1-12, or '*' for every month).",
		},
		"wday": {
			Description: "Day of week schedule field (0-7 where 0 and 7 are Sunday, or '*' for every day).",
		},
		"who": {
			Description: "User to run the command as (e.g. 'root').",
		},
		"command": {
			Description: "Command to execute.",
		},
	}
}

func (CronJobModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"minute":  types.StringType,
		"hour":    types.StringType,
		"mday":    types.StringType,
		"month":   types.StringType,
		"wday":    types.StringType,
		"who":     types.StringType,
		"command": types.StringType,
	}
}

func (m *CronJobModel) Set(_ context.Context, job pfsense.CronJob) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Minute = types.StringValue(job.Minute)
	m.Hour = types.StringValue(job.Hour)
	m.MDay = types.StringValue(job.MDay)
	m.Month = types.StringValue(job.Month)
	m.WDay = types.StringValue(job.WDay)
	m.Who = types.StringValue(job.Who)
	m.Command = types.StringValue(job.Command)

	return diags
}

func (m CronJobModel) Value(_ context.Context, job *pfsense.CronJob) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("minute"),
		"Minute cannot be parsed",
		job.SetMinute(m.Minute.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("hour"),
		"Hour cannot be parsed",
		job.SetHour(m.Hour.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("mday"),
		"Day of month cannot be parsed",
		job.SetMDay(m.MDay.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("month"),
		"Month cannot be parsed",
		job.SetMonth(m.Month.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("wday"),
		"Day of week cannot be parsed",
		job.SetWDay(m.WDay.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("who"),
		"User cannot be parsed",
		job.SetWho(m.Who.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("command"),
		"Command cannot be parsed",
		job.SetCommand(m.Command.ValueString()),
	)

	return diags
}
