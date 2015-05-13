package db

type Event struct {
	Id   int64
	Type int
	Text string
}

func AddEvent(profId int64, _type int, test string) (ok bool) {

	return
}

func EventList(profId int64) (list []Event) {

	return
}
