package provider

import (
	"terraform-provider-upkuapi/internal/kuma"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMonitorConvert(t *testing.T) {
	providerPlan := MonitorResourceModel{
		ID:                  types.Int64Value(1),
		Name:                types.StringValue("name"),
		Description:         types.StringValue("description test"),
		PathName:            types.StringValue("name"),
		Url:                 types.StringValue("https://example.com"),
		Method:              types.StringValue("GET"),
		Type:                types.StringValue("http"),
		Interval:            types.Int64Value(1),
		MaxRedirects:        types.Int64Value(1),
		RetryInterval:       types.Int64Value(1),
		ResendInterval:      types.Int64Value(1),
		MaxRetries:          types.Int64Value(1),
		ExpiryNotification:  types.BoolValue(true),
		UpsideDown:          types.BoolValue(true),
		IgnoreTls:           types.BoolValue(true),
		AcceptedStatusCodes: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("200")}),
		NotificationIDList:  types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(5)}),
		Tags: []MonitorTag{
			{
				TagId: types.Int64Value(5),
				Name:  types.StringValue("tag_name"),
				Value: types.StringValue("tag_value"),
			},
		},
	}

	plan, err := providerPlan.Convert()
	if err.HasError() {
		t.Fatal(err)
	}

	t.Logf("\n%+v\n", plan)
	if len(plan.AcceptedStatusCodes) == 0 || len(plan.NotificationIDList) == 0 {
		t.Fatal("failed to convert")
	}
	if len(plan.Tags) == 0 {
		t.Fatal("failed to convert")
	}
}

func TestMonitorConvertFrom(t *testing.T) {
	var providerPlan MonitorResourceModel

	plan := kuma.Monitor{
		ID:                  1,
		Name:                "name",
		Description:         "description test",
		PathName:            "name",
		Url:                 "https://example.com",
		Method:              "GET",
		Type:                "http",
		Interval:            1,
		MaxRedirects:        1,
		RetryInterval:       1,
		ResendInterval:      1,
		MaxRetries:          1,
		ExpiryNotification:  true,
		UpsideDown:          true,
		IgnoreTls:           true,
		AcceptedStatusCodes: []string{"200"},
		NotificationIDList:  []int64{5},
		Tags: []kuma.MonitorTag{
			{
				TagId: 5,
				Name:  "tag_name",
				Value: "tag_value",
			},
		},
	}

	err := providerPlan.ConvertFrom(plan)
	if err.HasError() {
		t.Fatal(err)
	}

	if len(providerPlan.Tags) == 0 {
		t.Fatal("failed to convert")
	}
	if len(providerPlan.AcceptedStatusCodes.Elements()) == 0 {
		t.Fatal("failed to convert")
	}
	if len(providerPlan.NotificationIDList.Elements()) == 0 {
		t.Fatal("failed to convert")
	}
}
