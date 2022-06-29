package block_user

type Users struct {
	UserFrom string
	UserTo   string
}

type BlockUserCommandType int8

const (
	RemoveConnectionFromUser BlockUserCommandType = iota
	RemoveConnectionToUser
	BlockUser
	UnknownCommand
)

type BlockUserCommand struct {
	Users Users
	Type  BlockUserCommandType
}

type BlockUserReplyType int8

const (
	RemoveConnectionFromUserUpdated BlockUserReplyType = iota
	RemoveConnectionToUserUpdated
	UserBlocked
	UnknownReply
)

type BlockUserReply struct {
	Users Users
	Type  BlockUserReplyType
}
