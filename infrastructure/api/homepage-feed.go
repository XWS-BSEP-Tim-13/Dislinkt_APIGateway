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

type HomepageFeedHandler struct {
	usersClientAddress string
	postsClientAddress string
	logger             *logger.Logger
}

func (handler *HomepageFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/user/feed", handler.HomepageFeed)
	if err != nil {
		panic(err)
	}
}

func NewHomepageFeedHandler(usersClientAddress, postsClientAddress string, logger *logger.Logger) Handler {
	return &HomepageFeedHandler{
		postsClientAddress: postsClientAddress,
		usersClientAddress: usersClientAddress,
		logger:             logger,
	}
}

func (handler *HomepageFeedHandler) HomepageFeed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
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
	rt, err := decodeHomepageFeedBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded body: %s\n", rt.Username)
	postsClient := services.NewPostsClient(handler.postsClientAddress)
	usersClient := services.NewUsersClient(handler.usersClientAddress)
	resp, err := usersClient.GetUsernames(context.TODO(), &userGw.ConnectionResponse{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	u := &postGw.Usernames{Username: resp.Usernames}
	respPosts, err := postsClient.GetFeedPosts(context.TODO(), &postGw.FeedRequest{Page: int64(rt.Page), Usernames: u})
	if err != nil {
		handler.logger.ErrorMessage("User: " + rt.Username + " | Action: HPF")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(respPosts)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		handler.logger.ErrorMessage("User: " + rt.Username + " | Action: HPF")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.logger.InfoMessage("User: " + rt.Username + " | Action: HPF")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
