package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	var wg sync.WaitGroup
	defer wg.Wait()

	ids := make(chan int64, pool)
	defer close(ids)

	users := make(chan user, n)
	defer close(users)

	for p := int64(0); p < pool && p < n; p++ {
		wg.Add(1)
		go func(ids <-chan int64, users chan<- user, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := range ids {
				users <- getOne(i)
			}
		}(ids, users, &wg)
	}

	for i := int64(0); i < n; i++ {
		ids <- i
	}

	for i := int64(0); i < n; i++ {
		res = append(res, <-users)
	}

	return res
}
