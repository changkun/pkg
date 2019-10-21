package lights

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/changkun/gobase/net"
)

// Bridge implements a hue light bridge
type Bridge struct {
	Hostname string
	Username string
}

// NewBridge creates a new hue bridge by given hostname and username.
func NewBridge(hostname, username string) *Bridge {
	return &Bridge{hostname, username}
}

// GetLights returns all lights in the given bridge
func (l *Bridge) GetLights() ([]Light, error) {
	var ll map[string]Light

	err := net.HTTPRequest(fmt.Sprintf(apiLights, l.Hostname, l.Username),
		http.MethodGet, nil, &net.RequestParams{Timeout: 100}, &ll)
	if err != nil {
		return nil, fmt.Errorf("hue: get lights went wrong, message: %v", err)
	}

	var ret []Light

	for id, v := range ll {
		iid, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("hue: convert light id error, message: %v", err)
		}
		v.ID = iid
		v.Bridge = l
		ret = append(ret, v)
	}

	return ret, nil
}
