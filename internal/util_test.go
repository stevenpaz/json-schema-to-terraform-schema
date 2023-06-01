package internal_test

import (
	"testing"

	"github.com/stevenpaz/tf-schema-gen/internal"
)

// TestOnlyOneTrue tests the OnlyOneTrue function.
func TestOnlyOneTrue(t *testing.T) {
	t.Parallel()

	type args struct {
		a bool
		b bool
		c bool
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "a is true",
			args: args{
				a: true,
				b: false,
				c: false,
			},
			want: true,
		},
		{
			name: "b is true",
			args: args{
				a: false,
				b: true,
				c: false,
			},
			want: true,
		},
		{
			name: "c is true",
			args: args{
				a: false,
				b: false,
				c: true,
			},
			want: true,
		},
		{
			name: "a and b are true",
			args: args{
				a: true,
				b: true,
				c: false,
			},
			want: false,
		},
		{
			name: "a and c are true",
			args: args{
				a: true,
				b: false,
				c: true,
			},
			want: false,
		},
		{
			name: "b and c are true",
			args: args{
				a: false,
				b: true,
				c: true,
			},
			want: false,
		},
		{
			name: "a, b, and c are true",
			args: args{
				a: true,
				b: true,
				c: true,
			},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := internal.OnlyOneTrue(
				test.args.a,
				test.args.b,
				test.args.c); got != test.want {
				t.Errorf("OnlyOneTrue() = %v, want %v", got, test.want)
			}
		})
	}
}
