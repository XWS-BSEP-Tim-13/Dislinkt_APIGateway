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
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

const (
	gatewayCertFile    = "cert/cert.pem"
	gatewayKeyFile     = "cert/key.pem"
	authCertFile       = "cert/auth-cert.pem"
	companyCertFile    = "cert/company-cert.pem"
	connectionCertFile = "cert/connection-cert.pem"
	postCertFile       = "cert/post-cert.pem"
	userCertFile       = "cert/user-cert.pem"
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
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	auth_creds, _ := credentials.NewClientTLSFromFile(authCertFile, "")
	opts_auth := []grpc.DialOption{grpc.WithTransportCredentials(auth_creds)}
	err := authGw.RegisterAuthenticationServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts_auth)
	if err != nil {
		panic(err)
	}

	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	company_creds, _ := credentials.NewClientTLSFromFile(companyCertFile, "")
	opts_company := []grpc.DialOption{grpc.WithTransportCredentials(company_creds)}
	err = companyGw.RegisterCompanyServiceHandlerFromEndpoint(context.TODO(), server.mux, companyEndpoint, opts_company)
	if err != nil {
		panic(err)
	}

	connectionEndpoint := fmt.Sprintf("%s:%s", "connection_service", "8000")
	connection_creds, _ := credentials.NewClientTLSFromFile(connectionCertFile, "")
	opts_connection := []grpc.DialOption{grpc.WithTransportCredentials(connection_creds)}
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndpoint, opts_connection)
	if err != nil {
		panic(err)
	}

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	post_creds, _ := credentials.NewClientTLSFromFile(postCertFile, "")
	opts_post := []grpc.DialOption{grpc.WithTransportCredentials(post_creds)}
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts_post)
	if err != nil {
		panic(err)
	}

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	user_creds, _ := credentials.NewClientTLSFromFile(userCertFile, "")
	opts_user := []grpc.DialOption{grpc.WithTransportCredentials(user_creds)}
	err = userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts_user)
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
	log.Fatal(http.ServeTLS(listener, mw.AuthMiddleware(server.mux, logger), gatewayCertFile, gatewayKeyFile))
}
