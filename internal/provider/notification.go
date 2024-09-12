package provider

import (
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Notification struct {
	ID        types.Int64  `tfsdk:"id"`
	UserId    types.Int64  `tfsdk:"user_id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Active    types.Bool   `tfsdk:"active"`
	IsDefault types.Bool   `tfsdk:"default"`
}

func (n *Notification) Convert() *kuma.Notification {
	return &kuma.Notification{
		ID:        n.ID.ValueInt64(),
		UserId:    n.UserId.ValueInt64(),
		Name:      n.Name.ValueString(),
		Type:      n.Type.ValueString(),
		Active:    n.Active.ValueBool(),
		IsDefault: n.IsDefault.ValueBool(),
	}
}

func (n *Notification) ConvertFrom(k kuma.Notification) {
	n.ID = types.Int64Value(int64(k.ID))
	n.UserId = types.Int64Value(int64(k.UserId))
	n.Name = types.StringValue(k.Name)
	n.Type = types.StringValue(k.Type)
	n.Active = types.BoolValue(k.Active)
	n.IsDefault = types.BoolValue(k.IsDefault)
}
