package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	auth "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	company "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	user "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
)

type RegistrationHandler struct {
	authClientAddress    string
	userClientAddress    string
	companyClientAddress string
	logger               *logger.Logger
	tracer               *opentracing.Tracer
}

func NewRegistrationHandler(authClientAddress, userClientAddress, companyClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) *RegistrationHandler {
	return &RegistrationHandler{
		authClientAddress:    authClientAddress,
		userClientAddress:    userClientAddress,
		companyClientAddress: companyClientAddress,
		logger:               logger,
		tracer:               tracer,
	}
}

func (handler *RegistrationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/registration", handler.HandleRegister)
	if err != nil {
		panic(err)
	}
}

func (handler *RegistrationHandler) HandleRegister(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("HandleRegister", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	registerRequestJson, err := decodeBodyToRegisterRequest(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Unable to decode request body!"))
		handler.logger.ErrorMessage("Action: RU")
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)

	var response []byte
	if !registerRequestJson.IsCompany {
		registerUserRequestPb := mapRegisterUserRequestPb(registerRequestJson)

		userClient := services.NewUserClient(handler.userClientAddress)
		user, err := userClient.CreateUser(ctx, registerUserRequestPb)
		if err != nil {
			handler.logger.ErrorMessage("Action: RU")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		response, err = json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Unable to register!"))
			return
		}
	} else {
		registerCompanyRequestPb := mapRegisterCompanyRequestPb(registerRequestJson)

		companyClient := services.NewCompanyClient(handler.companyClientAddress)
		newCompany, err := companyClient.CreateCompany(ctx, registerCompanyRequestPb)
		if err != nil {
			handler.logger.ErrorMessage("Action: RC")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		response, err = json.Marshal(newCompany)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Unable to register!"))
			return
		}
	}

	registerAuthRequestPb := mapRegisterAuthRequestPb(registerRequestJson)

	authClient := services.NewAuthClient(handler.authClientAddress)
	_, err = authClient.Register(ctx, registerAuthRequestPb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	handler.logger.InfoMessage("User: " + registerAuthRequestPb.User.Username + " | Action: Registration ")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func mapRegisterAuthRequestPb(registerRequestJson *domain.RegisterRequest) *auth.RegisterRequest {
	registerAuthRequestPb := &auth.RegisterRequest{
		User: &auth.User{
			Username: (*registerRequestJson).Username,
			Password: (*registerRequestJson).Password,
			Role:     (*registerRequestJson).Role,
			Email:    (*registerRequestJson).Email,
		},
	}
	return registerAuthRequestPb
}

func mapRegisterUserRequestPb(registerRequestJson *domain.RegisterRequest) *user.NewUser {
	registerUserRequestPb := &user.NewUser{
		User: &user.User{
			Username:    registerRequestJson.Username,
			FirstName:   registerRequestJson.FirstName,
			LastName:    registerRequestJson.LastName,
			Email:       registerRequestJson.Email,
			PhoneNumber: registerRequestJson.PhoneNumber,
			Gender:      user.User_Gender(registerRequestJson.Gender),
			DateOfBirth: timestamppb.New(registerRequestJson.DateOfBirth),
			Biography:   registerRequestJson.Biography,
			IsPrivate:   registerRequestJson.IsPrivate,
			Educations:  []*user.Education{},
			Experiences: []*user.Experience{},
			Skills:      []string{},
			Interests:   []string{},
			Connections: []string{},
		},
	}
	return registerUserRequestPb
}

func mapRegisterCompanyRequestPb(registerRequestJson *domain.RegisterRequest) *company.NewCompany {
	registerCompanyRequestPb := &company.NewCompany{
		Company: &company.Company{
			Username:    registerRequestJson.Username,
			Email:       registerRequestJson.Email,
			PhoneNumber: registerRequestJson.PhoneNumber,
			CompanyName: registerRequestJson.CompanyName,
			Description: registerRequestJson.Description,
			Location:    registerRequestJson.Location,
			Website:     registerRequestJson.Website,
			CompanySize: registerRequestJson.CompanySize,
			Industry:    registerRequestJson.Industry,
		},
	}
	return registerCompanyRequestPb
}

func decodeBodyToRegisterRequest(r io.Reader) (*domain.RegisterRequest, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var registerRequest domain.RegisterRequest
	if err := dec.Decode(&registerRequest); err != nil {
		return nil, err
	}
	return &registerRequest, nil
}
