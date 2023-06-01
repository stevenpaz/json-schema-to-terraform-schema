package tf_test

import (
	"reflect"
	"testing"

	"github.com/stevenpaz/tf-schema-gen/tf"
)

// TestTerraformSchema_Validate tests the Validate method of TerraformSchema.
func TestTerraformSchema_Validate(t *testing.T) {
	t.Parallel()

	type fields struct {
		Type           string
		Description    string
		Required       bool
		Computed       bool
		Optional       bool
		ValidationFunc string
	}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Type is required",
			fields: fields{
				Type:     "",
				Required: true,
			},
			want: []string{"Type is required"},
		},
		{
			name: "At least one of Required, Optional, or Computed must be true",
			fields: fields{
				Type: "string",
			},
			want: []string{"At least one of Required, Optional, or Computed must be true"},
		},
		{
			name: "Required, Optional, and Computed are mutually exclusive",
			fields: fields{
				Type:     "string",
				Required: true,
				Optional: true,
			},
			want: []string{
				"Required, Optional, and Computed are mutually exclusive",
			},
		},
		{
			name: "Required, Optional, and Computed are mutually exclusive",
			fields: fields{
				Type:     "string",
				Required: true,
				Computed: true,
			},
			want: []string{
				"Required, Optional, and Computed are mutually exclusive",
			},
		},
		{
			name: "Required, Optional, and Computed are mutually exclusive",
			fields: fields{
				Type:     "string",
				Optional: true,
				Computed: true,
			},
			want: []string{"Required, Optional, and Computed are mutually exclusive"},
		},
		{
			name: "Required, Optional, and Computed are mutually exclusive",
			fields: fields{
				Type:     "string",
				Required: true,
				Optional: true,
				Computed: true,
			},
			want: []string{"Required, Optional, and Computed are mutually exclusive"},
		},
		{
			name: "Multiple errors",
			fields: fields{
				Type:     "",
				Optional: true,
				Computed: true,
			},
			want: []string{
				"Type is required",
				"Required, Optional, and Computed are mutually exclusive",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			schema := tf.NewTerrformSchema("TestSchema")
			schema.SetType(test.fields.Type)
			schema.SetDescription(test.fields.Description)
			schema.SetRequired(test.fields.Required)
			schema.SetOptional(test.fields.Optional)
			schema.SetComputed(test.fields.Computed)
			schema.SetValidationFunc(test.fields.ValidationFunc)

			if got := schema.Validate(); !reflect.DeepEqual(got, test.want) {
				t.Errorf(
					"TerraformSchema.Validate() = %v, want %v",
					got,
					test.want)
			}
		})
	}
}
