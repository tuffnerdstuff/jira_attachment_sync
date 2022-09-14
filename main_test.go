package main

import (
	"fmt"
	"testing"
)

func TestGetSearchUriFromUri_IssueUri_SearchUriWithKeyJql(t *testing.T) {
	expectedSearchUri := "https://jira.example.com:1234/rest/api/2/search?fields=attachment,summary,description&jql=key+%3D+ISSUE-KEY"
	searchUri, err := getSearchUriFromUri("https://user:password@jira.example.com:1234/browse/ISSUE-KEY", 2)

	if err != nil || searchUri != expectedSearchUri {
		fmt.Printf("%s does not match %s\n", searchUri, expectedSearchUri)
		t.Fail()
	}

}

func TestGetSearchUriFromUri_FilterUri_SearchUriWithFilterJql(t *testing.T) {
	expectedSearchUri := "https://jira.example.com:1234/rest/api/2/search?fields=attachment,summary,description&jql=filter%3D4321"
	searchUri, err := getSearchUriFromUri("https://user:password@jira.example.com:1234/browse/ISSUE-KEY?filter=4321&somOpt=someVar", 2)

	if err != nil || searchUri != expectedSearchUri {
		fmt.Printf("%s does not match %s\n", searchUri, expectedSearchUri)
		t.Fail()
	}
}

func TestGetSearchUriFromUri_JqlUri_SearchUriWithJqlAsIs(t *testing.T) {
	expectedSearchUri := "https://jira.example.com:1234/rest/api/2/search?fields=attachment,summary,description&jql=assignee+%3D+currentUser%28%29+AND+createdDate+%3E+startOfYear%28%29"
	searchUri, err := getSearchUriFromUri("https://user:password@jira.example.com:1234/browse/ISSUE-KEY?jql=assignee%20%3D%20currentUser()%20AND%20createdDate%20%3E%20startOfYear()&filter=4321&somOpt=someVar", 2)

	if err != nil || searchUri != expectedSearchUri {
		fmt.Printf("%s does not match %s\n", searchUri, expectedSearchUri)
		t.Fail()
	}
}
