package validators_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/validators"
)

func TestEnsureNamePresent(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no error when name present",
			args:    args{name: "testname"},
			wantErr: false,
		},
		{
			name:    "has error when name is not present",
			args:    args{name: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validators.EnsureNamePresent(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("EnsureNamePresent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	type args struct {
		name string
		kind string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "correct name",
			args:    args{name: "correctname1"},
			wantErr: false,
		},
		{
			name:    "incorrect name",
			args:    args{name: "некорректноеимя"},
			wantErr: true,
		},
		{
			name:    "incorrect name",
			args:    args{name: "incorrect name"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validators.ValidateName(tt.args.name, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateKind(t *testing.T) {
	type args struct {
		kind string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "correct kind = counter",
			args:    args{kind: "counter"},
			wantErr: false,
		},
		{
			name:    "correct kind = gauge",
			args:    args{kind: "gauge"},
			wantErr: false,
		},
		{
			name:    "incorrect kind = gauger",
			args:    args{kind: "gauger"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validators.ValidateKind(tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("ValidateKind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
