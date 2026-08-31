package lights

import (
	"context"
	"fmt"
	"net/http"

	"changkun.de/x/pkg/net"
)

// requestTimeout bounds a single call to the bridge, in seconds.
const requestTimeout = 100

// Light represents an individual light component
type Light struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name"`
	State     State   `json:"state,omitempty"`
	Type      string  `json:"type,omitempty"`
	ModelID   string  `json:"modelid,omitempty"`
	SWVersion string  `json:"swversion,omitempty"`
	Bridge    *Bridge `json:"-"`
}

// State represents all states of a light light
type State struct {
	On             bool      `json:"on"`
	Hue            uint16    `json:"hue,omitempty"`
	Effect         string    `json:"effect,omitempty"`
	Bri            uint8     `json:"bri,omitempty"`
	Sat            uint8     `json:"sat,omitempty"`
	CT             uint16    `json:"ct,omitempty"`
	XY             []float32 `json:"xy,omitempty"`
	Alert          string    `json:"alert,omitempty"`
	TransitionTime uint16    `json:"transitiontime,omitempty"`
	Reachable      bool      `json:"reachable,omitempty"`
	ColorMode      string    `json:"colormode,omitempty"`
}

// Turn turns the light on or off. Cancelling ctx aborts the request to the
// bridge.
func (l *Light) Turn(ctx context.Context, on bool) (bool, error) {
	action := `{"on": false}`
	if on {
		action = `{"on": true}`
	}

	addr := fmt.Sprintf(apiLightState, l.Bridge.Hostname, l.Bridge.Username, l.ID)
	err := net.HTTPRequest(ctx, addr,
		http.MethodPut, []byte(action), &net.RequestParams{Timeout: requestTimeout}, &struct{}{})
	if err != nil {
		return false, fmt.Errorf("hue: turn lights went wrong, message: %w", err)
	}
	return true, nil
}
