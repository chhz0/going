package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigV is based on spf13/viper
type ConfigV struct {
	v   *viper.Viper
	ptr any // ptr must be a pointer.
}

func NewConfigV(ptr any) (*ConfigV, error) {
	if ptr == nil {
		return nil, fmt.Errorf("[NewConfigV] the param:ptr cannot be nil.")
	}

	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("[NewConfigV] the param:ptr must be a pointer, got %T.", ptr)
	}

	return &ConfigV{
		v:   viper.New(),
		ptr: ptr,
	}, nil
}

func (c *ConfigV) Load(configPath, configName, configType string) error {
	if configPath != "" {
		c.v.AddConfigPath(configPath)
	}

	c.v.SetConfigName(configName)
	c.v.SetConfigType(configType)

	if err := c.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("[ConfigV.Load] config file '%s.%s' not found in search paths:%s;",
				configName, configType, configPath)
		} else {
			return fmt.Errorf("[ConfigV.Load] failed to read config file: %w", err)
		}
	}

	if err := c.v.Unmarshal(c.ptr); err != nil {
		return fmt.Errorf("[ConfigV.Load] failed to unmarshal config to struct: %w", err)
	}

	return nil
}

func (c *ConfigV) Watch() {
	go func() {
		c.v.WatchConfig()
		c.v.OnConfigChange(func(in fsnotify.Event) {
			if err := c.v.ReadInConfig(); err != nil {
				log.Printf("[ConfigV.Watch] error re-reading config after change: %v", err)
				return
			}

			if err := c.v.Unmarshal(c.ptr); err != nil {
				log.Printf("[ConfigV.Watch] error unmarshaling new config: %v", err)
				return
			}
		})
	}()
}

func (c *ConfigV) Viper() *viper.Viper {
	return c.v
}

// ConfigVSafe is concurrency-safe. todo!
// type ConfigVSafe struct {
// 	rw sync.RWMutex
// 	v  *viper.Viper

// 	to any
// }
