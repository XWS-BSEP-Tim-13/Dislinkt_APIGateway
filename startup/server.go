package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/api"
	mw "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/middleware"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	saga "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging/nats"
	cfg "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	connectionGw "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
	tracer opentracing.Tracer
	closer io.Closer
}

const (
	gatewayCertFile    = "cert/cert.pem"
	gatewayKeyFile     = "cert/key.pem"
	authCertFile       = "cert/auth-cert.pem"
	companyCertFile    = "cert/company-cert.pem"
	connectionCertFile = "cert/connection-cert.pem"
	postCertFile       = "cert/post-cert.pem"
	userCertFile       = "cert/user-cert.pem"
	QueueGroup         = "create_post_service"
)

func NewServer(config *cfg.Config, logger *logger.Logger) *Server {
	tracer, closer := tracer.Init()
	opentracing.SetGlobalTracer(tracer)

	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
		tracer: tracer,
		closer: closer,
	}
	server.initHandlers()
	server.initRegistrationHandler(logger)

	commandPublisher := server.initPublisher(server.config.CreatePostCommandSubject)
	replySubscriber := server.initSubscriber(server.config.CreatePostReplySubject, QueueGroup)
	createPostOrchestrator := server.initCreatePostOrchestrator(commandPublisher, replySubscriber)

	commandSubscriber := server.initSubscriber(server.config.CreatePostCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreatePostReplySubject)
	server.initCreatePostHandler(replyPublisher, commandSubscriber)

	server.initAccountActivationHandler(logger)
	server.initUserPostsHandler(logger)
	server.initCreatePostApiHandler(logger, createPostOrchestrator)
	server.initHomepageFeedHandler(logger)
	server.initMessageSendHandler(logger)
	server.initUploadImageHandler(logger)
	server.initForgotPasswordHandler(logger)
	server.initChangePasswordPageHandler(logger)
	server.initReceiveJobOfferHandler(logger)
	return server
}

var grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || err == opentracing.ErrSpanContextNotFound {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}

func (server *Server) initHandlers() {
	//opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	//authCreds, _ := credentials.NewClientTLSFromFile(authCertFile, "")
	optsAuth := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}
	err := authGw.RegisterAuthenticationServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, optsAuth)
	if err != nil {
		panic(err)
	}

	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	//companyCreds, _ := credentials.NewClientTLSFromFile(companyCertFile, "")
	optsCompany := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}
	err = companyGw.RegisterCompanyServiceHandlerFromEndpoint(context.TODO(), server.mux, companyEndpoint, optsCompany)
	if err != nil {
		panic(err)
	}

	connectionEndpoint := fmt.Sprintf("%s:%s", "connection_service", "8000")
	//connectionCreds, _ := credentials.NewClientTLSFromFile(connectionCertFile, "")
	optsConnection := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndpoint, optsConnection)
	if err != nil {
		panic(err)
	}

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	//postCreds, _ := credentials.NewClientTLSFromFile(postCertFile, "")
	optsPost := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, optsPost)
	if err != nil {
		panic(err)
	}

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	//userCreds, _ := credentials.NewClientTLSFromFile(userCertFile, "")
	optsUser := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}
	err = userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, optsUser)
	if err != nil {
		panic(err)
	}
}

func (server *Server) initRegistrationHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	registrationHandler := api.NewRegistrationHandler(authEndpoint, userEndpoint, companyEndpoint, logger, &server.tracer)
	registrationHandler.Init(server.mux)
}

func (server *Server) initReceiveJobOfferHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	receive := api.NewReceiveJobOfferHandler(authEndpoint, companyEndpoint, logger, &server.tracer)
	receive.Init(server.mux)
}

func (server *Server) initAccountActivationHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", "auth_service", "8000")
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	companyEndpoint := fmt.Sprintf("%s:%s", "company_service", "8000")
	accountActivationHandler := api.NewAccountActivationHandler(authEndpoint, userEndpoint, companyEndpoint, logger, &server.tracer)
	accountActivationHandler.Init(server.mux)
}

func (server *Server) initMessageSendHandler(logger *logger.Logger) {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	accountActivationHandler := api.NewSendMessageHandler(postEndpoint, userEndpoint, logger, &server.tracer)
	accountActivationHandler.Init(server.mux)
}

func (server *Server) initUploadImageHandler(logger *logger.Logger) {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	uploadImageHandler := api.NewUploadImageHandler(postEndpoint, logger, &server.tracer)
	uploadImageHandler.Init(server.mux)
}

func (server *Server) initUserPostsHandler(logger *logger.Logger) {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	userPostsHandler := api.NewUsersPostsHandler(userEndpoint, postEndpoint, logger, &server.tracer)
	userPostsHandler.Init(server.mux)
}

func (server *Server) initChangePasswordPageHandler(logger *logger.Logger) {
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	changePasswordPage := api.NewChangePasswordPageHandlerHandler(authEndpoint, logger, &server.tracer)
	changePasswordPage.Init(server.mux)
}

func (server *Server) initForgotPasswordHandler(logger *logger.Logger) {
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	forgotPasswordHandler := api.NewForgotPasswordHandler(userEndpoint, authEndpoint, logger, &server.tracer)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) initHomepageFeedHandler(logger *logger.Logger) {
	connectionEndpoint := fmt.Sprintf("%s:%s", "connection_service", "8000")
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	forgotPasswordHandler := api.NewHomepageFeedHandler(connectionEndpoint, postEndpoint, logger, &server.tracer)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) initCreatePostApiHandler(logger *logger.Logger, orchestrator *application.CreateOrderOrchestrator) {
	forgotPasswordHandler := api.NewCreatePostHandler(logger, orchestrator)
	forgotPasswordHandler.Init(server.mux)
}

func (server *Server) initCreatePostHandler(publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := application.NewCreatePostCommandHandler(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initCreatePostOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *application.CreateOrderOrchestrator {
	orchestrator, err := application.NewCreatePostOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) Start(logger *logger.Logger) {
	address := fmt.Sprintf(":%s", server.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	log.Fatal(http.ServeTLS(listener, mw.AuthMiddleware(tracingWrapper(server.mux), logger), gatewayCertFile, gatewayKeyFile))
}

func (server *Server) GetTracer() opentracing.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) CloseTracer() error {
	return server.closer.Close()
}
