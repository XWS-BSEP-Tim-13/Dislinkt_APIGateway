package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/api"
	mw "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/middleware"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	cfg "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	connectionGw "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

const (
	serverCertFile = "cert/cert.pem"
	serverKeyFile  = "cert/key.pem"
)

func NewServer(config *cfg.Config, logger *logger.Logger) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initRegistrationHandler(logger)
	server.initAccountActivationHandler(logger)
	server.initUserPostsHandler(logger)
	server.initHomepageFeedHandler(logger)
	server.initUploadImageHandler(logger)
	server.initForgotPasswordHandler(logger)
	server.initChangePasswordPageHandler(logger)
	server.initReceiveJobOfferHandler(logger)
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	err := authGw.RegisterAuthenticationServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts)
	if err != nil {
		panic(err)
	}

	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	err = companyGw.RegisterCompanyServiceHandlerFromEndpoint(context.TODO(), server.mux, companyEndpoint, opts)
	if err != nil {
		panic(err)
	}

	connectionEndpoint := fmt.Sprintf("%s:%s", "connection_service", "8000")
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndpoint, opts)
	if err != nil {
		panic(err)
	}

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		panic(err)
	}

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err = userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initRegistrationHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	registrationHandler := api.NewRegistrationHandler(authEndpoint, userEndpoint, companyEndpoint, logger)
	registrationHandler.Init(server.mux)
}

func (server *Server) initReceiveJobOfferHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	receive := api.NewReceiveJobOfferHandler(authEndpoint, companyEndpoint, logger)
	receive.Init(server.mux)
}

func (server *Server) initAccountActivationHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	accountActivationHandler := api.NewAccountActivationHandler(authEndpoint, userEndpoint, companyEndpoint, logger)
	accountActivationHandler.Init(server.mux)
}

func (server *Server) initUploadImageHandler(logger *logger.Logger) {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	uploadImageHandler := api.NewUploadImageHandler(postEndpoint, logger)
	uploadImageHandler.Init(server.mux)
}

func (server *Server) initUserPostsHandler(logger *logger.Logger) {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	userPostsHandler := api.NewUsersPostsHandler(userEndpoint, postEndpoint, logger)
	userPostsHandler.Init(server.mux)
}

func (server *Server) initChangePasswordPageHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	changePasswordPage := api.NewChangePasswordPageHandlerHandler(authEndpoint, logger)
	changePasswordPage.Init(server.mux)
}

func (server *Server) initForgotPasswordHandler(logger *logger.Logger) {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	forgotPasswordHandler := api.NewForgotPasswordHandler(userEndpoint, authEndpoint, logger)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) initHomepageFeedHandler(logger *logger.Logger) {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	forgotPasswordHandler := api.NewHomepageFeedHandler(userEndpoint, postEndpoint, logger)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) Start(logger *logger.Logger) {
	address := fmt.Sprintf(":%s", server.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	log.Fatal(http.ServeTLS(listener, mw.AuthMiddleware(server.mux, logger), serverCertFile, serverKeyFile))
}
