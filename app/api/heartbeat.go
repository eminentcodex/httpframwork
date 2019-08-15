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

	return "heartbeat", "/heartbeat", http.MethodGet, func(writer http.ResponseWriter, request *http.Request) {
		h.Request = request
		h.Response = writer
		h.handler()
	}
}

// Perform the logic here
func (h *Heartbeat) handler() {
	res := map[string]interface{}{
		"Status":  1,
		"Message": "success",
	}

	h.ResponseJSON(res)
}
