package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"learn.oauth.billing/model"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

// Billing list of services to pay
type Billing struct {
	Services []string `json:"services"`
}

// BillingError error response
type BillingError struct {
	Error string `json:"error"`
}

// TokenIntrospect response
type TokenIntrospect struct {
	Jti      string `json:"jti"`
	Exp      int    `json:"exp"`
	Nbf      int    `json:"nbf"`
	Iat      interface{}    `json:"iat"`
	Aud      string `json:"aud"`
	Typ      string `json:"typ"`
	AuthTime int    `json:"auth_time"`
	Acr      string `json:"acr"`
	Active   bool   `json:"active"`
}

var config = struct {
	tokenIntroSpection string
}{
	tokenIntroSpection: "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/token/introspect",
}

func main() {
	http.HandleFunc("/billing/v1/services", enabledLog(services))
	http.ListenAndServe(":8081", nil)
}

func services(w http.ResponseWriter, r *http.Request) {
	token, err := getToken(r)

	if err != nil {
		log.Println(err)
		makeErrorMessage(w, err.Error())
		return
	}

	// Validate token
	if !validateToken(token) {
		makeErrorMessage(w, err.Error())
		return
	}

	claimBytes, err := getClaim(token)
	if err != nil {
		log.Println(err)
		makeErrorMessage(w, "Cannot parse token claim")
		return
	}

	tokenClaim := &model.TokenClaim{}
	err = json.Unmarshal(claimBytes, tokenClaim)
	if err != nil {
		log.Println(err)
		makeErrorMessage(w, err.Error())
		return
	}

	isValidAudience := false
	for _, v := tokenClaim.AudAsSlice() {
		if (v == "billingService" || v == "billingServiceV2") {
			isValidAudience = true
			break
		}
	}

	if !isValidAudience {
		fmt.Printf("Invalid audiences : %v", tokenClaim.AudAsSlice() )
		makeErrorMessage(w, "invalid token audience. Required audience [billingService, billingServiceV2]")
		return
	}

	if !strings.Contains(tokenClaim.Scope, "getBillingService") {
		makeErrorMessage(w, "invalid token scope. Required scope [getBillingService]")
		return
	}

	s := Billing{
		Services: []string{
			"electric",
			"phone",
			"internet",
			"water",
		},
	}
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	encoder.Encode(s)
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

func getToken(r *http.Request) (string, error) {
	// header
	token := r.Header.Get("Authorization")
	if token != "" {
		auths := strings.Split(token, " ")
		if len(auths) != 2 {
			return "", fmt.Errorf("invalid Authorization header format")
		}
		return auths[1], nil
	}
	// form body
	token = r.FormValue("access_token")
	if token != "" {
		return token, nil
	}

	// query
	token = r.URL.Query().Get("access_token")
	if token != "" {
		return token, nil
	}

	return token, fmt.Errorf("Access token is not presented")
}

func validateToken(token string) bool {
	// Request
	form := url.Values{}
	form.Add("token", token)
	form.Add("token_type_hint", "requesting_party_token")

	req, err := http.NewRequest("POST", config.tokenIntroSpection, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("tokenChecker", "2604a5a4-acbd-4762-914f-2168044b33b2")

	// client
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	byteBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	// process request.
	// if not 200
	if res.StatusCode != 200 {
		log.Println("Status is not 200 : ", res.StatusCode)
		return false
	}

	defer res.Body.Close()
	introSpect := &TokenIntrospect{}
	err = json.Unmarshal(byteBody, introSpect)
	if err != nil {
		log.Println(err)
		return false
	}

	return introSpect.Active
}

func makeErrorMessage(w http.ResponseWriter, errMsg string) {
	log.Println("Error message : ", errMsg)
	s := &BillingError{Error: errMsg} //error message
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(s)

}

func getClaim(token string) ([]byte, error) {
	tokenParts := strings.Split(token, ".")
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return []byte{}, err
	}
	return claim, nil
}
