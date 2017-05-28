package main
	
import (
	"testing"
	)
		
func TestAucklandToWellingtonDistance(t *testing.T) {
	knownAucklandToWellingtonDistanceInMeters := 491808.0
	accuracyInMeters := 5000.0
	
	aucklandTrkpt := trkpt{Latitude: -36.8485, Longitude: 174.7633}
	wellingtonTrkpt := trkpt{Latitude: -41.2865, Longitude: 174.7762}
	
	dist := distanceBetweenPointsInMeters(aucklandTrkpt, wellingtonTrkpt)
	
	minDistance := knownAucklandToWellingtonDistanceInMeters-accuracyInMeters
	maxDistance := knownAucklandToWellingtonDistanceInMeters+accuracyInMeters		
	if !(dist > minDistance && dist < maxDistance) {
		t.Error("Expected", knownAucklandToWellingtonDistanceInMeters, ", got ", dist);
	}
}

func TestOneDegreeAroundEquator(t *testing.T) {

	expectedDistance := 111226.0
	accuracyInMeters := 1000.0

	origin        := trkpt{Latitude: 0.0, Longitude: 0.0}
	originPlusOne := trkpt{Latitude: 1.0, Longitude: 0.0}
	
	dist := distanceBetweenPointsInMeters(origin, originPlusOne)

	minDistance := expectedDistance - accuracyInMeters
	maxDistance := expectedDistance + accuracyInMeters		
		
	if !(dist > minDistance && dist < maxDistance) {
		t.Error("Expected", expectedDistance, ", got ", dist);
	}
}

func TestOneDegreeAroundTropicish(t *testing.T) {

	expectedDistance := 77000.0
	accuracyInMeters := 2000.0

	origin        := trkpt{Latitude: 45.0, Longitude: 0.0}
	originPlusOne := trkpt{Latitude: 45.0, Longitude: 1.0}
	
	dist := distanceBetweenPointsInMeters(origin, originPlusOne)

	minDistance := expectedDistance - accuracyInMeters
	maxDistance := expectedDistance + accuracyInMeters		
		
	if !(dist > minDistance && dist < maxDistance) {
		t.Error("Expected", expectedDistance, ", got ", dist);
	}
}

func TestOneDegreeAroundAntarticish(t *testing.T) {

	expectedDistance := 2000.0
	accuracyInMeters := 200.0

	origin        := trkpt{Latitude: -89.0, Longitude: 0.0}
	originPlusOne := trkpt{Latitude: -89.0, Longitude: 1.0}
	
	dist := distanceBetweenPointsInMeters(origin, originPlusOne)

	minDistance := expectedDistance - accuracyInMeters
	maxDistance := expectedDistance + accuracyInMeters		
		
	if !(dist > minDistance && dist < maxDistance) {
		t.Error("Expected", expectedDistance, ", got ", dist);
	}
}	