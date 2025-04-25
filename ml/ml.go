package ml

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net/http"
)
var apiKey = "AIzaSyBfHli_XdpixsFNKXGWsghT4AmljGkBwnI"

type Candidate struct {
	Content Content `json:"content"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

func ReqML(req string) string {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	requestBody := RequestBody{
    Contents: []Content{
      {
        Parts: []Part{
          {Text: req},
        },
      },
    },
  }
  requestBodyBytes, err := json.Marshal(requestBody)
  if err != nil {
    log.Fatalf("Помилка перетворення JSON: %v", err)
  }
  reqq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
  if err != nil {
    log.Fatalf("Помилка створення запиту: %v", err)
  }
  reqq.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(reqq)
  if err != nil {
    log.Fatalf("Помилка виконання запиту: %v", err)
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatalf("Помилка читання відповіді: %v", err)
  }
  var jsonResponse Response
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Fatalf("Помилка розбору JSON: %v", err)
	}

	if len(jsonResponse.Candidates) > 0 && len(jsonResponse.Candidates[0].Content.Parts) > 0 {
		text := jsonResponse.Candidates[0].Content.Parts[0].Text
		return text
	}

	return ""
}