package gis

import (
	"fmt"
	"math"
)

type GeoMetric struct {
	DistanceM  float64   `json:"distance_m"`
	BearingDeg float64   `json:"bearing_deg"`
	Midpoint   FiberPoint `json:"midpoint"`
}

func GeoMetrics(a, b FiberPoint) (GeoMetric, error) {
	if math.Abs(a.Lat) > 90 || math.Abs(b.Lat) > 90 || math.Abs(a.Lng) > 180 || math.Abs(b.Lng) > 180 {
		return GeoMetric{}, fmt.Errorf("invalid coordinates")
	}
	return GeoMetric{
		DistanceM:  haversine(a, b),
		BearingDeg: bearing(a, b),
		Midpoint:   FiberPoint{Lat: (a.Lat + b.Lat) / 2, Lng: (a.Lng + b.Lng) / 2},
	}, nil
}

func bearing(a, b FiberPoint) float64 {
	lat1 := a.Lat * math.Pi / 180
	lat2 := b.Lat * math.Pi / 180
	dl := (b.Lng - a.Lng) * math.Pi / 180
	y := math.Sin(dl) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dl)
	d := math.Atan2(y, x) * 180 / math.Pi
	if d < 0 {
		d += 360
	}
	return d
}
