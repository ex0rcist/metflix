package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if server.storage == nil {
		t.Fatal("expected storage to not be nil")
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Config
		wantErr bool
	}{
		{
			name:    "shortcut",
			args:    []string{"-a0.0.0.0:8080"},
			want:    Config{Address: "0.0.0.0:8080"},
			wantErr: false,
		},
		{
			name:    "shortcut",
			args:    []string{"--address=127.0.0.1:81"},
			want:    Config{Address: "127.0.0.1:81"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Address: "default"}

			err := config.parseFlags("progname", tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected no error, got %v", err)
			}
			if tt.want != *config {
				t.Errorf("Expected %v, got %v", tt.want, config)
			}
		})
	}
}
