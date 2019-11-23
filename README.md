# crispy-pancake
A simple web service to auto-protect the MASTER branch of new GitHub repositories by listening for organization events; additionally tags a specified user with the details in a new issue.  Uses the golang/oauth2 and google/go-github libraries.  Keep your master branches safe and your cakes o' the pan crispy.

## Requirements
* A GitHub account and organization
* A GitHub Personal Access Token to access the GitHub API (Account Settings --> Developer Settings --> Personal Access Tokens); generate a new token and save it somewhere (you won't be able to view it later).
* Go version 1.9 or higher
* A publicly-accessible server location to run crispy-pancake and receive events (using a solution like [ngrok](https://ngrok.com/), this can be a machine without a public IP)
* The [golang/oauth2](https://github.com/golang/oauth2) library
* An organization webhook that gets triggered for repository events

## Usage ##
Install the golang/oauth2 library:
```go
go get golang.org/x/oauth2
```
Clone this repo:
```
git clone https://github.com/StruszSoft/crispy-pancake.git
```
Update the following variables in crispy-pancake.go:
```
var webhookSecret := "WEBHOOK_SECRET"
var personalAccessToken := "PERSONAL_ACCESS_TOKEN"
var entryPoint := "/webhook"
var listenPort := ":3000"
var userToTag := "@GITHUB_USER"
var issueTitle := "Branch Auto-Protected"	//Change as desired
```
Build crispy-pancake:
```go
go build
```
Run crispy-pancake:
```
crispy-pancake
```
Finally, create a new repository; it should be marked as protected, with an issue detailing the specific protections applied.


## Future Improvements ##
* Allow easy end-user customization of specific protection requirements
* Prompt for required customizations on first run
* crispier pancakes
