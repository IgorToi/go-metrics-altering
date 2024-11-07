// Package processmap provides helpful function for combination of two maps.
package processmap_test

import (
	"fmt"

	processmap "github.com/igortoigildin/go-metrics-altering/pkg/processMap"
)

func Example() {
	a := make(map[string]float64, 0)
	a["first"] = float64(10)
	b := make(map[string]int64, 0)
	b["second"] = int64(10)
	res := processmap.ConvertToSingleMap(a, b)

	fmt.Println(res)
	// Output:
	// map[first:10 second:10]
}
