package matil

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// === math sd ===

func Add[T Number](nums ...T) T {
	var sum T
	for _, n := range nums {
		sum += n
	}
	return sum
}

func Sub[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := nums[0]
	for _, n := range nums[1:] {
		res -= n
	}

	return res
}

func Mul[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := T(1)
	for _, n := range nums {
		res *= n
	}

	return res
}

func Div[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n == 0 {
			panic("you stupid! division by zero?!")
		}
		res /= n
	}

	return res
}

func Min[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n < res {
			res = n
		}
	}

	return res
}

func Max[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n > res {
			res = n
		}
	}

	return res
}

func Avg[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}
	return Add(nums...) / T(len(nums))
}

// === transform ===

func Abs[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		if n < 0 {
			res[i] = -n
		} else {
			res[i] = n
		}
	}
	return res
}

func Neg[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = -n
	}
	return res
}

func Square[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = n * n
	}
	return res
}

func Cube[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = n * n * n
	}
	return res
}

// === misc ===

func Mod(nums ...int) int {
	if len(nums) < 2 {
		return 0
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n == 0 {
			panic("mod by zero")
		}
		res %= n
	}

	return res
}

func Sign(nums ...float64) []float64 {
	res := make([]float64, len(nums))
	for i, n := range nums {
		switch {
		case n > 0:
			res[i] = 1
		case n < 0:
			res[i] = -1
		default:
			res[i] = 0
		}
	}
	return res
}

// does stdlib go have this tho?
func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
