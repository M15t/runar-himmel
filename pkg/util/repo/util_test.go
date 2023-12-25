package repoutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseConds(t *testing.T) {
	type args struct {
		conds []any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "empty conds",
			args: args{conds: []any{}},
			want: []any{},
		},
		{
			name: "empty string",
			args: args{conds: []any{""}},
			want: []any{},
		},
		{
			name: "empty map",
			args: args{conds: []any{map[string]any{}}},
			want: []any{},
		},
		{
			name: "empty array",
			args: args{conds: []any{[]any{}}},
			want: []any{},
		},
		{
			name: "nil",
			args: args{conds: []any{nil}},
			want: []any{},
		},
		{
			name: "ID-like string",
			args: args{conds: []any{"id; sql injection attempt--"}},
			want: []any{"id = ?", "id; sql injection attempt--"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseConds(tt.args.conds))
		})
	}
}
