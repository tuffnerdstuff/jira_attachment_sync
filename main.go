package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime/debug"
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

func handleBadResponse(resp *http.Response) {
	if !net.IsResponseOK(resp) {
		handleError(fmt.Errorf("could not get valid HTTP response, last response was %v", resp))
	}
}

func getIssue() model.Issue {
	issueURL, err := net.GetUrlForPath(conf.BaseUrl, model.ENDPOINT_ISSUE+"/"+args.IssueKey+"?fields=attachment,summary,description")
	handleError(err)
	resp, err := net.GetUrl(issueURL, conf.Username, conf.Password, conf.RetryCount)
	handleError(err)
	handleBadResponse(resp)
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

func runPostprocessingScript(scriptPath string, issue string, issueDir string) {
	if pathExists(scriptPath) {
		cmd := exec.Command(scriptPath, issue, issueDir)
		err := cmd.Run()
		if err != nil {
			fmt.Println("ERROR!")
		}
	} else {
		fmt.Printf("%s does not exist!\n", scriptPath)
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

func printHeader(title string, level int) {
	if !args.ShowProgress {
		fmt.Println(title)
	} else if level == 1 {
		printAsciiHeader(title, '║', '═', '╔', '╗', '╚', '╝')
	} else {
		printAsciiHeader(title, '│', '─', '┌', '┐', '├', '┘')
	}
}

func printAsciiHeader(title string, vertical rune, horizontal rune, luCorner rune, ruCorner rune, llCorner rune, rlCorner rune) {

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
	printHeader(issueTitle, 1)

	// Make sure issue dir exists
	issueDir := path.Join(conf.OutputDir, getPathSafeString(issueTitle))
	createDir(issueDir)

	if issue.Fields.Attachments == nil || len(issue.Fields.Attachments) == 0 {
		fmt.Println("Issue has no attachments!")
		return
	}

	// Download attachments
	printHeader("Downloading", 2)
	for i, attachment := range issue.Fields.Attachments {
		attachmentFileName := attachment.GetFilenameWithDatePrefix()
		filePath := path.Join(issueDir, attachmentFileName)
		prefix := getAttachmentProgressPrefix(i, issue.Fields.Attachments)
		if !pathExists(filePath) {
			resp, err := net.GetUrl(attachment.URL, conf.Username, conf.Password, conf.RetryCount)
			handleError(err)
			handleBadResponse(resp)
			err = net.DownloadFile(resp, filePath, prefix, args.ShowProgress)
			handleError(err)
		} else {
			fmt.Printf("%sSKIPPED\n", prefix)
		}

	}

	// Call post-processing script
	if args.Script != "" {
		printHeader("Post-Processing", 2)
		runPostprocessingScript(args.Script, issue.Key, issueDir)
	}

}

func catchPanic() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		debug.PrintStack()
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {

	defer catchPanic()

	// Parse Args
	config.ParseArgs(&args)

	// Load Config
	err := config.LoadConfig(&conf, args.ConfigPath)
	handleError(err)

	if args.ShowHelp || !conf.Validate() || !args.Validate() {
		flag.Usage()
	} else {
		downloadAttachments()
	}

}
