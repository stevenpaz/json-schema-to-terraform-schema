package tf

import "fmt"

// Constants for the primitive terraform types.
const (
	TypeString = "TypeString"
	TypeBool   = "TypeBool"
	TypeInt    = "TypeInt"
	TypeFloat  = "TypeFloat"
	TypeList   = "TypeList"
	TypeMap    = "TypeMap"
)

// Constants for Terraform SDKv2 validation functions.
const (
	ValidateFuncRFC3339Time  = "validation.IsRFC3339Time"
	ValidateFuncIntAtLeast   = "validation.IntAtLeast(%d)"
	ValidateFuncIntAtMost    = "validation.IntAtMost(%d)"
	ValidateFuncFloatAtLeast = "validation.FloatAtLeast(%f)"
	ValidateFuncFloatAtMost  = "validation.FloatAtMost(%f)"

	ValidateFuncFloatAtLeastExclusive = `func(i interface{}, p string) (s []string, es []error) {
	v, ok := i.(float64)
	if !ok {
		es = append(es, errors.New("expected type of float"))
		return
	}

	if v <= %f {
		es = append(es, errors.New("expected more than (%f)"))
		return
	}

	return
}`

	ValidateFuncFloatAtMostExclusive = `func(i interface{}, p string) (s []string, es []error) {
	v, ok := i.(float64)
	if !ok {
		es = append(es, errors.New("expected type of float"))
		return
	}

	if v >= %f {
		es = append(es, errors.New("expected less than (%f)"))
		return
	}

	return
}`
)

// Constants for Custom Validation Functions.
const (
	ValidateFuncMatchRegExPattern = "MatchRegExPattern"
)

func BuildValidateFuncFloatAtMostExclusive(max float64) string {
	return fmt.Sprintf(ValidateFuncFloatAtMostExclusive, max, max)
}

func BuildValidateFuncFloatAtLeastExclusive(min float64) string {
	return fmt.Sprintf(ValidateFuncFloatAtLeastExclusive, min, min)
}
