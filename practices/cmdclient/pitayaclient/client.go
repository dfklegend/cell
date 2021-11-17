package pitayaclient

import (
	"time"

	"github.com/topfreegames/pitaya/client2"
	"github.com/sirupsen/logrus"
)

var (
	theClient *client2.Client
)

func GetClient() *client2.Client {
	return theClient;
}

func Start() {
	theClient = client2.New(logrus.InfoLevel, 100*time.Millisecond)	
}
