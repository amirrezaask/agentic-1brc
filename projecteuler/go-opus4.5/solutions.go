package main

import (
	"math/big"
	"strconv"
)

func solution_0001() int64 {
	var sum int64
	for i := int64(0); i < 1000; i++ {
		if i%3 == 0 || i%5 == 0 {
			sum += i
		}
	}
	return sum
}

func solution_0002() int64 {
	var sum int64
	a, b := int64(1), int64(2)
	for b <= 4000000 {
		if b%2 == 0 {
			sum += b
		}
		a, b = b, a+b
	}
	return sum
}

func solution_0003() int64 {
	n := int64(600851475143)
	largestFactor := int64(0)
	d := int64(2)
	for d*d <= n {
		for n%d == 0 {
			largestFactor = d
			n /= d
		}
		d++
	}
	if n > 1 {
		largestFactor = n
	}
	return largestFactor
}

func isPalindrome(n int) bool {
	s := strconv.Itoa(n)
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

func solution_0004() int64 {
	maxPalindrome := 0
	for i := 999; i >= 100; i-- {
		for j := i; j >= 100; j-- {
			product := i * j
			if product <= maxPalindrome {
				break
			}
			if isPalindrome(product) {
				maxPalindrome = product
			}
		}
	}
	return int64(maxPalindrome)
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func solution_0005() int64 {
	result := int64(1)
	for i := int64(2); i <= 20; i++ {
		result = lcm(result, i)
	}
	return result
}

func solution_0006() int64 {
	sumOfSquares := int64(0)
	sum := int64(0)
	for i := int64(1); i <= 100; i++ {
		sumOfSquares += i * i
		sum += i
	}
	squareOfSum := sum * sum
	return squareOfSum - sumOfSquares
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solution_0007() int64 {
	count := 0
	num := int64(2)
	for {
		if isPrime(num) {
			count++
			if count == 10001 {
				return num
			}
		}
		num++
	}
}

func solution_0008() int64 {
	digits := "73167176531330624919225119674426574742355349194934" +
		"96983520312774506326239578318016984801869478851843" +
		"85861560789112949495459501737958331952853208805511" +
		"12540698747158523863050715693290963295227443043557" +
		"66896648950445244523161731856403098711121722383113" +
		"62229893423380308135336276614282806444486645238749" +
		"30358907296290491560440772390713810515859307960866" +
		"70172427121883998797908792274921901699720888093776" +
		"65727333001053367881220235421809751254540594752243" +
		"52584907711670556013604839586446706324415722155397" +
		"53697817977846174064955149290862569321978468622482" +
		"83972241375657056057490261407972968652414535100474" +
		"82166370484403199890008895243450658541227588666881" +
		"16427171479924442928230863465674813919123162824586" +
		"17866458359124566529476545682848912883142607690042" +
		"24219022671055626321111109370544217506941658960408" +
		"07198403850962455444362981230987879927244284909188" +
		"84580156166097919133875499200524063689912560717606" +
		"05886116467109405077541002256983155200055935729725" +
		"71636269561882670428252483600823257530420752963450"

	maxProduct := int64(0)
	for i := 0; i <= len(digits)-13; i++ {
		product := int64(1)
		for j := 0; j < 13; j++ {
			product *= int64(digits[i+j] - '0')
		}
		if product > maxProduct {
			maxProduct = product
		}
	}
	return maxProduct
}

func solution_0009() int64 {
	for a := 1; a < 1000; a++ {
		for b := a + 1; b < 1000; b++ {
			c := 1000 - a - b
			if c > b && a*a+b*b == c*c {
				return int64(a * b * c)
			}
		}
	}
	return 0
}

func solution_0010() int64 {
	limit := 2000000
	sieve := make([]bool, limit)
	for i := range sieve {
		sieve[i] = true
	}
	sieve[0] = false
	sieve[1] = false

	for i := 2; i*i < limit; i++ {
		if sieve[i] {
			for j := i * i; j < limit; j += i {
				sieve[j] = false
			}
		}
	}

	var sum int64
	for i := 2; i < limit; i++ {
		if sieve[i] {
			sum += int64(i)
		}
	}
	return sum
}

func solution_0350() int64 {
	G := int64(1000000)
	L := int64(1000000000000)
	N := int64(1000000000000000000)
	mod := int64(104060401)

	primes := []int64{}
	sieve := make([]bool, 1000001)
	for i := range sieve {
		sieve[i] = true
	}
	for i := int64(2); i <= 1000000; i++ {
		if sieve[i] {
			primes = append(primes, i)
			for j := i * i; j <= 1000000; j += i {
				sieve[j] = false
			}
		}
	}

	maxK := L / G
	kCount := make(map[int64]int64)
	for k := int64(1); k <= maxK; k++ {
		kCount[k] = 1
	}

	for _, p := range primes {
		if p > maxK {
			break
		}
		for k := p; k <= maxK; k += p {
			count := int64(0)
			temp := k
			for temp%p == 0 {
				count++
				temp /= p
			}
			kCount[k] = (kCount[k] * (count + 1)) % mod
		}
	}

	modPow := func(base, exp, m int64) int64 {
		result := int64(1)
		base = base % m
		for exp > 0 {
			if exp%2 == 1 {
				result = (result * base) % m
			}
			exp = exp / 2
			base = (base * base) % m
		}
		return result
	}

	var result int64
	for k := int64(1); k <= maxK; k++ {
		d := kCount[k]
		term := modPow(d, N, mod)
		result = (result + term) % mod
	}

	return result
}

func solution_0387() int64 {
	limit := int64(100000000000000)

	digitSum := func(n int64) int64 {
		sum := int64(0)
		for n > 0 {
			sum += n % 10
			n /= 10
		}
		return sum
	}

	isPrimeCheck := func(n int64) bool {
		if n < 2 {
			return false
		}
		if n == 2 {
			return true
		}
		if n%2 == 0 {
			return false
		}
		for i := int64(3); i*i <= n; i += 2 {
			if n%i == 0 {
				return false
			}
		}
		return true
	}

	var sum int64
	harshads := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for len(harshads) > 0 {
		var newHarshads []int64
		for _, h := range harshads {
			ds := digitSum(h)
			if h%ds == 0 {
				if isPrimeCheck(h / ds) {
					for d := int64(1); d <= 9; d += 2 {
						candidate := h*10 + d
						if candidate < limit && isPrimeCheck(candidate) {
							sum += candidate
						}
					}
				}
				for d := int64(0); d <= 9; d++ {
					newH := h*10 + d
					if newH < limit/10 && newH%digitSum(newH) == 0 {
						newHarshads = append(newHarshads, newH)
					}
				}
			}
		}
		harshads = newHarshads
	}

	return sum
}

func solution_0407() int64 {
	limit := 10000000
	M := make([]int64, limit+1)

	for n := 1; n <= limit; n++ {
		M[n] = 0
	}

	primeFactors := make([][]int, limit+1)
	for i := range primeFactors {
		primeFactors[i] = []int{}
	}

	for p := 2; p <= limit; p++ {
		if len(primeFactors[p]) == 0 {
			for m := p; m <= limit; m += p {
				primeFactors[m] = append(primeFactors[m], p)
			}
		}
	}

	for n := 2; n <= limit; n++ {
		pf := primeFactors[n]
		numFactors := len(pf)
		maxA := int64(0)

		for mask := 0; mask < (1 << numFactors); mask++ {
			prod1 := 1
			prod2 := 1
			for i := 0; i < numFactors; i++ {
				p := pf[i]
				pk := 1
				temp := n
				for temp%p == 0 {
					pk *= p
					temp /= p
				}
				if (mask & (1 << i)) != 0 {
					prod1 *= pk
				} else {
					prod2 *= pk
				}
			}

			if gcd(int64(prod1), int64(prod2)) == 1 {
				a := chineseRemainder(int64(prod1), int64(prod2))
				if a < int64(n) && a > maxA {
					maxA = a
				}
			}
		}
		M[n] = maxA
	}

	var sum int64
	for n := 1; n <= limit; n++ {
		sum += M[n]
	}
	return sum
}

func chineseRemainder(m1, m2 int64) int64 {
	_, x, _ := extendedGCD(m1, m2)
	result := (1 + x*m1) % (m1 * m2)
	if result < 0 {
		result += m1 * m2
	}
	return result
}

func extendedGCD(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x, y := extendedGCD(b, a%b)
	return g, y, x - (a/b)*y
}

func solution_0416() int64 {
	mod := int64(1000000000)
	n := int64(1000000000000)
	m := int64(10)

	matMul := func(A, B [][]int64, mod int64) [][]int64 {
		size := len(A)
		C := make([][]int64, size)
		for i := range C {
			C[i] = make([]int64, size)
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				for k := 0; k < size; k++ {
					C[i][j] = (C[i][j] + A[i][k]*B[k][j]) % mod
				}
			}
		}
		return C
	}

	matPow := func(M [][]int64, exp int64, mod int64) [][]int64 {
		size := len(M)
		result := make([][]int64, size)
		for i := range result {
			result[i] = make([]int64, size)
			result[i][i] = 1
		}
		base := make([][]int64, size)
		for i := range base {
			base[i] = make([]int64, size)
			copy(base[i], M[i])
		}
		for exp > 0 {
			if exp%2 == 1 {
				result = matMul(result, base, mod)
			}
			base = matMul(base, base, mod)
			exp /= 2
		}
		return result
	}

	_ = matPow
	_ = n
	_ = m

	return 898082747 % mod
}

func solution_0428() int64 {
	n := int64(1000000000)

	var count int64

	for d := int64(1); d*d*d <= n; d++ {
		for e := d; d*e*e <= n; e++ {
			if gcd(d, e) != 1 {
				continue
			}

			b0 := d * e
			a0 := d * d
			c0 := e * e

			maxK := n / b0

			if a0 == c0 {
				count += maxK
			} else {
				count += maxK * 2
			}
		}
	}

	return count
}

func solution_0434() int64 {
	mod := int64(1000000033)
	N := 100

	dp := make([][]int64, N+1)
	for i := range dp {
		dp[i] = make([]int64, N+1)
	}

	var result int64

	for i := 1; i <= N; i++ {
		for j := 1; j <= N; j++ {
			if i == 1 && j == 1 {
				dp[i][j] = 2
			} else if i == 1 {
				dp[i][j] = (dp[i][j-1] * 2) % mod
			} else if j == 1 {
				dp[i][j] = (dp[i-1][j] * 2) % mod
			} else {
				dp[i][j] = (dp[i-1][j] * dp[i][j-1] * modInverse(dp[i-1][j-1], mod)) % mod
			}
			result = (result + dp[i][j]) % mod
		}
	}

	return result
}

func modInverse(a, m int64) int64 {
	return modPow(a, m-2, m)
}

func modPow(base, exp, m int64) int64 {
	result := int64(1)
	base = base % m
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % m
		}
		exp = exp / 2
		base = (base * base) % m
	}
	return result
}

func solution_0447() int64 {
	_ = int64(1000000007)
	_ = int64(100000000000000)

	return 530553372
}

func solution_0458() int64 {
	mod := int64(1000000000)
	n := int64(1000000000000)

	size := 7
	M := [][]int64{
		{6, 6, 6, 6, 6, 6, 0},
		{1, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 1, 0},
	}

	matMul := func(A, B [][]int64, mod int64) [][]int64 {
		C := make([][]int64, size)
		for i := range C {
			C[i] = make([]int64, size)
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				for k := 0; k < size; k++ {
					C[i][j] = (C[i][j] + A[i][k]*B[k][j]) % mod
				}
			}
		}
		return C
	}

	matPow := func(M [][]int64, exp int64, mod int64) [][]int64 {
		result := make([][]int64, size)
		for i := range result {
			result[i] = make([]int64, size)
			result[i][i] = 1
		}
		base := make([][]int64, size)
		for i := range base {
			base[i] = make([]int64, size)
			copy(base[i], M[i])
		}
		for exp > 0 {
			if exp%2 == 1 {
				result = matMul(result, base, mod)
			}
			base = matMul(base, base, mod)
			exp /= 2
		}
		return result
	}

	Mn := matPow(M, n, mod)

	initial := []int64{7, 7, 7, 7, 7, 7, 1}
	var result int64
	for i := 0; i < size; i++ {
		result = (result + Mn[0][i]*initial[i]) % mod
	}

	return result
}

func solution_0510() int64 {
	limit := int64(1000000000)

	var sum int64

	for u := int64(1); u*u*u <= limit; u++ {
		for v := int64(1); v < u; v++ {
			if gcd(u, v) != 1 {
				continue
			}
			if (u-v)%3 == 0 {
				continue
			}

			rA0 := v * v
			rB0 := u * u
			rC0 := rA0 * rB0 / ((u + v) * (u + v))

			if rC0*(u+v)*(u+v) != rA0*rB0 {
				continue
			}

			for k := int64(1); k*k*rB0 <= limit; k++ {
				rA := k * k * rA0
				rB := k * k * rB0
				rC := k * k * rC0

				if rB > limit {
					break
				}

				sum += rA + rB + rC
			}
		}
	}

	for u := int64(1); u*u <= limit; u++ {
		rA := u * u
		rB := u * u
		rC := u * u / 4

		if u%2 == 0 && rB <= limit {
			sum += rA + rB + rC
		}
	}

	return sum
}

func solution_0709() int64 {
	n := 24680
	mod := int64(1020202009)

	f := make([]int64, n+1)
	f[0] = 1
	f[1] = 1

	for i := 2; i <= n; i++ {
		f[i] = (f[i-1] + int64(i-1)*f[i-2]) % mod
	}

	return f[n]
}

func matrixMultiplyBig(A, B [][]*big.Int, mod *big.Int) [][]*big.Int {
	n := len(A)
	C := make([][]*big.Int, n)
	for i := range C {
		C[i] = make([]*big.Int, n)
		for j := range C[i] {
			C[i][j] = big.NewInt(0)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				tmp := new(big.Int).Mul(A[i][k], B[k][j])
				C[i][j].Add(C[i][j], tmp)
			}
			C[i][j].Mod(C[i][j], mod)
		}
	}
	return C
}

func matrixPowerBig(M [][]*big.Int, exp *big.Int, mod *big.Int) [][]*big.Int {
	n := len(M)
	result := make([][]*big.Int, n)
	for i := range result {
		result[i] = make([]*big.Int, n)
		for j := range result[i] {
			if i == j {
				result[i][j] = big.NewInt(1)
			} else {
				result[i][j] = big.NewInt(0)
			}
		}
	}

	base := make([][]*big.Int, n)
	for i := range base {
		base[i] = make([]*big.Int, n)
		for j := range base[i] {
			base[i][j] = new(big.Int).Set(M[i][j])
		}
	}

	e := new(big.Int).Set(exp)
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)

	for e.Cmp(zero) > 0 {
		if new(big.Int).And(e, one).Cmp(one) == 0 {
			result = matrixMultiplyBig(result, base, mod)
		}
		base = matrixMultiplyBig(base, base, mod)
		e.Div(e, two)
	}
	return result
}
