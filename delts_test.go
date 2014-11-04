package delts

import (
	"fmt"
	"math"
	"testing"
)

func ExampleSortedDeltaStream() {

	s := NewSortedDeltaStream(1 << 20)

	for i := int64(1); i <= 20; i++ {
		s.Add(i)
	}

	s.Close()

	for d := range s.Output {
		fmt.Println(d)
	}

	// Output:
	// [1 20]

}

func ExampleReverseSortedDeltaStream() {

	s := NewSortedDeltaStream(1 << 20)

	for i := int64(20); i > 0; i-- {
		s.Add(i)
	}

	s.Close()

	for d := range s.Output {
		fmt.Println(d)
	}

	// Output:
	// [20]
	// [1 19]

}

func ExampleBreakInRangeSortedDeltaStream() {

	s := NewSortedDeltaStream(10)

	for i := int64(1); i <= 30; i++ {
		if i%10 == 0 {
			continue
		}
		s.Add(i)
	}

	s.Close()

	for d := range s.Output {
		fmt.Println(d)
	}

	// Output:
	// [1 9]
	// [11 19]
	// [21 29]

}

func BenchmarkSortedDeltaStream(b *testing.B) {

	s := NewSortedDeltaStream(100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Add(int64(i))
	}

}

func BenchmarkReverseSortedDeltaStream(b *testing.B) {

	s := NewSortedDeltaStream(100)

	go func() {
		for {
			<-s.Output
		}
	}()

	b.ResetTimer()

	var n int64 = math.MaxInt64

	for i := 0; i < b.N; i++ {
		s.Add(n)
		n--
	}

}

func BenchmarkBreakInRangeSortedDeltaStream(b *testing.B) {

	// assume 30k rps at 10 sec window
	s := NewSortedDeltaStream(300000)

	go func() {
		for {
			<-s.Output
		}
	}()

	b.ResetTimer()

	var n int64
	for i := 0; i < b.N; i++ {
		// skipping every 100th item
		if n%100 == 0 {
			n++
		}
		s.Add(n)
	}

}
