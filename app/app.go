package app

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"httpframwork/modules/constant"
	"httpframwork/modules/container"
	"httpframwork/modules/errorcache"
)

type Application struct {
	Config      *viper.Viper
	Container   *container.Container
	Domain      string
	Environment string
	Log         *logrus.Logger
	Routes      []*AppRoutes
}

// Create new application instance
func New() (app Application, err error) {

	app = Application{}

	if err = app.initConfig(); err != nil {
		return
	}

	if err = app.initContainer(); err != nil {
		return
	}

	app.initEnvironment()

	if err = app.initLogger(); err != nil {
		return
	}

	return app, nil
}

// Initializes application configuration
func (a *Application) initConfig() (err error) {
	path := os.Getenv(constant.EnvConfigPath)

	a.Config = viper.New()
	a.Config.SetConfigName("config")
	a.Config.SetConfigType("yaml")

	log.Println("Searching for application configuration file...")
	if path == "" {
		a.Config.AddConfigPath(constant.DefaultConfigPath) // look for config in the working directory
		log.Println("Loading configs from default location...")
	} else {
		a.Config.AddConfigPath(path)
		log.Printf("Loading configs from location: %s\n", path)
	}

	err = a.Config.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		return
	}
	return
}

// Init container
func (a *Application) initContainer() (err error) {

	// Register all the required services
	global := container.New().
		Register(errorcache.GetRegistry())

	cont := global.Duplicate()

	// Call service bootstrap methods
	if err = errorcache.PopulateErrorCodes(cont); err != nil {
		return
	}

	a.Container = cont

	return
}

// initEnvironment
func (a *Application) initEnvironment() {
	a.Domain = a.Config.GetString(constant.AppDomain)
	a.Environment = a.Config.GetString(constant.EnvConfigPath)
	return
}

// initLogger - Initializes default application logger
func (a *Application) initLogger() (err error) {
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.TraceLevel
	log.Formatter = new(logrus.JSONFormatter)

	a.Log = log

	return
}

// Start the application
func (a *Application) Run() (err error) {
	var (
		cert            tls.Certificate
		TLSCert, TLSKey string
		pURL            *url.URL
	)

	router := a.prepareRoutes()

	if pURL, err = url.ParseRequestURI(a.Domain); err != nil {
		return
	}

	ln, err := net.Listen("tcp4", ":"+pURL.Port())

	if err != nil {
		return fmt.Errorf("failed to start listener with error `%v`", err)
	}

	// if SSL enabled then init TLS

	if true == a.Config.GetBool(constant.SSLEnabled) {
		TLSCert = a.Config.GetString(constant.SSLCertFilePath)
		TLSKey = a.Config.GetString(constant.SSLKeyPath)

		if TLSCert != "" || TLSKey != "" {
			if cert, err = tls.LoadX509KeyPair(TLSCert, TLSKey); err != nil {
				return
			}
			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
			// update old listener with the upgraded tls one
			ln = tls.NewListener(ln, tlsConfig)
			log.Println("Switched to TLS")
		} else {
			return errors.New("ssl enabled but certificate missing")
		}
	}

	return http.Serve(ln, router)
}
