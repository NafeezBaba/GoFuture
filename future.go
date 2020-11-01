package GoFuture

import (
	"context"
	"errors"
	"fmt"
	"time"
)

/**
*   A Future represents the result of an asynchronous computation. Methods are provided to check if the computation is complete,
*	to wait for its completion, and to retrieve the result of the computation. The result can only be retrieved using method get
*	when the computation has completed, blocking if necessary until it is ready. Cancellation is performed by the cancel method.
*	Additional methods are provided to determine if the task completed normally or was cancelled. Once a computation has completed,
*	the computation cannot be cancelled.
*/
type Future interface {
	cancel() bool
	get() Result
	getWithTimeout(timeout time.Duration) Result
	isCancelled() bool
	isDone() bool
}


// Result structure returned by the future task, consists of contents of the result and errors of the computation.
type Result struct{
	value interface{}
	error error
}

/**
*	FutureTask that implements the future interface,
*	@done returns true if the async task execution is complete,
*	@result returns the Result structure,
*	@channel returns a channel to wait for the result.
 */
type FutureTask struct{
	done    bool
	result  Result
	channel chan Result
}

// Returns a reference to future object and results of the computations can be obtained by the functions associated with
// the reference.
func NewFuture(task func() Result) Future{

	commChannel := make(chan Result)

	futureObject := &FutureTask{
		done    : false,
		result  : Result{},
		channel:  commChannel,
	}

	go func(){
		fmt.Println("Inside go routine for new future instance")
		defer close(commChannel)
		resultObject := task()
		commChannel <- resultObject
		fmt.Println("Exit of routine")
	}()
	return futureObject
}


// Waits if necessary for the computation to complete, and then retrieves its result.
func (futureTask *FutureTask) get() Result{
	if futureTask.done {
		return futureTask.result
	}
	ctx := context.Background()
	return futureTask.getWithContext(ctx)
}

// Waits if necessary for at most the given time for the computation to complete, and then retrieves its result, if available.
func (futureTask *FutureTask) getWithTimeout(timeout time.Duration) Result{
	if futureTask.done {
		return futureTask.result
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return futureTask.getWithContext(ctx)
}

// Waits if necessary for the computation to complete, and then retrieves its result.
func (futureTask *FutureTask) getWithContext(ctx context.Context) Result{
	select {
	case <-ctx.Done():
		futureTask.done = true
		futureTask.result = Result{value: nil,error: errors.New("request timed out")}
		return futureTask.result
	case futureTask.result = <-futureTask.channel:
		futureTask.done = true
		return futureTask.result
	}
}

// Returns true if this task completed. Completion may be due to normal termination, an exception, or cancellation
// -- in all of these cases, this method will return true.
func (futureTask *FutureTask) isDone() bool{
	return futureTask.done
}

// Returns true if this task was cancelled before it completed normally.
func (futureTask *FutureTask) isCancelled() bool{
	if futureTask.done {
		if futureTask.result.error != nil && futureTask.result.error.Error() == "manually cancelled" {
			return true
		}
	}
	return false
}

// Attempts to cancel execution of this task. This attempt will fail if the task has already completed, has already been
// cancelled, or could not be cancelled for some other reason. If successful, and this task has not started when cancel
// is called, this task should never run.
func (futureTask *FutureTask) cancel() bool{
	if futureTask.isDone() || futureTask.isCancelled() {
		return false
	}
	futureTask.done = true
	futureTask.result = Result{value: nil,error: errors.New("manually cancelled")}
	futureTask.channel <- futureTask.result
	return true
}