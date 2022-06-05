package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/api"
	mw "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/middleware"
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

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initRegistrationHandler()
	server.initAccountActivationHandler()
	server.initUserPostsHandler()
	server.initHomepageFeedHandler()
	server.initUploadImageHandler()
	server.initForgotPasswordHandler()
	server.initChangePasswordPageHandler()
	server.initReceiveJobOfferHandler()
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

func (server *Server) initRegistrationHandler() {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	registrationHandler := api.NewRegistrationHandler(authEndpoint, userEndpoint, companyEndpoint)
	registrationHandler.Init(server.mux)
}

func (server *Server) initReceiveJobOfferHandler() {
	//authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	//companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	////receive := api.NewReceiveJobOfferHandler(authEndpoint, companyEndpoint)
	//handler := func(err error, msg string) { fmt.Println("AMQ MSG:", err, msg) }
	//if err := services.NewActiveMQ("activemq:61613").Subscribe("jobOffer.queue", handler); err != nil {
	//	fmt.Println("AMQ ERROR:", err)
	//}
	//receive.HandleReceive()
}

func (server *Server) initAccountActivationHandler() {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	accountActivationHandler := api.NewAccountActivationHandler(authEndpoint, userEndpoint, companyEndpoint)
	accountActivationHandler.Init(server.mux)
}

func (server *Server) initUploadImageHandler() {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	uploadImageHandler := api.NewUploadImageHandler(postEndpoint)
	uploadImageHandler.Init(server.mux)
}

func (server *Server) initUserPostsHandler() {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	userPostsHandler := api.NewUsersPostsHandler(userEndpoint, postEndpoint)
	userPostsHandler.Init(server.mux)
}

func (server *Server) initChangePasswordPageHandler() {
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	changePasswordPage := api.NewChangePasswordPageHandlerHandler(authEndpoint)
	changePasswordPage.Init(server.mux)
}

func (server *Server) initForgotPasswordHandler() {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	forgotPasswordHandler := api.NewForgotPasswordHandler(userEndpoint, authEndpoint)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) initHomepageFeedHandler() {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	forgotPasswordHandler := api.NewHomepageFeedHandler(userEndpoint, postEndpoint)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) Start() {
	address := fmt.Sprintf(":%s", server.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	log.Fatal(http.ServeTLS(listener, mw.AuthMiddleware(server.mux), serverCertFile, serverKeyFile))
}
