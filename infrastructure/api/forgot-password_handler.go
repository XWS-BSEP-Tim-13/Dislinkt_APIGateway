package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"mime"
	"net/http"
)

type ForgotPasswordHandler struct {
	usersClientAddress          string
	authenticationClientAddress string
}

func (handler *ForgotPasswordHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/forgot-password/{email}", handler.ForgotPassword)
	if err != nil {
		panic(err)
	}
}

func NewForgotPasswordHandler(usersClientAddress, authenticationClientAddress string) Handler {
	return &ForgotPasswordHandler{
		authenticationClientAddress: authenticationClientAddress,
		usersClientAddress:          usersClientAddress,
	}
}

func (handler *ForgotPasswordHandler) ForgotPassword(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Println("Request started")
	email := pathParams["email"]
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	fmt.Printf("Decoded body: %s\n", email)
	usersClient := services.NewUsersClient(handler.usersClientAddress)
	_, err = usersClient.GetByEmail(context.TODO(), &userGw.GetRequest{Id: email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authClient := services.NewAuthClient(handler.authenticationClientAddress)
	_, err = authClient.ForgotPassword(context.TODO(), &authGw.ForgotPasswordRequest{Email: email})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
