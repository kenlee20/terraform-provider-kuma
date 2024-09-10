package provider

import (
	"context"
	"fmt"
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &monitorsDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorsDataSource{}
)

func NewMonitorsDataSource() datasource.DataSource {
	return &monitorsDataSource{}
}

type MonitorsDataSourceModel struct {
	Monitors []monitorsModel `tfsdk:"monitors"`
}

type monitorsModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type monitorsDataSource struct {
	client *kuma.Client
}

func (d *monitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

// Schema defines the schema for the data source.
func (d *monitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Example identifier",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Example identifier",
						},
					},
				},
			},
		},
	}
}

func (d *monitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var states MonitorsDataSourceModel

	monitors, err := d.client.GetMonitors()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read monotors",
			err.Error(),
		)
		return
	}

	for _, monitor := range monitors {
		state := monitorsModel{
			ID:   types.Int64Value(int64(monitor.ID)),
			Name: types.StringValue(monitor.Name),
		}
		states.Monitors = append(states.Monitors, state)
	}

	// Set state
	diags := resp.State.Set(ctx, &states)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *monitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kuma.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
