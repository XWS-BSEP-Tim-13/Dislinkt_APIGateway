package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"mime"
	"net/http"
)

type UsersPostsHandler struct {
	usersClientAddress string
	postsClientAddress string
	logger             *logger.Logger
}

func NewUsersPostsHandler(usersClientAddress, postsClientAddress string, logger *logger.Logger) Handler {
	return &UsersPostsHandler{
		postsClientAddress: postsClientAddress,
		usersClientAddress: usersClientAddress,
		logger:             logger,
	}
}

func (handler *UsersPostsHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/user/get-posts", handler.GetUsersPosts)
	if err != nil {
		panic(err)
	}
}

func (handler *UsersPostsHandler) GetUsersPosts(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
	rt, err := decodePostDtoBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded body: %s,%s\n", rt.IdTo, rt.Username)
	usersClient := services.NewUsersClient(handler.usersClientAddress)
	connection := &userGw.Connection{IdTo: rt.IdTo, IdFrom: rt.IdFrom}
	resp, err1 := usersClient.CheckIfUserCanReadPosts(context.TODO(), &userGw.ConnectionBody{Connection: connection})
	fmt.Printf("First response: \n")
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !resp.IsReadable {
		response, err := json.Marshal(nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}
	postsClient := services.NewPostsClient(handler.postsClientAddress)
	posts, err2 := postsClient.GetByUser(context.TODO(), &postGw.GetByUserRequest{Username: rt.Username})
	fmt.Printf("Second response: \n")
	if err2 != nil {
		handler.logger.ErrorMessage("Action: GUP")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(posts)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		handler.logger.ErrorMessage("Action: GUP")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.logger.InfoMessage("Action: GUP")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
