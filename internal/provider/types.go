package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MonitorResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PathName    types.String `tfsdk:"path_name"`
	Url         types.String `tfsdk:"url"`
	Method      types.String `tfsdk:"method"`
	// Active      types.Bool   `tfsdk:"active"`
	Type types.String `tfsdk:"type"`

	// Timeout        types.Int64 `tfsdk:"timeout"`
	Interval       types.Int64 `tfsdk:"interval"`
	RetryInterval  types.Int64 `tfsdk:"retry_interval"`
	ResendInterval types.Int64 `tfsdk:"resend_interval"`
	MaxRetries     types.Int64 `tfsdk:"max_retries"`
	MaxRedirects   types.Int64 `tfsdk:"max_redirects"`

	AcceptedStatusCodes types.List `tfsdk:"accepted_statuscodes"`
	NotificationIDList  types.List `tfsdk:"notification_list"`
	ExpiryNotification  types.Bool `tfsdk:"expiry_notification"`
	IgnoreTls           types.Bool `tfsdk:"ignore_tls"`
	UpsideDown          types.Bool `tfsdk:"upside_down"`
}

type Notification struct {
	ID        types.Int64  `tfsdk:"id"`
	UserId    types.Int64  `tfsdk:"user_id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Active    types.Bool   `tfsdk:"active"`
	IsDefault types.Bool   `tfsdk:"default"`
}

type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}
