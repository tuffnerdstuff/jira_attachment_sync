package net

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

const RETRY_SLEEP = 2 * time.Second

func IsResponseOK(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusOK
}

func GetUrl(url string, username string, password string, retries int) (*http.Response, error) {
	var resp *http.Response = nil
	var err error = nil
	padLen := len(fmt.Sprintf("%d", retries))
	decPlaceholder := fmt.Sprintf("%%0%dd", padLen)
	countFormatString := fmt.Sprintf("[%s/%s]", decPlaceholder, decPlaceholder)
	for i := 1; i <= retries; i++ {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(username, password)
		resp, err = http.DefaultClient.Do(req)
		if err == nil && IsResponseOK(resp) {
			break
		}
		fmt.Printf(countFormatString+" Retrying %s\n", i, retries, url)
		time.Sleep(RETRY_SLEEP)
	}
	return resp, err
}

func GetUrlForPath(base_url string, path string) (string, error) {
	baseURL, err := url.Parse(base_url)
	if err != nil {
		return "", err
	}
	issuePath, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(issuePath).String(), err
}

func DownloadFile(resp *http.Response, filepath string, prefix string, showProgress bool) error {

	// Get the data
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Progress Bar
	var writer io.Writer = out
	if showProgress {
		bar := progressbar.DefaultBytes(resp.ContentLength, prefix)
		writer = io.MultiWriter(out, bar)
	}

	// Write the body to progressbar and file
	_, err = io.Copy(writer, resp.Body)
	return err
}
