package min

import (
	"fmt"
	"testing"
)

func TestGetMin(t *testing.T) {
	ans := GetMin(10, 3)
	if ans != 3 {
		t.Errorf("GetMin(10,3) = %d, want 3", ans)
	}
	ans = GetMin(3, 10)
	if ans != 3 {
		t.Errorf("GetMin(3,10) = %d, want 3", ans)
	}
}

func TestGetMinTable(t *testing.T) {
	var tests = []struct {
		a, b, want int
	}{{10, 3, 3}, {5, 50, 5}, {-4, -10, -10}}
	for _, item := range tests {
		testName := fmt.Sprintf("%d,%d", item.a, item.b)
		testFunc := func(t *testing.T) {
			ans := GetMin(item.a, item.b)
			if ans != item.want {
				t.Errorf("GetMin(%d,%d) = %d, want %d", item.a, item.b, ans, item.want)
			}
		}
		t.Run(testName, testFunc)
	}
}

func BenchmarkGetMin(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMin(100, 53)
	}
}

func benchmarkGetMin(x, y int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetMin(x, y)
	}
}

func BenchmarkGetMinA(b *testing.B) { benchmarkGetMin(10, 3, b) }

func BenchmarkGetMinB(b *testing.B) { benchmarkGetMin(-10, -3, b) }
