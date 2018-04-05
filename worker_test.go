package worker

import (
	"testing"
	"fmt"
)

const integer = 5

type test struct {
	number uint
}

func start(t *testing.T)  {
	worker := Worker {
		MaxWorkers:integer,
		MaxQueue:integer,
		Testing: t,
	}
	worker.Start()
}

func TestStartWorkers(t *testing.T) {
	start(t)

	if len(workers) != integer {
		t.Errorf("Expected %d workers got %d", integer, len(workers))
	}
}

func TestAddJob(t *testing.T) {
	for i := 0; i < 100; i++ {
		c := test{number: 5}
		AddJob(c)
	}

	GracefulShutdown()
	start(nil)

	for i := 0; i < 100; i++ {
		c := test{number: 5}
		AddJob(c)
	}
}

func (t test) perform() {
	if integer << integer != t.number << t.number{
		fmt.Errorf("Expected %d workers got %d", integer * integer, t.number << t.number)
	}
}

func (t test) performTest(test *testing.T) {
	if integer << integer != t.number << t.number{
		test.Errorf("Expected %d workers got %d", integer * integer, t.number << t.number)
	}
}

func TestGracefulShutdown(t *testing.T) {
	GracefulShutdown()
	if len(workers) != 0 {
		t.Errorf("Expected %d workers got %d", 0, len(workers))
	}
}