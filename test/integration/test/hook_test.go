/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/prow/pkg/github"
	"sigs.k8s.io/prow/pkg/phony"
)

func TestHook(t *testing.T) {
	const (
		commentFile = "./testdata/test_comment.json"
		url         = "http://localhost/hook"
		hmac        = "abcde12345"
		org         = "fake-org-hook"
		repo        = "fake-repo-hook"
		label       = "area/kubectl"
	)

	t.Parallel()

	githubClient, err := github.NewClient(func() []byte { return nil }, func(b []byte) []byte { return b }, "", "http://localhost/fakeghserver")
	if err != nil {
		t.Fatalf("failed to construct GitHub client: %v", err)
	}

	issueID, err := githubClient.CreateIssue(org, repo, "Dummy PR, do not merge", "", 0, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed creating issue: %v", err)
	}
	if err := githubClient.CreateComment(org, repo, issueID, "this is an important work"); err != nil {
		t.Fatalf("Failed creating comment: %v", err)
	}
	comments, err := githubClient.ListIssueComments(org, repo, issueID)
	if err != nil {
		t.Fatalf("Failed listing comments: %v", err)
	}
	if len(comments) == 0 {
		t.Fatal("This shouldn't happen, comment created cannot be found")
	}
	if err := githubClient.AddRepoLabel(org, repo, label, "", ""); err != nil {
		t.Fatalf("Failed add label: %v", err)
	}

	d, err := os.ReadFile(commentFile)
	if err != nil {
		t.Fatalf("Could not read payload file: %v", err)
	}

	d = []byte(strings.ReplaceAll(strings.ReplaceAll(string(d), "{ISSUE_ID_PLACEHOLDER}", strconv.Itoa(issueID)), "{COMMENT_ID_PLACEHOLDER}", strconv.Itoa(comments[0].ID)))

	// Intentionally separate webhook from fakeghserver, to avoid the hassle of
	// supporting webhooks for all faked gh events, as hook is the only place
	// where webhook events are relevant
	t.Log("Send webhook")
	if err := phony.SendHook(url, "issue_comment", d, []byte(hmac)); err != nil {
		t.Fatalf("Error sending hook: %v", err)
	}

	if err := wait.PollUntilContextTimeout(context.Background(), 500*time.Millisecond, 1*time.Minute, true, func(ctx context.Context) (bool, error) {
		gotLabels, err := githubClient.GetIssueLabels(org, repo, issueID)
		if err != nil {
			return false, fmt.Errorf("failed listing issue labels: %w", err)
		}
		for _, l := range gotLabels {
			if l.Name == label {
				t.Log("Found label")
				return true, nil
			}
		}
		return false, nil
	}); err != nil {
		t.Fatalf("Didn't get label: %v", err)
	}
}
