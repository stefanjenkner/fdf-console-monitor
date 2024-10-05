package events

type DataEvent struct {
	ElapsedTime         uint16
	Level               uint8
	Distance            *uint16
	RemainingDistance   *uint16
	Time500mSplit       *uint16
	Time500mAverage     *uint16
	Strokes             *uint16
	StrokesPerMinute    *uint8
	WattsPreviousStroke *uint16
	WattsAverage        *uint16
	CaloriesPerHour     *uint16
	CaloriesTotal       *uint16
}
type DataEventOption func(*DataEvent)

func NewDataEvent(elapsedTime uint16, level uint8, opts ...DataEventOption) *DataEvent {
	e := &DataEvent{
		ElapsedTime: elapsedTime,
		Level:       level,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithDistance(distance uint16) DataEventOption {
	return func(e *DataEvent) {
		e.Distance = &distance
	}
}

func WithRemainingDistance(remainingDistance uint16) DataEventOption {
	return func(e *DataEvent) {
		e.RemainingDistance = &remainingDistance
	}
}

func WithTime500mSplit(time500mSplit uint16) DataEventOption {
	return func(e *DataEvent) {
		e.Time500mSplit = &time500mSplit
	}
}

func WithTime500mAverage(time500mAverage uint16) DataEventOption {
	return func(e *DataEvent) {
		e.Time500mAverage = &time500mAverage
	}
}

func WithStrokes(strokes uint16) DataEventOption {
	return func(e *DataEvent) {
		e.Strokes = &strokes
	}
}

func WithStrokesPerMinute(strokesPerMinute uint8) DataEventOption {
	return func(e *DataEvent) {
		e.StrokesPerMinute = &strokesPerMinute
	}
}

func WithWattsPreviousStroke(wattsPreviousStroke uint16) DataEventOption {
	return func(e *DataEvent) {
		e.WattsPreviousStroke = &wattsPreviousStroke
	}
}

func WithWattsAverage(wattsAverage uint16) DataEventOption {
	return func(e *DataEvent) {
		e.WattsAverage = &wattsAverage
	}
}

func WithCaloriesPerHour(caloriesPerHour uint16) DataEventOption {
	return func(e *DataEvent) {
		e.CaloriesPerHour = &caloriesPerHour
	}
}

func WithCaloriesTotal(caloriesTotal uint16) DataEventOption {
	return func(e *DataEvent) {
		e.CaloriesTotal = &caloriesTotal
	}
}
