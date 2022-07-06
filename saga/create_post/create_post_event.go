package create_post

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostDto struct {
	Id       primitive.ObjectID
	Content  string
	Image    string
	Date     time.Time
	Username string
}

type PostFront struct {
	Content string
	Image   string
}

type UserPost struct {
	Dto      PostDto
	Username string
}

type SavePostCommantType int8

const (
	SavePost SavePostCommantType = iota
	SaveNotification
	RollbackPost
	UnknownCommand
	FinishSaga
)

type SavePostReplyType int8

const (
	SavePostFinished SavePostReplyType = iota
	SaveNotificationsFinished
	ErrorOccured
	UnknownReply
)

type SavePostCommand struct {
	Dto  UserPost
	Type SavePostCommantType
}

type SavePostReply struct {
	Dto  UserPost
	Type SavePostReplyType
}
