package http_client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HttpRequestWithHeaders interface {
	AddHeader(key, value string)
}

type HttpRequest struct {
	Url             string
	QueryParameters map[string]string
	Headers         map[string]string
}

func (r *HttpRequest) AddHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}

	r.Headers[key] = value
}

type HttpRequestWithBody[T any] struct {
	HttpRequest
	Body T
}

type ErrorBody struct {
	Problem string `json:"problem"`
}

type HttpResponse struct {
	StatusCode int
	Headers    map[string][]string
	HasError   bool
	ErrorBody  ErrorBody
}

type HttpResponseWithBody[T any] struct {
	HttpResponse
	Body T
}

const UserAgent = "otp-service-client-go"

func Get[T any](request HttpRequest) (HttpResponseWithBody[T], error) {
	req, err := http.NewRequest(http.MethodGet, request.Url, nil)
	if err != nil {
		return HttpResponseWithBody[T]{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponseWithBody[T]{}, err
	}

	response := HttpResponseWithBody[T]{
		HttpResponse: HttpResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		},
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	jsonBody, err := parseJson[T](body)
	if err != nil {
		return response, err
	}
	response.Body = jsonBody

	return response, nil
}

func Post[T any](request HttpRequest) (HttpResponseWithBody[T], error) {
	req, err := http.NewRequest(http.MethodPost, request.Url, nil)
	if err != nil {
		return HttpResponseWithBody[T]{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponseWithBody[T]{}, err
	}

	response := HttpResponseWithBody[T]{
		HttpResponse: HttpResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		},
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	jsonBody, err := parseJson[T](body)
	if err != nil {
		return response, err
	}
	response.Body = jsonBody

	return response, nil
}

func PostWithBody[T any, T1 any](request HttpRequestWithBody[T]) (HttpResponseWithBody[T1], error) {
	byteData, err := json.Marshal(request.Body)
	if err != nil {
		return HttpResponseWithBody[T1]{}, err
	}

	byteReader := bytes.NewReader(byteData)

	req, err := http.NewRequest(http.MethodPost, request.Url, byteReader)
	if err != nil {
		return HttpResponseWithBody[T1]{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponseWithBody[T1]{}, err
	}

	response := HttpResponseWithBody[T1]{
		HttpResponse: HttpResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		},
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	jsonBody, err := parseJson[T1](body)
	if err != nil {
		return response, err
	}
	response.Body = jsonBody

	return response, nil
}

func PostWithNoContent(request HttpRequest) (HttpResponse, error) {
	req, err := http.NewRequest(http.MethodPost, request.Url, nil)
	if err != nil {
		return HttpResponse{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponse{}, err
	}

	response := HttpResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	return response, nil
}

func PostWithBodyWithNoContent[T any](request HttpRequestWithBody[T]) (HttpResponse, error) {
	byteData, err := json.Marshal(request.Body)
	if err != nil {
		return HttpResponse{}, err
	}

	byteReader := bytes.NewReader(byteData)

	req, err := http.NewRequest(http.MethodPost, request.Url, byteReader)
	if err != nil {
		return HttpResponse{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponse{}, err
	}

	response := HttpResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	return response, nil
}

func DeleteWithNoContent(request HttpRequest) (HttpResponse, error) {
	req, err := http.NewRequest(http.MethodDelete, request.Url, nil)
	if err != nil {
		return HttpResponse{}, err
	}

	q := req.URL.Query()

	for queryParameter, queryValue := range request.QueryParameters {
		q.Add(queryParameter, queryValue)
	}

	req.URL.RawQuery = q.Encode()

	for headerKey, headerValue := range request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	req.Header.Add("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpResponse{}, err
	}

	response := HttpResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}

	if (response.StatusCode < 200 || response.StatusCode > 299) && response.StatusCode != http.StatusNotFound {
		errorJson, err := parseJson[ErrorBody](body)
		if err != nil {
			return response, err
		}

		response.ErrorBody = errorJson
		response.HasError = true
	}

	if response.StatusCode == http.StatusNotFound {
		return response, nil
	}

	return response, nil
}

func parseJson[T any](s []byte) (T, error) {
	var body T

	err := json.Unmarshal(s, &body)
	return body, err
}
