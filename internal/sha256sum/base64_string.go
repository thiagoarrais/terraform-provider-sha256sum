package sha256sum

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	Base64String Base64StringType                           = Base64StringType{}
	_            basetypes.StringTypable                    = Base64StringType{}
	_            basetypes.StringValuable                   = Base64StringValue{}
	_            basetypes.StringValuableWithSemanticEquals = Base64StringValue{}
)

type Base64StringType struct {
	basetypes.StringType
}

func (t Base64StringType) Equal(o attr.Type) bool {
	other, ok := o.(Base64StringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t Base64StringType) String() string {
	return "Base64StringType"
}

func (t Base64StringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := Base64StringValue{
		StringValue: in,
	}

	return value, nil
}

func (t Base64StringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t Base64StringType) ValueType(ctx context.Context) attr.Value {
	return t.NullValue()
}

func (t Base64StringType) NullValue() Base64StringValue {
	return Base64StringValue{StringValue: types.StringNull()}
}

func (t Base64StringType) Value(v string) Base64StringValue {
	return Base64StringValue{StringValue: types.StringValue(v)}
}

type Base64StringValue struct {
	basetypes.StringValue
}

func (v Base64StringValue) Equal(o attr.Value) bool {
	other, ok := o.(Base64StringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v Base64StringValue) Type(ctx context.Context) attr.Type {
	return Base64String
}

func (v Base64StringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(Base64StringValue)

	if !ok {
		return false, diags
	}

	old := v.StringValue.ValueString()
	new := newValue.StringValue.ValueString()

	diags.AddError("This is a forced error", "Nothing is wrong, but we're returning an error anyway")

	if old == new {
		return true, diags
	}

	if len(new) == 64 {
		oldDecoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(old))
		oldHash := sha256.New()
		_, err := io.Copy(oldHash, oldDecoder)
		if err != nil {
			return false, diags
		}

		return fmt.Sprintf("%x", oldHash.Sum(nil)) == new, diags
	} else if len(old) == 64 {
		newDecoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(new))
		newHash := sha256.New()
		_, err := io.Copy(newHash, newDecoder)
		if err != nil {
			return false, diags
		}
		return fmt.Sprintf("%x", newHash.Sum(nil)) == old, diags

	}

	return false, diags
}
