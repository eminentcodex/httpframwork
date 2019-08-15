package api

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
	"httpframwork/modules/container"
)

type Api struct {
	Container *container.Container
	Config    *viper.Viper
	Output    string
	Request   *http.Request
	Response  http.ResponseWriter
}

func (api *Api) ResponseJSON(res interface{}) {
	b, _ := json.Marshal(res)
	api.Response.WriteHeader(http.StatusOK)
	api.Response.Header().Add("Content-Type", "application/json")
	api.Response.Write(b)
	return
}

func (api *Api) ResponseText(res interface{}) {
	api.Response.WriteHeader(http.StatusOK)
	api.Response.Header().Add("Content-Type", "application/text")
	api.Response.Write([]byte(res.(string)))
	return
}

func (api *Api) Redirect(url string) {
	http.Redirect(api.Response, api.Request, url, http.StatusMovedPermanently)
}

func (api *Api) ResponseHeaders() {

}
