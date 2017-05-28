package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"
)

type gpx struct {
	TrackList []trk `xml:"trk"`
}

type trk struct {
	SegmentList []trkseg `xml:"trkseg"`
}

type trkseg struct {
	PointList []trkpt `xml:"trkpt"`
}

type trkpt struct {
	Latitude  float64    `xml:"lat,attr"`
	Longitude float64    `xml:"lon,attr"`
	Timestamp *time.Time `xml:"time,omitempty"`
}

func readGpx(filename string) (gpx, error) {

	fmt.Printf("Opening %s...\n", filename)

	empty := gpx{}

	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return empty, err
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	var q gpx
	xml.Unmarshal(b, &q)

	// fmt.Println(q)

	return q, nil

}

// func isEqual(x trkpt, test trkpt) bool {
// 	return test.Timestamp.Equal(x.Timestamp)
// }
//
// func isBefore(x trkpt, test trkpt) bool {
// 	return test.Timestamp.Before(x.Timestamp)
// }
//
// func isAfter(x trkpt, test trkpt) bool {
// 	return test.Timestamp.After(x.Timestamp)
// }
//
// func isInside(x trkpt, y trkpt, test trkpt) bool {
// 	if isAfter(x, test) && isBefore(y, test) {
// 		return true
// 	} else {
// 		return false
// 	}
// }

func distanceBetweenPointsInMeters(a trkpt, b trkpt) float64 {
	meanRadiusOfEarth := 6372797.0
	aLatMeters := a.Latitude * (math.Pi / 180)
	bLatMeters := b.Latitude * (math.Pi / 180)

	deltaLatitudeMeters := (b.Latitude - a.Latitude) * (math.Pi / 180)
	deltaLongitudeMeters := (b.Longitude - a.Longitude) * (math.Pi / 180)

	aaa := math.Sin(deltaLatitudeMeters/2)*math.Sin(deltaLatitudeMeters/2) + math.Cos(aLatMeters)*math.Cos(bLatMeters)*math.Sin(deltaLongitudeMeters/2)*math.Sin(deltaLongitudeMeters/2)
	ccc := 2 * math.Atan2(math.Sqrt(aaa), math.Sqrt(1-aaa))

	distance := meanRadiusOfEarth * ccc
	return distance
}

func main() {

	aFile := flag.String("a", "", "A File")
	bFile := flag.String("b", "", "B File")
	outFilename := flag.String("o", "", "Out File")
	flag.Parse()

	aSet, aErr := readGpx(*aFile)
	if aErr != nil {
		fmt.Println(aErr)
		os.Exit(1)
	}
	bSet, bErr := readGpx(*bFile)
	if bErr != nil {
		fmt.Println(bErr)
		os.Exit(2)
	}

	for _, knownTimeTrackList := range aSet.TrackList {
		for _, knownTimeSegment := range knownTimeTrackList.SegmentList {
			for _, knownPoint := range knownTimeSegment.PointList {
				//find nearest point in bSet
				minDistance := 0.0
				foundSomething := false
				closestTrackIndex := 0
				closestSegmentIndex := 0
				closestPointIndex := 0

				for a, bSetTrack := range bSet.TrackList {
					for b, bSetSegment := range bSetTrack.SegmentList {
						for c, bSetVaguePoint := range bSetSegment.PointList {
							dist := distanceBetweenPointsInMeters(knownPoint, bSetVaguePoint)
							if (dist < minDistance) || (foundSomething == false) {
								foundSomething = true
								minDistance = dist
								closestTrackIndex = a
								closestSegmentIndex = b
								closestPointIndex = c
							}
						}
					}
				}

				if foundSomething {
					// fmt.Println("knownPoint: ", knownPoint)
					// fmt.Println(closestTrackIndex, closestSegmentIndex, closestPointIndex)

					//set the timestamp of the bSet point to the closest point in aSet
					bSet.TrackList[closestTrackIndex].SegmentList[closestSegmentIndex].PointList[closestPointIndex].Timestamp = knownPoint.Timestamp
				} else {
					fmt.Println("bSet appears to be empty")
					os.Exit(1)
				}

			}
		}
	}

	//interpolate times in bSet
	var previousPointWithTimestamp trkpt
	previousTrackIndex := 0
	previousSegmentIndex := 0
	previousPointIndex := 0

	for a, bSetTrack := range bSet.TrackList {
		for b, bSetSegment := range bSetTrack.SegmentList {
			for c, bSetPoint := range bSetSegment.PointList {

				if (bSetPoint.Timestamp != nil) && (previousPointWithTimestamp.Timestamp != nil) {
					//not the first time coming through
					fmt.Println("happening", previousTrackIndex, previousSegmentIndex, previousPointIndex, a, b, c)

					startTimestamp := *bSetPoint.Timestamp
					endTimestamp := *previousPointWithTimestamp.Timestamp
					duration := startTimestamp.Sub(endTimestamp)
					fmt.Println("duration:", duration.Minutes())
					//iterate between first known point and second known point
					for aa, bbSetTrack := range bSet.TrackList {
						for bb, bbSetSegment := range bbSetTrack.SegmentList {
							for cc, _ := range bbSetSegment.PointList {
								if (aa >= previousTrackIndex && aa <= a) && (bb >= previousSegmentIndex && bb <= b) && (cc > previousPointIndex && cc < c) {
									//interpolate time for this point

									fmt.Println("zzzz", aa, bb, cc)
								}
							}
						}
					}

				}

				//prepare for next loop
				if bSetPoint.Timestamp != nil {
					previousPointWithTimestamp = bSetPoint
					previousTrackIndex = a
					previousSegmentIndex = b
					previousPointIndex = c
					fmt.Println("I made it", a, b, c, previousPointWithTimestamp)
				}
			}
		}
	}

	fmt.Println("Marshalling up...")

	bOut, _ := xml.MarshalIndent(bSet, "", "\t")
	//fmt.Printf("%s", bOut)

	fmt.Println("Writing out...")

	out, err := os.Create(*outFilename)
	fmt.Println(err)
	fmt.Fprintf(out, "%s%s", xml.Header, bOut)
	defer out.Close()

}
