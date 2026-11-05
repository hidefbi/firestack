package dnsx

import (
	"testing"

	"github.com/celzero/firestack/intra/settings"
)

// UnproxiablePid must never return the base proxy under global lockdown.
func TestUnproxiablePid(t *testing.T) {
	defer settings.GlobalLockdown.Store(false)
	defer settings.Loopingback.Store(false)

	tests := []struct {
		name     string
		lockdown bool
		loop     bool
		want     string
	}{
		{"default", false, false, NetBaseProxy},
		{"loopback", false, true, NetExitProxy},
		{"lockdown", true, false, NetNoProxy},
		{"lockdown+loopback", true, true, NetNoProxy},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			settings.GlobalLockdown.Store(tc.lockdown)
			settings.Loopingback.Store(tc.loop)
			if got := UnproxiablePid(); got != tc.want {
				t.Fatalf("UnproxiablePid() = %q, want %q", got, tc.want)
			}
		})
	}
}
