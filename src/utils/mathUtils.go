//Package utils 工具类
package utils

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

//MaxBigint 获取最大素数可能值
func MaxBigint(n int) (string, int) {
	var nCount int
	result := big.NewInt(1)

	tPrime := time.Now()
	primeArr := make([]int64, 100)
	for i := 1; i <= n; i++ {
		if prime(i) {
			if nCount < len(primeArr) {
				primeArr[nCount] = int64(i)
			} else {
				primeArr = append(primeArr, int64(i))

			}
			nCount++
		}
	}

	fmt.Println()
	elapsedPrime := time.Since(tPrime)

	//can be del--------打印素数(这里可以做个txt的保存)-------
	// for i, v := range primeArr {
	// 	fmt.Printf("%8d", v)
	// 	if i%8 == 0 {
	// 		fmt.Println()
	// 	}
	// }

	// fmt.Println()
	//can be del--------打印素数(这里可以做个txt的保存)-------
	fmt.Println("获取质数耗时:", elapsedPrime)
	tStart := time.Now()

	for index := 0; index < len(primeArr); index++ {
		i := primeArr[index]
		if i == 0 {
			break
		}

		temp := big.NewInt(i)
		result.Mul(result, temp)
	}

	elapsed := time.Since(tStart)
	fmt.Println("所有质数相乘耗时:", elapsed)
	return result.String(), nCount
}

//CreatePrimeMul 根据用户传入的int数组计算md5哈希值
func CreatePrimeMul(nums ...int) string {
	result := big.NewInt(1)
	for _, num := range nums {
		temp := big.NewInt(int64(num))
		result.Mul(result, temp)
	}
	return Md5ByteArr(result.Bytes())
}

var tenToAny = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}

//DecimalToAny 10进制转任意进制
func DecimalToAny(num, n int) string {
	var newNumstr string
	var remainder int
	var remainderString string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remainderString = tenToAny[remainder]
		} else {
			remainderString = strconv.Itoa(remainder)
		}
		newNumstr = remainderString + newNumstr
		num = num / n
	}
	return newNumstr
}

//AnyToDecimal 任意进制转10进制
func AnyToDecimal(num string, n int) int {
	var newNum float64
	newNum = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(findkey(value))
		if tmp != -1 {
			newNum = newNum + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(newNum)
}

//------------私有方法------------------

func prime(value int) bool {
	if value <= 1 {
		return false
	}
	if value == 2 || value == 3 || value == 5 || value == 7 {
		return true
	}
	if value%2 == 0 || value%3 == 0 || value%5 == 0 || value%7 == 0 {
		return false
	}
	factor := 7
	c := []int{4, 2, 4, 2, 4, 6, 2, 6}
	max := int(math.Sqrt(float64(value)))
	if max*max == value {
		return false
	}
	for factor < max {
		for i := 0; i < len(c); i++ {
			factor += c[i]
			if value%factor == 0 {
				return false
			}
		}
	}
	return true
}

// map根据value找key
func findkey(in string) int {
	result := -1
	for k, v := range tenToAny {
		if in == v {
			result = k
		}
	}
	return result
}
