package geojson

/*Multipolygon is a representation of a GeoJSON Multipolygon,
an array of distinct polygons.*/
type Multipolygon struct {
	Type string		       `json:"type"`
	Coordinates [][][]Coordinate       `json:"coordinates"`
}

func NewMultipolygon(coordinates [][][]Coordinate) Multipolygon {
	return Multipolygon{
		Type:"multipolygon",
		Coordinates:coordinates,
	}
}

/*Polygon is a representation of a GeoJSON point*/
type Point struct {
	Type string		       `json:"type"`
	Coordinates Coordinate       `json:"coordinates"`
}

func NewPoint(coordinates Coordinate) Point {
	return Point{
		Type:"point",
		Coordinates:coordinates,
	}
}

/*Coordinate is a pair of floats representing the latitude and longitude
  of a single geographic point. In GeoJSON, these ar given as longitude,
  latitude pairs*/

type Coordinate [2]float64
