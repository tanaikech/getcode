/*
Package getcode (doc.go) :
This is a Golang library to automatically get an authorization code for retrieving access token using OAuth2.

When it retrieves an access token and refresh token using OAuth2, the code for retrieving them has to be got by authorization on own browser. In order to retrieve the code, in generally, users have to click the authorization button and copy the code on the browser. This library can be automatically got the code by launching HTML server as a redirected server. At first, I have used this for retrieving the code from Google. But recently I noticed that this can be used for other sites. They are Google, GitHub, Slack and so on. This library can be used for creating such applications.

# Install
You can get this by
$ go get -u github.com/tanaikech/getcode

# Usage

code := getcode.Init(CodeURL, ServerPort, waitingTime, showing, manual).Do()

- CodeURL (string) : URL to get code
- ServerPort (int) : Server port
- waitingTime(second) (int) : waiting time for request
- showing (bool) : Display message for authorization
- manual (bool) : If you use manual mode, it's true.


# Samples

1. In the case of Google

- CodeURL : https://accounts.google.com/o/oauth2/auth?client_id=###&redirect_uri=http://localhost:8080&scope=###&response_type=code&approval_prompt=force&access_type=offline
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false

2. In the case of GitHub

- CodeURL : https://github.com/login/oauth/authorize?client_id=###&scope=###
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false

3. In the case of Slack

- CodeURL : https://slack.com/oauth/authorize?client_id=###&scope=###
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false

More information is https://github.com/tanaikech/getcode

---------------------------------------------------------------
*/
package getcode
