package application

import (
	"fmt"
	events "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/create_post"
	saga "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateOrderOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreatePostOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateOrderOrchestrator, error) {
	o := &CreateOrderOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *CreateOrderOrchestrator) Start(order *events.PostFront, username string) error {
	PostDto := events.PostDto{
		Id:       primitive.NewObjectID(),
		Image:    order.Image,
		Content:  order.Content,
		Date:     time.Now(),
		Username: username,
	}
	UserPost := events.UserPost{
		Dto:      PostDto,
		Username: username,
	}
	event := &events.SavePostCommand{
		Type: events.SavePost,
		Dto:  UserPost,
	}
	fmt.Println(UserPost)
	return o.commandPublisher.Publish(event)
}

func (o *CreateOrderOrchestrator) handle(reply *events.SavePostReply) {
	command := events.SavePostCommand{Dto: reply.Dto}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateOrderOrchestrator) nextCommandType(reply events.SavePostReplyType) events.SavePostCommantType {
	switch reply {
	case events.SavePostFinished:
		return events.SaveNotification
	case events.SaveNotificationsFinished:
		return events.FinishSaga
	case events.ErrorOccured:
		return events.RollbackPost
	default:
		return events.UnknownCommand
	}
}
