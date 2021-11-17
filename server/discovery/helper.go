package discovery

import (
	"testing"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/integration"
)

// GetTestEtcd gets a test in memory etcd server
func GetTestEtcd(t *testing.T) (*integration.ClusterV3, *clientv3.Client) {
	t.Helper()
	c := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	cli := c.RandClient()
	return c, cli
}
