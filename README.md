[![CircleCI](https://circleci.com/gh/alexj50/worker/tree/master.svg?style=svg)](https://circleci.com/gh/alexj50/worker/tree/master)
[![codecov.io Code Coverage](https://codecov.io/gh/alexj50/worker/branch/master/graph/badge.svg)](https://codecov.io/github/alexj50/worker?branch=master)
[![Code Climate](https://codeclimate.com/github/alexj50/worker/badges/gpa.svg)](https://codeclimate.com/github/alexj50/worker)

# Worker

A Go routine buffer to speed up async jobs

## Getting Started

```
go get github.com/alexj50/worker
```

### Initializing

```
import "github.com/alexj50/worker"

worker := &Worker {
    maxWorkers: 10,
    maxQueue:   10,
    testing:    t, // *testing.T (optional), only for tests
}
worker.Start()
```

### Add to the Queue

```
func someFunc(){
    job := test{number: 5}
    AddJob(job)
}
```

implement the interface methods

```
type test struct {
    worker.Job // add if out of module
    number int
}

// perform will get called
func (t test) perform() {
    // do something
}

func (t test) performTest(test *testing.T) {
    // do the test
}
```
## Stopping the Queue

This will stop the workers and wait until all jobs have finished

```
GracefulShutdown()
```

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
