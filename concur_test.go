package concurex

import (
	// "log"
	"sync"
	"testing"
)

type h struct {
	f func(a address)
}

func (h *h) ship(a address) {
	h.f(a)
}

func Test0(t *testing.T) {
	wg := new(sync.WaitGroup)

	wg.Add(100)

	f := func(a address) {
		wg.Done()
	}

	s := newServer0(&h{f: f})

	for i := address(0); i < 100; i++ {
		go s.order(i)
	}

	wg.Wait()
}

func Test1(t *testing.T) {
	sum := 0

	wg := new(sync.WaitGroup)

	wg.Add(100)

	f := func(a address) {
		sum += int(a)
		wg.Done()
	}

	s := newServer1(&h{f: f})

	for i := address(0); i < 100; i++ {
		i := i
		go s.order(i)
	}

	wg.Wait()

	if sum != 4950 {
		t.Fatalf("sum was %v; expected %v", sum, 4950)
	}
}

func Test2(t *testing.T) {
	sum := 0

	orderWg := new(sync.WaitGroup)
	orderWg.Add(100)
	shipDoneWg := new(sync.WaitGroup)
	shipDoneWg.Add(100)
	unfreezeShipWg := new(sync.WaitGroup)
	unfreezeShipWg.Add(1)

	f := func(a address) {
		unfreezeShipWg.Wait()
		sum += int(a)
		shipDoneWg.Done()
	}

	s := newServer2(&h{f: f})
	for i := address(0); i < 100; i++ {
		i := i
		go func() {
			s.order(i)
			orderWg.Done()
		}()
	}

	orderWg.Wait()
	unfreezeShipWg.Done()
	orderWg.Wait()
	shipDoneWg.Wait()

	if sum != 4950 {
		t.Fatalf("sum was %v; expected %v", sum, 4950)
	}
}

type outOfOrder struct {
	a    address
	last address
}

func Test3(t *testing.T) {
	sum := 0
	last := address(-1)
	var ooo *outOfOrder

	orderWg := new(sync.WaitGroup)
	orderWg.Add(100)
	shipDoneWg := new(sync.WaitGroup)
	shipDoneWg.Add(100)
	unfreezeShipWg := new(sync.WaitGroup)
	unfreezeShipWg.Add(1)

	f := func(a address) {
		unfreezeShipWg.Wait()
		if ooo == nil && !(a > last) {
			ooo = &outOfOrder{a, last}
		}
		last = a
		sum += int(a)
		shipDoneWg.Done()
	}

	s := newServer3(&h{f: f})
	for i := address(0); i < 100; i++ {
		i := i
		go func() {
			s.order(i)
			orderWg.Done()
		}()
	}

	orderWg.Wait()
	unfreezeShipWg.Done()
	orderWg.Wait()
	shipDoneWg.Wait()

	if sum != 4950 {
		t.Fatalf("sum was %v; expected %v", sum, 4950)
	}
	if ooo != nil {
		t.Fatalf("elem out of order: %v is not less than %v", ooo.a, ooo.last)
	}
}
