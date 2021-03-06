package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/jwt"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"mime"
	"net/http"
)

type ReceiveJobOffer struct {
	authClientAddress    string
	companyClientAddress string
	logger               *logger.Logger
	tracer               *opentracing.Tracer
}

func NewReceiveJobOfferHandler(authClientAddress, companyClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) *ReceiveJobOffer {
	return &ReceiveJobOffer{
		authClientAddress:    authClientAddress,
		companyClientAddress: companyClientAddress,
		logger:               logger,
		tracer:               tracer,
	}
}

func (handler *ReceiveJobOffer) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/receive-job-offer", handler.ReceiveJobOfferHandler)
	if err != nil {
		panic(err)
	}
}

func (handler *ReceiveJobOffer) ReceiveJobOfferHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("ReceiveJobOfferHandler", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Media type error")
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		fmt.Println("App type error")
		return
	}
	fmt.Println(r.Body)
	jobTokenDto, err := decodeJobOfferDtoBody(r.Body)
	if err != nil {
		handler.logger.ErrorMessage("Action: DcdJO")
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	_, claims, err := jwt.ParseJwtWithEmail(jobTokenDto.Token)
	if err != nil {
		handler.logger.ErrorMessage("Action: PJWTC")
		fmt.Println("Parse claims error")
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)

	authClient := services.NewAuthClient(handler.authClientAddress)
	resp, err := authClient.CheckIfUserExist(ctx, &authGw.CheckIfUserExistsRequest{Username: claims.Email})
	if !resp.Exists {
		return
	}
	fmt.Printf("Received first response \n")
	companyClient := services.NewCompanyClient(handler.companyClientAddress)
	pb := mapJobDomainToPb(jobTokenDto.JobOffer)
	_, err = companyClient.CreateJobOffer(ctx, &companyGw.JobOfferRequest{Dto: pb})
	if err != nil {
		handler.logger.ErrorMessage("Company: " + claims.Email + " | Action: CJO")
		return
	}

	response, err := json.Marshal(true)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.logger.InfoMessage("Company: " + claims.Email + " | Action: CJO")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
