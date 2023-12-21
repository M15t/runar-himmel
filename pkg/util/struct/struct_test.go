package structutil_test

import (
	"reflect"
	structutil "runar-himmel/pkg/util/struct"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

type updateData struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,phone"`
}

var testData = updateData{
	FirstName: strPtr("Test"),
	LastName:  strPtr("Gopher"),
	Email:     strPtr("test.gopher@mail.com"),
	Phone:     strPtr("+84989898989"),
}
var resData = map[string]interface{}{
	"first_name": testData.FirstName,
	"last_name":  testData.LastName,
	"email":      testData.Email,
	"phone":      testData.Phone,
}

func TestToMap(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantOut map[string]interface{}
	}{
		{
			name: "Success",
			args: args{
				in: testData,
			},
			wantOut: resData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := structutil.ToMap(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ToMap() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
