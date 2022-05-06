package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	auth "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	"github.com/golang/protobuf/jsonpb"
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
	registerRequestJson, err := decodeBodyToRegisterRequest(r.Body)
	fmt.Println(registerRequestJson)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Unable to decode request body!"))
		return
	}

	//var registerRequestPb *auth.RegisterRequest
	registerRequestPb := auth.RegisterRequest{}
	jsonpb.Unmarshal(r.Body, &registerRequestPb)
	
	authClient := services.NewAuthClient(handler.authClientAddress)
	username, err := authClient.Register(context.TODO(), &registerRequestPb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Unable to connect on auth service!"))
		return
	}

	response, err := json.Marshal(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Unable to register!"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func decodeBodyToRegisterRequest(r io.Reader) (*User, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var registerRequest User
	if err := dec.Decode(&registerRequest); err != nil {
		return nil, err
	}
	return &registerRequest, nil
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
