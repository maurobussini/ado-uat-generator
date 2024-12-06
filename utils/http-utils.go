package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Execute a GET on specific url
func HttpGet[T any](url string,
	authorizationHeaderValue string) (
	T, error) {

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", authorizationHeaderValue)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	var response T

	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return response, err
	}

	jsonErr := json.Unmarshal([]byte(body), &response)

	if jsonErr != nil {
		return response, jsonErr
	}

	return response, nil
}

func HttpPost[TRequest any, TResponse any](url string,
	authorizationHeaderValue string,
	request *TRequest) (
	TResponse, error) {

	return httpWithBody[TRequest, TResponse](
		http.MethodPost,
		url,
		authorizationHeaderValue,
		request)
}

func HttpPatch[TRequest any, TResponse any](url string,
	authorizationHeaderValue string,
	request *TRequest) (
	TResponse, error) {

	return httpWithBody[TRequest, TResponse](
		http.MethodPatch,
		url,
		authorizationHeaderValue,
		request)
}

func httpWithBody[TRequest any, TResponse any](verb string,
	url string,
	authorizationHeaderValue string,
	request *TRequest) (
	TResponse, error) {

	var response TResponse

	jsonBody, errJson := json.Marshal(request)
	if errJson != nil {
		return response, errJson
	}

	req, errReq := http.NewRequest(
		verb,
		url,
		bytes.NewBuffer(jsonBody))
	if errReq != nil {
		return response, errReq
	}

	req.Header.Set("Authorization", authorizationHeaderValue)
	req.Header.Set("Content-Type", "application/json-patch+json")

	client := &http.Client{}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return response, errResp
	}

	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)

	if errRead != nil {
		return response, errRead
	}

	errDeserialize := json.Unmarshal([]byte(body), &response)

	if errDeserialize != nil {
		return response, errDeserialize
	}

	return response, nil
}
