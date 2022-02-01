package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/schollz/progressbar/v3"
)

const ENDPOINT_ISSUE = "issue"
const ENDPOINT_ATTACHMENT = "attachment"

var username string
var password string
var base_url string
var output_dir string = "."
var issueKey string
var http_retries = 0

type Issue struct {
	ID     int    `json:"id"`
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Attachments []Attachment `json:"attachment"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
}

type Attachment struct {
	ID          string `json:"id"`
	URL         string `json:"content"`
	Filename    string `json:"filename"`
	MimeType    string `json:"mimeType"`
	CreatedTime string `json:"created"`
	Size        int    `json:"size"`
}

func parseArgs() {
	flag.StringVar(&username, "username", "myuser", "Your Jira username")
	flag.StringVar(&password, "password", "mypass", "Your Jira password")
	flag.StringVar(&base_url, "base_url", "https://example.com/rest/api/2", "Your Jira password")
	flag.StringVar(&output_dir, "output", ".", "The base path to your downloaded attachments")
	flag.StringVar(&issueKey, "issue", "WHATEVER-42", "The isse ID")
	flag.Parse()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getUrl(url string) *http.Response {
	var resp *http.Response = nil
	for i := 0; i < http_retries+1; i++ {

		req, err := http.NewRequest("GET", url, nil)
		handleError(err)
		req.SetBasicAuth(username, password)
		resp, err = http.DefaultClient.Do(req)
		handleError(err)
		if resp.StatusCode == http.StatusOK {
			break
		}
	}
	return resp
}

func getIssue() Issue {
	issueURL := getUrlForPath(ENDPOINT_ISSUE + "/" + issueKey + "?fields=attachment,summary,description")
	resp := getUrl(issueURL)
	defer resp.Body.Close()
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	var issue Issue
	json.Unmarshal([]byte(jsonBytes), &issue)
	return issue
}

func getUrlForPath(path string) string {
	baseURL, err := url.Parse(base_url)
	handleError(err)
	issuePath, err := url.Parse(path)
	handleError(err)
	return baseURL.ResolveReference(issuePath).String()
}

func getPathSafeString(str string) string {
	slug.CustomSub = map[string]string{
		" ": "_",
		"[": "",
		"]": "",
	}
	slug.Lowercase = false
	return slug.Make(str)
}

func createDir(dirPath string) {
	os.MkdirAll(dirPath, os.ModePerm)
}

func downloadFile(filepath string, url string, prefix string) error {

	// Get the data
	resp := getUrl(url)
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

func isCompressed(attachment Attachment) bool {

	extensions := []string{".zip", ".rar", ".7z", ".001"}
	for _, extension := range extensions {
		if strings.HasSuffix(strings.ToLower(attachment.Filename), extension) {
			return true
		}

	}
	return false
}

func extractFile(filePath string, outputDir string, prefix string) {
	// TODO: Use Temp Dir and only copy if extraction successful
	fmt.Printf("%sEXTRACTING ...", prefix)
	if !pathExists(outputDir) {
		createDir(outputDir)
		cmd := exec.Command("7z", "x", "-aos", "-o"+outputDir, filePath)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("\r%sERROR!        \n", prefix)
		} else {
			fmt.Printf("\r%sOK            \n", prefix)
		}
	} else {
		fmt.Printf("\r%sSKIPPED       \n", prefix)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getAttachmentFileName(attachment Attachment) string {
	prefix := attachment.ID
	// golang is whack! These are "magic numbers" in the pattern string ...
	createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", attachment.CreatedTime)
	if err == nil {
		prefix = createdTime.Format("2006-01-02")
	}
	return fmt.Sprintf("%s_%s", prefix, attachment.Filename)
}

func getAttachmentProgressPrefix(index int, attachments []Attachment) string {
	bullet := '├'
	if index+1 == len(attachments) {
		bullet = '└'
	}
	return fmt.Sprintf("%c─ %s: ", bullet, attachments[index].Filename)
}

func downloadAttachments() {
	// Make sure output dir exists
	createDir(output_dir)

	// Retrieve issue
	issue := getIssue()
	issueTitle := fmt.Sprintf("%s %s", issue.Key, issue.Fields.Summary)
	fmt.Println(issueTitle)
	//fmt.Println(issue.Fields.Description)

	// Make sure issue dir exists
	issueDir := path.Join(output_dir, getPathSafeString(issueTitle))
	createDir(issueDir)

	// Download attachments
	fmt.Println("┌─────────────┐")
	fmt.Println("│ Downloading │")
	fmt.Println("├─────────────┘")
	var compressedAttachments []Attachment
	for i, attachment := range issue.Fields.Attachments {
		//fmt.Printf("%c┬─ %s (%d bytes)", bullet1, attachment.Filename, attachment.Size)
		attachmentFileName := getAttachmentFileName(attachment)
		filePath := path.Join(issueDir, attachmentFileName)
		prefix := getAttachmentProgressPrefix(i, issue.Fields.Attachments)
		if !pathExists(filePath) {
			err := downloadFile(filePath, attachment.URL, prefix)
			handleError(err)
		} else {
			fmt.Printf("%sSKIPPED\n", prefix)
		}

		if isCompressed(attachment) {
			compressedAttachments = append(compressedAttachments, attachment)
		}
	}

	// Extract compressed attachments

	fmt.Println("┌─────────────┐")
	fmt.Println("│ Extracting  │")
	fmt.Println("├─────────────┘")
	extractedDir := path.Join(issueDir, "__extracted__")
	for i, attachment := range compressedAttachments {
		prefix := getAttachmentProgressPrefix(i, compressedAttachments)
		createDir(extractedDir)
		attachmentFileName := getAttachmentFileName(attachment)
		filePath := path.Join(issueDir, attachmentFileName)
		extractFile(filePath, path.Join(extractedDir, getPathSafeString(attachmentFileName)), prefix)
	}

}

func main() {

	parseArgs()
	downloadAttachments()

}
