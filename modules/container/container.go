package container

import "sync"

type Container struct {
	sync.Mutex
	bag map[string]interface{}
}

type GlobalContainer struct {
	Container
}

type Registry struct {
	Key   string
	Value interface{}
}

type Registries []Registry

var globalContainerInstance *GlobalContainer

// New function - returns new instance of global container.
func New() *GlobalContainer {
	globalContainerInstance = &GlobalContainer{Container{sync.Mutex{}, make(map[string]interface{})}}
	return globalContainerInstance.getDefaults()
}

// GetInstance function - returns current instance of global container.
func GetInstance() *GlobalContainer {
	return globalContainerInstance
}

// Register method - triggers init for components for global container...
func (me *GlobalContainer) Register(r Registries) *GlobalContainer {
	for _, v := range r {
		globalContainerInstance.Container.bag[v.Key] = v.Value
	}

	return me
}

// Register method - triggers init for components for container...
func (me *Container) Register(r Registries) *Container {
	for _, v := range r {
		me.bag[v.Key] = v.Value
	}

	return me
}

// Duplicate method - returns duplicate instance of container, while reinitializing required components passed.
func (me *Container) Duplicate(r ...Registries) *Container {
	instance := Container{sync.Mutex{}, make(map[string]interface{})}

	for k, v := range globalContainerInstance.Container.bag {
		instance.bag[k] = v
	}

	if len(r) > 0 {
		for _, v := range r {
			instance.Register(v)
		}
	}

	return &instance
}

// StoreToGlobal method - store a value to global container ...
func (me *Container) StoreToGlobal(key string, value interface{}) {
	globalContainerInstance.Container.Lock()
	globalContainerInstance.Container.bag[key] = value
	globalContainerInstance.Container.Unlock()
	me.Set(key, value)
}

// StoreToContainer ...
func (me *Container) Set(key string, value interface{}) *Container {
	me.Lock()
	me.bag[key] = value
	me.Unlock()
	return me
}

// RetrieveFromContainer ...
func (me *Container) Get(key string) (value interface{}) {
	me.Lock()
	value = me.bag[key]
	me.Unlock()
	return value
}

// GetAll method - get all bag values.
func (me *Container) GetAll() map[string]interface{} {
	return me.bag
}

// ReSync method - sync values from global.
func (me *Container) ReSync() *Container {
	for k, v := range globalContainerInstance.Container.bag {
		me.bag[k] = v
	}

	return me
}

// getDefaults method - register default global container values.
func (me *GlobalContainer) getDefaults() *GlobalContainer {
	return me
}