package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequestRatePerMinute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rpm     types.Int64
		want    int
		wantErr bool
	}{
		{name: "omitted", rpm: types.Int64Null(), want: 300},
		{name: "unknown", rpm: types.Int64Unknown(), want: 300},
		{name: "minimum", rpm: types.Int64Value(1), want: 1},
		{name: "maximum", rpm: types.Int64Value(6000), want: 6000},
		{name: "zero", rpm: types.Int64Value(0), wantErr: true},
		{name: "above maximum", rpm: types.Int64Value(6001), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, diags := requestRatePerMinute(tt.rpm)
			if tt.wantErr {
				if !diags.HasError() {
					t.Fatalf("requestRatePerMinute() error = false, want an error for %s", tt.rpm)
				}
				return
			}
			if diags.HasError() {
				t.Fatalf("requestRatePerMinute() unexpected error: %v", diags)
			}
			if got != tt.want {
				t.Fatalf("requestRatePerMinute() = %d, want %d", got, tt.want)
			}
		})
	}
}
