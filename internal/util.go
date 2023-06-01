package internal

import "github.com/iancoleman/strcase"

// OnlyOneTrue returns true if only one of the three arguments is true.
func OnlyOneTrue(a, b, c bool) bool {
	return (a && !b && !c) || (!a && b && !c) || (!a && !b && c)
}

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(b bool) *bool {
	return &b
}

// Float64Ptr returns a pointer to the given float64.
func Float64Ptr(f float64) *float64 {
	return &f
}

// ToSnakeCase converts the given string to snake case.
func ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
}

// ToCamelCase converts the given string to camel case.
func ToCamelCase(s string) string {
	return strcase.ToCamel(s)
}
