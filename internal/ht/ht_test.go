package ht

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/franciscorojas27/HT/internal/ui"

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

func TestParseEnv(t *testing.T) {
	data := []byte("token=${TOKEN}\nmissing=${MISSING}\n")
	vars := map[string]any{"TOKEN": "abc"}

	out := ParseEnv(data, vars)

	require.Contains(t, string(out), "token=abc")
	require.Contains(t, string(out), "missing=${MISSING}")
}

func TestLoadEnvFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, ".env")
	content := "A=1\n# comment\nB='two'\nC=\"three\"\nEMPTY=\n"
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	vars, err := LoadEnvFile(filePath)
	require.NoError(t, err)
	require.Equal(t, "1", vars["A"])
	require.Equal(t, "two", vars["B"])
	require.Equal(t, "three", vars["C"])
	require.Equal(t, "", vars["EMPTY"])
}

func TestParseYMLWithEnvAndVars(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("user=envUser\ntoken=envToken\n"), 0o600))

	yml := strings.Join([]string{
		"config:",
		"  headers:",
		"    X-Token: ${token}",
		"vars:",
		"  token: varToken",
		"env_file: '" + filepath.ToSlash(envPath) + "'",
		"requests:",
		"  - name: get-user",
		"    method: GET",
		"    url: /users/${user}",
	}, "\n")

	requests, err := ParseYML([]byte(yml))
	require.NoError(t, err)
	require.Equal(t, "varToken", requests.Config.Headers["X-Token"])
	require.Len(t, requests.Requests, 1)
	require.Equal(t, "/users/envUser", requests.Requests[0].URL)
}

func TestFindRequest(t *testing.T) {
	color.NoColor = true
	reqs := HT{Requests: []HTRequest{{Name: "one"}, {Name: "two"}}}

	found, ok := reqs.FindRequest("two")
	require.True(t, ok)
	require.Equal(t, "two", found.Name)

	out := captureOutput(t, func() {
		_, ok := reqs.FindRequest("missing")
		require.False(t, ok)
	})
	require.Contains(t, out, "not found")

	first, ok := reqs.FindRequest("")
	require.True(t, ok)
	require.Equal(t, "one", first.Name)

	out = captureOutput(t, func() {
		_, ok := (HT{}).FindRequest("")
		require.False(t, ok)
	})
	require.Contains(t, out, "No requests found")
}

func TestListRequests(t *testing.T) {
	color.NoColor = true
	reqs := HT{Requests: []HTRequest{{Method: "GET", URL: "/one"}, {Method: "POST", URL: "/two"}}}

	out := captureOutput(t, func() {
		reqs.ListRequests()
	})

	require.Contains(t, out, "Request #1")
	require.Contains(t, out, "Request #2")
}

func TestPrepareBody(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "body.txt")
	require.NoError(t, os.WriteFile(filePath, []byte("file-body"), 0o600))

	reader, contentType, err := (HTRequest{BodyFile: filePath}).PrepareBody()
	require.NoError(t, err)
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "file-body", string(data))
	require.Equal(t, "", contentType)

	reader, contentType, err = (HTRequest{Body: "plain"}).PrepareBody()
	require.NoError(t, err)
	data, err = io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "plain", string(data))
	require.Equal(t, "", contentType)

	reader, contentType, err = (HTRequest{Body: map[string]any{"a": 1}}).PrepareBody()
	require.NoError(t, err)
	data, err = io.ReadAll(reader)
	require.NoError(t, err)
	require.Contains(t, string(data), "\"a\"")
	require.Equal(t, "application/json", contentType)

	reader, contentType, err = (HTRequest{}).PrepareBody()
	require.NoError(t, err)
	require.Nil(t, reader)
	require.Equal(t, "", contentType)
}

func TestLoadHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	require.NoError(t, err)

	LoadHeaders(req, map[string]string{"X-Test": "1"})
	require.Equal(t, "1", req.Header.Get("X-Test"))
}

func TestParseHeaderLines(t *testing.T) {
	headers, err := parseHeaderLines([]string{"A: 1", "B:two"})
	require.NoError(t, err)
	require.Equal(t, "1", headers["A"])
	require.Equal(t, "two", headers["B"])

	headers, err = parseHeaderLines(nil)
	require.NoError(t, err)
	require.Empty(t, headers)

	_, err = parseHeaderLines([]string{"invalid"})
	require.Error(t, err)
}

func TestQuickMethod(t *testing.T) {
	require.Equal(t, "HEAD", quickMethod("", "", true))
	require.Equal(t, "PUT", quickMethod("put", "", false))
	require.Equal(t, "POST", quickMethod("", "a=1", false))
	require.Equal(t, "GET", quickMethod("", "", false))
}

func TestShouldUseQuick(t *testing.T) {
	require.True(t, shouldUseQuick(ui.Flags{}))
	require.False(t, shouldUseQuick(ui.Flags{YML: "file.yml"}))
	require.True(t, shouldUseQuick(ui.Flags{YML: "file.yml", URL: "https://example.com"}))
	require.True(t, shouldUseQuick(ui.Flags{YML: "file.yml", Method: "GET"}))
}

func TestRunQuickDryRun(t *testing.T) {
	color.NoColor = true
	flags := ui.Flags{
		URL:    "https://example.com",
		Data:   "a=1",
		DryRun: true,
		User:   "user:pass",
		Headers: []string{
			"X-Test: 1",
		},
	}

	out := captureOutput(t, func() {
		err := runQuick(flags)
		require.NoError(t, err)
	})

	require.Contains(t, out, "Dry run mode enabled")
}

func TestLoadYAMLRequest(t *testing.T) {
	dir := t.TempDir()
	ymlPath := filepath.Join(dir, "req.yml")
	yml := strings.Join([]string{
		"config:",
		"  base_url: https://api.example.com",
		"  timeout: 5",
		"requests:",
		"  - name: get-user",
		"    method: GET",
		"    url: /users/1",
	}, "\n")
	require.NoError(t, os.WriteFile(ymlPath, []byte(yml), 0o600))

	reqs, target, timeout, err := loadYAMLRequest(ui.Flags{YML: ymlPath, Name: "get-user"})
	require.NoError(t, err)
	require.Equal(t, "get-user", target.Name)
	require.Equal(t, "https://api.example.com/users/1", target.URL)
	require.Equal(t, 5*time.Second, timeout)
	require.Len(t, reqs.Requests, 1)

	_, _, _, err = loadYAMLRequest(ui.Flags{YML: ymlPath, Name: "missing"})
	require.Error(t, err)

	_, _, _, err = loadYAMLRequest(ui.Flags{YML: ymlPath, LS: true})
	require.True(t, errors.Is(err, errListDone))
}

func TestExecuteRequestDryRun(t *testing.T) {
	color.NoColor = true
	flags := ui.Flags{DryRun: true}
	target := HTRequest{Method: "GET", URL: "https://example.com"}

	out := captureOutput(t, func() {
		err := executeRequest(flags, HT{}, target, time.Second, requestOptions{}, true)
		require.NoError(t, err)
	})

	require.Contains(t, out, "Dry run mode enabled")
}

func TestParseHeaderLinesTrimsSpaces(t *testing.T) {
	headers, err := parseHeaderLines([]string{"  X-Test  :  value  "})
	require.NoError(t, err)
	require.Equal(t, "value", headers["X-Test"])
}

func TestParseHeaderLinesEmptyKey(t *testing.T) {
	_, err := parseHeaderLines([]string{": value"})
	require.Error(t, err)
}

func TestPrepareBodyError(t *testing.T) {
	reader, contentType, err := (HTRequest{BodyFile: filepath.Join(t.TempDir(), "missing.txt")}).PrepareBody()
	require.Error(t, err)
	require.Nil(t, reader)
	require.Equal(t, "", contentType)
}

func TestLoadEnvFileMissing(t *testing.T) {
	_, err := LoadEnvFile(filepath.Join(t.TempDir(), "missing.env"))
	require.Error(t, err)
}

func TestParseYMLInvalid(t *testing.T) {
	_, err := ParseYML([]byte("{not: yaml"))
	require.Error(t, err)
}

func TestParseYMLInvalidEnvFile(t *testing.T) {
	yml := "env_file: 'missing.env'\nrequests: []\n"
	_, err := ParseYML([]byte(yml))
	require.Error(t, err)
}

func TestPrepareBodyJSONError(t *testing.T) {
	body := map[string]any{"ch": make(chan int)}
	reader, contentType, err := (HTRequest{Body: body}).PrepareBody()
	require.Error(t, err)
	require.Nil(t, reader)
	require.Equal(t, "", contentType)
}

func TestLoadHeadersDoesNotPanic(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	require.NoError(t, err)

	out := captureOutput(t, func() {
		LoadHeaders(req, map[string]string{"X-A": "1", "X-B": "2"})
	})

	require.Empty(t, strings.TrimSpace(out))
	require.Equal(t, "1", req.Header.Get("X-A"))
	require.Equal(t, "2", req.Header.Get("X-B"))
}

func TestPrepareBodyStringReader(t *testing.T) {
	reader, _, err := (HTRequest{Body: "hello"}).PrepareBody()
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	require.NoError(t, err)
	require.Equal(t, "hello", buf.String())
}
