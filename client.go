package main

import (
	"fmt"
	"net/http"

	"github.com/osuAkatsuki/otp-service-client-go/internal/http_client"
)

type OtpClient struct {
	BaseUrl string
	Secret  string
}

func NewOtpClient(baseUrl, secret string) OtpClient {
	return OtpClient{baseUrl, secret}
}

func handleResponse(resp http_client.HttpResponse) error {
	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if resp.HasError {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return &BadRequestError{resp.ErrorBody.Problem}
		case http.StatusConflict:
			return &ConflictError{resp.ErrorBody.Problem}
		default:
			return &UnknownError{resp.ErrorBody.Problem}
		}
	}

	return nil
}

func handleResponseWithBody[T any](resp http_client.HttpResponseWithBody[T]) (T, error) {
	var def T
	err := handleResponse(resp.HttpResponse)
	if err != nil {
		return def, err
	}

	return resp.Body, nil
}

func prepareRequest(oc *OtpClient, request http_client.HttpRequestWithHeaders) {
	request.AddHeader("X-Secret", oc.Secret)
}

func getRequest[T any](oc *OtpClient, request http_client.HttpRequest) (T, error) {
	prepareRequest(oc, &request)

	var def T
	resp, err := http_client.Get[T](request)
	if err != nil {
		return def, err
	}

	return handleResponseWithBody[T](resp)
}

func postRequest[T any](oc *OtpClient, request http_client.HttpRequest) (T, error) {
	prepareRequest(oc, &request)

	var def T
	resp, err := http_client.Post[T](request)
	if err != nil {
		return def, err
	}

	return handleResponseWithBody[T](resp)
}

func postRequestWithNoContent(oc *OtpClient, request http_client.HttpRequest) error {
	prepareRequest(oc, &request)

	resp, err := http_client.PostWithNoContent(request)
	if err != nil {
		return err
	}

	return handleResponse(resp)
}

func postRequestWithBodyWithNoContent[T any](oc *OtpClient, request http_client.HttpRequestWithBody[T]) error {
	prepareRequest(oc, &request)

	resp, err := http_client.PostWithBodyWithNoContent(request)
	if err != nil {
		return err
	}

	return handleResponse(resp)
}

func deleteRequestWithNoContent(oc *OtpClient, request http_client.HttpRequest) error {
	prepareRequest(oc, &request)

	resp, err := http_client.DeleteWithNoContent(request)
	if err != nil {
		return err
	}

	return handleResponse(resp)
}

type GetUserOtpResponse struct {
	Verified bool   `json:"verified"`
	Enabled  bool   `json:"enabled"`
	Secret   string `json:"secret"`
	AuthUrl  string `json:"auth_url"`
}

func (oc *OtpClient) GetUserOtp(userId int) (GetUserOtpResponse, error) {
	req := http_client.HttpRequest{
		Url: oc.BaseUrl + fmt.Sprintf("/users/%d/otp", userId),
	}

	resp, err := getRequest[GetUserOtpResponse](oc, req)
	if err != nil {
		return GetUserOtpResponse{}, err
	}

	return resp, nil
}

type CreateUserOtpResponse struct {
	Secret  string `json:"secret"`
	AuthUrl string `json:"auth_url"`
}

func (oc *OtpClient) CreateUserOtp(userId int) (CreateUserOtpResponse, error) {
	req := http_client.HttpRequest{
		Url: oc.BaseUrl + fmt.Sprintf("/users/%d/otp", userId),
	}

	resp, err := postRequest[CreateUserOtpResponse](oc, req)
	if err != nil {
		return CreateUserOtpResponse{}, err
	}

	return resp, nil
}

func (oc *OtpClient) DisableUserOtp(userId int) error {
	req := http_client.HttpRequest{
		Url: oc.BaseUrl + fmt.Sprintf("/users/%d/otp/disable", userId),
	}

	err := postRequestWithNoContent(oc, req)
	if err != nil {
		return err
	}

	return nil
}

func (oc *OtpClient) DeleteUserOtp(userId int) error {
	req := http_client.HttpRequest{
		Url: oc.BaseUrl + fmt.Sprintf("/users/%d/otp", userId),
	}

	err := deleteRequestWithNoContent(oc, req)
	if err != nil {
		return err
	}

	return nil
}

type VerifyOtpRequest struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (oc *OtpClient) VerifyOtp(userId int, token string) error {
	req := http_client.HttpRequestWithBody[VerifyOtpRequest]{
		HttpRequest: http_client.HttpRequest{
			Url: oc.BaseUrl + "/otp/verify",
		},
		Body: VerifyOtpRequest{
			UserId: userId,
			Token:  token,
		},
	}

	err := postRequestWithBodyWithNoContent[VerifyOtpRequest](oc, req)
	if err != nil {
		return err
	}

	return nil
}

type ValidateOtpRequest struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func (oc *OtpClient) ValidateOtp(userId int, token string) error {
	req := http_client.HttpRequestWithBody[ValidateOtpRequest]{
		HttpRequest: http_client.HttpRequest{
			Url: oc.BaseUrl + "/otp/validate",
		},
		Body: ValidateOtpRequest{
			UserId: userId,
			Token:  token,
		},
	}

	err := postRequestWithBodyWithNoContent[ValidateOtpRequest](oc, req)
	if err != nil {
		return err
	}

	return nil
}
