package templates

import (
	"testing"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   string
		want string
	}{
		{
			name: "UTC",
			tm:   "2021-02-24T13:35:50.028603852Z",
			want: "2021-02-24 at 13:35:50",
		},
		{
			name: "Empty",
			tm:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
