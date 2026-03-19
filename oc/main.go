package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	opencode "github.com/sst/opencode-sdk-go"
	"github.com/sst/opencode-sdk-go/option"
)

const baseURL = "http://127.0.0.1:4096"

type Review struct {
	Status string  `json:"status"`
	Issues []Issue `json:"issues"`
}

type Issue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func main() {
	// Get staged diff
	out, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal").Output()
	if err != nil {
		fatal("unable to get git diff: %v", err)
		return
	}
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		fmt.Println("no staged changes to review")
		return
	}

	// Build prompt
	prompt := `You are a code reviewer. Review the staged git diff below.

Look for bugs, typos, security issues, and code style problems.

Respond ONLY with a JSON object (no markdown fences, no extra text):
{"status":"pass|fail|warn","issues":[{"file":"...","line":0,"severity":"error|warning|info","message":"..."}]}
If everything looks good, return {"status":"pass","issues":[]}.

` + "```diff\n" + diff + "\n```"

	client := opencode.NewClient(
		option.WithBaseURL(baseURL),
		option.WithMaxRetries(1),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// create opencode session
	session, err := client.Session.New(ctx, opencode.SessionNewParams{
		Title: opencode.F("pre-commit review"),
	})
	if err != nil {
		fatal("unable to create session: %v", err)
		return
	}
	defer client.Session.Delete(context.Background(), session.ID, opencode.SessionDeleteParams{})

	fmt.Fprintln(os.Stderr, "Reviewing staged changes...")

	// send prompt
	resp, err := client.Session.Prompt(ctx, session.ID, opencode.SessionPromptParams{
		Parts: opencode.F([]opencode.SessionPromptParamsPartUnion{
			opencode.TextPartInputParam{
				Type: opencode.F(opencode.TextPartInputTypeText),
				Text: opencode.F(prompt),
			},
		}),
	})
	if err != nil {
		fatal("unable to prompt: %v", err)
		return
	}

	var text string
	for _, part := range resp.Parts {
		if tp, ok := part.AsUnion().(opencode.TextPart); ok {
			text += tp.Text
		}
	}

	var review Review
	if err := json.Unmarshal([]byte(text), &review); err != nil {
		fatal("unable to parse json: %v\nraw response:\n%s", err, text)
		return
	}

	fmt.Printf("Review status: %s\n", review.Status)
	for _, issue := range review.Issues {
		fmt.Printf("  [%s] %s:%d — %s\n", issue.Severity, issue.File, issue.Line, issue.Message)
	}
	if len(review.Issues) == 0 {
		fmt.Println("No issues found!")
	}

	if review.Status == "fail" {
		os.Exit(1)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "opencode-pre-commit: "+format+"\n", args...)
	os.Exit(1)
}
