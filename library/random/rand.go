package random

import (
	"math/rand"
	"startkit/library/systems"
	"time"
)

func RandInt64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	if systems.GetGOOS() == "linux" {
		rand.Seed(time.Now().UnixNano())
	}
	return rand.Int63n(max-min) + min
}

func RandInt64Array(min, max, length int64) (array []int64) {
	if min >= max || max == 0 {
		return []int64{}
	}
	if systems.GetGOOS() == "linux" {
		for i := int64(0); i < length; i++ {
			rand.Seed(time.Now().UnixNano())
			array = append(array, rand.Int63n(max-min)+min)
		}
	} else {
		for i := int64(0); i < length; i++ {
			array = append(array, rand.Int63n(max-min)+min)
		}
	}
	return
}

func RandInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	if systems.GetGOOS() == "linux" {
		rand.Seed(time.Now().UnixNano())
	}
	return rand.Intn(max-min) + min
}

func RandIntArray(min, max, length int) (array []int) {
	if min >= max || max == 0 {
		return []int{}
	}
	if systems.GetGOOS() == "linux" {
		for i := 0; i < length; i++ {
			rand.Seed(time.Now().UnixNano())
			array = append(array, rand.Intn(max-min)+min)
		}
	} else {
		for i := 0; i < length; i++ {
			array = append(array, rand.Intn(max-min)+min)
		}
	}
	return
}

func ShuffleArray(src []interface{}) (array []interface{}) {
	array = make([]interface{}, len(src))
	for i, v := range RandIntArray(0, len(src), len(src)) {
		array[i] = src[v]
	}
	return
}

func ShuffleStringArray(src []string) (array []string) {
	array = make([]string, len(src))
	for i, v := range RandIntArray(0, len(src), len(src)) {
		array[i] = src[v]
	}
	return
}
