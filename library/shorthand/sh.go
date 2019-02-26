package sh

func For(count int, function func(index int)) func() {
	return func() {
		for i := 0; i < count; i++ {
			function(i)
		}
	}
}

func If(b bool, function func()) {
	if b {
		function()
	}
}

func IfElse(b bool, ifFunc, elseFunc func()) {
	if b {
		ifFunc()
	} else {
		elseFunc()
	}
}
