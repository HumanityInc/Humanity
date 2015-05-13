package db

type Message struct {
	Id     int64
	ProfId int64
	CTime  int64
	Text   string
}

func SendMessage(fromProfId, toProfId int64, text string) (ok bool) {

	return
}

func MessageList(profId int64) (list []Message) {

	return
}
