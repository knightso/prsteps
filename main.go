package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/v55/github"
)

const version = "1.0.0"

var (
	showVersion bool
	format      string
	pat         string
)

func init() {
	flag.StringVar(&format, "format", "csv", "Format output(csv, tsv)")
	flag.StringVar(&pat, "pat", "", "GitHub Private Access Tokens(for accessing private repositories)")
	flag.BoolVar(&showVersion, "version", false, "Print version information and quit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: prsteps [OPTIONS] REPO PR\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}
	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "requires exactly 2 argument.\n")
		flag.Usage()
		os.Exit(1)
	}
	owner, repo, err := parseRepositoryName(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pr, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client := github.NewClient(nil)
	if pat := getPAT(); pat != "" {
		client = client.WithAuthToken(pat)
	}

	commitFiles, err := getAllCommitFiles(context.Background(), client, owner, repo, pr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch format {
	case "csv":
		printTable(commitFiles, ',')
	case "tsv":
		printTable(commitFiles, '\t')
	default:
		fmt.Fprintf(os.Stderr, "unsupported format type: %s\n", format)
		os.Exit(1)
	}
}

func printTable(commitFiles []*github.CommitFile, delimiter rune) {
	w := csv.NewWriter(os.Stdout)
	w.Comma = delimiter
	_ = w.Write([]string{"file", "additions", "deletions"})
	for _, file := range commitFiles {
		additions := strconv.Itoa(file.GetAdditions())
		deletions := strconv.Itoa(file.GetDeletions())
		_ = w.Write([]string{file.GetFilename(), additions, deletions})
	}
	w.Flush()
}

func parseRepositoryName(name string) (owner, repo string, err error) {
	owner, repo, ok := strings.Cut(name, "/")
	if !ok {
		return "", "", errors.New("the format of repository name is invalid(expect: {owner}/{repo})")
	}
	return owner, repo, nil
}

func getAllCommitFiles(ctx context.Context, client *github.Client, owner, repo string, pr int) ([]*github.CommitFile, error) {
	allCommitFiles := make([]*github.CommitFile, 0)
	opts := &github.ListOptions{
		PerPage: 100,
		Page:    1,
	}
	for {
		commitFiles, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pr, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in PR%d: %v", pr, err)
		}
		if len(commitFiles) == 0 {
			break
		}
		opts.Page++
		allCommitFiles = append(allCommitFiles, commitFiles...)
	}
	return allCommitFiles, nil
}

func getPAT() string {
	if pat != "" {
		return pat
	}
	return getPATFromGitHubCLI()
}

func getPATFromGitHubCLI() string {
	b, _ := exec.Command("gh", "auth", "token").Output()
	return strings.TrimSpace(string(b))
}
