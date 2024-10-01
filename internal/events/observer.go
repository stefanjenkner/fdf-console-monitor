package events

type Observer interface {
	OnData(event DataEvent)
	OnStatusChange(event StatusChangeEvent)
}
