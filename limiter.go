package go_limit

import (
	"errors"
	"fmt"
	"time"
)

type Limiter struct {
	ticker       *time.Ticker
	limit        chan int
	req          chan int
	reached      bool
	timeToUnlock time.Time
	killLimiter  chan bool
	active       bool
	killedTime   time.Time
}

type LimiterResponse struct {
	IsReached        bool
	TimeLeftToUnlock time.Duration
}


func NewLimiter(timeInterval time.Duration, limitNumber int) *Limiter {
	limit := &Limiter{}
	limit.ticker = time.NewTicker(timeInterval)
	limit.limit = make(chan int, limitNumber)
	limit.req = make(chan int,1)
	limit.timeToUnlock = time.Now().Add(timeInterval)
	limit.active = true
	limit.killLimiter = make(chan bool, 1)
	go limit.run(timeInterval, limitNumber)
	return limit
}

func (limit *Limiter) run(timeInterval time.Duration, limitNumber int) {
	for {
		select {
		case <-limit.ticker.C:
			limit.limit = make(chan int, limitNumber)
			limit.timeToUnlock = time.Now().Add(10 * time.Second)
		case <-limit.req:
			if len(limit.limit) < cap(limit.limit) {
				limit.limit <- 1
				limit.reached = false
				continue
			}
			limit.reached = true
		case <-limit.killLimiter:
			return
		}
	}
}

func (limit *Limiter) Check() (*LimiterResponse, error) {
	limiterRes := LimiterResponse{}
	if !limit.active {
		errMsg := fmt.Sprintf("limiter was killed at %s", limit.killedTime)
		return nil, errors.New(errMsg)
	}
	limit.req <- 1
	time.Sleep(time.Millisecond)
	limiterRes.IsReached = limit.reached
	limiterRes.TimeLeftToUnlock = limit.timeToUnlock.Sub(time.Now())
	return &limiterRes, nil
}

func (limit *Limiter) Kill() error {
	if !limit.active {
		return errors.New("limiter is already killed")
	}
	limit.active = false
	limit.killedTime = time.Now()
	limit.killLimiter <- true
	return nil
}