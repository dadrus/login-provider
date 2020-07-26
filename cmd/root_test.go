package cmd_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"login-provider/cmd"
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

func waitForServiceToStart(t *testing.T, port int) bool {
	client := &http.Client{}

	url := fmt.Sprintf("http://127.0.0.1:%d%s", port, "/health/alive")

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

type stringWriter struct {
	value string
}

// The only function required by the io.Writer interface.  Will append
// written data to the String.value string.
func (s *stringWriter) Write(p []byte) (n int, err error) {
	s.value += string(p)
	return len(p), nil
}

func TestStartService(t *testing.T) {
	// GIVEN
	cmd.RootCmd.SetOut(os.Stdout)
	port := freePort()
	os.Setenv("PORT", fmt.Sprintf("%d", port))

	// WHEN
	go func() {
		assert.Nil(t, cmd.RootCmd.Execute())
	}()

	// THEN
	// service shall be up and running
	var count = 1
	for waitForServiceToStart(t, port) {
		t.Logf("Ports are not yet open, retrying attempt #%d...", count)
		count++
		if count > 5 {
			t.FailNow()
		}
		time.Sleep(time.Second)
	}
}

func TestVersion(t *testing.T) {
	// GIVEN
	cmd.RootCmd.SetArgs([]string{"--version"})
	w := &stringWriter{}
	cmd.RootCmd.SetOut(w)

	// WHEN
	err := cmd.RootCmd.Execute()

	// THEN
	require.NoError(t, err)
	require.Contains(t, w.value, "version master", "Default version must be master")
}
