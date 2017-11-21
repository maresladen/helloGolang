//Package utils 工具类
package utils

import (
	"crypto/md5"
	"encoding/hex"
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
	primeArr := make([]int64, 78498)
	for i := 1; i <= n; i++ {
		if prime(i) {
			primeArr[nCount] = int64(i)
			nCount++
			// temp := big.NewInt(int64(i))
			// result.Mul(result, temp)
		}
	}

	elapsedPrime := time.Since(tPrime)
	fmt.Println("获取质数耗时:", elapsedPrime)

	tStart := time.Now()

	for i := range primeArr {

		temp := big.NewInt(int64(i))
		result.Mul(result, temp)
	}

	elapsed := time.Since(tStart)
	fmt.Println("所有质数相乘耗时:", elapsed)
	return result.String(), nCount
}

//Md5fun 字符串转换为MD5
func Md5fun(s string) string {
	signByte := []byte(s)
	hash := md5.New()
	hash.Write(signByte)
	return hex.EncodeToString(hash.Sum(nil))
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
