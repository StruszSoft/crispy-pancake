/*
	crispy-pancake

	Receive organization events from GitHub; watch for and act upon repo creation
	by protecting the master branch for the new repo and creating an issue describing
	the protections applied, tagging a user in the process.

	This is my first foray into Go (hey, I figure if I'm supposed to show that I can
	learn new things quickly, what better way, right?), so I apologize if I've violated
	ALL the conventions and offended thine eyes.  That doesn't sound sincere, but it is.

*/
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var webhookSecret = "WEBHOOK_SECRET"
var personalAccessToken = "PERSONAL_ACCESS_TOKEN"
var entryPoint = "/webhook"
var listenPort = ":3000"
var userToTag = "@GITHUB_USER"
var issueTitle = "Branch Auto-Protected" // Change as desired

// Listen for calls to the server and act on organization events from GitHub
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Use go-github to validate the payload against our secret.
	payload, err := github.ValidatePayload(r, []byte(webhookSecret))
	if err != nil {
		log.Printf("Error validating payload: err=%s\n", err)
		return
	}
	defer r.Body.Close()
	// Get details about this event from the payload.
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("Could not parse webhook: err=%s\n", err)
		return
	}
	// When we detect a repository has been created, protect the MASTER branch and create an issue to detail the protection; tag a specified user in that issue.
	switch e := event.(type) {
	case *github.RepositoryEvent:
		if e.Action != nil && *e.Action == "created" {
			log.Printf("%s created repository %s; protecting the master branch\n", *e.Sender.Login, *e.Repo.FullName)
			time.Sleep(10 * time.Second) // It seems that the master branch takes just a moment to be available; tread water for just a bit.  10 seconds is probably excessive.
			protectMasterBranch(*e.Repo.Owner.Login, *e.Repo.Name)
			createIssue(*e.Repo.Owner.Login, *e.Repo.Name)
		}
	default:
		return
	}
}

// Apply delete protection to master branch of a repository
func protectMasterBranch(orgName string, repoName string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personalAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// Configure variables for our protection request
	protectionRequest := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   true,
			Contexts: []string{"continuous-integration"},
		},
		EnforceAdmins:              false,
		RequiredPullRequestReviews: nil,
		Restrictions:               nil,
	}
	// Update protection on the master branch
	log.Printf("--Protecting master branch")
	protection, _, err := client.Repositories.UpdateBranchProtection(context.Background(), orgName, repoName, "master", protectionRequest)
	if err != nil {
		// Handle error.
		log.Fatal(err)
	}
	log.Printf("\n%v\n", github.Stringify(protection))
}

// Create a new issue on a repository detailing protections applied to the master branch; tag the user specified by the userToTag variable.
func createIssue(orgName string, repoName string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personalAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	// Configure variables for our issue request
	issueRequest := &github.IssueRequest{
		Title: github.String(issueTitle),
		Body:  github.String(userToTag + " - Protections applied:  Require status checks to pass before merging (continuous-integration), Require branches to be up to date before merging."),
	}
	// Create the new issue
	log.Printf("--Creating Issue")
	issue, _, err := client.Issues.Create(context.Background(), orgName, repoName, issueRequest)
	if err != nil {
		// Handle error.
		log.Fatal(err)
	}
	log.Printf("\n%v\n", github.Stringify(issue))
}

func main() {
	log.Println("Starting server.")
	http.HandleFunc(entryPoint, handleWebhook)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
