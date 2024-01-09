package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"

	"github.com/adriancrafter/todoapp/internal/am"
)

// NewGeoPoint creates a new GeoPoint.
func NewGeoPoint(lng, lat float64) am.GeoPoint {
	return am.GeoPoint{
		Lng: lng,
		Lat: lat,
	}
}

// NullGeoPoint represents a GeoPoint that may be null.
type NullGeoPoint struct {
	am.GeoPoint
	Valid bool
}

// NewNullGeoPoint creates a new NullGeoPoint.
func NewNullGeoPoint(point am.GeoPoint) NullGeoPoint {
	return NullGeoPoint{
		GeoPoint: point,
		Valid:    true,
	}
}

// Scan implements the sql.Scanner interface.
func (ngp *NullGeoPoint) Scan(value interface{}) error {
	if value == nil {
		ngp.Valid = false
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	geom, err := wkb.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("error unmarshalling WKB data: %v", err)
	}

	point, ok := geom.(orb.Point)
	if !ok {
		return fmt.Errorf("decoded type is not a Point")
	}

	ngp.GeoPoint = NewGeoPoint(point[0], point[1]) // orb.Point is [lng, lat]
	ngp.Valid = true
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (ngp NullGeoPoint) Value() (driver.Value, error) {
	if !ngp.Valid {
		return nil, nil
	}

	wkbData, err := wkb.Marshal(orb.Point{ngp.GeoPoint.Lng, ngp.GeoPoint.Lat})
	if err != nil {
		return nil, err
	}
	return wkbData, nil
}
