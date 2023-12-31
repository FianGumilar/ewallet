package utils

import "math"

func GetDistance(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6371.0 // Radius of the Earth in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lon1Ran := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Ran := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Ran - lon1Ran

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}
