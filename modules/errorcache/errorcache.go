package errorcache

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"httpframwork/modules/constant"
	"httpframwork/modules/container"
)

type (
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"msg"`
	}

	ErrorsConfig map[string]Error

	errorsCache struct {
		sync.Mutex
		bag ErrorsConfig
	}
)

const (
	InstanceKey = "Errors"
)

// GetRegistry function ...
func GetRegistry() container.Registries {
	obj := &errorsCache{
		bag: make(ErrorsConfig),
	}

	return container.Registries{
		container.Registry{
			Key:   InstanceKey,
			Value: obj,
		},
	}
}

// GetInstance function ...
func GetInstance(c *container.Container) *errorsCache {
	return c.Get(InstanceKey).(*errorsCache)
}

// PopulateErrorCodes function ...
func PopulateErrorCodes(c *container.Container) (err error) {
	ins := GetInstance(c)
	if err = ins.RetrieveErrors(); err != nil {
		return
	}
	c.StoreToGlobal(InstanceKey, ins)

	return
}

// GetError method - gets the stored error object on Container
func (me *errorsCache) GetError(errorCode string) Error {

	if err, ok := me.bag[errorCode]; ok {
		return err
	}

	return Error{}
}

// RetrieveErrors method ... Error storage
func (me *errorsCache) RetrieveErrors() (err error) {

	var (
		conf      *viper.Viper
		path      = os.Getenv("CONFIG_PATH")
		errorLang = os.Getenv("ERROR_LANG")
	)

	conf = viper.New()
	conf.SetConfigName("errors")
	conf.SetConfigType("yaml")

	log.Println("Searching for error configuration file...")
	if path == "" {
		conf.AddConfigPath(constant.DefaultConfigPath) // look for config in the working directory
		log.Println("Loading error configs from default location...")
	} else {
		conf.AddConfigPath(path)
		log.Printf("Loading error configs from location: %s\n", path)
	}

	err = conf.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		return
	}

	errs := conf.Get(errorLang)
	if errs == nil {
		return errors.New("no errors defined")
	}

	newBag := make(ErrorsConfig)

	for identifier, e := range errs.(map[string]interface{}) {

		newBag[identifier] = Error{
			Status:  e.(map[string]interface{})["status"].(int),
			Message: e.(map[string]interface{})["msg"].(string),
		}
	}

	me.bag = newBag

	return
}
