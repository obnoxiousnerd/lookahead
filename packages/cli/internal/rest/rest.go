package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
	"lookahead.web.app/cli/internal/credential"
	"lookahead.web.app/cli/internal/version"
)

func getAPIEndpoint() string {
	if version.Version == "Dev" {
		return "http://localhost:3000/todos"
	}
	return "https://lookahead-api.vercel.app/todos"
}

type dataSchema struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Id         string `json:"id"`
	LastEdited string `json:"last_edited"`
}

type restClientStruct struct {
	//IdToken The token used for making transactions with the database
	IdToken    string
	httpClient http.Client
}

type restClientMethods interface {
	Add(title string, content string, last_edited string) error
	Delete()
	GetAll() ([]dataSchema, error)
	Get(id string) (dataSchema, error)
	Set(id string, title string, content string, last_edited string) error
}

func (c restClientStruct) Add(title string, content string, last_edited string) error {
	reqBodyData := map[string]string{
		"title":       title,
		"content":     content,
		"last_edited": last_edited,
	}
	reqBody := map[string]interface{}{
		"data": reqBodyData,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, getAPIEndpoint(), bytes.NewBuffer(reqBodyJSON))
	req.Header.Add("Authorization", "Bearer "+c.IdToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	resMsg := gjson.GetBytes(resBody, "message").Str
	if resMsg != "OK" {
		return errors.New("There was some problem on our side. Sorry for incovenience!!")
	}
	return nil
}
func (c restClientStruct) Delete(id string) error {
	reqBody := map[string]string{"id": id}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodDelete, getAPIEndpoint(), bytes.NewBuffer(reqBodyJSON))
	req.Header.Add("Authorization", "Bearer "+c.IdToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	resMsg := gjson.GetBytes(resBody, "message").Str
	if resMsg != "OK" {
		return errors.New("There was some problem on our side. Sorry for incovenience!!")
	}
	return nil
}

func (c restClientStruct) Get(id string) (dataSchema, error) {
	req, _ := http.NewRequest(http.MethodGet, getAPIEndpoint(), nil)
	req.Header.Add("Authorization", "Bearer "+c.IdToken)
	req.URL.Query().Add("id", id)

	res, err := c.httpClient.Do(req)
	if err == nil {
		todoBytes, _ := ioutil.ReadAll(res.Body)
		todoParsed := gjson.ParseBytes(todoBytes)
		todo := todoParsed.Value().(dataSchema)

		return todo, nil
	} else {
		return dataSchema{}, err
	}
}

// GetAll fetch all the user todos from the serverless API.
// To identify the user, the user-authenticated OAuth2 token or a Firebase ID token.
func (c restClientStruct) GetAll() ([]dataSchema, error) {
	req, _ := http.NewRequest(http.MethodGet, getAPIEndpoint(), nil)
	req.Header.Add("Authorization", "Bearer "+c.IdToken)

	res, err := c.httpClient.Do(req)
	if err == nil {
		todos := make([]dataSchema, 0)
		todosBytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(todosBytes, &todos)

		return todos, nil
	} else {
		return nil, err
	}
}
func (c restClientStruct) Set(id string, title string, content string, last_edited string) error {
	reqBodyData := map[string]string{
		"title":   title,
		"content": content,
	}
	reqBody := map[string]interface{}{
		"id":          id,
		"data":        reqBodyData,
		"last_edited": last_edited,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPatch, getAPIEndpoint(), bytes.NewBuffer(reqBodyJSON))
	req.Header.Add("Authorization", "Bearer "+c.IdToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	resBody, _ := ioutil.ReadAll(res.Body)
	resMsg := gjson.GetBytes(resBody, "message").Str
	if resMsg != "OK" {
		return errors.New("There was some problem on our side. Sorry for incovenience!!")
	}
	return nil
}

var RestClient restClientStruct = restClientStruct{
	IdToken: credential.ReadCredentials().IdToken,
}
