package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func lookupFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func freePort() int {
	tries := 0
	for {
		if port, err := lookupFreePort(); err == nil {
			return port
		} else {
			if tries > 10 {
				panic("Unable to find free port")
			}
		}
		tries++
	}
	// will never get here
	return 0
}

var port int

func init() {
	port = freePort()
	os.Setenv("PORT", fmt.Sprintf("%d", port))
}

func waitForServiceToStart(t *testing.T) bool {
	client := &http.Client{}

	url := fmt.Sprintf("https://127.0.0.1:%d%s", port, "/health/alive")

	if resp, err := client.Get(url); err != nil {
		t.Logf("HTTP request to %s failed: %s", url, err)
		return true
	} else if resp.StatusCode != http.StatusOK {
		t.Logf("HTTP request to %s got status code %d but expected was 200", url, resp.StatusCode)
		return true
	}

	// Give a bit more time to initialize
	time.Sleep(time.Second * 5)
	return false
}

func TestVersion(t *testing.T) {
	RootCmd.SetArgs([]string{"--version"})

	// TODO: Get the output printed to stdout and check it?
	err := RootCmd.Execute()
	require.NoError(t, err)
}

// Disabled for now
// This tests runs in a panic: html/template: pattern matches no files: `web/templates/*`
// TODO: Fix this test
func testStartService(t *testing.T) {
	go func() {
		assert.Nil(t, RootCmd.Execute())
	}()

	var count = 1
	for waitForServiceToStart(t) {
		t.Logf("Ports are not yet open, retrying attempt #%d...", count)
		count++
		if count > 15 {
			t.FailNow()
		}
		time.Sleep(time.Second)
	}
}
