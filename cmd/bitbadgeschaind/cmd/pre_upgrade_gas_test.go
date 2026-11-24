package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateAppTomlMinGasPrices(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		wantChanged bool
		wantVal     string
	}{
		{
			name:        "rewrites 10ubadge to 10000000000abadge",
			in:          "minimum-gas-prices = \"10ubadge\"\n\n[api]\nenable = true\n",
			wantChanged: true,
			wantVal:     `minimum-gas-prices = "10000000000abadge"`,
		},
		{
			name:        "rewrites fractional 0.025ubadge cleanly",
			in:          "minimum-gas-prices = '0.025ubadge'\n",
			wantChanged: true,
			wantVal:     `minimum-gas-prices = "25000000abadge"`,
		},
		{
			name:        "preserves foreign denoms alongside ubadge",
			in:          "minimum-gas-prices = \"10ubadge,0.01ibc/E1116484B327AEE59CDC3DA73D319834781A13DB2A7DFC1F38A30CD45ABF58B8\"\n",
			wantChanged: true,
			wantVal:     `minimum-gas-prices = "10000000000abadge,0.01ibc/E1116484B327AEE59CDC3DA73D319834781A13DB2A7DFC1F38A30CD45ABF58B8"`,
		},
		{
			name:        "already abadge is a no-op",
			in:          "minimum-gas-prices = \"10000000000abadge\"\n",
			wantChanged: false,
			wantVal:     `minimum-gas-prices = "10000000000abadge"`,
		},
		{
			name:        "empty value is a no-op",
			in:          "minimum-gas-prices = \"\"\n",
			wantChanged: false,
			wantVal:     `minimum-gas-prices = ""`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			p := filepath.Join(dir, "app.toml")
			if err := os.WriteFile(p, []byte(tc.in), 0o644); err != nil {
				t.Fatal(err)
			}

			changed, err := migrateAppTomlMinGasPrices(p)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}

			gotBytes, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			got := string(gotBytes)
			if !strings.Contains(got, tc.wantVal) {
				t.Fatalf("got:\n%s\nwant substring:\n%s", got, tc.wantVal)
			}

			// Re-running must always be a no-op.
			changed2, err := migrateAppTomlMinGasPrices(p)
			if err != nil {
				t.Fatalf("second run error: %v", err)
			}
			if changed2 {
				t.Fatalf("second run unexpectedly modified file")
			}
		})
	}
}
