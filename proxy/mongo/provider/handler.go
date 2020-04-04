package provider

type Handler interface {
	HandleQuery(msg Message) (r *ReplyMessage, err error)

	HandleOpCommand(msg Message) (r *CommandReplyMessage, err error)

	HandleOpMsg(msg Message) (r *MessageMessage, err error)

	GetUserDbList(user interface{}) error
}
