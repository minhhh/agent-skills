package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "git-reword",
		Short: "Rewrites git commit messages with conventional format support",
	}

	applyCmd := &cobra.Command{
		Use:   "apply <commit-hash> <new-message>",
		Short: "Rewrite a specific commit message",
		Args:  cobra.ExactArgs(2),
		Run:   runApply,
	}
	applyCmd.Flags().String("author", "", "New author (Name <email>)")

	rootCmd.AddCommand(
		applyCmd,
		&cobra.Command{
			Use:   "analyze <commit-hash>",
			Short: "Get context for a commit to help generate a conventional message",
			Args:  cobra.ExactArgs(1),
			Run:   runAnalyze,
		},
		&cobra.Command{
			Use:   "lint <message>",
			Short: "Check if a message follows conventional commit format",
			Args:  cobra.ExactArgs(1),
			Run:   runLint,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check if the working tree is clean for rebasing",
			Run:   runStatus,
		},
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runApply(cmd *cobra.Command, args []string) {
	hash := args[0]
	msg := args[1]

	author, _ := cmd.Flags().GetString("author")

	if err := checkWorkingTreeClean(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rewriting commit %s with message: %s\n", hash[:7], msg)
	if author != "" {
		fmt.Printf("Updating author to: %s\n", author)
	}

	if err := rewordCommit(hash, msg, author); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully rewrote commit")
}

func runAnalyze(cmd *cobra.Command, args []string) {
	hash := args[0]
	if err := analyzeCommit(hash); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runLint(cmd *cobra.Command, args []string) {
	msg := args[0]
	if err := lintMessage(msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Message follows conventional commit format.")
}

func runStatus(cmd *cobra.Command, args []string) {
	if err := checkWorkingTreeClean(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Working tree is clean. Ready to rebase.")
}

func rewordCommit(hash, newMsg, newAuthor string) error {
	parent, err := getParent(hash)
	if err != nil {
		return fmt.Errorf("getting parent: %w", err)
	}

	commits, err := getCommitList(parent)
	if err != nil {
		return fmt.Errorf("getting commit list: %w", err)
	}

	targetIdx := -1
	for i, c := range commits {
		if strings.HasPrefix(c, hash) || c == hash {
			targetIdx = i
			break
		}
	}

	if targetIdx == -1 {
		return fmt.Errorf("commit %s not found in history", hash[:7])
	}

	tmpDir, err := os.MkdirTemp("", "git-reword")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write new message to a temp file
	msgFile := filepath.Join(tmpDir, "commit-msg.txt")
	if err := os.WriteFile(msgFile, []byte(newMsg), 0644); err != nil {
		return fmt.Errorf("writing message file: %w", err)
	}

	seqEditor := filepath.Join(tmpDir, "sequence-editor.sh")
	if err := writeSequenceEditor(seqEditor, targetIdx, newAuthor, msgFile); err != nil {
		return fmt.Errorf("writing sequence editor: %w", err)
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("GIT_SEQUENCE_EDITOR=%s", seqEditor))
	// We no longer need GIT_EDITOR because we use 'exec git commit --amend'

	rebaseCmd := exec.Command("git", "rebase", "-i", parent)
	rebaseCmd.Env = env
	rebaseCmd.Stdin = os.Stdin
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr

	if err := rebaseCmd.Run(); err != nil {
		return fmt.Errorf("git rebase failed: %w", err)
	}

	return nil
}

func getParent(hash string) (string, error) {
	cmd := exec.Command("git", "rev-parse", hash+"^")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func getCommitList(parentHash string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", parentHash+"..HEAD", "--reverse")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	commits := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(commits) > 0 && commits[0] == "" {
		return commits[:len(commits)-1], nil
	}
	return commits, nil
}

func writeSequenceEditor(path string, targetIdx int, newAuthor, msgFile string) error {
	// Escape quotes for bash
	authorEscaped := strings.ReplaceAll(newAuthor, `"`, `\"`)

	var execCmd string
	if newAuthor != "" {
		execCmd = fmt.Sprintf("exec git commit --amend --author=\\\"%s\\\" --file=\\\"%s\\\" --no-edit", authorEscaped, msgFile)
	} else {
		execCmd = fmt.Sprintf("exec git commit --amend --file=\\\"%s\\\" --no-edit", msgFile)
	}

	content := fmt.Sprintf(`#!/bin/bash
PICK_FILE=$(mktemp)
grep "^pick " "$1" > "$PICK_FILE"
REWORD_LINE=%d
awk -v n="$REWORD_LINE" -v cmd="%s" 'NR == n { print "pick " $2; print cmd; next } { print }' "$PICK_FILE" > "$1"
exit 0
`, targetIdx+1, execCmd)

	return os.WriteFile(path, []byte(content), 0755)
}

func checkWorkingTreeClean() error {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check git status: %w", err)
	}
	if len(out) > 0 {
		return fmt.Errorf("working tree is dirty; please commit or stash changes before rewording")
	}
	return nil
}

func analyzeCommit(hash string) error {
	fmt.Printf("--- Analysis for commit %s ---\n", hash[:7])

	// Current Message
	msgCmd := exec.Command("git", "log", "-1", "--format=%B", hash)
	msgOut, err := msgCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get commit message: %w", err)
	}
	fmt.Printf("Current Message:\n%s\n", string(msgOut))

	// Diff Summary
	diffCmd := exec.Command("git", "show", "--stat", "--patch", hash)
	diffCmd.Stdout = os.Stdout
	diffCmd.Stderr = os.Stderr
	fmt.Println("Diff:")
	return diffCmd.Run()
}

func lintMessage(msg string) error {
	// Simple conventional commit regex
	// type(scope): description
	re := regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?: .+$`)
	if !re.MatchString(msg) {
		return fmt.Errorf("invalid format; expected: <type>(<scope>): <description>")
	}
	return nil
}