package main

import (
	"encoding/base64"
	"log"
	"strings"
	"testing"
)

func TestCliaimDecoder(t *testing.T) {
	tokenString := `eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI0TVVSUGlEVzlfUVJKQklybGZYclBYcVhiSERDUUZWZ2M5UlBNeTdSeDRJIn0.eyJqdGkiOiIxNzJiOWZkNi0yY2MyLTQyN2UtODczOC01YWRkOTJhMGNmOWEiLCJleHAiOjE1NzkzMTYyMjksIm5iZiI6MCwiaWF0IjoxNTc5MzE1OTI5LCJpc3MiOiJodHRwOi8vMTAuMTAwLjE5Ni42MDo4MDgwL2F1dGgvcmVhbG1zL2xlYXJuaW5nQXBwIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6IjQwMDkyOGI0LWMwMGItNGVjNS1hNTIyLTlmOTg3Y2NiNzZiMiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImJpbGxpbmdBcHAiLCJhdXRoX3RpbWUiOjE1NzkzMTQ4NzAsInNlc3Npb25fc3RhdGUiOiI4MjljMWRiMS1kNTFhLTQxMjEtYmQxMC01MTk1YTI5MTkyNzMiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODA4MCJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwibmFtZSI6ImJvYiBib2IiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib2IiLCJnaXZlbl9uYW1lIjoiYm9iIiwiZmFtaWx5X25hbWUiOiJib2IiLCJlbWFpbCI6ImJvYkBib2IuY29tIn0.R6xa4vhAADT7N7E2RESfrggK9fxMv6cRUMJkKs6qQIhPH3v0njCbYnGIvwz5_YqIclNDyLEVBplqMG7FOKcs5Xc01rgGKTMg3JRsuMQMpjY0aSkfm0vrfn1P4EZvCxltbNL0iLfntdmtEls-pVU7IAINSDc7CPTfbOFcIv85aN1XHoEObJxpU23MqeFloXNHV0wMBkrXVJE1zbkp4sgMKr3GnWTN1H9FUSSZu5wSbe9qrNoitm8vttijz_AjXSyGnCMEZ-Ow4Aq9vWvaSNAlWhCJ5bOqU01OQWVp_Cj0l3CbNeYkBRmtu09hxM9d5rcJoF2YPMTItDQC7Icx39LA3w`
	tokenParts := strings.Split(tokenString, ".")
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	log.Println("tokenParts[1]\n", tokenParts[1])

	if err != nil {
		t.Error(err)
	}
	log.Println("Claim : ", string(claim))
}
