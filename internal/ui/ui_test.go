package ui

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/franciscorojas27/ht/internal/runner"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	oldColorOut := color.Output
	oldColorErr := color.Error
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	color.Output = w
	color.Error = w

	fn()

	require.NoError(t, w.Close())
	os.Stdout = oldStdout
	color.Output = oldColorOut
	color.Error = oldColorErr

	out, err := io.ReadAll(r)
	require.NoError(t, err)
	require.NoError(t, r.Close())

	return string(out)
}

func TestFormatJSON(t *testing.T) {
	out, err := FormatJSON([]byte("{\"a\":1}"))
	require.NoError(t, err)
	require.Contains(t, out, "\"a\"")
}

func TestFormatJSONInvalid(t *testing.T) {
	_, err := FormatJSON([]byte("not-json"))
	require.Error(t, err)
}

func TestTimeTaken(t *testing.T) {
	color.NoColor = true
	out := captureOutput(t, func() {
		TimeTaken(3 * time.Second)
		TimeTaken(1500 * time.Millisecond)
		TimeTaken(500 * time.Millisecond)
	})
	require.Contains(t, out, "Time taken")
}

func TestDryRun(t *testing.T) {
	color.NoColor = true
	req, err := http.NewRequest("POST", "https://example.com", bytes.NewBufferString("hello"))
	require.NoError(t, err)

	out := captureOutput(t, func() {
		DryRun(req)
	})

	require.Contains(t, out, "Dry run mode enabled")
	require.Contains(t, out, "Body")
}

func TestRenderResponseJSONVerbose(t *testing.T) {
	color.NoColor = true
	res := &runner.RequestResult{
		FinalTime:   1500 * time.Millisecond,
		DumpHeaders: []byte("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n"),
		Body:        []byte("{\"ok\":true}"),
		DNSDur:      10 * time.Millisecond,
		ConnDur:     20 * time.Millisecond,
		TLSDur:      30 * time.Millisecond,
		FirstByte:   40 * time.Millisecond,
		TLSVersion:  0x0303,
		SawTLS:      true,
		Headers:     http.Header{"Content-Type": []string{"application/json"}},
	}

	out := captureOutput(t, func() {
		err := RenderResponse("GET", "https://example.com", res, true)
		require.NoError(t, err)
	})

	require.Contains(t, out, "DNS Lookup")
	require.Contains(t, out, "TLS Handshake")
	require.Contains(t, out, "\"ok\": true")
}

func TestRenderResponsePlain(t *testing.T) {
	color.NoColor = true
	res := &runner.RequestResult{
		FinalTime:   200 * time.Millisecond,
		DumpHeaders: []byte("HTTP/1.1 204 No Content\r\nContent-Type: text/plain\r\n\r\n"),
		Body:        []byte("ok"),
		Headers:     http.Header{"Content-Type": []string{"text/plain"}},
	}

	out := captureOutput(t, func() {
		err := RenderResponse("GET", "https://example.com", res, false)
		require.NoError(t, err)
	})

	require.True(t, strings.Contains(out, "HTTP/1.1") || strings.Contains(out, "ok"))
}
