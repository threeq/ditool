package ratectrl

import (
	"fmt"
	"testing"
	"time"
)

func TestRectifier(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		var arr = make([]int, 100)
		for i := 0; i < len(arr); i++ {
			arr[i] = i
		}
		testHelper(t, arr)
	})

	t.Run("string", func(t *testing.T) {
		var arr = make([]string, 100)
		for i := 0; i < len(arr); i++ {
			arr[i] = fmt.Sprintf("s%v", i)
		}
		testHelper(t, arr)
	})

}

func testHelper[T any](t *testing.T, arr []T) {
	ch := make(chan T, 10)
	go func() {
		for i := 0; i < len(arr); i++ {
			ch <- arr[i]
		}
		close(ch)
	}()

	ch = Rectifier(ch, 3*time.Second, 15)
	ch = Rectifier(ch, 1*time.Second, 10)

	for data := range ch {
		t.Logf("%v", data)
	}
}
