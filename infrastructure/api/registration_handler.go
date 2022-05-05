package api

import (
	"context"
	"encoding/json"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	auth "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io"
	"net/http"
)

type RegistrationHandler struct {
	authClientAddress    string
	userClientAddress    string
	companyClientAddress string
}

func NewRegistrationHandler(authClientAddress, userClientAddress, companyClientAddress string) *RegistrationHandler {
	return &RegistrationHandler{
		authClientAddress:    authClientAddress,
		userClientAddress:    userClientAddress,
		companyClientAddress: companyClientAddress,
	}
}

func (handler *RegistrationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/registration", handler.GetDetails)
	if err != nil {
		panic(err)
	}
}

func (handler *RegistrationHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	registerRequest, err := decodeBodyToRegisterRequest(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authClient := services.NewAuthClient(handler.authClientAddress)
	username, err := authClient.Register(context.TODO(), registerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Unable to register!"))
		return
	}

	response, err := json.Marshal(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func decodeBodyToRegisterRequest(r io.Reader) (*auth.RegisterRequest, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var registerRequest auth.RegisterRequest
	if err := dec.Decode(&registerRequest); err != nil {
		return nil, err
	}
	return &registerRequest, nil
}
