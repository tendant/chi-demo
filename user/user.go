package user

import "net/http"

type UserImpl struct{}

func (*UserImpl) PostLogin(w http.ResponseWriter, r *http.Request) *Response {
	// Implement me
	return &Response{
		body:        "ok",
		Code:        200,
		contentType: "application/json",
	}
}

func (*UserImpl) PostRegister(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{
		body:        "ok",
		Code:        200,
		contentType: "application/json",
	}
}

func (*UserImpl) PostResetPassword(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{
		body:        "ok",
		Code:        200,
		contentType: "application/json",
	}
}
