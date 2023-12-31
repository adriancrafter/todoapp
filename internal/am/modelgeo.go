package am

type GeoPoint struct {
	Lng float64
	Lat float64
	Alt float64
}

func NewGeoPoint(lng, lat, alt float64) GeoPoint {
	return GeoPoint{
		Lng: lng,
		Lat: lat,
		Alt: alt,
	}
}
