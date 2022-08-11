package ratectrl

import (
	"time"
)

func chanSize(total int) int {
	l := total / 5
	if l > 100 {
		l = 100
	} else if l < 5 {
		l = 5
	}
	return l
}

//Rectifier 整流器
func Rectifier[T any](src chan T, window time.Duration, speed int) chan T {

	ch := make(chan T, chanSize(speed))
	go func() {
		defer close(ch)
		var count = 0
		barrier := time.Now()
		for d := range src {
			if time.Now().Sub(barrier) > window {
				count = 0
				barrier = time.Now()
			}
			ch <- d
			count++
			if count >= speed {
				cc := window - time.Now().Sub(barrier)
				if cc > 0 {
					<-time.After(cc)
				}
				count = 0
				barrier = time.Now()
			}
		}
	}()

	return ch
}
