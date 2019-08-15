package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"httpframwork/modules/constant"
	"httpframwork/modules/container"
	"httpframwork/modules/logger"
)

type Api struct {
	Name      string
	Container *container.Container
	Config    *viper.Viper
	RawBody   interface{}
	Request   *http.Request
	Response  http.ResponseWriter
	Vars      map[string]string
	Log       *logs.Log
	Status    int
}

func (api *Api) GetHandler(name string, path string, method string, handler func()) (string, string, string, http.HandlerFunc) {
	return name, path, method, func(writer http.ResponseWriter, request *http.Request) {
		api.Name = name
		api.Request = request
		api.Response = writer
		api.Init()
		handler()
		api.Defer()
	}
}

// Calls other init method to initialize resources
func (api *Api) Init() {
	api.Vars = mux.Vars(api.Request)
	api.initLogger()
}

// initialize request level logging
func (api *Api) initLogger() {
	// obtain log folder form config
	lPath := api.Config.GetString(constant.AppLogFolder)
	api.Log = logs.New(lPath)

	api.Log.Print("Start ", time.Now().UTC().Format(constant.DefaultDateTimeFormat))
	api.Log.Print("IP ", api.GetClientIP())
	api.Log.Print("Resource ", api.Name)
	api.Log.Print("Method ", api.Request.Method)
	api.Log.Print("URL ", api.Request.URL.String())
	api.Log.Print("Param ", api.Vars)
	api.Log.Print("Request ", strings.Replace(string(api.Body()), "\n", "", -1))
}

// ResponseJSON
func (api *Api) ResponseJSON(res interface{}) {
	b, _ := json.Marshal(res)
	api.Response.WriteHeader(http.StatusOK)
	api.Response.Header().Add("Content-Type", "application/json")
	api.Response.Write(b)
	return
}

// ResponseTest
func (api *Api) ResponseText(res interface{}) {
	api.Response.WriteHeader(http.StatusOK)
	api.Response.Header().Add("Content-Type", "application/text")
	api.Response.Write([]byte(res.(string)))
	return
}

// Redirect - redirects to the specified URL
func (api *Api) Redirect(url string) {
	http.Redirect(api.Response, api.Request, url, http.StatusMovedPermanently)
}

// ResponseHeader
func (api *Api) ResponseHeaders(headers map[string]string) {
	for h, v := range headers {
		api.Response.Header().Set(h, v)
	}
}

// GetClientIP
func (api *Api) GetClientIP() (string) {
	req := api.Request
	ipAddress := req.RemoteAddr

	if ip := req.Header.Get("X-Forwarded-For"); "" != ip {
		ipAddress = ip

		// X-Forwarded-For might contain multiple IPs. Get the last one.
		if strings.Contains(ipAddress, ",") {
			ips := strings.Split(ipAddress, ",")
			ipAddress = strings.Trim(ips[len(ips)-1], " ")
		}
	}

	var (
		ip  net.IP
		err error
	)

	if -1 != strings.Index(ipAddress, ":") {
		if ipAddress, _, err = net.SplitHostPort(ipAddress); nil != err {
			return ""
		}
	}

	if err := ip.UnmarshalText([]byte(ipAddress)); nil != err {
		return ""
	}

	return ipAddress
}

// Body returns the body from the request
func (api *Api) Body() []byte {
	body, _ := ioutil.ReadAll(api.Request.Body)
	// Restore the io.ReadCloser to its original state
	api.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// Handles primary response
func (api *Api) Defer() {
	var b bytes.Buffer

	if api.Status == 0 {
		api.Log.Print("Status", http.StatusInternalServerError)
		api.Log.Print("Stack trace", string(debug.Stack()))
	} else {
		api.Log.Print("Status", api.Status)
	}

	if api.RawBody != nil {
		if fmt.Sprint(api.RawBody) == "[]" {
			emptyResponse, _ := json.Marshal(make([]int64, 0))
			api.Log.Print("Response", string(emptyResponse))
		} else {
			enc := json.NewEncoder(&b)
			enc.SetEscapeHTML(false)
			enc.Encode(api.RawBody)
			api.Log.Print("Response", strings.Replace(string(b.Bytes()), "\n", "", -1))
		}
	}

	api.Log.Print("End", time.Now().UTC().Format(constant.DefaultDateTimeFormat))
	api.Log.Dump()
}
