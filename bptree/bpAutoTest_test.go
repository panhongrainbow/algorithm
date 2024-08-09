package bpTree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var (
	// randomQuantity represents the number of elements to be generated for random testing.
	randomQuantity = 20000

	// randomMax represents the maximum value for generating random numbers.
	randomMax = 150000

	// randomMin represents the minimum value for generating random numbers.
	randomMin = 10
)

// GenerateUniqueNumbers generates a slice of unique integers
// count: the number of unique integers to generate
// max: the maximum value for the integers
// min: the minimum value for the integers
func GenerateUniqueNumbers(count, max, min int) ([]int64, error) {
	// Ensure the range [min, max] is large enough to generate the required count of unique numbers
	if max-min+1 < count {
		return nil, fmt.Errorf("not enough numbers in the range [%d, %d] to generate %d unique values", min, max, count)
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Use a map to keep track of unique numbers
	numbers := make(map[int64]struct{})
	// Result slice to store the unique numbers
	result := make([]int64, 0, count)

	// Generate unique numbers until the required count is reached
	for len(result) < count {
		num := int64(rand.Intn(max-min+1) + min)
		// Check if the number is already in the map (i.e., it's unique)
		if _, exists := numbers[num]; !exists {
			numbers[num] = struct{}{}
			result = append(result, num)
		}
	}

	return result, nil
}

// Test_Check_BpTree_Automatic is used for automated testing, generating test data with random numbers for B+ tree insertion and deletion.
func Test_Check_BpTree_Automatic(t *testing.T) {
	// Automated random testing for B+ tree.
	t.Run("Automated Testing Section", func(t *testing.T) {

		numbersForAdding, _ := GenerateUniqueNumbers(randomQuantity, randomMax, randomMin)

		// Set up randomization.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Generate random data for insertion.
		/*numbersForAdding := make([]int64, randomQuantity)
		for i := 0; i < randomQuantity; i++ {
			num := int64(random.Intn(randomMax-randomMin+1) + randomMin)
			numbersForAdding[i] = num
		}*/
		fmt.Println("Random data for insertion:", numbersForAdding)

		// Generate random data for deletion.
		numbersForDeleting := make([]int64, randomQuantity)
		copy(numbersForDeleting, numbersForAdding)
		shuffleSlice(numbersForDeleting, random)
		fmt.Println("Random data for deletion:", numbersForDeleting)

		// Generate sorted data.
		sortedNumbers := make([]int64, randomQuantity)
		copy(sortedNumbers, numbersForAdding)
		sort.Slice(sortedNumbers, func(i, j int) bool {
			return sortedNumbers[i] < sortedNumbers[j]
		})
		fmt.Println("Sorted data:", sortedNumbers)

		// Initialize B-tree.
		root := NewBpTree(5)

		// Start inserting data.
		for i := 0; i < randomQuantity; i++ {
			// Insert data entries continuously.
			root.InsertValue(BpItem{Key: numbersForAdding[i]})
		}

		// Start deleting data.
		for i := 0; i < randomQuantity; i++ {

			// 显示目前的删除值
			value := numbersForDeleting[i]
			fmt.Println(i, value)

			// Deleting data entries continuously.
			deleted, _, _, err := root.RemoveValue(BpItem{Key: numbersForDeleting[i]})
			// deleted, _, _, err := root.RemoveValueLeft(BpItem{Key: numbersForDeleting[i]})

			if deleted == false {
				fmt.Println("Breakpoint: Data deletion not successful. 💢 The number is ", numbersForDeleting[i], i)
				panic("Breakpoint: Data deletion not successful.")
			}
			if err != nil {
				fmt.Println("Breakpoint: Deletion encountered an error. 💢 The number is ", numbersForDeleting[i], i)
				panic("Breakpoint: Deletion encountered an error.")
			}
		}

		fmt.Println()
	})
	// Automated random testing for B+ tree.
	t.Run("Manually Identify B Plus Tree Operation Errors", func(t *testing.T) {
		// 数量二百的例子
		// Generate random data for insertion.
		var randomNumbers = []int64{1538, 1064, 249, 1966, 778, 1046, 1764, 1797, 847, 726, 1212, 119, 1063, 266, 1622, 511, 1609, 450, 1011, 707, 1425, 1045, 821, 1294, 1154, 1723, 1349, 1499, 230, 1320, 312, 917, 845, 1738, 1462, 236, 320, 1381, 409, 1805, 1709, 943, 1879, 69, 211, 367, 898, 1700, 1234, 395, 710, 1196, 1526, 384, 509, 1962, 1456, 205, 1830, 576, 587, 419, 1252, 1091, 346, 1066, 1876, 1088, 351, 1031, 1568, 1233, 761, 715, 691, 1368, 1739, 1314, 1197, 224, 1049, 1060, 1036, 1420, 567, 1305, 1618, 1557, 919, 115, 155, 1601, 874, 540, 260, 892, 1423, 794, 362, 484, 868, 1945, 958, 969, 1977, 905, 229, 1914, 376, 736, 156, 1105, 530, 1629, 405, 1398, 1706, 1691, 1683, 743, 1971, 279, 1256, 247, 1745, 785, 1119, 1513, 1078, 879, 1556, 1804, 1873, 388, 1418, 1880, 1362, 1392, 611, 930, 1240, 1571, 502, 1013, 439, 1581, 1881, 114, 235, 1703, 1341, 591, 393, 488, 538, 1925, 624, 1975, 1536, 1654, 89, 1902, 495, 1944, 674, 942, 1248, 1250, 660, 1928, 527, 1017, 1161, 1682, 71, 1807, 158, 1768, 435, 1623, 458, 1630, 125, 151, 67, 1087, 1820, 1009, 1506, 1823, 494, 863, 1439, 887, 765, 1600, 1428, 671, 1608, 1394, 583, 1288, 1824, 1737, 1180, 416, 1350, 1867, 714, 687, 138, 1535, 701, 614, 411, 1694, 522, 1913, 1328, 23, 1767, 87, 1365, 1441, 314, 1868, 262, 857, 54, 1055, 1765, 282, 1681, 43, 1325, 363, 1577, 396, 1302, 1023, 1190, 760, 417, 276, 1194, 1489, 1220, 1806, 1487, 275, 1659, 1859, 1777, 1163, 1204, 575, 1175, 947, 870, 163, 60, 104, 753, 1217, 748, 1244, 1758, 1658, 440, 1071, 888, 1592, 1040, 639, 601, 1417, 222, 118, 1862, 73, 605, 1003, 1177, 1152, 291, 1665, 1126, 330, 1464, 830, 762, 677, 933, 1301, 473, 1385, 652, 193, 1552, 621, 1888, 323, 1293, 1626, 597, 818, 343, 307, 940, 574, 446, 584, 1891, 1249, 97, 1690, 820, 918, 1313, 521, 935, 1183, 1229, 1253, 474, 1108, 1831, 1840, 582, 829, 1616, 34, 1363, 498, 1095, 675, 953, 1632, 1263, 103, 1653, 669, 453, 1751, 1495, 1171, 1704, 113, 1118, 1168}
		// Generate random data for deletion.
		var shuffledNumbers = []int64{236, 1571, 671, 1880, 669, 674, 715, 396, 119, 1806, 1758, 388, 363, 1420, 701, 193, 1626, 1777, 863, 691, 1709, 320, 1305, 857, 1196, 1556, 905, 1313, 458, 1234, 1428, 1204, 346, 1071, 151, 249, 918, 947, 1623, 343, 118, 1608, 97, 114, 887, 1618, 1629, 892, 1040, 743, 1171, 1745, 1423, 494, 113, 1925, 1694, 488, 276, 1737, 1658, 1064, 1609, 1240, 1363, 820, 1055, 1977, 611, 1046, 1060, 942, 760, 230, 1881, 1487, 247, 23, 1820, 1600, 1706, 540, 1088, 1738, 1513, 639, 235, 1031, 1190, 1108, 1249, 1263, 1723, 1876, 1462, 222, 1581, 1003, 1557, 1879, 1293, 1368, 1456, 1830, 1392, 34, 1506, 1180, 821, 266, 1441, 314, 940, 587, 1017, 1683, 1212, 1301, 829, 1690, 830, 1764, 1536, 1653, 1183, 1256, 1862, 1091, 262, 502, 1768, 1700, 1341, 1294, 384, 614, 1868, 1013, 205, 419, 1385, 943, 1233, 1804, 435, 1823, 60, 710, 1681, 687, 930, 376, 707, 511, 498, 919, 601, 71, 393, 440, 762, 1418, 958, 1632, 1797, 1320, 1654, 1049, 509, 1499, 591, 1288, 1244, 1194, 446, 1867, 933, 1126, 1105, 584, 367, 1568, 395, 583, 69, 411, 1577, 1161, 473, 1622, 484, 453, 1439, 1066, 1751, 1045, 43, 785, 953, 1592, 530, 1119, 156, 163, 1250, 158, 1944, 291, 527, 621, 879, 409, 870, 1767, 474, 417, 522, 73, 1220, 1425, 115, 1394, 1163, 1036, 736, 1914, 675, 888, 1975, 495, 1362, 362, 1859, 761, 576, 935, 868, 1704, 1739, 1349, 1023, 1552, 969, 1526, 1253, 597, 67, 1913, 765, 351, 1302, 575, 1177, 1325, 103, 605, 1118, 1252, 224, 818, 794, 1962, 1009, 1011, 138, 1945, 1175, 1971, 1152, 1495, 574, 1682, 714, 748, 1807, 624, 521, 677, 1902, 260, 1538, 1398, 1154, 229, 1840, 1616, 1873, 1805, 282, 778, 1229, 567, 652, 1765, 323, 450, 898, 1966, 1703, 1381, 1691, 1888, 874, 1314, 87, 753, 845, 1095, 1328, 439, 330, 1350, 1168, 1197, 1659, 726, 307, 405, 104, 125, 1489, 1217, 279, 1824, 1365, 660, 1601, 1417, 1063, 582, 155, 89, 1665, 1087, 917, 275, 1248, 847, 1891, 211, 1464, 1630, 312, 1535, 1928, 54, 416, 1078, 538, 1831}

		// Initialize B plus tree.
		root := NewBpTree(3)

		// Start inserting data.
		for i := 0; i < randomQuantity; i++ {
			// Insert data entries continuously.
			root.InsertValue(BpItem{Key: randomNumbers[i]})
		}

		// Start deleting data.
		for i := 0; i < randomQuantity; i++ {

			// 中断检验
			value := shuffledNumbers[i]
			fmt.Println(i, value)
			if shuffledNumbers[i] == 1824 {
				fmt.Println(">>>>> !")
				fmt.Print()
			}

			/*if shuffledNumbers[i] == 1365 {
				root.root.Index[1] = 1078
			}*/

			deleted, _, _, err := root.RemoveValue(BpItem{Key: shuffledNumbers[i]})

			if deleted == false {
				fmt.Println("Breakpoint: Data deletion not successful. 💢 The number is ", shuffledNumbers[i], i)
			}
			if err != nil {
				fmt.Println("Breakpoint: Deletion encountered an error. 💢 The number is ", shuffledNumbers[i], i)
			}
		}

		fmt.Println()
	})
}

// shuffleSlice randomly shuffles the elements in the slice.
func shuffleSlice(slice []int64, rng *rand.Rand) {
	// Iterate through the slice in reverse order, starting from the last element.
	for i := len(slice) - 1; i > 0; i-- {
		// Generate a random index 'j' between 0 and i (inclusive).
		j := rng.Intn(i + 1)

		// Swap the elements at indices i and j.
		slice[i], slice[j] = slice[j], slice[i]
	}
}
