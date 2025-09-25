package matil

// Number is a type constraint that includes all integer and floating-point types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Add returns the sum of all given numbers.
// If no numbers are provided, it returns the zero value of T.
func Add[T Number](nums ...T) T {
	var sum T
	for _, n := range nums {
		sum += n
	}
	return sum
}

// Sub subtracts all subsequent numbers from the first.
// If no numbers are provided, it returns the zero value of T.
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

// Mul returns the product of all given numbers.
// If no numbers are provided, it returns the zero value of T.
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

// Div divides the first number by all subsequent numbers.
// It panics if division by zero is attempted.
// If no numbers are provided, it returns the zero value of T.
func Div[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n == 0 {
			panic("division by zero")
		}
		res /= n
	}
	return res
}

// Min returns the smallest number among the provided values.
// If no numbers are provided, it returns the zero value of T.
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

// Max returns the largest number among the provided values.
// If no numbers are provided, it returns the zero value of T.
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

// Avg returns the arithmetic mean of the given numbers.
// If no numbers are provided, it returns the zero value of T.
func Avg[T Number](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}
	return Add(nums...) / T(len(nums))
}

// Abs returns a slice containing the absolute values of the input numbers.
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

// Neg returns a slice containing the negated values of the input numbers.
func Neg[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = -n
	}
	return res
}

// Square returns a slice containing the square of each input number.
func Square[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = n * n
	}
	return res
}

// Cube returns a slice containing the cube of each input number.
func Cube[T Number](nums ...T) []T {
	res := make([]T, len(nums))
	for i, n := range nums {
		res[i] = n * n * n
	}
	return res
}

// Mod performs modulo operation sequentially from left to right.
// Example: Mod(20, 7, 3) => ((20 % 7) % 3).
// It panics if modulo by zero is attempted.
// If fewer than 2 numbers are provided, it returns 0.
func Mod(nums ...int) int {
	if len(nums) < 2 {
		return 0
	}

	res := nums[0]
	for _, n := range nums[1:] {
		if n == 0 {
			panic("modulo by zero")
		}
		res %= n
	}
	return res
}

// Sign returns a slice with the sign of each input number.
// The result is 1 for positive numbers, -1 for negative numbers, and 0 for zero.
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
