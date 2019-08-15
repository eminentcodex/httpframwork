package api

import (
	"net/http"

	"github.com/spf13/viper"
	"httpframwork/modules/container"
)

type Heartbeat struct {
	Api
}

// Registers handle function with the router
func RegisterHeartbeat(cont *container.Container, conf *viper.Viper) (string, string, string, http.HandlerFunc) {
	h := &Heartbeat{
		Api{
			Container: cont,
			Config:    conf,
		},
	}

	return h.GetHandler("heartbeat", "/heartbeat", http.MethodGet, h.handler)
}

// Perform the logic here
func (h *Heartbeat) handler() {
	res := map[string]interface{}{
		"Status":  1,
		"Message": "success",
	}

	h.Status = http.StatusOK
	h.ResponseJSON(res)
}
