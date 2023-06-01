package openapi_test

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stevenpaz/tf-schema-gen/internal"
	"github.com/stevenpaz/tf-schema-gen/openapi"
	"github.com/stevenpaz/tf-schema-gen/tf"
)

func TestBuildValidationFunc(t *testing.T) {
	tests := []struct {
		name string
		arg  *openapi3.Schema
		want string
	}{
		{
			name: "null case",
			arg:  nil,
			want: "",
		},
		{
			name: "no format",
			arg:  &openapi3.Schema{Type: "string"},
			want: "",
		},
		{
			name: "date format",
			arg:  &openapi3.Schema{Type: "string", Format: "date"},
			want: "validation.IsRFC3339Time",
		},
		{
			name: "date-time format",
			arg:  &openapi3.Schema{Type: "string", Format: "date-time"},
			want: "validation.IsRFC3339Time",
		},
		{
			name: "int inclusive minimum",
			arg: &openapi3.Schema{
				Type: "integer",
				Min:  internal.Float64Ptr(1),
			},
			want: "validation.IntAtLeast(1)",
		},
		{
			name: "int inclusive maximum",
			arg: &openapi3.Schema{
				Type: "integer",
				Max:  internal.Float64Ptr(1),
			},
			want: "validation.IntAtMost(1)",
		},
		{
			name: "int exclusive minimum",
			arg: &openapi3.Schema{
				Type:         "integer",
				Min:          internal.Float64Ptr(1),
				ExclusiveMin: true,
			},
			want: "validation.IntAtLeast(2)",
		},
		{
			name: "int exclusive maximum",
			arg: &openapi3.Schema{
				Type:         "integer",
				Max:          internal.Float64Ptr(5),
				ExclusiveMax: true,
			},
			want: "validation.IntAtMost(4)",
		},
		{
			name: "float32 inclusive minimum",
			arg: &openapi3.Schema{
				Type:   "number",
				Format: "float",
				Min:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtLeast(1.000000)",
		},
		{
			name: "float32 inclusive maximum",
			arg: &openapi3.Schema{
				Type:   "number",
				Format: "float",
				Max:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtMost(1.000000)",
		},
		// TF doesn't support int64, so we convert to float64.
		{
			name: "int64 inclusive minimum",
			arg: &openapi3.Schema{
				Type:   "integer",
				Format: "int64",
				Min:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtLeast(1.000000)",
		},
		{
			name: "int64 inclusive maximum",
			arg: &openapi3.Schema{
				Type:   "integer",
				Format: "int64",
				Max:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtMost(1.000000)",
		},
		{
			name: "int64 exclusive minimum",
			arg: &openapi3.Schema{
				Type:         "integer",
				Format:       "int64",
				Min:          internal.Float64Ptr(1),
				ExclusiveMin: true,
			},
			want: "validation.FloatAtLeast(2.000000)",
		},
		{
			name: "int64 exclusive maximum",
			arg: &openapi3.Schema{
				Type:         "integer",
				Format:       "int64",
				Max:          internal.Float64Ptr(5),
				ExclusiveMax: true,
			},
			want: "validation.FloatAtMost(4.000000)",
		},
		{
			name: "float64 inclusive minimum",
			arg: &openapi3.Schema{
				Type:   "number",
				Format: "double",
				Min:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtLeast(1.000000)",
		},
		{
			name: "float64 inclusive maximum",
			arg: &openapi3.Schema{
				Type:   "number",
				Format: "double",
				Max:    internal.Float64Ptr(1),
			},
			want: "validation.FloatAtMost(1.000000)",
		},
		{
			name: "float64 exclusive minimum",
			arg: &openapi3.Schema{
				Type:         "number",
				Format:       "double",
				Min:          internal.Float64Ptr(1),
				ExclusiveMin: true,
			},
			want: tf.BuildValidateFuncFloatAtLeastExclusive(1),
		},
		{
			name: "float64 exclusive maximum",
			arg: &openapi3.Schema{
				Type:         "number",
				Format:       "double",
				Max:          internal.Float64Ptr(1),
				ExclusiveMax: true,
			},
			want: tf.BuildValidateFuncFloatAtMostExclusive(1),
		},
		{
			name: "compound validation",
			arg: &openapi3.Schema{
				Type:         "integer",
				Min:          internal.Float64Ptr(1),
				ExclusiveMin: true,
				Max:          internal.Float64Ptr(5),
				ExclusiveMax: true,
			},
			want: "validation.All(validation.IntAtLeast(2),validation.IntAtMost(4))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openapi.BuildValidationFunc(tt.arg); got != tt.want {
				t.Errorf("BuildValidationFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTFValidationFunc(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "null case",
			arg:  "",
			want: "",
		},
		{
			name: "date format",
			arg:  "date",
			want: "validation.IsRFC3339Time",
		},
		{
			name: "date-time format",
			arg:  "date-time",
			want: "validation.IsRFC3339Time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openapi.GetTFValidationFunc(tt.arg); got != tt.want {
				t.Errorf("GetTFValidationFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTFType(t *testing.T) {
	tests := []struct {
		name    string
		arg     *openapi3.Schema
		want    string
		wantErr bool
	}{
		{
			name:    "null case",
			arg:     nil,
			want:    "",
			wantErr: true,
		},
		{
			name:    "unknown type",
			arg:     &openapi3.Schema{Type: "unknown"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "string",
			arg:     &openapi3.Schema{Type: "string"},
			want:    tf.TypeString,
			wantErr: false,
		},
		{
			name:    "int",
			arg:     &openapi3.Schema{Type: "integer"},
			want:    tf.TypeInt,
			wantErr: false,
		},
		{
			name:    "int8",
			arg:     &openapi3.Schema{Type: "integer", Format: "int8"},
			want:    tf.TypeInt,
			wantErr: false,
		},
		{
			name:    "int16",
			arg:     &openapi3.Schema{Type: "integer", Format: "int16"},
			want:    tf.TypeInt,
			wantErr: false,
		},
		{
			name:    "int32",
			arg:     &openapi3.Schema{Type: "integer", Format: "int32"},
			want:    tf.TypeInt,
			wantErr: false,
		},
		{
			name:    "int64",
			arg:     &openapi3.Schema{Type: "integer", Format: "int64"},
			want:    tf.TypeFloat,
			wantErr: false,
		},
		{
			name:    "bool",
			arg:     &openapi3.Schema{Type: "boolean"},
			want:    tf.TypeBool,
			wantErr: false,
		},
		{
			name:    "float32",
			arg:     &openapi3.Schema{Type: "number", Format: "float"},
			want:    tf.TypeFloat,
			wantErr: false,
		},
		{
			name:    "float64",
			arg:     &openapi3.Schema{Type: "number", Format: "double"},
			want:    tf.TypeFloat,
			wantErr: false,
		},
		{
			name:    "array",
			arg:     &openapi3.Schema{Type: "array"},
			want:    tf.TypeList,
			wantErr: false,
		},
		{
			name:    "object",
			arg:     &openapi3.Schema{Type: "object"},
			want:    tf.TypeMap,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openapi.GetTFType(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTFType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTFType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToTFSchema(t *testing.T) {
	tests := []struct {
		name       string
		schemaName string
		arg        *openapi3.Schema
		want       *tf.TerraformSchema
		wantErr    bool
	}{
		{
			name:    "null case",
			arg:     nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:       "unknown type",
			schemaName: "TestSchema",
			arg: &openapi3.Schema{
				Type: "unknown",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "Happy path",
			schemaName: "TestSchema",
			arg: &openapi3.Schema{
				Type:        "string",
				Format:      "date-time",
				Description: "test",
			},
			want: &tf.TerraformSchema{
				Name: "TestSchema",
				Properties: map[string]interface{}{
					"Type":           tf.TypeString,
					"ValidationFunc": "validation.IsRFC3339Time",
					"Description":    "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openapi.ConvertToTFSchema(tt.schemaName, tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToTFSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToTFSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
