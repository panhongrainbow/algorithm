package randhub

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// =====================================================================================================================
//                  ⚗️ Functional Testing  (No Pool)
// =====================================================================================================================
// 🧪 Functional testing focuses on verifying the core functionalities of the system.
// 🧪 It involves inputting a large volume of data to ensure that all features work
// as expected and can handle a variety of possible inputs.

// Number 🧫 is a type constraint that allows int64 and float64 types.
type Number interface {
	int64 | float64
}

// GenerateNumbers 🧫 generates a slice of numbers of type T.
// The numbers may contain duplicates.
// count: the number of numbers to generate.
// minNum: the minimum value for the numbers.
// maxNum: the maximum value for the numbers.
func GenerateNumbers[T Number](count int64, minNum, maxNum T) ([]T, error) {
	// Convert minNum and maxNum to float64 for comparison
	minFloat := float64(minNum)
	maxFloat := float64(maxNum)

	// Ensure minNum is less than or equal to maxNum.
	if minFloat > maxFloat {
		return nil, errors.New("minNum must be less than or equal to maxNum")
	}

	// Create a new random number generator with a seed based on the current time.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Result slice to store the generated numbers.
	result := make([]T, 0, count)

	// Generate numbers until the required count is reached.
	for int64(len(result)) < count {
		num := T(minFloat + (maxFloat-minFloat)*rnd.Float64())
		result = append(result, num)
	}

	return result, nil
}

// GenerateUniqueNumbers 🧫 generates a slice of unique numbers of type T.
// count: the number of unique numbers to generate.
// minNum: the minimum value for the numbers.
// maxNum: the maximum value for the numbers.
func GenerateUniqueNumbers[T Number](count int64, minNum, maxNum T) ([]T, error) {
	// Convert minNum and maxNum to float64 for comparison
	minFloat := float64(minNum)
	maxFloat := float64(maxNum)

	// Ensure minNum is less than or equal to maxNum.
	if minFloat > maxFloat {
		return nil, errors.New("minNum must be less than or equal to maxNum")
	}

	// Ensure the range [minNum, maxNum] is large enough to generate the required count of unique numbers.
	switch v := any(maxNum).(type) {
	case int64:
		if v-int64(minNum)+1 < count {
			return nil, fmt.Errorf("not enough numbers in the range [%v, %v] to generate %d unique values", minNum, maxNum, count)
		}
		// For float64 types, no range validation is performed, as floating-point numbers can represent an infinite range and may involve irrational numbers.
	}

	// Create a new random number generator with a seed based on the current time.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Use a map to keep track of unique numbers.
	numbers := make(map[T]struct{})
	// Result slice to store the unique numbers.
	result := make([]T, 0, count)

	// Generate unique numbers until the required count is reached.
	for int64(len(result)) < count {
		num := T(minFloat + (maxFloat-minFloat)*rnd.Float64())
		// Check if the number is already in the map (i.e., it's unique).
		if _, exists := numbers[num]; !exists {
			numbers[num] = struct{}{}
			result = append(result, num)
		}
	}

	return result, nil
}

// GenerateNumbersWithOption 🧫 chooses between GenerateNumbers and GenerateUniqueNumbers based on the unique flag.
// count: the number of numbers to generate.
// minNum: the minimum value for the numbers.
// maxNum: the maximum value for the numbers.
// unique: a flag indicating whether to generate unique numbers or not.
func GenerateNumbersWithOption[T Number](count int64, minNum, maxNum T, unique bool) ([]T, error) {
	// Choose the appropriate function based on the unique flag.
	if unique {
		return GenerateUniqueNumbers(count, minNum, maxNum)
	}
	return GenerateNumbers(count, minNum, maxNum)
}
