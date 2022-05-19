package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/api"
	mw "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/middleware"
	cfg "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	server.initUserPostsHandler()
	server.initHomepageFeedHandler()
	server.initUploadImageHandler()
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

func (server *Server) initCustomHandlers() {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	registrationHandler := api.NewRegistrationHandler(authEndpoint, userEndpoint, companyEndpoint)
	registrationHandler.Init(server.mux)
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

func (server *Server) initHomepageFeedHandler() {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	userPostsHandler := api.NewHomepageFeedHandler(userEndpoint, postEndpoint)
	userPostsHandler.Init(server.mux)
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), mw.AuthMiddleware(server.mux)))
}
