package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type ChangePasswordPageHandler struct {
	authenticationClientAddress string
}

func (handler *ChangePasswordPageHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/change-password/{token}", handler.ChangePasswordPage)
	if err != nil {
		panic(err)
	}
}

func NewChangePasswordPageHandlerHandler(authenticationClientAddress string) Handler {
	return &ChangePasswordPageHandler{
		authenticationClientAddress: authenticationClientAddress,
	}
}

func (handler *ChangePasswordPageHandler) ChangePasswordPage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Println("Request started")
	token := pathParams["token"]
	//contentType := r.Header.Get("Content-Type")
	//mediatype, _, err := mime.ParseMediaType(contentType)
	//if err != nil {
	//	fmt.Printf("#################################### token: %s\n", token)
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//if mediatype != "application/json" {
	//	http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
	//	return
	//}

	authClient := services.NewAuthClient(handler.authenticationClientAddress)
	_, err := authClient.ChangePasswordPage(context.TODO(), &authGw.ChangePasswordPageRequest{Token: token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "https://localhost:3000/change-password/"+token, http.StatusTemporaryRedirect)
	//w.WriteHeader(http.StatusOK)
}
