package routines

import (
	"fmt"
	"time"
)

var (
	TotalRunTimes int
)

type RunTimes struct {
	runTimes int
}

func (s *RunTimes) run(endChan chan int) {
	for true {
		s.runTimes++

		// wait end
		select {
		case <-endChan:
			return
		default:
		}
	}
}

func (s *RunTimes) print() {
	fmt.Printf("times:%d\n", s.runTimes)
}

func RunCounter(num int) {
	datas := make([]RunTimes, num)
	endChan := make(chan int)

	for i := 0; i < num; i++ {
		one := &datas[i]

		go one.run(endChan)
	}

	time.Sleep(8 * time.Second)
	close(endChan)

	var total int
	for i := 0; i < num; i++ {
		one := &datas[i]
		//one.print()
		total += one.runTimes
	}

	fmt.Printf("total:%d\n", total)
	TotalRunTimes = total
}

// func main() {
// 	test(32)
// }
