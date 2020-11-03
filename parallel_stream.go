package stream

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type parallelStage struct {
	baseStage
	pumper   chan interface{}
	done     chan bool
	routines int
}

const (
	parallelBase = 16
)

func (p *parallelStage) begin(size int) {
	if size > 0 {
		p.pumper = make(chan interface{}, size) // buffered channel
	} else {
		p.pumper = make(chan interface{})
	}
	p.startLoops(size)
}

func (p *parallelStage) accept(t interface{}) {
	p.pumper <- t
}

func (p *parallelStage) end() {
	close(p.pumper)
	for i := 0; i < p.routines; i++ {
		<-p.done
	}
	p.downStream.end()
}

func (p *parallelStage) startLoops(size int) {
	routines := size/parallelBase + 1
	p.routines = routines
	p.done = make(chan bool)
	for i := 0; i < routines; i++ {
		go p.looper()
	}
}

func (p *parallelStage) looper() {
	for {
		v := <-p.pumper
		if v == nil || p.cancellationRequested() {
			break
		}
		//fmt.Printf("goroutine %d\n", GoID())
		p.downStream.accept(v)
	}
	p.done <- true
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
