package doi

import (
	"encoding/json"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"10.1000/182", "10.1000/182", false},
		{"https://doi.org/10.1000/182", "10.1000/182", false},
		{"http://doi.org/10.1000/182", "10.1000/182", false},
		{"  10.1000/182  ", "10.1000/182", false},
		{"invalid/doi", "", true},
		{"https://google.com/10.1000/182", "", true},
		{"10.1002/abc", "10.1002/abc", false}, // Old Wiley format
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("Parse() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestDoiJSON(t *testing.T) {
	type Container struct {
		D DOI `json:"doi"`
	}

	t.Run("Marshal", func(t *testing.T) {
		d, _ := Parse("10.1000/182")
		c := Container{D: d}
		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		expected := `{"doi":"10.1000/182"}`
		if string(data) != expected {
			t.Errorf("Marshal output = %s, want %s", string(data), expected)
		}
	})

	t.Run("Unmarshal Valid", func(t *testing.T) {
		jsonStr := `{"doi":"10.1000/182"}`
		var c Container
		if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if c.D.String() != "10.1000/182" {
			t.Errorf("Unmarshal value = %s, want 10.1000/182", c.D.String())
		}
	})

	t.Run("Unmarshal Invalid", func(t *testing.T) {
		jsonStr := `{"doi":"invalid"}`
		var c Container
		if err := json.Unmarshal([]byte(jsonStr), &c); err == nil {
			t.Error("Unmarshal should have failed for invalid DOI")
		}
	})
}
