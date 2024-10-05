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

type DataEventBuilder struct {
	dataEvent *DataEvent
}

func NewDataEventBuilder(ElapsedTime uint16, Level uint8) *DataEventBuilder {
	return &DataEventBuilder{
		&DataEvent{
			ElapsedTime: ElapsedTime,
			Level:       Level,
		},
	}
}

func (b *DataEventBuilder) SetDistance(distance uint16) *DataEventBuilder {
	b.dataEvent.Distance = &distance
	return b
}

func (b *DataEventBuilder) SetRemainingDistance(remainingDistance uint16) *DataEventBuilder {
	b.dataEvent.RemainingDistance = &remainingDistance
	return b
}

func (b *DataEventBuilder) SetTime500mSplit(time500mSplit uint16) *DataEventBuilder {
	b.dataEvent.Time500mSplit = &time500mSplit
	return b
}

func (b *DataEventBuilder) SetTime500mAverage(time500mAverage uint16) *DataEventBuilder {
	b.dataEvent.Time500mAverage = &time500mAverage
	return b
}

func (b *DataEventBuilder) SetStrokes(strokes uint16) *DataEventBuilder {
	b.dataEvent.Strokes = &strokes
	return b
}

func (b *DataEventBuilder) SetStrokesPerMinute(strokesPerMinute uint8) *DataEventBuilder {
	b.dataEvent.StrokesPerMinute = &strokesPerMinute
	return b
}

func (b *DataEventBuilder) SetWattsPreviousStroke(wattsPreviousStroke uint16) *DataEventBuilder {
	b.dataEvent.WattsPreviousStroke = &wattsPreviousStroke
	return b
}

func (b *DataEventBuilder) SetWattsAverage(wattsAverage uint16) *DataEventBuilder {
	b.dataEvent.WattsAverage = &wattsAverage
	return b
}

func (b *DataEventBuilder) SetCaloriesPerHour(caloriesPerHour uint16) *DataEventBuilder {
	b.dataEvent.CaloriesPerHour = &caloriesPerHour
	return b
}

func (b *DataEventBuilder) SetCaloriesTotal(caloriesTotal uint16) *DataEventBuilder {
	b.dataEvent.CaloriesTotal = &caloriesTotal
	return b
}

func (b *DataEventBuilder) Build() *DataEvent {
	return b.dataEvent
}
