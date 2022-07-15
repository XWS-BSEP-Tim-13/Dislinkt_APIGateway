package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	connectionGw "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
)

type HomepageFeedHandler struct {
	connectionClientAddress string
	postsClientAddress      string
	logger                  *logger.Logger
	tracer                  *opentracing.Tracer
}

func (handler *HomepageFeedHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/user/feed/{page}", handler.HomepageFeed)
	if err != nil {
		panic(err)
	}
}

func NewHomepageFeedHandler(connectionClientAddress, postsClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) Handler {
	return &HomepageFeedHandler{
		postsClientAddress:      postsClientAddress,
		connectionClientAddress: connectionClientAddress,
		logger:                  logger,
		tracer:                  tracer,
	}
}

func (handler *HomepageFeedHandler) HomepageFeed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("HomepageFeed", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("Request started")
	page := pathParams["page"]
	pageNum, _ := strconv.Atoi(page)
	postsClient := services.NewPostsClient(handler.postsClientAddress)
	fmt.Println(postsClient)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	resp, err := connectionClient.GetConnectionUsernamesForUser(ctx, &connectionGw.ConnectionResponse{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	u := &postGw.Usernames{Username: resp.Usernames}
	respPosts, err := postsClient.GetFeedPosts(ctx, &postGw.FeedRequest{Page: int64(pageNum), Usernames: u})

	if err != nil {
		fmt.Println(err, "Error while getting posts")
		handler.logger.ErrorMessage("User: | Action: HPF")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(respPosts)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		handler.logger.ErrorMessage("User:  Action: HPF")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.logger.InfoMessage("User:  Action: HPF")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
