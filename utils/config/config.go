package config

import (
    "time"

    "github.com/spf13/viper"
)


type Config struct {
    config *viper.Viper
}

// 公共配置
func LoadConfig(path string) *Config {
    v := viper.New()
    v.SetConfigName("config")
    v.AddConfigPath(path)
    v.SetConfigType("yaml")
    v.ReadInConfig()    
    return &Config{
        config: v,
    }
}

// GetDuration returns a duration from the inner config
func (c *Config) GetDuration(s string, def ...time.Duration) time.Duration {
    if !c.config.IsSet(s) {
        if len(def) != 0 {
            return def[0]
        }        
    }
    return c.config.GetDuration(s)
}

// GetString returns a string from the inner config
func (c *Config) GetString(s string, def ...string) string {
    if !c.config.IsSet(s) {
        if len(def) != 0 {
            return def[0]
        }        
    }
    return c.config.GetString(s)
}

// GetInt returns an int from the inner config
func (c *Config) GetInt(s string, def ...int) int {
    if !c.config.IsSet(s) {
        if len(def) != 0 {
            return def[0]
        }        
    }
    return c.config.GetInt(s)
}

// GetBool returns an boolean from the inner config
func (c *Config) GetBool(s string, def ...bool) bool {
    if !c.config.IsSet(s) {
        if len(def) != 0 {
            return def[0]
        }        
    }
    return c.config.GetBool(s)
}

// GetStringSlice returns a string slice from the inner config
func (c *Config) GetStringSlice(s string) []string {
    return c.config.GetStringSlice(s)
}

// Get returns an interface from the inner config
func (c *Config) Get(s string, def ...interface{}) interface{} {
    if !c.config.IsSet(s) {
        if len(def) != 0 {
            return def[0]
        }        
    }
    return c.config.Get(s)
}

// GetStringMapString returns a string map string from the inner config
func (c *Config) GetStringMapString(s string) map[string]string {
    return c.config.GetStringMapString(s)
}

// UnmarshalKey unmarshals key into v
func (c *Config) UnmarshalKey(s string, v interface{}) error {
    return c.config.UnmarshalKey(s, v)
}

// Unmarshal unmarshals config into v
func (c *Config) Unmarshal(v interface{}) error {
    return c.config.Unmarshal(v)
}



