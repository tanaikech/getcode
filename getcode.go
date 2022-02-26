// Package getcode (getcode.go) :
// Get authorization code
package getcode

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// authCode : For getting authorization code
type authCode struct {
	Code string
	Err  error
}

// serverInfToGetCode : For getting authorization code
type serverInfToGetCode struct {
	Response chan authCode
	Start    chan bool
	End      chan bool
}

// GcParams : Parameters for authorization.
type GcParams struct {
	AuthURL string
	Port    int
	Twait   int
	Msg     bool
	Manual  bool
}

// Init : Inputs parameters for authorization.
func Init(authURL string, port, wait int, msg, manual bool) *GcParams {
	return &GcParams{
		AuthURL: authURL,
		Port:    port,
		Twait:   wait,
		Msg:     msg,
		Manual:  manual,
	}
}

// getAutoCode // Retrieves code by launching server.
func (p *GcParams) getAutoCode() (string, error) {
	if p.Msg {
		fmt.Printf("\n### This is a automatic input mode.\n### Please follow opened browser, login and click authentication.\n### Moves to a manual mode if you wait for %d seconds under this situation.\n", p.Twait)
	}
	s := &serverInfToGetCode{
		Response: make(chan authCode, 1),
		Start:    make(chan bool, 1),
		End:      make(chan bool, 1),
	}
	defer func() {
		s.End <- true
	}()
	go func(port int) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if len(code) == 0 {
				fmt.Fprintf(w, `<html><head><title>Status</title></head><body><p>Erorr.</p></body></html>`)
				s.Response <- authCode{Err: fmt.Errorf("not found code")}
				return
			}
			fmt.Fprintf(w, `<html><head><title>Status</title></head><body><p>The authentication was done. Please close this page.</p></body></html>`)
			s.Response <- authCode{Code: code}
		})
		var err error
		Listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			s.Response <- authCode{Err: err}
			return
		}
		server := http.Server{}
		server.Handler = mux
		go server.Serve(Listener)
		s.Start <- true
		<-s.End
		Listener.Close()
		s.Response <- authCode{Err: err}
	}(p.Port)
	<-s.Start
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", strings.Replace(p.AuthURL, "&", `\&`, -1))
	case "linux":
		cmd = exec.Command("xdg-open", strings.Replace(p.AuthURL, "&", `\&`, -1))
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", strings.Replace(p.AuthURL, "&", `^&`, -1))
	default:
		return "", fmt.Errorf("go manual mode")
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("go manual mode")
	}
	var result authCode
	select {
	case result = <-s.Response:
	case <-time.After(time.Duration(p.Twait) * time.Second):
		return "", fmt.Errorf("go manual mode")
	}
	if result.Err != nil {
		return "", fmt.Errorf("go manual mode")
	}
	return result.Code, nil
}

// Do : Retrieve code for authorization.
func (p *GcParams) Do() string {
	var code string
	var err error
	if !p.Manual {
		code, err = p.getAutoCode()
		if err != nil {
			fmt.Printf("\n### This is a manual input mode.\n### Please input code retrieved by importing following URL to your browser.\n\n"+
				"[URL]==> %v\n"+
				"[CODE]==>", p.AuthURL)
			if _, err := fmt.Scan(&code); err != nil {
				log.Fatalf("Error: %v.\n", err)
			}
		}
	} else {
		fmt.Printf("%v\n", p.AuthURL)
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatalf("Error: %v.\n", err)
		}
	}
	return code
}
