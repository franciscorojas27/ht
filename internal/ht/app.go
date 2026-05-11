package ht

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/franciscorojas27/HT/internal/runner"
	"github.com/franciscorojas27/HT/internal/ui"
)

type basicAuth struct {
	User string
	Pass string
}

type requestOptions struct {
	BasicAuth *basicAuth
	Insecure  bool
}

var errListDone = errors.New("list done")

func Run(flags ui.Flags) error {
	quickMode := shouldUseQuick(flags)
	if quickMode {
		return runQuick(flags)
	}

	requests, targetReq, timeout, err := loadYAMLRequest(flags)
	if err != nil {
		if errors.Is(err, errListDone) {
			return nil
		}
		return err
	}

	opts := requestOptions{}
	return executeRequest(flags, requests, targetReq, timeout, opts, false)
}

func shouldUseQuick(flags ui.Flags) bool {
	if flags.YML == "" {
		return true
	}
	if flags.URL != "" || flags.Method != "" {
		return true
	}
	return false
}

func runQuick(flags ui.Flags) error {
	headers, err := parseHeaderLines(flags.Headers)
	if err != nil {
		return err
	}

	method := quickMethod(flags.Method, flags.Data, flags.HeadOnly)
	if flags.Data != "" {
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	var auth *basicAuth
	if flags.User != "" {
		user, pass, _ := strings.Cut(flags.User, ":")
		auth = &basicAuth{User: user, Pass: pass}
	}

	targetReq := HTRequest{
		Method:  method,
		URL:     flags.URL,
		Headers: headers,
	}
	if flags.Data != "" {
		targetReq.Body = flags.Data
	}

	opts := requestOptions{BasicAuth: auth, Insecure: flags.Insecure}
	return executeRequest(flags, HT{}, targetReq, 30*time.Second, opts, true)
}

func loadYAMLRequest(flags ui.Flags) (HT, HTRequest, time.Duration, error) {
	data, err := os.ReadFile(flags.YML)
	if err != nil {
		return HT{}, HTRequest{}, 0, fmt.Errorf("error reading YAML file: %w", err)
	}

	requests, err := ParseYML(data)
	if err != nil {
		return HT{}, HTRequest{}, 0, fmt.Errorf("error parsing YAML file: %w", err)
	}
	if flags.LS {
		requests.ListRequests()
		return HT{}, HTRequest{}, 0, errListDone
	}

	targetReq, found := requests.FindRequest(flags.Name)
	if !found {
		return HT{}, HTRequest{}, 0, fmt.Errorf("request not found")
	}

	if requests.Config.BaseURL != "" {
		targetReq.URL = requests.Config.BaseURL + targetReq.URL
	}

	var timeout time.Duration
	if requests.Config.Timeout > 0 {
		timeout = time.Duration(requests.Config.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}

	return requests, targetReq, timeout, nil
}

func executeRequest(flags ui.Flags, requests HT, targetReq HTRequest, timeout time.Duration, opts requestOptions, quickMode bool) error {
	bodyReader, contentType, err := targetReq.PrepareBody()
	if err != nil {
		return fmt.Errorf("error preparing request body: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqClient, err := http.NewRequestWithContext(ctx, targetReq.Method, targetReq.URL, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	if !quickMode {
		LoadHeaders(reqClient, requests.Config.Headers)
	}
	LoadHeaders(reqClient, targetReq.Headers)
	if contentType != "" && reqClient.Header.Get("Content-Type") == "" {
		reqClient.Header.Set("Content-Type", contentType)
	}
	if opts.BasicAuth != nil {
		reqClient.SetBasicAuth(opts.BasicAuth.User, opts.BasicAuth.Pass)
	}

	if flags.DryRun {
		ui.DryRun(reqClient)
		return nil
	}

	initTime := time.Now()
	res, err := runner.DoRequest(reqClient, initTime, flags.Verbose, runner.NewClient(opts.Insecure))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	return ui.RenderResponse(targetReq.Method, targetReq.URL, res, flags.Verbose)
}

func parseHeaderLines(lines []string) (map[string]string, error) {
	if len(lines) == 0 {
		return map[string]string{}, nil
	}
	headers := make(map[string]string, len(lines))
	for _, line := range lines {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("invalid header %q, expected 'Key: Value'", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return nil, fmt.Errorf("invalid header %q, empty key", line)
		}
		headers[key] = value
	}
	return headers, nil
}

func quickMethod(method string, data string, headOnly bool) string {
	if headOnly {
		return "HEAD"
	}
	if method != "" {
		return strings.ToUpper(method)
	}
	if data != "" {
		return "POST"
	}
	return "GET"
}
