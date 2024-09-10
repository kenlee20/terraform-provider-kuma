package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type testTFStruect struct {
	Num  types.Int64
	Str  types.String
	Bol  types.Bool
	StrL types.List
	NumL types.List
	BolL types.List
}

type testStruect struct {
	Num  int
	Str  string
	Bol  bool
	StrL []string
	NumL []int
	BolL []bool
}

func TestConvertFromModel(t *testing.T) {
	var output testStruect
	planTF := testTFStruect{
		Num:  types.Int64Value(10),
		Str:  types.StringValue("test"),
		Bol:  types.BoolValue(true),
		StrL: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("test1"), types.StringValue("test2")}),
		NumL: types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(20), types.Int64Value(30)}),
		BolL: types.ListValueMust(types.BoolType, []attr.Value{types.BoolValue(true), types.BoolValue(false)}),
	}

	if err := ConvertStruct(planTF, &output, false); err != nil {
		t.Errorf("Convert() error = %v", err)
	}

	if output.Num != int(planTF.Num.ValueInt64()) {
		t.Errorf("Convert() error = %v", output.Num)
	}
	if output.Str != planTF.Str.ValueString() {
		t.Errorf("Convert() error = %v", output.Str)
	}
	if output.Bol != planTF.Bol.ValueBool() {
		t.Errorf("Convert() error = %v", output.Bol)
	}
	if len(output.StrL) != len(planTF.StrL.Elements()) {
		t.Errorf("Convert() error = %v", output.StrL)
	}
	if len(output.NumL) != len(planTF.NumL.Elements()) {
		t.Errorf("Convert() error = %v", output.NumL)
	}
	if len(output.BolL) != len(planTF.BolL.Elements()) {
		t.Errorf("Convert() error = %v", output.BolL)
	}
}

func TestConvertToModel(t *testing.T) {
	var output testTFStruect
	planS := testStruect{
		Num:  10,
		Str:  "test",
		Bol:  true,
		StrL: []string{"test1", "test2"},
		NumL: []int{20, 30},
		BolL: []bool{true, false},
	}

	if err := ConvertStruct(planS, &output, true); err != nil {
		t.Errorf("Convert() error = %v", err)
	}

	if output.Num.ValueInt64() != int64(planS.Num) {
		t.Errorf("Convert() error = %v", output.Num.ValueInt64())
	}
	if output.Str.ValueString() != planS.Str {
		t.Errorf("Convert() error = %v", output.Str.ValueString())
	}
	if output.Bol.ValueBool() != planS.Bol {
		t.Errorf("Convert() error = %v", output.Bol.ValueBool())
	}
	if len(output.StrL.Elements()) != 2 {
		t.Errorf("Convert() error = %v", output.StrL)
	}
	if len(output.NumL.Elements()) != len(planS.NumL) {
		t.Errorf("Convert() error = %v", output.NumL)
	}
	if len(output.BolL.Elements()) != len(planS.BolL) {
		t.Errorf("Convert() error = %v", output.BolL)
	}
}
