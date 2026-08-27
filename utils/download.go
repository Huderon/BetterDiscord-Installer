package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func DownloadFile(url string, destPath string) (response *http.Response, err error) {
	// Setup the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "BetterDiscord/cli")
	req.Header.Add("Accept", "application/octet-stream")

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Check the response before touching the destination so a failed download
	// (bad status, dropped connection) can't truncate an existing file such as
	// a previously-valid betterdiscord.asar.
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("bad status code: %s", resp.Status)
	}

	// Stream into a temp file alongside the destination, then atomically rename
	// into place so an interrupted or partial download never leaves a corrupt
	// file behind.
	out, err := os.CreateTemp(filepath.Dir(destPath), ".download-*")
	if err != nil {
		return resp, err
	}
	tmpPath := out.Name()
	// Best-effort cleanup; a no-op once the file has been renamed into place.
	defer os.Remove(tmpPath)

	if _, err = io.Copy(out, resp.Body); err != nil {
		out.Close()
		return resp, err
	}
	if err = out.Close(); err != nil {
		return resp, err
	}

	if err = os.Rename(tmpPath, destPath); err != nil {
		return resp, err
	}

	return resp, nil
}

func DownloadJSON[T any](url string) (T, error) {
	var data T

	// Setup the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("User-Agent", "BetterDiscord/cli")

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("bad status: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
