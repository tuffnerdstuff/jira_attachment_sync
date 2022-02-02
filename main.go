package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gosimple/slug"
	"github.com/tuffnerdstuff/jira-attachment-sync/config"
	"github.com/tuffnerdstuff/jira-attachment-sync/model"
	"github.com/tuffnerdstuff/jira-attachment-sync/net"
)

var conf config.Config
var args config.Arguments

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getIssue() model.Issue {
	issueURL, err := net.GetUrlForPath(conf.BaseUrl, model.ENDPOINT_ISSUE+"/"+args.IssueKey+"?fields=attachment,summary,description")
	handleError(err)
	resp, err := net.GetUrl(issueURL, conf.Username, conf.Password, conf.RetryCount)
	handleError(err)
	defer resp.Body.Close()
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	var issue model.Issue
	err = json.Unmarshal([]byte(jsonBytes), &issue)
	handleError(err)
	return issue
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

func getAttachmentProgressPrefix(index int, attachments []model.Attachment) string {
	bullet := '├'
	if index+1 == len(attachments) {
		bullet = '└'
	}
	return fmt.Sprintf("%c─ %s: ", bullet, attachments[index].Filename)
}

func printHeader(title string, vertical rune, horizontal rune, luCorner rune, ruCorner rune, llCorner rune, rlCorner rune) {

	repeatCount := len(title) + 3
	horizontalLine := strings.Repeat(fmt.Sprintf("%c", horizontal), repeatCount)
	fmt.Printf("%c%s%c\n", luCorner, horizontalLine, ruCorner)
	fmt.Printf("%c %s  %c\n", vertical, title, vertical)
	fmt.Printf("%c%s%c\n", llCorner, horizontalLine, rlCorner)
}

func downloadAttachments() {

	// Retrieve issue
	issue := getIssue()
	issueTitle := issue.GetTitle()
	printHeader(issueTitle, '║', '═', '╔', '╗', '╚', '╝')

	// Make sure issue dir exists
	issueDir := path.Join(conf.OutputDir, getPathSafeString(issueTitle))
	createDir(issueDir)

	if issue.Fields.Attachments == nil || len(issue.Fields.Attachments) == 0 {
		fmt.Println("Issue has no attachments!")
		return
	}

	// Download attachments
	printHeader("Downloading", '│', '─', '┌', '┐', '├', '┘')
	var compressedAttachments []model.Attachment
	for i, attachment := range issue.Fields.Attachments {
		attachmentFileName := attachment.GetFilenameWithDatePrefix()
		filePath := path.Join(issueDir, attachmentFileName)
		prefix := getAttachmentProgressPrefix(i, issue.Fields.Attachments)
		if !pathExists(filePath) {
			resp, err := net.GetUrl(attachment.URL, conf.Username, conf.Password, conf.RetryCount)
			handleError(err)
			err = net.DownloadFile(resp, filePath, prefix)
			handleError(err)
		} else {
			fmt.Printf("%sSKIPPED\n", prefix)
		}

		if attachment.IsCompressed() {
			compressedAttachments = append(compressedAttachments, attachment)
		}
	}

	// Extract compressed attachments
	printHeader("Extracting", '│', '─', '┌', '┐', '├', '┘')
	extractedDir := path.Join(issueDir, "__extracted__")
	for i, attachment := range compressedAttachments {
		prefix := getAttachmentProgressPrefix(i, compressedAttachments)
		createDir(extractedDir)
		attachmentFileName := attachment.GetFilenameWithDatePrefix()
		filePath := path.Join(issueDir, attachmentFileName)
		extractFile(filePath, path.Join(extractedDir, getPathSafeString(attachmentFileName)), prefix)
	}

}

func catchPanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {

	defer catchPanic()

	// Parse Args
	config.ParseArgs(&args)

	// Parse Config
	err := config.LoadConfig(&conf, &args)
	handleError(err)

	// Override Config from file with Config from command-line
	config.ParseConfig(&conf)

	if args.ShowHelp || !conf.Validate() || !args.Validate() {
		flag.Usage()
	} else {
		downloadAttachments()
	}

}
