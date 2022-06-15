package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type ChangePasswordPageHandler struct {
	authenticationClientAddress string
	logger                      *logger.Logger
}

func (handler *ChangePasswordPageHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/change-password/{token}", handler.ChangePasswordPage)
	if err != nil {
		panic(err)
	}
}

func NewChangePasswordPageHandlerHandler(authenticationClientAddress string, logger *logger.Logger) Handler {
	return &ChangePasswordPageHandler{
		authenticationClientAddress: authenticationClientAddress,
		logger:                      logger,
	}
}

func (handler *ChangePasswordPageHandler) ChangePasswordPage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Println("Request started")
	token := pathParams["token"]
	authClient := services.NewAuthClient(handler.authenticationClientAddress)
	_, err := authClient.ChangePasswordPage(context.TODO(), &authGw.ChangePasswordPageRequest{Token: token})
	if err != nil {
		handler.logger.ErrorMessage("Action: CP")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "https://localhost:3000/change-password/"+token, http.StatusTemporaryRedirect)
	handler.logger.InfoMessage("Action: CP")
	//w.WriteHeader(http.StatusOK)
}
