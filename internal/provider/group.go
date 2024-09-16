package provider

import (
	"context"
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Tags        types.Map    `tfsdk:"tags"`
}

func (g *groupModel) Convert(ctx context.Context) (*kuma.Monitor, diag.Diagnostics) {
	var diags diag.Diagnostics

	tmpTag := make(map[string]string)
	newTag := make([]kuma.MonitorTag, 0)

	if !g.Tags.IsNull() && !g.Tags.IsUnknown() {
		diags = g.Tags.ElementsAs(ctx, &tmpTag, false)
		if diags.HasError() {
			return nil, diags
		}
	}

	for k, v := range tmpTag {
		newTag = append(newTag, kuma.MonitorTag{
			Name:  k,
			Value: v,
		})
	}

	return &kuma.Monitor{
		ID:          g.ID.ValueInt64(),
		Name:        g.Name.ValueString(),
		Type:        g.Type.ValueString(),
		Description: g.Description.ValueString(),
		Tags:        newTag,
	}, nil
}

func (g *groupModel) ConvertFrom(ctx context.Context, m *kuma.Monitor) diag.Diagnostics {
	var diags diag.Diagnostics
	newTag := make(map[string]string)

	for _, tag := range m.Tags {
		newTag[tag.Name] = tag.Value
	}

	g.ID = types.Int64Value(m.ID)
	g.Name = types.StringValue(m.Name)
	g.Type = types.StringValue(m.Type)
	g.Description = types.StringValue(m.Description)
	if len(newTag) == 0 {
		g.Tags = types.MapNull(types.StringType)
		return diags
	}
	g.Tags, diags = types.MapValueFrom(ctx, types.StringType, newTag)

	return diags
}
