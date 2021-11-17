package discovery

import (
	"time"
	"github.com/dfklegend/cell/utils/config"
)

// EtcdServiceDiscoveryConfig Etcd service discovery config
type EtcdServiceDiscoveryConfig struct {
	Endpoints   []string
	User        string
	Pass        string
	DialTimeout time.Duration
	Prefix      string
	Heartbeat   struct {
		TTL time.Duration
		Log bool
	}
	SyncServers struct {
		Interval    time.Duration
		Parallelism int
	}
	Revoke struct {
		Timeout time.Duration
	}
	GrantLease struct {
		Timeout       time.Duration
		MaxRetries    int
		RetryInterval time.Duration
	}
	Shutdown struct {
		Delay time.Duration
	}
	ServerTypesBlacklist []string
}

// NewDefaultEtcdServiceDiscoveryConfig Etcd service discovery default config
func NewDefaultEtcdServiceDiscoveryConfig() *EtcdServiceDiscoveryConfig {
	return &EtcdServiceDiscoveryConfig{
		Endpoints:   []string{"localhost:2379"},
		User:        "",
		Pass:        "",
		DialTimeout: time.Duration(5 * time.Second),
		Prefix:      "cell/",
		Heartbeat: struct {
			TTL time.Duration
			Log bool
		}{
			TTL: time.Duration(60 * time.Second),
			Log: false,
		},
		SyncServers: struct {
			Interval    time.Duration
			Parallelism int
		}{
			Interval:    time.Duration(120 * time.Second),
			Parallelism: 10,
		},
		Revoke: struct {
			Timeout time.Duration
		}{
			Timeout: time.Duration(5 * time.Second),
		},
		GrantLease: struct {
			Timeout       time.Duration
			MaxRetries    int
			RetryInterval time.Duration
		}{
			Timeout:       time.Duration(60 * time.Second),
			MaxRetries:    15,
			RetryInterval: time.Duration(5 * time.Second),
		},
		Shutdown: struct {
			Delay time.Duration
		}{
			Delay: time.Duration(300 * time.Millisecond),
		},
		ServerTypesBlacklist: nil,
	}
}

// NewEtcdServiceDiscoveryConfig Etcd service discovery config with default config paths
func NewEtcdServiceDiscoveryConfig(cfgFile *config.Config) *EtcdServiceDiscoveryConfig {
	conf := NewDefaultEtcdServiceDiscoveryConfig()
	// override
	if err := cfgFile.UnmarshalKey("etcd", &conf); err != nil {
		panic(err)
	}
	//
	return conf
}
