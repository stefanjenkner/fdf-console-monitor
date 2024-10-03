package serialmonitor

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	line := "A8000060001410243028105065904"
	want := capture{
		elapsedTime:      6,
		distance:         14,
		time500m:         163,
		strokesPerMinute: 28,
		watts:            105,
		cals:             659,
		level:            4,
	}
	if got := parse(line); !reflect.DeepEqual(got, want) {
		t.Errorf("parse() = %v, want %v", got, want)
	}
}
