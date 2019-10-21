package lights

import (
	"fmt"
	"net/http"

	"github.com/changkun/gobase/net"
)

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

// Turn turns the lights on or off.
func (l *Light) Turn(on bool) (bool, error) {
	var action string
	if on {
		action = "{\"on\": true}"
	} else {
		action = "{\"on\": false}"
	}

	addr := fmt.Sprintf(apiLightState, l.Bridge.Hostname, l.Bridge.Username, l.ID)
	err := net.HTTPRequest(addr,
		http.MethodPut, []byte(action), &net.RequestParams{Timeout: 100}, &struct{}{})
	if err != nil {
		return false, fmt.Errorf("hue: turn off lights went wrong, message: %v", err)
	}
	return true, nil
}
