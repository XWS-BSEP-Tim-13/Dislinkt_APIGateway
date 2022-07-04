package api

import (
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/application"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"mime"
	"net/http"
)

type CreatePostHandler struct {
	usersClientAddress string
	postsClientAddress string
	logger             *logger.Logger
	orchestrator       *application.CreateOrderOrchestrator
}

func (handler *CreatePostHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/forgot-password/{email}", handler.CreatePost)
	if err != nil {
		panic(err)
	}
}

func NewCreatePostHandler(logger *logger.Logger, orchestrator *application.CreateOrderOrchestrator) Handler {
	return &CreatePostHandler{
		logger:       logger,
		orchestrator: orchestrator,
	}
}

func (handler *CreatePostHandler) CreatePost(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Println("Request started")
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
	postDto, err := decodeCreatePostBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, claims, err := jwt.ParseJwt(r.Header.Get("Authorization"))
	if err != nil {
		handler.logger.ErrorMessage("Action: PJWTC")
		fmt.Println("Parse claims error")
		return
	}
	err = handler.orchestrator.Start(postDto, claims.Username)
	if err != nil {
		fmt.Println("Bad request jwt")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(true)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
