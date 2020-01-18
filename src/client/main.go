package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"learn.oauth.client/model"
)

var config = struct {
	appID               string
	appPassword         string
	authURL             string
	logout              string
	afterLogoutRedirect string
	authCodeCallback    string
	tokenEndpoint       string
	servicesEndpoint    string
}{
	appID:               "billingApp",
	appPassword:         "9399482a-24de-4543-9c23-38fe0de45357",
	authURL:             "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logout:              "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutRedirect: "http://localhost:8080/home",
	authCodeCallback:    "http://localhost:8080/authCodeRedirect",
	tokenEndpoint:       "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/token",
	servicesEndpoint:    "http://localhost:8081/billing/v1/services",
}

var t = template.Must(template.ParseFiles("template/index.html"))
var tServices = template.Must(template.ParseFiles("template/index.html", "template/services.html"))

// App Application private variables
type AppVar struct {
	AuthCode     string
	SessionState string
	AccessToken  string
	RefreshToken string
	Scope        string
	Services     []string
	State        map[string]struct{}
}

func newAppVar() AppVar {
	return AppVar{State: make(map[string]struct{})}
}

var appVar = newAppVar()

func init() {
	log.SetFlags(log.Ltime)
}

func main() {
	http.HandleFunc("/home", enabledLog(home))
	http.HandleFunc("/login", enabledLog(login))
	http.HandleFunc("/logout", enabledLog(logout))
	http.HandleFunc("/services", enabledLog(services))
	http.HandleFunc("/authCodeRedirect", enabledLog(authCodeRedirect))
	http.ListenAndServe(":8080", nil)
}

func enabledLog(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.SetPrefix(handlerName + " ")

		log.Println("--> " + handlerName)
		log.Printf("request : %v\n", r.RequestURI)
		// log.Printf("response : %v\n", w)
		handler(w, r)
		log.Println("<-- " + handlerName + "\n")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("template/index.html"))
	t.Execute(w, appVar)
}

func login(w http.ResponseWriter, r *http.Request) {
	// create a redirect URL for authentication endpoint
	req, err := http.NewRequest("GET", config.authURL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	qs := url.Values{}
	state := uuid.New().String()
	appVar.State[state] = struct{}{}
	qs.Add("state", state)
	qs.Add("client_id", config.appID)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", config.authCodeCallback)

	req.URL.RawQuery = qs.Encode()
	http.Redirect(w, r, req.URL.String(), http.StatusFound)
}

func authCodeRedirect(w http.ResponseWriter, r *http.Request) {
	appVar.AuthCode = r.URL.Query().Get("code")
	callBackState := r.URL.Query().Get("state")

	if _, ok := appVar.State[callBackState]; !ok {
		fmt.Fprintf(w, "Error")
		return
	}

	delete(appVar.State, callBackState)

	appVar.SessionState = r.URL.Query().Get("session_state")
	r.URL.RawQuery = ""
	fmt.Printf("Request queries : %+v\n", appVar)
	// http.Redirect(w, r, "http://localhost:8080", http.StatusFound)

	// exchange token here
	exchangeToken()
	t.Execute(w, appVar)
}

func logout(w http.ResponseWriter, r *http.Request) {
	q := url.Values{}
	q.Add("redirect_uri", config.afterLogoutRedirect)

	logoutURL, err := url.Parse(config.logout)
	if err != nil {
		log.Println(err)
	}

	logoutURL.RawQuery = q.Encode()
	appVar = newAppVar()
	http.Redirect(w, r, logoutURL.String(), http.StatusFound)
}

func exchangeToken() {

	// Request
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", appVar.AuthCode)
	form.Add("redirect_uri", config.authCodeCallback)
	form.Add("client_id", config.appID)
	req, err := http.NewRequest("POST", config.tokenEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Println(err)
		return
	}

	req.SetBasicAuth(config.appID, config.appPassword)

	// Client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println("couldn't get access token", err)
		return
	}

	// Process response
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Println(err)
		return
	}

	accessTokenResponse := &model.AccessTokenResponse{}
	json.Unmarshal(byteBody, accessTokenResponse)

	// 取得できたaccess_tokenはjwt形式になっている
	appVar.AccessToken = accessTokenResponse.AccessToken
	appVar.RefreshToken = accessTokenResponse.RefreshToken
	appVar.Scope = accessTokenResponse.Scope
	log.Println(string(byteBody))

}

func services(w http.ResponseWriter, r *http.Request) {

	// request
	req, err := http.NewRequest("GET", config.servicesEndpoint, nil)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	req.Header.Add("Authorization", "Bearer "+appVar.AccessToken)

	// client
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()

	c := http.Client{}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	byteBody, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	// process response
	if res.StatusCode != 200 {
		log.Println(string(byteBody))
	}

	billingResponse := &model.Billing{}

	err = json.Unmarshal(byteBody, billingResponse)
	if err != nil {
		log.Println(err)
		tServices.Execute(w, appVar)
		return
	}

	appVar.Services = billingResponse.Services
	tServices.Execute(w, appVar)
}