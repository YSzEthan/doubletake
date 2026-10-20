package airplay

import (
	"strings"
	"testing"
)

const xrandrDualHead = `Screen 0: minimum 320 x 200, current 3840 x 1080, maximum 16384 x 16384
eDP-1 connected (normal left inverted right x axis y axis)
   1920x1080     60.05 +  60.05
DP-1 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 0mm x 0mm
   1280x720R     59.86 +
HDMI-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 530mm x 300mm
   1920x1080     60.00*+  74.97    59.94
DP-2 disconnected (normal left inverted right x axis y axis)
`

func TestSelectXrandrMonitor(t *testing.T) {
	tests := []struct {
		name    string
		xrandr  string
		output  string
		want    [4]int
		wantErr string
	}{
		{name: "primary by default", xrandr: xrandrDualHead, want: [4]int{0, 0, 1920, 1080}},
		{
			name: "first active when no primary",
			xrandr: "eDP-1 connected (normal left inverted right x axis y axis)\n" +
				"DP-1 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 0mm x 0mm\n",
			want: [4]int{1920, 0, 3840, 1080},
		},
		{name: "scaled virtual output", xrandr: xrandrDualHead, output: "DP-1", want: [4]int{1920, 0, 3840, 1080}},
		{
			name:    "connected but inactive",
			xrandr:  xrandrDualHead,
			output:  "eDP-1",
			wantErr: `x11 output "eDP-1" not found or not active (active: DP-1, HDMI-1)`,
		},
		{
			name:    "disconnected",
			xrandr:  xrandrDualHead,
			output:  "DP-2",
			wantErr: `x11 output "DP-2" not found or not active`,
		},
		{
			name:    "unknown",
			xrandr:  xrandrDualHead,
			output:  "DP-9",
			wantErr: `x11 output "DP-9" not found or not active`,
		},
		{
			name: "exact name only",
			xrandr: "DP-10 connected 1920x1080+0+0 (normal) 0mm x 0mm\n" +
				"DP-1 connected 1280x720+1920+0 (normal) 0mm x 0mm\n",
			output: "DP-1",
			want:   [4]int{1920, 0, 3200, 720},
		},
		{
			name:   "rotated output reports rotated geometry",
			xrandr: "DP-1 connected 1080x1920+1920+0 left (normal left inverted right x axis y axis) 0mm x 0mm\n",
			output: "DP-1",
			want:   [4]int{1920, 0, 3000, 1920},
		},
		{
			name:   "narrow output",
			xrandr: "DP-1 connected 320x240+1920+0 (normal) 0mm x 0mm\n",
			output: "DP-1",
			want:   [4]int{1920, 0, 2240, 240},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x0, y0, x1, y1, err := selectXrandrMonitor(tt.xrandr, tt.output)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := [4]int{x0, y0, x1, y1}; got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
