package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestGroupConvert(t *testing.T) {
	plan := groupModel{
		ID:   types.Int64Unknown(),
		Name: types.StringValue("test"),
		Tags: types.MapUnknown(types.StringType),
	}

	item, diag := plan.Convert(context.TODO())
	if diag.HasError() {
		t.Fatal(diag)
	}

	t.Logf("%+v", item)
}
