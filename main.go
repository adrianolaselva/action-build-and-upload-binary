package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	_ = godotenv.Load(".env")

	accessToken := ""
	if value, ok := syscall.Getenv("API_TOKEN_GITHUB"); ok {
		accessToken = value
	}

	owner := ""
	if value, ok := syscall.Getenv("OWNER"); ok {
		owner = value
	}

	repository := ""
	if value, ok := syscall.Getenv("REPOSITORY"); ok {
		repSplitted := strings.Split(value, "/")
		if len(repSplitted) > 0 {
			repository = repSplitted[1]
		}
	}

	tag := ""
	if value, ok := syscall.Getenv("TAG"); ok && strings.Contains(value, "refs/tags/") {
		tag = strings.ReplaceAll(value, "refs/tags/", "")
	}

	sha := ""
	if value, ok := syscall.Getenv("SHA"); ok {
		sha = value
	}

	if repository == "" {
		fmt.Fprintf(os.Stderr, "Repository not defined\n")
		os.Exit(-10)
	}

	if owner == "" {
		fmt.Fprintf(os.Stderr, "Owner not defined\n")
		os.Exit(-11)
	}

	if tag == "" && sha == "" {
		fmt.Fprintf(os.Stderr, "TAG and SHA not defined\n")
		os.Exit(-11)
	}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	githubClient := github.NewClient(oauth2.NewClient(ctx, tokenSource))

	if tag != "" {
		retrieveReleaseByTag(ctx, githubClient, owner, repository, tag)
		return
	}

	if sha != "" {
		retrieveCommitBySha(ctx, githubClient, owner, repository, sha)
		return
	}
}

func retrieveReleaseByTag(ctx context.Context, githubClient *github.Client, owner string, repository string, currentTag string) {
	tagResp, resp, err := githubClient.Repositories.GetReleaseByTag(ctx, owner, repository, currentTag)
	if err != nil {
		defer resp.Response.Body.Close()
		fmt.Fprintf(os.Stderr, "Failed To Load Release Tag Properties: %v\n", err.Error())
		os.Exit(-1)
	}

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("\n- Id: %v\n", *tagResp.ID))
	message.WriteString(fmt.Sprintf("- Release: %v\n", *tagResp.TagName))
	message.WriteString(fmt.Sprintf("- Url: %v\n", *tagResp.AssetsURL))
	message.WriteString(fmt.Sprintf("- User: %v\n", *tagResp.Author.Login))
	message.WriteString(fmt.Sprintf("- Data: %v\n", *tagResp.PublishedAt))
	message.WriteString(fmt.Sprintf("- Changelog: \n%v\n", *tagResp.Body))

	fmt.Fprintf(os.Stdout, "%v\n", message.String())
}

func retrieveCommitBySha(ctx context.Context, githubClient *github.Client, owner string, repository string, sha string) {
	commit, resp, err := githubClient.Repositories.GetCommit(ctx, owner, repository, sha)
	if err != nil {
		defer resp.Response.Body.Close()
		fmt.Fprintf(os.Stderr, "Failed To Load Release Tag Properties: %v\n", err.Error())
		os.Exit(-1)
	}

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("\n- SHA: %v\n", *commit.SHA))
	message.WriteString(fmt.Sprintf("- Url: %v\n", *commit.URL))
	message.WriteString(fmt.Sprintf("- User: %v\n", *commit.Author.Login))

	fmt.Fprintf(os.Stdout, "%v\n", message.String())
}

func executeCommand(cmd string, args string) {
	out, err := exec.Command(cmd, args).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed To Execute Command: %v\n", err.Error())
		os.Exit(-2)
	}

	fmt.Fprintf(os.Stdout, "Output: %v", string(out))
}
