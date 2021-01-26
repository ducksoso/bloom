package bloom

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

func TestBasic1(t *testing.T) {
	bloom, _ := Create(1000, 0.01)
	bloom.Put([]byte("10086"))
	if bloom.MightContain([]byte("10086")) {
		fmt.Println("should be in")
	}

	for i := 0; i < 1000; i++ {
		n := make([]byte, 4)
		binary.BigEndian.PutUint32(n, uint32(i))
		bloom.Put(n)
	}

	count := 0
	for i := 950; i < 1100; i++ {
		n := make([]byte, 4)
		binary.BigEndian.PutUint32(n, uint32(i))
		if bloom.MightContain(n) {
			count++
			fmt.Println("误判了")
		}
	}
	fmt.Println("总共误判：", count)
}

func TestBasic(t *testing.T) {
	f := initBloomFilter(1000, 4)
	n1 := []byte("Bess")
	n2 := []byte("Jane")
	n3 := []byte("Emma")
	f.Put(n1)
	n3a := f.TestAndAdd(n3)
	n1b := f.MightContain(n1)
	n2b := f.MightContain(n2)
	n3b := f.MightContain(n3)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3a {
		t.Errorf("%v should not be in the first time we look.", n3)
	}
	if !n3b {
		t.Errorf("%v should be in the second time we look.", n3)
	}
}

func TestBasicUint32(t *testing.T) {
	f := initBloomFilter(1000, 4)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	binary.BigEndian.PutUint32(n4, 103)
	f.Put(n1)
	n3a := f.TestAndAdd(n3)
	n1b := f.MightContain(n1)
	n2b := f.MightContain(n2)
	n3b := f.MightContain(n3)
	f.MightContain(n4)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3a {
		t.Errorf("%v should not be in the first time we look.", n3)
	}
	if !n3b {
		t.Errorf("%v should be in the second time we look.", n3)
	}
}

func TestNewWithLowNumbers(t *testing.T) {
	f := initBloomFilter(0, 0)
	if f.k != 1 {
		t.Errorf("%v should be 1", f.k)
	}
	if f.m != 1 {
		t.Errorf("%v should be 1", f.m)
	}
}

func TestString(t *testing.T) {
	f := NewWithEstimates(1000, 0.001)
	n1 := "Love"
	n2 := "is"
	n3 := "in"
	n4 := "bloom"
	f.AddString(n1)
	n3a := f.TestAndAddString(n3)
	n1b := f.TestString(n1)
	n2b := f.TestString(n2)
	n3b := f.TestString(n3)
	f.TestString(n4)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3a {
		t.Errorf("%v should not be in the first time we look.", n3)
	}
	if !n3b {
		t.Errorf("%v should be in the second time we look.", n3)
	}

}

func testEstimated(n uint, maxFp float64, t *testing.T) {
	m, k := EstimateParameters(n, maxFp)
	f := NewWithEstimates(n, maxFp)
	fpRate := f.EstimateFalsePositiveRate(n)
	if fpRate > 1.5*maxFp {
		t.Errorf("False positive rate too high: n: %v; m: %v; k: %v; maxFp: %f; fpRate: %f, fpRate/maxFp: %f", n, m, k, maxFp, fpRate, fpRate/maxFp)
	}
}

func TestEstimated1000_0001(t *testing.T)   { testEstimated(1000, 0.000100, t) }
func TestEstimated10000_0001(t *testing.T)  { testEstimated(10000, 0.000100, t) }
func TestEstimated100000_0001(t *testing.T) { testEstimated(100000, 0.000100, t) }

func TestEstimated1000_001(t *testing.T)   { testEstimated(1000, 0.001000, t) }
func TestEstimated10000_001(t *testing.T)  { testEstimated(10000, 0.001000, t) }
func TestEstimated100000_001(t *testing.T) { testEstimated(100000, 0.001000, t) }

func TestEstimated1000_01(t *testing.T)   { testEstimated(1000, 0.010000, t) }
func TestEstimated10000_01(t *testing.T)  { testEstimated(10000, 0.010000, t) }
func TestEstimated100000_01(t *testing.T) { testEstimated(100000, 0.010000, t) }

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// The following function courtesy of Nick @turgon
// This helper function ranges over the input data, applying the hashing
// which returns the bit locations to set in the filter.
// For each location, increment a counter for that bit address.
//
// If the Bloom Filter's location() method distributes locations uniformly
// at random, a property it should inherit from its hash function, then
// each bit location in the filter should end up with roughly the same
// number of hits.  Importantly, the value of k should not matter.
//
// Once the results are collected, we can run a chi squared goodness of fit
// test, comparing the result histogram with the uniform distribition.
// This yields a test statistic with degrees-of-freedom of m-1.
func chiTestBloom(m, k, rounds uint, elements [][]byte) (succeeds bool) {
	f := initBloomFilter(m, k)
	results := make([]uint, m)
	chi := make([]float64, m)

	for _, data := range elements {
		h := baseHashes(data)
		for i := uint(0); i < f.k; i++ {
			results[f.location(h, i)]++
		}
	}

	// Each element of results should contain the same value: k * rounds / m.
	// Let's run a chi-square goodness of fit and see how it fares.
	var chiStatistic float64
	e := float64(k*rounds) / float64(m)
	for i := uint(0); i < m; i++ {
		chi[i] = math.Pow(float64(results[i])-e, 2.0) / e
		chiStatistic += chi[i]
	}

	// this tests at significant level 0.005 up to 20 degrees of freedom
	table := [20]float64{
		7.879, 10.597, 12.838, 14.86, 16.75, 18.548, 20.278,
		21.955, 23.589, 25.188, 26.757, 28.3, 29.819, 31.319, 32.801, 34.267,
		35.718, 37.156, 38.582, 39.997}
	df := min(m-1, 20)

	succeeds = table[df-1] > chiStatistic
	return

}
