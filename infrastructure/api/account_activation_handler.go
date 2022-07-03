package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	auth "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	company "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	user "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type AccountActivationHandler struct {
	authClientAddress    string
	userClientAddress    string
	companyClientAddress string
	logger               *logger.Logger
	tracer               *opentracing.Tracer
}

func NewAccountActivationHandler(authClientAddress, userClientAddress, companyClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) *AccountActivationHandler {
	return &AccountActivationHandler{
		authClientAddress:    authClientAddress,
		userClientAddress:    userClientAddress,
		companyClientAddress: companyClientAddress,
		logger:               logger,
		tracer:               tracer,
	}
}

func (handler *AccountActivationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/activate/{code}", handler.HandleActivateAccount)
	if err != nil {
		panic(err)
	}
}

func (handler *AccountActivationHandler) HandleActivateAccount(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("HandleActivateAccount", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	code := pathParams["code"]

	authClient := services.NewAuthClient(handler.authClientAddress)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	response, err := authClient.ActivateAccount(ctx, &auth.ActivateAccountRequest{
		Code: code,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if response.ActivatedAccount.Role == "USER" {
		userClient := services.NewUserClient(handler.userClientAddress)
		_, err := userClient.ActivateAccount(ctx, &user.ActivateAccountRequest{
			Email: response.ActivatedAccount.Email,
		})
		if err != nil {
			handler.logger.ErrorMessage("User: " + response.ActivatedAccount.Email + " | Action: AA")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	} else if response.ActivatedAccount.Role == "COMPANY" {
		companyClient := services.NewCompanyClient(handler.companyClientAddress)
		_, err := companyClient.ActivateAccount(ctx, &company.ActivateAccountRequest{
			Email: response.ActivatedAccount.Email,
		})
		if err != nil {
			handler.logger.ErrorMessage("Company: " + response.ActivatedAccount.Email + " | Action: AA")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	handler.logger.InfoMessage("User: " + response.ActivatedAccount.Email + " | Action: AA")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response.ActivatedAccount.Message))
}
