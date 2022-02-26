# getcode

<a name="top"></a>

# Overview

This is a Golang library to automatically get an authorization code for retrieving access token using OAuth2.

# Description

When it retrieves an access token and refresh token using OAuth2, the code for retrieving them has to be got by authorization on own browser. In order to retrieve the code, in generally, users have to click the authorization button and copy the code on the browser. This library can be automatically got the code by launching HTML server as a redirected server. At first, I have used this for retrieving the code from Google. But recently I noticed that this can be used for other sites. They are Google, GitHub, Slack and so on. This library can be used for creating such applications.

This method was used for [gogauth](https://github.com/tanaikech/gogauth). Next, it was used for [ggsrun](https://github.com/tanaikech/ggsrun) And, this was recreated as a library.

# Install

You can get this by

```bash
$ go get -u github.com/tanaikech/getcode
```

- If "go get" cannot be used, please use this library with `GO111MODULE=on`. And `import("github.com/tanaikech/getcode")`

# Usage

```
code := getcode.Init(CodeURL, ServerPort, waitingTime, showing, manual).Do()
```

- CodeURL (string) : URL to get code
- ServerPort (int) : Server port
- waitingTime(second) (int) : waiting time for request
- showing (boolean) : Display message for authorization
- manual (boolean) : If you use manual mode, it's true.

- In this case, please set `http://localhost:8080` as the redirect URI.

## Flow

1. `getcode` is run, your browser is launched and waits for login to the site you want to retrieve access token.
1. After logged in the site, you will see the button for authorizing. When push it, the HTML server by this library retrieves the authorization code.
1. The authorization code can be retrieved automatically. And `Done`. is displayed on your terminal.

If your browser isn't launched or spends for 30 seconds from the wait of authorization, it becomes the input work queue. This is a manual mode. Please copy displayed URL and paste it to your browser, and login the site.

**By above flow, the authorization code can be retrieved.**

# Samples

These are the samples for retrieving the code from Google, GitHub and Slack.

### 1. In the case of Google

```go
code := getcode.Init(CodeURL, 8080, 30, true, false).Do()
```

- CodeURL : `https://accounts.google.com/o/oauth2/auth?client_id=###&redirect_uri=http://localhost:8080&scope=###&response_type=code&approval_prompt=force&access_type=offline`
- ServerPort : 8080
- waitingTime : 30
- showing message : true
- manual : false

In this case, please create the OAuth2 credential as Web application. please set `http://localhost:8080` as the redirect URI.

As a sample script, when this library is used for [the sample script of Quickstart for Go of Drive API](https://developers.google.com/drive/api/v3/quickstart/go), it becomes as follows. When you run this script, a browser is automatically opened and when you authorize the scopes, the authorization code is automatically retrieved, and `token.json` is created.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/tanaikech/getcode"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// Here, this library is used.
    // In this case, the redirect URI is http://localhost:8080
	code := getcode.Init(authURL, 8080, 30, true, false).Do()

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}
```

### 2. In the case of GitHub

- CodeURL : `https://github.com/login/oauth/authorize?client_id=###&scope=###`
- ServerPort : 8080
- waitingTime : 30
- showing message : true
- manual : false

In this case, please set `http://localhost:8080` as the redirect URI.

### 3. In the case of Slack

- CodeURL : `https://slack.com/oauth/authorize?client_id=###&scope=###`
- ServerPort : 8080
- waitingTime : 30
- showing message : true
- manual : false

In this case, please set `http://localhost:8080` as the redirect URI.

<a name="Update_History"></a>

# Update History

- v1.0.0 (June 21, 2017)

  Initial release.

- v1.0.1 (February 26, 2022)

  Latest libraries are reflected to this library. And, the sample script is shown using the Quickstart for go of Drive API.

[TOP](#top)
