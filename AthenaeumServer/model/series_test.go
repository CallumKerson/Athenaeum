package model_test

import (
	"testing"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func TestShouldFormatSeries(t *testing.T) {
	//given
	expected := "The Stormlight Archive Book 2.5"
	series := &model.Series{2.500, "The Stormlight Archive"}

	//when
	actual := series.String()

	//then
	if actual != expected {
		t.Errorf("Generated wrong entry\nExpected: %s\nWas: %s", expected, actual)
	}
}

func TestShouldFormatEntryForWholeEntryNumber(t *testing.T) {
	//given
	expected := "2"
	series := &model.Series{Entry: 2}

	//when
	actual := series.EntryString()

	//then
	if actual != expected {
		t.Errorf("Generated wrong entry\nExpected: %s\nWas: %s", expected, actual)
	}
}

func TestShouldFormatEntryForIntermidiateEntryNumber(t *testing.T) {
	//given
	expected := "2.5"
	series := &model.Series{Entry: 2.500}

	//when
	actual := series.EntryString()

	//then
	if actual != expected {
		t.Errorf("Generated wrong entry\nExpected: %s\nWas: %s", expected, actual)
	}
}
