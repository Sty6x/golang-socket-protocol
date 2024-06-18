package message

type Header struct {
	Protocol       string
	ConnectionType string
}

type Request struct {
	Header
	Namespace       string
	DateEstablished string
	UserId          string
}

type Response struct {
	Header
	ConnectionId    string
	DateEstablished string
	Status          string
}

type PushMessage struct {
	Header
	Status          string
	UserId          string
	Namespace       string
	DateEstablished string
	Payload         string
}
