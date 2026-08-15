package geo

import (
	"errors"
	"math"
)

// Point represents a geographic location in decimal degrees.
type Point struct {
	Lat float64
	Lon float64
}

// Target is a routable FTN resource with an optional geographic position.
type Target struct {
	ID     string
	Point  Point
	Active bool
	Metric uint32
}

var ErrInvalidPoint = errors.New("invalid geographic point")

func validPoint(p Point) bool {
	return p.Lat >= -90 && p.Lat <= 90 && p.Lon >= -180 && p.Lon <= 180 && !math.IsNaN(p.Lat) && !math.IsNaN(p.Lon)
}

// RankByDistance returns active targets ordered by great-circle distance.
// Geographic ranking is advisory; it never overrides explicit routing policy.
func RankByDistance(origin Point, targets []Target) ([]Target, error) {
	if !validPoint(origin) {
		return nil, ErrInvalidPoint
	}
	out := make([]Target, 0, len(targets))
	for _, t := range targets {
		if t.Active && validPoint(t.Point) {
			out = append(out, t)
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && haversine(origin, out[j].Point) < haversine(origin, out[j-1].Point); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

func haversine(a, b Point) float64 {
	const earthRadiusKm = 6371.0088
	lat1, lat2 := a.Lat*math.Pi/180, b.Lat*math.Pi/180
	dLat := (b.Lat - a.Lat) * math.Pi / 180
	dLon := (b.Lon - a.Lon) * math.Pi / 180
	h := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}
