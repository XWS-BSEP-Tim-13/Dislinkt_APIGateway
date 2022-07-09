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
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

type EventsHandler struct {
	connectionClientAddress string
	postsClientAddress      string
	logger                  *logger.Logger
	tracer                  *opentracing.Tracer
}

func (handler *EventsHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/events", handler.Events)
	if err != nil {
		panic(err)
	}
}

func NewEventsHandler(connectionClientAddress, postsClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) Handler {
	return &EventsHandler{
		postsClientAddress:      postsClientAddress,
		connectionClientAddress: connectionClientAddress,
		logger:                  logger,
		tracer:                  tracer,
	}
}

func (handler *EventsHandler) Events(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("Events", *handler.tracer, r)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling event sourcing at %s\n", r.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	postsClient := services.NewPostsClient(handler.postsClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)

	connEvents, err := connectionClient.GetEvents(ctx, &connectionGw.EventRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	postEvents, err := postsClient.GetEvents(ctx, &postGw.EventRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	for _, event := range postEvents.Events {
		connEvents.Events = append(connEvents.Events, &connectionGw.Event{
			Id:        event.Id,
			Action:    event.Action,
			User:      event.User,
			Published: timestamppb.New(event.Published.AsTime()),
		})
		fmt.Printf("event: %s\n", event)
	}

	response, err := json.Marshal(connEvents)
	fmt.Printf("json response: %s\n", response)
	if err != nil {
		handler.logger.ErrorMessage("Action: ES")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	handler.logger.InfoMessage("Action: ES")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
