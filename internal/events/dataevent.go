package events

type DataEvent struct {
	ElapsedTime         uint64
	Level               uint64
	Distance            *uint64
	RemainingDistance   *uint64
	Time500mSplit       *uint64
	Time500mAverage     *uint64
	Strokes             *uint64
	StrokesPerMinute    *uint64
	WattsPreviousStroke *uint64
	WattsAverage        *uint64
	CaloriesPerHour     *uint64
	CaloriesTotal       *uint64
}

type DataEventBuilder struct {
	dataEvent *DataEvent
}

func NewDataEventBuilder(ElapsedTime uint64, Level uint64) *DataEventBuilder {
	return &DataEventBuilder{
		&DataEvent{
			ElapsedTime: ElapsedTime,
			Level:       Level,
		},
	}
}

func (b *DataEventBuilder) SetDistance(distance uint64) *DataEventBuilder {
	b.dataEvent.Distance = &distance
	return b
}

func (b *DataEventBuilder) SetRemainingDistance(remainingDistance uint64) *DataEventBuilder {
	b.dataEvent.RemainingDistance = &remainingDistance
	return b
}

func (b *DataEventBuilder) SetTime500mSplit(time500mSplit uint64) *DataEventBuilder {
	b.dataEvent.Time500mSplit = &time500mSplit
	return b
}

func (b *DataEventBuilder) SetTime500mAverage(time500mAverage uint64) *DataEventBuilder {
	b.dataEvent.Time500mAverage = &time500mAverage
	return b
}

func (b *DataEventBuilder) SetStrokes(strokes uint64) *DataEventBuilder {
	b.dataEvent.Strokes = &strokes
	return b
}

func (b *DataEventBuilder) SetStrokesPerMinute(strokesPerMinute uint64) *DataEventBuilder {
	b.dataEvent.StrokesPerMinute = &strokesPerMinute
	return b
}

func (b *DataEventBuilder) SetWattsPreviousStroke(wattsPreviousStroke uint64) *DataEventBuilder {
	b.dataEvent.WattsPreviousStroke = &wattsPreviousStroke
	return b
}

func (b *DataEventBuilder) SetWattsAverage(wattsAverage uint64) *DataEventBuilder {
	b.dataEvent.WattsAverage = &wattsAverage
	return b
}

func (b *DataEventBuilder) SetCaloriesPerHour(caloriesPerHour uint64) *DataEventBuilder {
	b.dataEvent.CaloriesPerHour = &caloriesPerHour
	return b
}

func (b *DataEventBuilder) SetCaloriesTotal(caloriesTotal uint64) *DataEventBuilder {
	b.dataEvent.CaloriesTotal = &caloriesTotal
	return b
}

func (b *DataEventBuilder) Build() *DataEvent {
	return b.dataEvent
}
