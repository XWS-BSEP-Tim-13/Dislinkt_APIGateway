package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type ChangePasswordPageHandler struct {
	authenticationClientAddress string
	logger                      *logger.Logger
	tracer                      *opentracing.Tracer
}

func (handler *ChangePasswordPageHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/change-password/{token}", handler.ChangePasswordPage)
	if err != nil {
		panic(err)
	}
}

func NewChangePasswordPageHandlerHandler(authenticationClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) Handler {
	return &ChangePasswordPageHandler{
		authenticationClientAddress: authenticationClientAddress,
		logger:                      logger,
		tracer:                      tracer,
	}
}

func (handler *ChangePasswordPageHandler) ChangePasswordPage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("ChangePasswordPage", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	token := pathParams["token"]
	authClient := services.NewAuthClient(handler.authenticationClientAddress)
	_, err := authClient.ChangePasswordPage(ctx, &authGw.ChangePasswordPageRequest{Token: token})
	if err != nil {
		handler.logger.ErrorMessage("Action: CP")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "https://localhost:3000/change-password/"+token, http.StatusTemporaryRedirect)
	handler.logger.InfoMessage("Action: CP")
	//w.WriteHeader(http.StatusOK)
}
