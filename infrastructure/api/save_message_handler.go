package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/tracer"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type SendMessageHandler struct {
	postsClientAddress string
	logger             *logger.Logger
	tracer             *opentracing.Tracer
	usersClientAddress string
}

func (handler *SendMessageHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/message", handler.SendMessage)
	if err != nil {
		panic(err)
	}
}

func NewSendMessageHandler(postsClientAddress, usersClientAddress string, logger *logger.Logger, tracer *opentracing.Tracer) Handler {
	return &SendMessageHandler{
		postsClientAddress: postsClientAddress,
		usersClientAddress: usersClientAddress,
		logger:             logger,
		tracer:             tracer,
	}
}

func (handler *SendMessageHandler) SendMessage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	span := tracer.StartSpanFromRequest("UploadImage", *handler.tracer, r)
	defer span.Finish()
	fmt.Println("Send message started")
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling post create at %s\n", r.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)

	//err := r.ParseMultipartForm(32 << 20)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	messageDto, err := decodeMessageDtoBody(r.Body)
	if err != nil {
		handler.logger.ErrorMessage("Action: DcdJO")
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Println(messageDto)

	postsClient := services.NewPostsClient(handler.postsClientAddress)
	messagePb := &postGw.MessageDto{
		Content:     messageDto.Content,
		MessageFrom: messageDto.MessageFrom,
		MessageTo:   messageDto.MessageTo,
	}
	resp, err := postsClient.SaveMessage(ctx, &postGw.SaveMessageRequest{Message: messagePb})
	if err != nil {
		fmt.Println(err)
	}
	usersClient := services.NewUsersClient(handler.usersClientAddress)
	_, err = usersClient.MessageNotification(ctx, &userGw.Connection{IdTo: resp.MessageTo, IdFrom: resp.MessageFrom})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	response, _ := json.Marshal(resp)

	handler.logger.InfoMessage("FU")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
