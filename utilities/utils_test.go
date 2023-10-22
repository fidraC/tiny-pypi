package utilities_test

import (
	"pypi/utilities"
	"testing"
)

func TestNormalize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no change",
			in:   "foo",
			want: "foo",
		},
		{
			name: "lowercase",
			in:   "Foo",
			want: "foo",
		},
		{
			name: "spaces",
			in:   "Foo Bar",
			want: "foo-bar",
		},
		{
			name: "double dash",
			in:   "Foo--Bar",
			want: "foo-bar",
		},
		{
			name: "double underscore",
			in:   "Foo__Bar",
			want: "foo-bar",
		},
		{
			name: "sepcial characters",
			in:   "Foo!@#$%^&*()Bar",
			want: "foo-bar",
		},
		{
			name: "double period",
			in:   "foo..bar.txt",
			want: "foo-bar.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilities.Normalize(tt.in)
			if got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
