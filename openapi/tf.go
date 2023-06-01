package openapi

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stevenpaz/tf-schema-gen/internal"
	"github.com/stevenpaz/tf-schema-gen/tf"
)

func OpenAPI3ToTerraform(filePath string) (*tf.TerraformScope, error) {
	doc, err := openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	scope := tf.NewTerrformScope(doc.Info.Title)

	for name, schema := range doc.Components.Schemas {
		if ts, err := ConvertToTFSchema(name, scope, schema.Value); err == nil {
			scope.AddSchema(ts)
		} else {
			return nil,
				fmt.Errorf("failed to convert schema '%s': %w", name, err)
		}
	}

	return scope, nil
}

// ConvertToTFSchema converts an OpenAPI schema to a Terraform schema.
func ConvertToTFSchema(
	name string,
	scope *tf.TerraformScope,
	s *openapi3.Schema,
) (*tf.TerraformSchema, error) {
	if s == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	tfSchema := tf.NewTerrformSchema(name, scope)

	for name, prop := range s.Properties {
		propSchema := prop.Value

		tfProp := tf.NewTerraformProperty()

		if vf := BuildValidationFunc(propSchema); len(vf) > 0 {
			tfProp.SetValidateFunc(vf)
			tfSchema.HasValidateFuncs = true
		}

		tfProp.SetDescription(propSchema.Description)

		if t, err := GetTFType(propSchema); err == nil {
			tfProp.Type = t
		} else {
			return nil, err
		}

		if propSchema.Nullable || propSchema.AllowEmptyValue {
			tfProp.SetOptional(true)
		}

		if propSchema.ReadOnly {
			tfProp.SetComputed(true)
		}

		// Check if the property is required.
		for _, required := range s.Required {
			if required == name {
				tfProp.SetRequired(true)
				break
			}
		}

		// If not marked as Read-Only or Required, mark as optional.
		if !tfProp.IsComputed() && !tfProp.IsRequired() {
			tfProp.SetOptional(true)
		}

		// validate tfProp before adding and throw error if invalid
		if errs := tfProp.Validate(); len(errs) > 0 {
			return nil, fmt.Errorf(strings.Join(errs, "\n"))
		}

		tfSchema.AddProp(internal.ToSnakeCase(name), tfProp)
	}

	return tfSchema, nil
}

// BuildValidationFunc builds a Terraform validation from an OpenAPI schema.
func BuildValidationFunc(s *openapi3.Schema) string {
	if s == nil {
		return ""
	}

	t := s.Type

	f := []string{}

	if len(s.Format) > 0 {
		temp := GetTFValidationFunc(s.Format)
		if len(temp) > 0 {
			f = append(f, temp)
		}
	}

	if s.Min != nil {
		bound := *s.Min

		if t == TypeInteger && s.Format != FormatInt64 {
			if s.ExclusiveMin {
				bound++
			}

			f = append(f, fmt.Sprintf(tf.ValidateFuncIntAtLeast, int(bound)))
		} else if t == TypeInteger && s.Format == FormatInt64 {
			if s.ExclusiveMin {
				bound++
			}

			f = append(f, fmt.Sprintf(tf.ValidateFuncFloatAtLeast, bound))
		} else if t == TypeNumber {
			if s.ExclusiveMin {
				f = append(f, tf.BuildValidateFuncFloatAtLeastExclusive(bound))
			} else {
				f = append(f, fmt.Sprintf(tf.ValidateFuncFloatAtLeast, bound))
			}
		}
	}

	if s.Max != nil {
		bound := *s.Max

		if t == TypeInteger && s.Format != FormatInt64 {
			if s.ExclusiveMax {
				bound--
			}

			f = append(f, fmt.Sprintf(tf.ValidateFuncIntAtMost, int(bound)))
		} else if t == TypeInteger && s.Format == FormatInt64 {
			if s.ExclusiveMax {
				bound--
			}

			f = append(f, fmt.Sprintf(tf.ValidateFuncFloatAtMost, bound))
		} else if t == TypeNumber {
			if s.ExclusiveMax {
				f = append(f, tf.BuildValidateFuncFloatAtMostExclusive(bound))
			} else {
				f = append(f, fmt.Sprintf(tf.ValidateFuncFloatAtMost, bound))
			}
		}
	}

	if len(f) == 0 {
		return ""
	}

	if len(f) == 1 {
		return fmt.Sprintf("validation.ToDiagFunc(%s)", f[0])
	}

	return fmt.Sprintf("validation.ToDiagFunc(validation.All(%s))", strings.Join(f, ","))
}

// GetTFType returns the Terraform type that corresponds to the given OpenAPI
// type.
func GetTFType(s *openapi3.Schema) (string, error) {
	if s == nil {
		return "", fmt.Errorf("schema is nil")
	}

	t := s.Type
	switch t {
	case "string":
		return tf.TypeString, nil
	case "boolean":
		return tf.TypeBool, nil
	case "integer":
		if s.Format == FormatInt64 {
			// Terraform does not have an Int64 type, so use Float instead,
			// which is float64 under the hood.
			return tf.TypeFloat, nil
		}

		return tf.TypeInt, nil
	case "number":
		return tf.TypeFloat, nil
	case "array":
		return tf.TypeList, nil
	case "object":
		return tf.TypeMap, nil
	default:
		return "", fmt.Errorf("unsupported type %s", t)
	}
}

// GetTFValidationFunc returns a Terraform validation function that corresponds
// to the given format.
func GetTFValidationFunc(format string) string {
	switch format {
	case "date-time", "date":
		return tf.ValidateFuncRFC3339Time
	}

	return ""
}
