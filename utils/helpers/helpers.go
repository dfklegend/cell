package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"path/filepath"
	"reflect"	
	"strings"
	"testing"
	"time"
)

// GetFreePort returns a free port
func GetFreePort(t testing.TB) int {
	t.Helper()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// GetMapKeys returns a string slice with the map keys
func GetMapKeys(t *testing.T, m interface{}) []string {
	if reflect.ValueOf(m).Kind() != reflect.Map {
		t.Fatal(errors.New("GetMapKeys should receive a map"))
	}
	if reflect.TypeOf(m).Key() != reflect.TypeOf("bla") {
		t.Fatal(errors.New("GetMapKeys should receive a map with string keys"))
	}
	t.Helper()
	res := make([]string, 0)
	for _, k := range reflect.ValueOf(m).MapKeys() {
		res = append(res, k.String())
	}
	return res
}

// WriteFile test helper
func WriteFile(t *testing.T, filepath string, bytes []byte) {
	t.Helper()
	if err := ioutil.WriteFile(filepath, bytes, 0644); err != nil {
		t.Fatalf("failed writing file: %s", err)
	}
}

// ReadFile test helper
func ReadFile(t *testing.T, filepath string) []byte {
	t.Helper()
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed reading file: %s", err)
	}
	return b
}

// StartProcess starts a process
func StartProcess(t testing.TB, program string, args ...string) *exec.Cmd {
	t.Helper()
	return exec.Command(program, args...)
}

func waitForServerToBeReady(t testing.TB, out *bufio.Reader) {
	t.Helper()
	ShouldEventuallyReturn(t, func() bool {
		line, _, err := out.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		return strings.Contains(string(line), "all modules started!")
	}, true, 100*time.Millisecond, 30*time.Second)
}

// FixtureGoldenFileName returns the golden file name on fixtures path
func FixtureGoldenFileName(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("fixtures", name+".golden")
}

func vetExtras(extras []interface{}) (bool, string) {
	for i, extra := range extras {
		if extra != nil {
			zeroValue := reflect.Zero(reflect.TypeOf(extra)).Interface()
			if !reflect.DeepEqual(zeroValue, extra) {
				message := fmt.Sprintf("unexpected non-nil/non-zero extra argument at index %d:\n\t<%T>: %#v", i+1, extra, extra)
				return false, message
			}
		}
	}
	return true, ""
}

func pollFuncReturn(f interface{}) (interface{}, error) {
	values := reflect.ValueOf(f).Call([]reflect.Value{})

	extras := []interface{}{}
	for _, value := range values[1:] {
		extras = append(extras, value.Interface())
	}

	success, message := vetExtras(extras)

	if !success {
		return nil, errors.New(message)
	}

	return values[0].Interface(), nil
}

// ShouldEventuallyReceive should asserts that eventually channel c receives a value
func ShouldEventuallyReceive(t testing.TB, c interface{}, timeouts ...time.Duration) interface{} {
	t.Helper()
	if !isChan(c) {
		t.Fatal("ShouldEventuallyReceive c argument should be a channel")
	}
	v := reflect.ValueOf(c)

	timeout := time.After(500 * time.Millisecond)

	if len(timeouts) > 0 {
		timeout = time.After(timeouts[0])
	}

	recvChan := make(chan reflect.Value)

	go func() {
		v, ok := v.Recv()
		if ok {
			recvChan <- v
		}
	}()

	select {
	case <-timeout:
		t.Fatal(errors.New("timed out waiting for channel to receive"))
	case a := <-recvChan:
		return a.Interface()
	}

	return nil
}

// ShouldAlwaysReturn asserts that the return of f should always be v, timeouts: 0 - evaluation interval, 1 - timeout
func ShouldAlwaysReturn(t testing.TB, f interface{}, v interface{}, timeouts ...time.Duration) {
	t.Helper()
	interval := 10 * time.Millisecond
	timeout := time.After(50 * time.Millisecond)
	switch len(timeouts) {
	case 1:
		interval = timeouts[0]
		break
	case 2:
		interval = timeouts[0]
		timeout = time.After(timeouts[1])
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if isFunction(f) {
		for {
			select {
			case <-timeout:
				return
			case <-ticker.C:
				val, err := pollFuncReturn(f)
				if err != nil {
					t.Fatal(err)
				}
				if v != val {
					t.Fatalf("function f returned wrong value %s", val)
				}
			}
		}
	} else {
		t.Fatal("ShouldAlwaysReturn should receive a function with no args and more than 0 outs")
		return
	}
}

// ShouldEventuallyReturn asserts that eventually the return of f should be v, timeouts: 0 - evaluation interval, 1 - timeout
func ShouldEventuallyReturn(t testing.TB, f interface{}, v interface{}, timeouts ...time.Duration) {
	t.Helper()
	interval := 10 * time.Millisecond
	timeout := time.After(500 * time.Millisecond)
	switch len(timeouts) {
	case 1:
		interval = timeouts[0]
		break
	case 2:
		interval = timeouts[0]
		timeout = time.After(timeouts[1])
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if isFunction(f) {
		for {
			select {
			case <-timeout:
				t.Fatalf("function f never returned value %s", v)
			case <-ticker.C:
				val, err := pollFuncReturn(f)
				if err != nil {
					t.Fatal(err)
				}
				if v == val {
					return
				}
			}
		}
	} else {
		t.Fatal("ShouldEventuallyEqual should receive a function with no args and more than 0 outs")
		return
	}
}
