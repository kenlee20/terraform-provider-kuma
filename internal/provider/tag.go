package provider

import (
	"terraform-provider-kuma/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Tag struct {
	ID    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Color types.String `tfsdk:"color"`
}

func (t *Tag) Convert() *kuma.Tag {
	return &kuma.Tag{
		ID:    t.ID.ValueInt64(),
		Name:  t.Name.ValueString(),
		Color: t.Color.ValueString(),
	}
}

func (t *Tag) ConvertFrom(tag kuma.Tag) {
	t.ID = types.Int64Value(tag.ID)
	t.Name = types.StringValue(tag.Name)
	t.Color = types.StringValue(tag.Color)
}
