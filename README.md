# crispy-pancake
crispy-pancake is a simple web service to auto-protect the `master` branches of new GitHub repositories by listening for Organization events; additionally tags a specified user with the details in a new issue.  Uses the golang/oauth2 and google/go-github libraries.  Keep your `master` branches safe and your cakes o' the pan crispy.

## Requirements
* A GitHub account and Organization
* A GitHub [Personal Access Token](#personal-access-tokens) to access the GitHub API
* [Go](https://golang.org/doc/install) version 1.9 or higher (`go version` to check your version)
* A publicly-accessible server location to run crispy-pancake and receive events (using a solution like [ngrok](https://ngrok.com/), this can be a machine without a public IP)
* The [golang/oauth2](https://github.com/golang/oauth2) library
* An Organization Webhook that gets triggered for repository events; [more details](#installing-a-webhook)

## Usage ##
Install the golang/oauth2 library:
```go
go get golang.org/x/oauth2
```
Clone this repo:
```
git clone https://github.com/StruszSoft/crispy-pancake.git
```
crispy-pancake uses the google/github-go library.  Depending on how your Go installation is configured, you may need to change the line that imports the github-go package.  Find the line that starts with `import "github.com/google/go-github/` and modify accordingly; from that project's README:
```go
import "github.com/google/go-github/v28/github"	// with go modules enabled (GO111MODULE=on or outside GOPATH)
import "github.com/google/go-github/github" // with go modules disabled
```
Update the following variables in crispy-pancake.go:
```
var webhookSecret := "WEBHOOK_SECRET"                 //Must include your Webhook Secret here
var personalAccessToken := "PERSONAL_ACCESS_TOKEN"    //Must include your PAT here
var userToTag := "@GITHUB_USER"                       //Must include the user to tag in new issues here
var entryPoint := "/webhook"                          //Change as desired
var listenPort := ":3000"                             //Change as desired
var issueTitle := "Branch Auto-Protected"	            //Change as desired
```
Build crispy-pancake:
```go
go build
```
Run crispy-pancake.  For Windows (from within the repository directory; use full path to executable otherwise):
```
crispy-pancake
```
For OS/X and Unix/Linux (from within the repository directory; use full path to executable otherwise):
```
./crispy-pancake
```
Finally, create and initialize a new repository; the `master` branch should be marked as protected, with an issue detailing the specific protections applied.  To verify this, go to 'Settings' for the new repository, select 'Branches', and check for 'Branch protection rule.  You can click 'Edit' to view the details of the protection(s) applied.

## Personal Access Tokens ##
To create a personal access token, go to 'Settings' for your user account, select 'Developer Settings', then 'Personal Access Tokens').  Generate a new token with the following scopes:  'repo', 'admin:org', 'admin:org_hook', and 'write:discussion'.  Paste the new token into a temporary document (you won't be able to view it later). Do not save this document (it's got secrets!), just leave it to hang around for a few minutes.  Once we've pasted this token into crispy-pancake, you won't need it.  [Full PAT documentation](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line).
 
## Installing a Webhook ##
From the Organization page, go to 'Settings'.  Click 'Webhooks'.  For the 'Payload URL', enter the URL of the server that will host crispy-pancake.  Use the variable you defined in `entryPoint` after the name/IP of the public site (ie, https://www.example.com/webhook).  For 'Secret', you may enter any string you like; be sure to enter the same string for the `webhookSecret` variable.  Choose the option to "Let me select individual events", and select 'Repositories'.  Make sure it's 'Active', and click 'Update webhook'.

## Future Improvements ##
* Allow easy end-user customization of specific protection requirements
* Prompt for required customizations on first run
* crispier pancakes
