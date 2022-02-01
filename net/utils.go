package net

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/schollz/progressbar/v3"
)

func GetUrl(url string, username string, password string, retries int) (*http.Response, error) {
	var resp *http.Response = nil
	var err error
	for i := 0; i < retries+1; i++ {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(username, password)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
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

func DownloadFile(resp *http.Response, filepath string, prefix string) error {

	// Get the data
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Progress Bar
	bar := progressbar.DefaultBytes(resp.ContentLength, prefix)

	// Write the body to progressbar and file
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return err
}
