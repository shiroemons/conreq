package requestid

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	id1 := Generate()
	id2 := Generate()

	if id1 == "" {
		t.Error("Generate() returned empty string")
	}

	if id2 == "" {
		t.Error("Generate() returned empty string")
	}

	if id1 == id2 {
		t.Error("Generate() returned the same ID twice")
	}

	if !IsValid(id1) {
		t.Errorf("Generated ID %s is not valid", id1)
	}

	if !IsValid(id2) {
		t.Errorf("Generated ID %s is not valid", id2)
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "valid UUID v4",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			want: true,
		},
		{
			name: "valid generated UUID",
			id:   Generate(),
			want: true,
		},
		{
			name: "invalid UUID - wrong format",
			id:   "not-a-uuid",
			want: false,
		},
		{
			name: "invalid UUID - empty",
			id:   "",
			want: false,
		},
		{
			name: "invalid UUID - wrong length",
			id:   "550e8400-e29b-41d4-a716",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.id); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
