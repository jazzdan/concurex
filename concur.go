package concurex

// interfaces
type address int

type server interface {
	order(a address)
}

type shipper interface {
	ship(a address)
}

// exercise 0
type server0 struct {
	sh shipper
}

func newServer0(sh shipper) server {
	return &server0{
		sh: sh,
	}
}

func (s *server0) order(a address) {
	s.sh.ship(a)
}

// exercise 1
type server1 struct {
	sh    shipper
	queue chan address
}

func newServer1(sh shipper) server {
	queue := make(chan address)
	s := &server1{
		sh:    sh,
		queue: queue,
	}

	go s.loop()

	return s
}

func (s *server1) loop() {
	for a := range s.queue {
		s.sh.ship(a)
	}
}

func (s *server1) order(a address) {
	s.queue <- a
}

type server2 struct {
	sh    shipper
	ch    chan address
	queue []address
}

func newServer2(sh shipper) server {
	ch := make(chan address)
	s := &server2{
		sh: sh,
		ch: ch,
	}

	go s.accept(s.ch)
	ch = nil

	return s
}

func (s *server2) accept(ch chan address) {
	var workerDone chan struct{}
	for {
		if len(s.queue) == 0 && workerDone == nil && ch == nil {
			return
		}

		select {
		case a, ok := <-ch:
			if !ok {
				ch = nil
			} else {
				s.queue = append(s.queue, a)
			}
		case <-workerDone:
			workerDone = nil
		}

		if len(s.queue) > 0 && workerDone == nil {
			workerDone = make(chan struct{})
			var a address
			a, s.queue = s.queue[0], s.queue[1:]
			go s.work(workerDone, a)
		}

	}
}

func (s *server2) work(doneCh chan struct{}, a address) {
	s.sh.ship(a)
	doneCh <- struct{}{}
}

func (s *server2) order(a address) {
	s.ch <- a
}

// exercise 3
type server3 struct {
	sh shipper
}

func newServer3(sh shipper) server {
	return &server3{
		sh: sh,
	}
}

func (s *server3) order(a address) {
	s.sh.ship(a)
}
