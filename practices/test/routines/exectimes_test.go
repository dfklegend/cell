package routines

import (
	"fmt"
	"testing"
)

func Test_Simple(t *testing.T) {
	fmt.Println("--------")
	RunCounter(32)
	fmt.Printf("total:%d", TotalRunTimes)
	fmt.Println("--------")
}
