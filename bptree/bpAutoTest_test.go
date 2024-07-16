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
	randomQuantity = 200

	// randomMax represents the maximum value for generating random numbers.
	randomMax = 2000

	// randomMin represents the minimum value for generating random numbers.
	randomMin = 10
)

// Test_Check_BpTree_Automatic is used for automated testing, generating test data with random numbers for B+ tree insertion and deletion.
func Test_Check_BpTree_Automatic(t *testing.T) {
	// Automated random testing for B+ tree.
	t.Run("Automated Testing Section", func(t *testing.T) {
		// Set up randomization.
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)

		// Generate random data for insertion.
		numbersForAdding := make([]int64, randomQuantity)
		for i := 0; i < randomQuantity; i++ {
			num := int64(random.Intn(randomMax-randomMin+1) + randomMin)
			numbersForAdding[i] = num
		}
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
			if deleted == false {
				fmt.Println("Breakpoint: Data deletion not successful. 💢 The number is ", numbersForDeleting[i], i)
			}
			if err != nil {
				fmt.Println("Breakpoint: Deletion encountered an error. 💢 The number is ", numbersForDeleting[i], i)
			}
		}

		fmt.Println()
	})
	// Automated random testing for B+ tree.
	t.Run("Manually Identify B Plus Tree Operation Errors", func(t *testing.T) {
		// 数量二百的例子
		// Generate random data for insertion.
		var randomNumbers = []int64{1176, 1546, 1945, 1951, 76, 1323, 1726, 1536, 155, 1556, 1075, 787, 435, 83, 1921, 132, 442, 623, 1617, 343, 98, 1870, 1881, 678, 739, 759, 237, 1140, 1347, 1729, 1124, 747, 453, 1161, 1407, 1764, 1373, 1867, 547, 973, 901, 1804, 805, 844, 1568, 970, 1975, 1869, 836, 1427, 423, 1535, 386, 1988, 1486, 549, 35, 1377, 752, 1018, 1156, 1719, 1313, 1568, 896, 1642, 560, 1648, 1806, 641, 1275, 1670, 1400, 272, 1981, 1118, 1388, 1907, 910, 1507, 1904, 330, 474, 1911, 223, 1517, 1803, 632, 498, 408, 1446, 1084, 1646, 1176, 1789, 1495, 177, 718, 1138, 647, 1373, 1713, 162, 1320, 1847, 1271, 1873, 1290, 41, 314, 1892, 873, 1410, 1656, 1457, 1511, 851, 822, 748, 1053, 1996, 1121, 901, 382, 677, 766, 1325, 1024, 97, 153, 1993, 1632, 901, 84, 723, 1732, 1346, 1712, 35, 1324, 263, 1439, 280, 221, 909, 1206, 272, 1220, 1006, 312, 400, 789, 403, 690, 418, 262, 975, 481, 1952, 613, 1663, 1429, 1835, 1012, 1862, 1213, 1872, 218, 290, 1551, 1101, 340, 1114, 316, 1676, 1823, 866, 637, 90, 1982, 1006, 569, 338, 1052, 449, 745, 548, 524, 1764, 1640, 315, 497, 432, 994, 1496, 1784, 750, 715, 652, 1145}
		// Generate random data for deletion.
		var shuffledNumbers = []int64{1988, 1290, 1347, 1640, 221, 789, 994, 836, 1712, 1407, 1220, 162, 153, 343, 1535, 1975, 97, 290, 1869, 1617, 1114, 548, 1904, 481, 338, 442, 901, 1006, 1862, 1313, 805, 1121, 1656, 1632, 1719, 1124, 677, 1140, 1024, 1161, 1012, 1410, 1907, 1052, 1764, 315, 272, 1145, 1507, 1732, 223, 340, 715, 498, 844, 98, 35, 745, 1324, 1373, 1118, 272, 1018, 1446, 1429, 1803, 1835, 632, 1325, 652, 1138, 1952, 177, 1427, 418, 132, 314, 896, 1320, 637, 1806, 547, 1804, 787, 1568, 1676, 851, 750, 678, 690, 873, 237, 1323, 1867, 35, 1373, 1911, 1873, 453, 822, 1646, 718, 1993, 1951, 752, 759, 1921, 280, 76, 1789, 1872, 1213, 1377, 432, 1729, 560, 262, 1457, 1556, 423, 975, 449, 1713, 435, 1006, 330, 524, 901, 316, 1275, 497, 766, 613, 155, 403, 1981, 1495, 641, 386, 1670, 1496, 1206, 970, 1075, 1996, 1546, 1517, 910, 1439, 218, 1551, 1346, 866, 1271, 1486, 1642, 1388, 1084, 973, 1648, 647, 1881, 1784, 1176, 1823, 623, 1176, 382, 84, 739, 901, 400, 748, 1536, 1945, 747, 1764, 1156, 1053, 909, 408, 1400, 723, 90, 1726, 1892, 41, 1870, 569, 83, 1568, 1847, 312, 1982, 549, 1663, 1511, 474, 1101, 263}

		// Initialize B plus tree.
		root := NewBpTree(5)

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
			if shuffledNumbers[i] == 960 {
				fmt.Print()
			}

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
