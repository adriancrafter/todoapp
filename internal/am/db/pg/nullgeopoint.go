package pg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adriancrafter/todoapp/internal/am"
)

func NewGeoPoint(lng, lat, alt float64) am.GeoPoint {
	return am.GeoPoint{
		Lng: lng,
		Lat: lat,
		Alt: alt,
	}
}

type NullGeoPoint struct {
	am.GeoPoint
	Valid bool
}

func (n *NullGeoPoint) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}

	var pointString string
	switch v := value.(type) {
	case string:
		pointString = v
	// case []byte: // If it's in binary format
	// pointString = string(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	pointString = strings.TrimPrefix(pointString, "POINT(")
	pointString = strings.TrimSuffix(pointString, ")")
	coords := strings.Split(pointString, " ")
	if len(coords) != 2 {
		return fmt.Errorf("invalid POINT format")
	}

	lng, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return fmt.Errorf("invalid longitude value: %v", err)
	}

	lat, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return fmt.Errorf("invalid latitude value: %v", err)
	}

	n.GeoPoint = NewGeoPoint(lng, lat, 0) // Assume altitude is 0
	n.Valid = true
	return nil
}
