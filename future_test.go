package GoFuture

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestNewFutureWithGet(t *testing.T) {
	futureInstance1:= NewFuture(func() Result {
		var temp interface{}
		temp = 1
		time.Sleep(2 * time.Second)
		return Result{value: temp}
	})

	actual := futureInstance1.get()
	expected := &Result{error: nil, value: 1}

	if actual.error != expected.error{
		t.Errorf("Error actual = %v, and Expected = %v.", actual.error, expected.error)
	}

	if actual.value != expected.value{
		t.Errorf("Value actual = %v, and Expected = %v.", actual.value, expected.value)
	}

}

func TestNewFutureWithGetTimeout(t *testing.T) {
	futureInstance1:= NewFuture(func() Result {
		var temp interface{}
		temp = 1
		time.Sleep(2 * time.Second)
		return Result{value: temp}
	})

	actual := futureInstance1.getWithTimeout(1 * time.Second)
	expected := Result{value: nil,error: errors.New("request timed out")}

	if actual.error.Error() != expected.error.Error(){
		t.Errorf("Error actual = %v, and Expected = %v.", actual.error, expected.error)
	}

	if actual.value != expected.value{
		t.Errorf("Value actual = %v, and Expected = %v.", actual.value, expected.value)
	}

}

func TestNewFutureWithGetTimeout_2(t *testing.T) {
	futureInstance1:= NewFuture(func() Result {
		var temp interface{}
		temp = 1
		time.Sleep(1 * time.Second)
		return Result{value: temp}
	})

	actual := futureInstance1.getWithTimeout(2 * time.Second)
	expected := Result{value: 1,error: nil}

	if actual.error != expected.error{
		t.Errorf("Error actual = %v, and Expected = %v.", actual.error, expected.error)
	}

	if actual.value != expected.value{
		t.Errorf("Value actual = %v, and Expected = %v.", actual.value, expected.value)
	}

}

func TestNewFutureIsDone(t *testing.T) {
	futureInstance1:= NewFuture(func() Result {
		var temp interface{}
		temp = 1
		time.Sleep(1 * time.Second)
		return Result{value: temp}
	})

	fmt.Println(futureInstance1.get().error)
	actual := futureInstance1.isDone()
	expected := true

	if actual != expected{
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}

}

func TestNewFutureCancel(t *testing.T) {
	futureInstance1:= NewFuture(func() Result {
		var temp interface{}
		temp = 1
		time.Sleep(2 * time.Second)
		return Result{value: temp}
	})

	go futureInstance1.cancel()
	futureInstance1.get()
	actual := futureInstance1.isCancelled()
	expected := true

	if actual != expected{
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}

}