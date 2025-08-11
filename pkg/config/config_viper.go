package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigV is based on spf13/viper.
type ConfigV struct {
	v *viper.Viper

	// configPath      []string
	// configName      string
	// configType      string
	configUnmarshal any    // ptr must be a pointer.
	onChange        func() // optional callback for config change.
}

// NewConfigV creates a new ConfigV instance.
func NewConfigV(ptr any) (*ConfigV, error) {
	if ptr == nil {
		return nil, fmt.Errorf("[NewConfigV] the param:ptr cannot be nil.")
	}

	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("[NewConfigV] the param:ptr must be a pointer, got %T.", ptr)
	}

	return &ConfigV{
		v:               viper.New(),
		configUnmarshal: ptr,
	}, nil
}

// SetOnChange sets the callback function when the configuration changes.
func (c *ConfigV) SetOnChange(fn func()) {
	c.onChange = fn
}

// Load loads the configuration from the specified path.
func (c *ConfigV) Load(configPath, configName, configType string) error {
	if configPath != "" {
		c.v.AddConfigPath(configPath)
	}

	c.v.SetConfigName(configName)
	c.v.SetConfigType(configType)

	if err := c.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("[ConfigV.Load] config file '%s.%s' not found in search paths:%s.",
				configName, configType, configPath)
		}
		return fmt.Errorf("[ConfigV.Load] failed to read config file: %w.", err)
	}

	if err := c.v.Unmarshal(c.configUnmarshal); err != nil {
		return fmt.Errorf("[ConfigV.Load] failed to unmarshal config to struct: %w.", err)
	}

	return nil
}

// Watch watches for changes to the config file and reloads it.
func (c *ConfigV) Watch() {
	c.v.WatchConfig()
	c.v.OnConfigChange(func(in fsnotify.Event) {
		if err := c.v.ReadInConfig(); err != nil {
			log.Printf("[ConfigV.Watch] error re-reading config after change: %v.", err)
			return
		}

		if err := c.v.Unmarshal(c.configUnmarshal); err != nil {
			log.Printf("[ConfigV.Watch] error unmarshaling new config: %v.", err)
			return
		}

		if c.onChange != nil {
			c.onChange()
		}
	})
}

// Viper returns the underlying viper instance.
func (c *ConfigV) Viper() *viper.Viper {
	return c.v
}

// // ConfigVSafe is concurrency-safe.
// type ConfigVSafe struct {
// 	rw sync.RWMutex
// 	v  *viper.Viper

// 	configUnmarshal any
// 	decoder         func(any) error
// }

// // NewConfigVSafe create a new thread-safe ConfigVSafe instance.
// func NewConfigVSafe(ptr any) (*ConfigVSafe, error) {
// 	if ptr == nil {
// 		return nil, fmt.Errorf("[NewConfigVSafe] the param:ptr cannot be nil.")
// 	}

// 	val := reflect.ValueOf(ptr)
// 	if val.Kind() != reflect.Pointer {
// 		return nil, fmt.Errorf("[NewConfigVSafe] the param:ptr must be a pointer, got %T.", ptr)
// 	}

// 	// create a zero value of the pointed-to type
// 	configType := reflect.TypeOf(ptr).Elem()
// 	configValue := reflect.New(configType).Elem().Interface()

// 	return &ConfigVSafe{
// 		v:               viper.New(),
// 		configUnmarshal: configValue,
// 		decoder:         func(dst any) error { return nil },
// 	}, nil
// }

// // SetDecoder sets a custom decoder function (e.g., json.Unmarshal, yaml.Unmarshal, toml.Unmarshal)
// func (c *ConfigVSafe) SetDecoder(decoder func(any) error) {
// 	c.decoder = decoder
// }

// // Load loads configuration from file
// func (c *ConfigVSafe) Load(configPath, configName, configType string) error {
// 	if configPath != "" {
// 		c.v.AddConfigPath(configPath)
// 	}
// 	c.v.SetConfigName(configName)
// 	c.v.SetConfigType(configType)

// 	if err := c.v.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			return fmt.Errorf("[ConfigVSafe.Load] config file '%s.%s' not found in search paths:%s.",
// 				configName, configType, configPath)
// 		}
// 		return fmt.Errorf("[ConfigVSafe.Load] failed to read config file: %w.", err)
// 	}

// 	rawData := c.v.AllSettings()
// 	jsonData, err := json.Marshal(rawData)
// 	if err != nil {
// 		return fmt.Errorf("[ConfigVSafe.Load] failed to marshal config to json: %w.", err)
// 	}

// 	// Create new config instance
// 	uType := reflect.TypeOf(c.configUnmarshal)
// 	newConfig := reflect.New(uType).Interface()

// 	if err := json.Unmarshal(jsonData, newConfig); err != nil {
// 		return fmt.Errorf("[ConfigVSafe.Load] failed to unmarshal config to struct: %w.", err)
// 	}

// 	c.rw.Lock()
// 	c.configUnmarshal = c.configUnmarshal
// 	c.rw.Unlock()

// 	return nil
// }

// // Watch watches for config changes and updates the config
// func (c *ConfigVSafe) Watch() {
// 	c.v.WatchConfig()
// 	c.v.OnConfigChange(func(in fsnotify.Event) {
// 		// Get raw config data
// 		rawData := c.v.AllSettings()
// 		jsonData, err := json.Marshal(rawData)
// 		if err != nil {
// 			log.Printf("[ConfigVSafe.Watch] failed to marshal config: %v", err)
// 			return
// 		}

// 		// Create new config instance
// 		configType := reflect.TypeOf(c.configUnmarshal)
// 		newConfig := reflect.New(configType).Interface()

// 		// Decode into new instance
// 		if err := json.Unmarshal(jsonData, newConfig); err != nil {
// 			log.Printf("[ConfigVSafe.Watch] failed to decode config: %v", err)
// 			return
// 		}

// 		// Atomically update the config
// 		c.rw.Lock()
// 		c.configUnmarshal = reflect.ValueOf(newConfig).Elem().Interface()
// 		c.rw.Unlock()
// 	})
// }

// // Get returns a deep copy of the current configuration
// func (c *ConfigVSafe) Get() (any, error) {
// 	c.rw.RLock()
// 	defer c.rw.RUnlock()

// 	// Marshal current config
// 	jsonData, err := json.Marshal(c.configUnmarshal)
// 	if err != nil {
// 		return nil, fmt.Errorf("[ConfigVSafe.Get] failed to marshal config: %w", err)
// 	}

// 	// Create new instance
// 	configType := reflect.TypeOf(c.configUnmarshal)
// 	newConfig := reflect.New(configType).Interface()

// 	// Unmarshal into new instance
// 	if err := json.Unmarshal(jsonData, newConfig); err != nil {
// 		return nil, fmt.Errorf("[ConfigVSafe.Get] failed to unmarshal config: %w", err)
// 	}

// 	return newConfig, nil
// }

// // View executes a function with read-locked access to the config.
// func (c *ConfigVSafe) View(fn func(any)) {
// 	c.rw.RLock()
// 	fn(c.configUnmarshal)
// 	c.rw.RUnlock()
// }

// // Viper returns the underlying viper instance.
// func (c *ConfigVSafe) Viper() *viper.Viper {
// 	return c.v
// }
