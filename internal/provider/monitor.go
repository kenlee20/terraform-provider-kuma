package provider

import (
	"context"

	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MonitorResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PathName    types.String `tfsdk:"path_name"`
	Url         types.String `tfsdk:"url"`
	Type        types.String `tfsdk:"type"`
	// Active      types.Bool   `tfsdk:"active"`
	// Timeout        types.Int64 `tfsdk:"timeout"`
	Interval       types.Int64 `tfsdk:"interval"`
	RetryInterval  types.Int64 `tfsdk:"retry_interval"`
	ResendInterval types.Int64 `tfsdk:"resend_interval"`
	MaxRetries     types.Int64 `tfsdk:"max_retries"`
	MaxRedirects   types.Int64 `tfsdk:"max_redirects"`

	Method           types.String `tfsdk:"http_option_method"`
	HTTPBodyEncoding types.String `tfsdk:"http_option_body_encoding"`
	Body             types.String `tfsdk:"http_option_body"`
	Headers          types.String `tfsdk:"http_option_headers"`

	AcceptedStatusCodes types.List `tfsdk:"accepted_statuscodes"`
	NotificationIDList  types.List `tfsdk:"notification_list"`
	ExpiryNotification  types.Bool `tfsdk:"expiry_notification"`
	IgnoreTls           types.Bool `tfsdk:"ignore_tls"`
	UpsideDown          types.Bool `tfsdk:"upside_down"`
	Tags                types.Map  `tfsdk:"tags"`
}

type MonitorTag struct {
	TagId types.Int64  `tfsdk:"tag_id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (m *MonitorResourceModel) Convert() (*kuma.Monitor, diag.Diagnostics) {
	ctx := context.Background()
	var tmpTag map[string]string
	var newTag []kuma.MonitorTag
	var newAcceptedStatusCodes []string
	var newNotificationIDList []int64

	err := m.Tags.ElementsAs(ctx, &tmpTag, true)
	if err.HasError() {
		return nil, err
	}

	for key, value := range tmpTag {
		newTag = append(newTag, kuma.MonitorTag{
			Name:  key,
			Value: value,
		})
	}

	err = m.AcceptedStatusCodes.ElementsAs(ctx, &newAcceptedStatusCodes, true)
	if err.HasError() {
		return nil, err
	}
	err = m.NotificationIDList.ElementsAs(ctx, &newNotificationIDList, true)
	if err.HasError() {
		return nil, err
	}

	return &kuma.Monitor{
		ID:                  m.ID.ValueInt64(),
		Name:                m.Name.ValueString(),
		Description:         m.Description.ValueString(),
		PathName:            m.PathName.ValueString(),
		Url:                 m.Url.ValueString(),
		Method:              m.Method.ValueString(),
		HTTPBodyEncoding:    m.HTTPBodyEncoding.ValueString(),
		Body:                m.Body.ValueString(),
		Headers:             m.Headers.ValueString(),
		Type:                m.Type.ValueString(),
		Interval:            m.Interval.ValueInt64(),
		RetryInterval:       m.RetryInterval.ValueInt64(),
		ResendInterval:      m.ResendInterval.ValueInt64(),
		MaxRetries:          m.MaxRetries.ValueInt64(),
		MaxRedirects:        m.MaxRedirects.ValueInt64(),
		AcceptedStatusCodes: newAcceptedStatusCodes,
		NotificationIDList:  newNotificationIDList,
		ExpiryNotification:  m.ExpiryNotification.ValueBool(),
		IgnoreTls:           m.IgnoreTls.ValueBool(),
		UpsideDown:          m.UpsideDown.ValueBool(),
		Tags:                newTag,
	}, nil
}

func (m *MonitorResourceModel) ConvertFrom(stu kuma.Monitor) (err diag.Diagnostics) {
	ctx := context.Background()

	m.ID = types.Int64Value(stu.ID)
	m.Name = types.StringValue(stu.Name)
	m.Description = types.StringValue(stu.Description)
	m.PathName = types.StringValue(stu.PathName)
	m.Url = types.StringValue(stu.Url)
	m.Method = types.StringValue(stu.Method)
	m.HTTPBodyEncoding = types.StringValue(stu.HTTPBodyEncoding)
	m.Body = types.StringValue(stu.Body)
	m.Headers = types.StringValue(stu.Headers)
	m.Type = types.StringValue(stu.Type)
	m.Interval = types.Int64Value(stu.Interval)
	m.RetryInterval = types.Int64Value(stu.RetryInterval)
	m.ResendInterval = types.Int64Value(stu.ResendInterval)
	m.MaxRedirects = types.Int64Value(stu.MaxRedirects)
	m.MaxRetries = types.Int64Value(stu.MaxRetries)
	m.ExpiryNotification = types.BoolValue(stu.ExpiryNotification)
	m.IgnoreTls = types.BoolValue(stu.IgnoreTls)
	m.UpsideDown = types.BoolValue(stu.UpsideDown)

	m.AcceptedStatusCodes, err = types.ListValueFrom(ctx, types.StringType, stu.AcceptedStatusCodes)
	if err.HasError() {
		return err
	}

	m.NotificationIDList, err = types.ListValueFrom(ctx, types.Int64Type, stu.NotificationIDList)
	if err.HasError() {
		return err
	}

	newTag := make(map[string]string)
	for _, tag := range stu.Tags {
		newTag[tag.Name] = tag.Value
	}

	m.Tags, err = types.MapValueFrom(ctx, types.StringType, newTag)

	return err
}
