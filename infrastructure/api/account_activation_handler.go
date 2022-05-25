package api

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	auth "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	company "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	user "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type AccountActivationHandler struct {
	authClientAddress    string
	userClientAddress    string
	companyClientAddress string
}

func NewAccountActivationHandler(authClientAddress, userClientAddress, companyClientAddress string) *AccountActivationHandler {
	return &AccountActivationHandler{
		authClientAddress:    authClientAddress,
		userClientAddress:    userClientAddress,
		companyClientAddress: companyClientAddress,
	}
}

func (handler *AccountActivationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/activate/{code}", handler.HandleActivateAccount)
	if err != nil {
		panic(err)
	}
}

func (handler *AccountActivationHandler) HandleActivateAccount(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	code := pathParams["code"]

	authClient := services.NewAuthClient(handler.authClientAddress)

	response, err := authClient.ActivateAccount(context.TODO(), &auth.ActivateAccountRequest{
		Code: code,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if response.ActivatedAccount.Role == "USER" {
		userClient := services.NewUserClient(handler.userClientAddress)
		_, err := userClient.ActivateAccount(context.TODO(), &user.ActivateAccountRequest{
			Email: response.ActivatedAccount.Email,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else if response.ActivatedAccount.Role == "COMPANY" {
		companyClient := services.NewCompanyClient(handler.companyClientAddress)
		_, err := companyClient.ActivateAccount(context.TODO(), &company.ActivateAccountRequest{
			Email: response.ActivatedAccount.Email,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.ActivatedAccount.Message))
}
