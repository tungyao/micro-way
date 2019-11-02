package gate_way

import (
	"log"
)

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	t := &Task{f: f}
	go t.Execute()
	return t
}
func (t *Task) Execute() {
	err := t.f()
	if err != nil {
		log.Panic(err)
	}
}

type FPool struct {
	EntryChannel  chan *Task
	InnerChanel   chan *Task
	MaxWorkNumber int
	Stop          chan bool
}

func NewPool(cap int) *FPool {
	p := FPool{
		EntryChannel:  make(chan *Task),
		InnerChanel:   make(chan *Task),
		Stop:          make(chan bool),
		MaxWorkNumber: cap,
	}
	return &p
}
func (p *FPool) Worker(worked int) {
	for task := range p.InnerChanel {
		task.Execute()
	}
}
func (p *FPool) Close() {
	close(p.EntryChannel)
	close(p.InnerChanel)
}

func (p *FPool) Run() {
	for i := 0; i < p.MaxWorkNumber; i++ {
		go p.Worker(i)
	}
	for task := range p.EntryChannel {
		p.InnerChanel <- task
	}

}
