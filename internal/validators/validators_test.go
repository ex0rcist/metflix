package validators_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/validators"
)

func TestValidateMetric(t *testing.T) {
	type args struct {
		name string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "correct counter", args: args{name: "testname", kind: "counter"}, wantErr: false},
		{name: "correct gauge", args: args{name: "testname", kind: "gauge"}, wantErr: false},

		{name: "name not present", args: args{name: "", kind: "counter"}, wantErr: true},
		{name: "incorrect name", args: args{name: "некорректноеимя", kind: "counter"}, wantErr: true},
		{name: "incorrect name", args: args{name: "incorrect name", kind: "gauge"}, wantErr: true},
		{name: "incorrect kind", args: args{kind: "gauger"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validators.ValidateMetric(tt.args.name, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
