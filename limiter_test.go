package go_limiter

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter1 := NewLimiter(10*time.Second, 2)
	for i:=0; i<20;i++ {
		res, err := limiter1.Check()
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.IsReached {
			fmt.Printf("Wait %.2f seconds to unlock\n", res.TimeLeftToUnlock.Seconds())
		} else {
			fmt.Println("Limit not reached")
		}
		time.Sleep(1*time.Second)
	}
	err := limiter1.Kill()

	if err != nil {
		t.Fail()
	}

	res, err := limiter1.Check()

	if err == nil {
		t.Fail()
	}

	if res != nil {
		t.Fail()
	}

}
