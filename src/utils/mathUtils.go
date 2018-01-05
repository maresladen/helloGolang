//Package utils 工具类
package utils

import (
	"fmt"
	"math"
	"math/big"
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
