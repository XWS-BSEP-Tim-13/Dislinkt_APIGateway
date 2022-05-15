package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"mime"
	"net/http"
)

type HomepageFeedHandler struct {
	usersClientAddress string
	postsClientAddress string
}

func (handler *HomepageFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/user/feed", handler.HomepageFeed)
	if err != nil {
		panic(err)
	}
}

func NewHomepageFeedHandler(usersClientAddress, postsClientAddress string) Handler {
	return &HomepageFeedHandler{
		postsClientAddress: postsClientAddress,
		usersClientAddress: usersClientAddress,
	}
}

func (handler *HomepageFeedHandler) HomepageFeed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
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
	usersClient := services.NewUsersClient(handler.usersClientAddress)
	resp, err := usersClient.GetConnectionUsernamesForUser(context.TODO(), &userGw.UserUsername{Username: rt.Username})
	fmt.Printf("First response: \n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postsClient := services.NewPostsClient(handler.postsClientAddress)
	u := &postGw.Usernames{Username: resp.Usernames}
	respPosts, err := postsClient.GetFeedPosts(context.TODO(), &postGw.FeedRequest{Page: int64(rt.Page), Usernames: u})
	fmt.Printf("Second response: %s\n", respPosts.LastPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(respPosts)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
