package events

type StatusChange uint64

const (
	Started         StatusChange = 0
	PausedOrStopped StatusChange = 1
	Resumed         StatusChange = 2
	Reset           StatusChange = 3
	LevelChanged    StatusChange = 4
)

type StatusChangeEvent struct {
	StatusChange StatusChange
}
