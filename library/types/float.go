package types

func ToFloat64(value interface{}) (float64, bool) {
	switch value.(type) {
	case int:
		float := float64(value.(int))
		return float, true
	case string:
		number, _ := ToInt(value)
		float, _ := ToFloat64(number)
		return float, true
	default:
		float, ok := value.(float64)
		return float, ok
	}
}
