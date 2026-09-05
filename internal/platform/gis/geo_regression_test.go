package gis

import "testing"

func TestGeoMetricsRejectsOutOfRangeCoordinates(t *testing.T) {
	_, err := GeoMetrics(FiberPoint{Lat:91,Lng:0}, FiberPoint{Lat:0,Lng:0})
	if err == nil { t.Fatal("latitude outside [-90,90] must be rejected") }
}

func TestGeoMetricsAcceptsValidCoordinates(t *testing.T) {
	m, err := GeoMetrics(FiberPoint{Lat:0,Lng:0}, FiberPoint{Lat:0,Lng:1})
	if err != nil || m.DistanceM <= 0 { t.Fatalf("expected positive distance, metric=%#v err=%v", m, err) }
}
