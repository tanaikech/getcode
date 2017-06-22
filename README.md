getcode
=====

[![Build Status](https://travis-ci.org/tanaikech/getcode.svg?branch=master)](https://travis-ci.org/tanaikech/getcode)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENCE)

<a name="TOP"></a>
# Overview
This is a Golang library to automatically get an authorization code for retrieving access token using OAuth2.

# Description
When it retrieves an access token and refresh token using OAuth2, the code for retrieving them has to be got by authorization on own browser. In order to retrieve the code, in generally, users have to click the authorization button and copy the code on the browser. This library can be automatically got the code by launching HTML server as a redirected server. At first, I have used this for retrieving the code from Google. But recently I noticed that this can be used for other sites. They are Google, GitHub, Slack and so on. This library can be used for creating such applications.

This method was used for [gogauth](https://github.com/tanaikech/gogauth). Next, it was used for [ggsrun](https://github.com/tanaikech/ggsrun) And, this was recreated as a library.

# Install
You can get this by

~~~bash
$ go get -u github.com/tanaikech/getcode
~~~

# Usage

~~~
code := getcode.Init(CodeURL, ServerPort, waitingTime, showing, manual).Do()
~~~

- CodeURL (string) : URL to get code
- ServerPort (int) : Server port
- waitingTime(second) (int) : waiting time for request
- showing (boolean) : Display message for authorization
- manual (boolean) : If you use manual mode, it's true.

## Flow
1. ``getcode`` is run, your browser is launched and waits for login to the site you want to retrieve access token.
1. After logged in the site, you will see the button for authorizing. When push it, the HTML server by this library retrieves the authorization code.
1. The authorization code can be retrieved automatically. And ``Done``. is displayed on your terminal.

If your browser isn't launched or spends for 30 seconds from the wait of authorization, it becomes the input work queue. This is a manual mode. Please copy displayed URL and paste it to your browser, and login the site.

**By above flow, the authorization code can be retrieved.**

# Samples
These are the samples for retrieving the code from Google, GitHub and Slack.

### 1. In the case of Google

- CodeURL : ``https://accounts.google.com/o/oauth2/auth?client_id=###&redirect_uri=http://localhost:8080&scope=###&response_type=code&approval_prompt=force&access_type=offline``
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false

### 2. In the case of GitHub

- CodeURL : ``https://github.com/login/oauth/authorize?client_id=###&scope=###``
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false

### 3. In the case of Slack

- CodeURL : ``https://slack.com/oauth/authorize?client_id=###&scope=###``
- ServerPort : 8080
- waitingTime : 10
- showing message : true
- manual : false


<a name="Update_History"></a>
# Update History

* v1.0.0 (June 21, 2017)

    Initial release.

[TOP](#TOP)
