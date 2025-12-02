package main

func solution_0001() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		if i%3 == 0 || i%5 == 0 {
			sum += i
		}
	}
	return sum
}

func solution_0002() int {
	sum := 0
	a, b := 1, 2
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
	factor := int64(2)
	for factor*factor <= n {
		for n%factor == 0 {
			n = n / factor
		}
		factor++
	}
	return n
}

func isPalindrome(n int) bool {
	original := n
	reversed := 0
	for n > 0 {
		reversed = reversed*10 + n%10
		n = n / 10
	}
	return original == reversed
}

func solution_0004() int {
	maxPalindrome := 0
	for i := 100; i < 1000; i++ {
		for j := i; j < 1000; j++ {
			product := i * j
			if isPalindrome(product) && product > maxPalindrome {
				maxPalindrome = product
			}
		}
	}
	return maxPalindrome
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return (a * b) / gcd(a, b)
}

func solution_0005() int64 {
	result := int64(1)
	for i := int64(2); i <= 20; i++ {
		result = lcm(result, i)
	}
	return result
}

func solution_0006() int {
	sumOfSquares := 0
	sum := 0
	for i := 1; i <= 100; i++ {
		sumOfSquares += i * i
		sum += i
	}
	squareOfSum := sum * sum
	return squareOfSum - sumOfSquares
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solution_0007() int {
	count := 0
	n := 2
	for {
		if isPrime(n) {
			count++
			if count == 10001 {
				return n
			}
		}
		n++
	}
}

func solution_0008() int64 {
	digits := "7316717653133062491922511967442657474235534919493496983520312774506326239578318016984801869478851843858615607891129494954595017379583319528532088055111254069874715852386305071569329096329522744304355766896648950445244523161731856403098711121722383113622298934233803081353362766142828064444866452387493035890729629049156044077239071381051585930796086670172427121883998797908792274921901699720888093776657273330010533678812202354218097512545405947522435258490771167055601360483958644670632441572215539753697817977846174064955149290862569321978468622482839722413756570560574902614079729686524145351004748216637048440319989000889524345065854122758866688116427171479924442928230863465674813919123162824586178664583591245665294765456828489128831426076900422421902267105562632111110937054421750694165896040807198403850962455444362981230987879927244284909188845801561660979191338754992005240636899125607176060588611646710940507754100225698315520005593572972571636269561882670428252483600823257530420752963450"
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

func solution_0009() int {
	for a := 1; a < 1000; a++ {
		for b := a + 1; b < 1000-a; b++ {
			c := 1000 - a - b
			if a*a+b*b == c*c {
				return a * b * c
			}
		}
	}
	return 0
}

func sieveOfEratosthenes(limit int) []bool {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func solution_0010() int64 {
	limit := 2000000
	isPrime := sieveOfEratosthenes(limit - 1)
	sum := int64(0)
	for i := 2; i < limit; i++ {
		if isPrime[i] {
			sum += int64(i)
		}
	}
	return sum
}

func isPrimeLarge(n int64) bool {
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

func solution_0387() int64 {
	limit := int64(100000000000000)
	sum := int64(0)
	
	type HarshadInfo struct {
		num int64
		ds  int64
	}
	
	queue := []HarshadInfo{}
	for i := int64(1); i <= 9; i++ {
		queue = append(queue, HarshadInfo{i, i})
	}
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		num := current.num
		digitSum := current.ds
		
		isStrong := false
		quotient := num / digitSum
		if quotient > 1 {
			isStrong = isPrimeLarge(quotient)
		}
		
		if isStrong {
			for d := int64(1); d <= 9; d += 2 {
				prime := num*10 + d
				if prime < limit && isPrimeLarge(prime) {
					sum += prime
				}
			}
		}
		
		for d := int64(0); d <= 9; d++ {
			next := num*10 + d
			nextDS := digitSum + d
			if next < limit/10 && next%nextDS == 0 {
				queue = append(queue, HarshadInfo{next, nextDS})
			}
		}
	}
	
	return sum
}

func gcdInt64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solution_0407() int64 {
	limit := int64(10000000)
	M := make([]int64, limit+1)
	
	for n := int64(1); n <= limit; n++ {
		M[n] = 0
	}
	
	for d := int64(1); d <= limit; d++ {
		for n := d; n <= limit; n += d {
			if n%d == 0 {
				candidate := n - d
				if candidate >= 0 && candidate < n {
					a := candidate
					if (a*a)%n == a%n && a > M[n] {
						M[n] = a
					}
				}
			}
		}
	}
	
	for d := int64(1); d <= limit; d++ {
		for n := d; n <= limit; n += d {
			if n%d == 0 {
				a := d
				if a < n && (a*a)%n == a%n && a > M[n] {
					M[n] = a
				}
			}
		}
	}
	
	sum := int64(0)
	for n := int64(1); n <= limit; n++ {
		sum += M[n]
	}
	return sum
}

func gcdInt(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solution_0510() int64 {
	return 0
}

func solution_0350() int64 {
	return 0
}

func solution_0416() int64 {
	return 0
}

func solution_0428() int64 {
	return 0
}

func solution_0434() int64 {
	return 0
}

func solution_0447() int64 {
	return 0
}

func solution_0458() int64 {
	return 0
}

func solution_0709() int64 {
	return 0
}
