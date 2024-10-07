// Package processmap provides helpful function for combination of two maps.
package processmap

func ConvertToSingleMap(a map[string]float64, b map[string]int64) map[string]interface{} {
	c := make(map[string]interface{}, 33)
	for i, v := range a {
		c[i] = v
	}
	for j, l := range b {
		c[j] = l
	}
	return c
}
