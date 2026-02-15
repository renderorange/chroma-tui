package tui

import "testing"

func TestPanelDimensions(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		wantLeftWidth  int
		wantRightWidth int
		wantListHeight int
	}{
		{
			name:           "standard 80x40 terminal",
			width:          80,
			height:         40,
			wantLeftWidth:  17,
			wantRightWidth: 53,
			wantListHeight: 34,
		},
		{
			name:           "wide 120x40 terminal",
			width:          120,
			height:         40,
			wantLeftWidth:  27,
			wantRightWidth: 83,
			wantListHeight: 34,
		},
		{
			name:           "minimum size 60x20",
			width:          60,
			height:         20,
			wantLeftWidth:  12,
			wantRightWidth: 38,
			wantListHeight: 14,
		},
		{
			name:           "below minimum 50x15",
			width:          50,
			height:         15,
			wantLeftWidth:  10,
			wantRightWidth: 10,
			wantListHeight: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftWidth, rightWidth, listHeight := panelDimensions(tt.width, tt.height)

			if leftWidth != tt.wantLeftWidth {
				t.Errorf("leftWidth = %d, want %d", leftWidth, tt.wantLeftWidth)
			}
			if rightWidth != tt.wantRightWidth {
				t.Errorf("rightWidth = %d, want %d", rightWidth, tt.wantRightWidth)
			}
			if listHeight != tt.wantListHeight {
				t.Errorf("listHeight = %d, want %d", listHeight, tt.wantListHeight)
			}

			// Verify non-negative
			if leftWidth < 0 || rightWidth < 0 || listHeight < 0 {
				t.Errorf("got negative dimension: left=%d, right=%d, height=%d",
					leftWidth, rightWidth, listHeight)
			}
		})
	}
}
