package tf

import (
	"github.com/stevenpaz/tf-schema-gen/internal"
)

// TerraformScope represents a collection of Terraform Schemas.
type TerraformScope struct {
	Name          string
	NameCamelCase string
	NameSnakeCase string
	Schemas       []*TerraformSchema
}

// NewTerrformScope creates a new TerraformScope.
func NewTerrformScope(name string) *TerraformScope {
	return &TerraformScope{
		Name:          name,
		NameCamelCase: internal.ToCamelCase(name),
		NameSnakeCase: internal.ToSnakeCase(name),
		Schemas:       make([]*TerraformSchema, 0),
	}
}

// AddSchema adds a TerraformSchema to the TerraformScope.
func (ts *TerraformScope) AddSchema(schema *TerraformSchema) {
	if ts.Schemas == nil {
		ts.Schemas = make([]*TerraformSchema, 0)
	}

	ts.Schemas = append(ts.Schemas, schema)
}

// Schema represents a Terraform Schema.
type TerraformSchema struct {
	Scope            *TerraformScope
	Name             string
	NameCamelCase    string
	NameSnakeCase    string
	Properties       map[string]TerraformProperty
	HasValidateFuncs bool
}

// TerraformProperty represents a property of a Terraform Schema.
type TerraformProperty struct {
	Type         string
	Required     *bool
	Optional     *bool
	Computed     *bool
	Description  *string
	ValidateFunc *string
}

// NewTerrformSchema creates a new TerraformSchema.
func NewTerrformSchema(name string, scope *TerraformScope) *TerraformSchema {
	return &TerraformSchema{
		Scope:         scope,
		Name:          name,
		NameCamelCase: internal.ToCamelCase(name),
		NameSnakeCase: internal.ToSnakeCase(name),
		Properties:    make(map[string]TerraformProperty),
	}
}

// NewTerraformProperty creates a new TerraformProperty.
func NewTerraformProperty() *TerraformProperty {
	return &TerraformProperty{}
}

// SetDescription sets the description of the TerraformProperty.
func (tp *TerraformProperty) SetDescription(desc string) {
	if len(desc) > 0 {
		tp.Description = internal.StringPtr(desc)
	}
}

func (tp *TerraformProperty) SetRequired(required bool) {
	tp.Required = internal.BoolPtr(required)
}

func (tp *TerraformProperty) SetOptional(optional bool) {
	tp.Optional = internal.BoolPtr(optional)
}

func (tp *TerraformProperty) SetComputed(computed bool) {
	tp.Computed = internal.BoolPtr(computed)
}

func (tp *TerraformProperty) SetValidateFunc(validateFunc string) {
	if len(validateFunc) > 0 {
		tp.ValidateFunc = internal.StringPtr(validateFunc)
	}
}

func (ts *TerraformSchema) AddProp(name string, prop *TerraformProperty) {
	ts.Properties[name] = *prop
}

// IsRequired returns true if the TerraformProperty is required.
func (tp TerraformProperty) IsRequired() bool {
	return tp.Required != nil && *tp.Required
}

// IsOptional returns true if the TerraformProperty is optional.
func (tp TerraformProperty) IsOptional() bool {
	return tp.Optional != nil && *tp.Optional
}

// IsComputed returns true if the TerraformProperty is computed.
func (tp TerraformProperty) IsComputed() bool {
	return tp.Computed != nil && *tp.Computed
}

// Validate validates the TerraformSchema.
func (ts TerraformSchema) Validate() []string {
	var errs []string

	if ts.Name == "" {
		errs = append(errs, "name is required")
	}

	if len(ts.Properties) == 0 {
		errs = append(errs, "at least one property is required")
	} else {
		for name, prop := range ts.Properties {
			errs = append(errs, prop.Validate()...)
			if len(errs) > 0 {
				errs = append(errs, "property: "+name)
			}
		}
	}

	return errs
}

// Validate validates the TerraformProperty.
func (tp TerraformProperty) Validate() []string {
	var errs []string

	if tp.Type == "" {
		errs = append(errs, "Type is required")
	}

	req := tp.IsRequired()
	opt := tp.IsOptional()
	comp := tp.IsComputed()

	if !req && !opt && !comp {
		errs = append(
			errs,
			"At least one of Required, Optional, or Computed must be true")
	} else if !internal.OnlyOneTrue(comp, opt, req) {
		errs = append(
			errs,
			"Required, Optional, and Computed are mutually exclusive")
	}

	return errs
}
