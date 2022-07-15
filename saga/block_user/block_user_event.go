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
	RollbackUpdates
	FinnishFunction
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
	ErrorOccured
	UnknownReply
)

type BlockUserReply struct {
	Users Users
	Type  BlockUserReplyType
}
