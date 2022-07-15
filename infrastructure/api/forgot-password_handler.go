package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"mime"
	"net/http"
)

type ForgotPasswordHandler struct {
	usersClientAddress          string
	authenticationClientAddress string
	logger                      *logger.Logger
	tracer                      *opentracing.Tracer
}

func (handler *ForgotPasswordHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/forgot-password/{email}", handler.ForgotPassword)
	if err != nil {
		panic(err)
	}
}

func NewForgotPasswordHandler(usersClientAddress, authenticationClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) Handler {
	return &ForgotPasswordHandler{
		authenticationClientAddress: authenticationClientAddress,
		usersClientAddress:          usersClientAddress,
		logger:                      logger,
		tracer:                      tracer,
	}
}

func (handler *ForgotPasswordHandler) ForgotPassword(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("ForgotPassword", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

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
	_, err = usersClient.GetByEmail(ctx, &userGw.GetRequest{Id: email})
	if err != nil {
		handler.logger.WarningMessage("User " + email + " | Action: FP")
		handler.logger.ErrorMessage("User " + email + " | Action: FP")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authClient := services.NewAuthClient(handler.authenticationClientAddress)
	_, err = authClient.ForgotPassword(ctx, &authGw.ForgotPasswordRequest{Email: email})
	if err != nil {
		handler.logger.ErrorMessage("User " + email + " | Action: FP")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	handler.logger.InfoMessage("User " + email + " | Action: FP")
}
