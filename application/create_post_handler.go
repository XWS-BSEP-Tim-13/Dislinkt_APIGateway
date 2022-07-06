package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	events "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/create_post"
	saga "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	userGw "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/infrastructure/grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreatePostCommandHandler struct {
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreatePostCommandHandler(publisher saga.Publisher, subscriber saga.Subscriber) (*CreatePostCommandHandler, error) {
	o := &CreatePostCommandHandler{
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreatePostCommandHandler) handle(command *events.SavePostCommand) {
	reply := events.SavePostReply{Dto: command.Dto}
	fmt.Println("Orch started")
	switch command.Type {
	case events.SavePost:
		postsEndpoint := fmt.Sprintf("%s:%s", "post_service", "8000")
		postsClient := services.NewPostsClient(postsEndpoint)
		dto := mapDtoToPostPb(&command.Dto.Dto)
		_, err := postsClient.CreatePost(context.TODO(), &postGw.NewPostRequest{Post: dto})
		if err != nil {
			fmt.Println(err, "create post error")
			return
		}
		reply.Type = events.SavePostFinished
		fmt.Println("Step 1")
	case events.SaveNotification:
		userEndpoint := fmt.Sprintf("%s:%s", "user_service", "8000")
		userClient := services.NewUserClient(userEndpoint)
		notification := &userGw.NotificationDto{
			Username: command.Dto.Username,
			Type:     0,
		}
		_, err := userClient.CreateNotification(context.TODO(), &userGw.NotificationRequest{Notification: notification})
		if err != nil {
			fmt.Println(err, "create notification error")
			return
		}
		reply.Type = events.SaveNotificationsFinished
		fmt.Println("Step 2")
	case events.FinishSaga:
		fmt.Println("Step 3")
		return
	case events.RollbackPost:
		postsEndpoint := fmt.Sprintf("%s:%s", "post_service", "8000")
		postsClient := services.NewPostsClient(postsEndpoint)
		postsClient.DeletePost(context.TODO(), &postGw.GetRequest{Id: command.Dto.Dto.Id.Hex()})
		fmt.Println("Step 4")
		return
	default:
		fmt.Println("Step unknown")
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}

func mapDtoToPostPb(dto *events.PostDto) *postGw.PostDto {
	post := &postGw.PostDto{
		Image:   dto.Image,
		Content: dto.Content,
		Date:    timestamppb.New(dto.Date),
		Id:      dto.Id.Hex(),
	}
	return post
}
