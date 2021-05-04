package errors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/payfazz/go-errors/v2"
)

func TestWrappingNil(t *testing.T) {
	wrappedErr := errors.Wrap(nil)
	if wrappedErr != nil {
		t.Errorf("Wrap(nil) should be nil")
	}
}

func TestIndempotentWrap(t *testing.T) {
	errOri := errors.New("testerr")
	errWrapped1 := errors.Wrap(errOri)
	errWrapped2 := errors.Wrap(errWrapped1)
	errWrapped3 := errors.Wrap(errWrapped2)

	if errWrapped1 != errWrapped2 || errWrapped2 != errWrapped3 {
		t.Errorf("wrapped error must be indempotent")
	}
}

func TestWrapMessage(t *testing.T) {
	errOri := errors.Errorf("testerr")
	errWrapped := errors.Wrap(errOri)
	if errWrapped.Error() != "testerr" {
		t.Errorf("Wrapped error should not change error message")
	}
}

func TestNew(t *testing.T) {
	err0 := fmt.Errorf("err1")
	err1 := errors.Wrap(err0)
	err2 := errors.NewWithCause("err2", err1)
	err3 := errors.NewWithCause("err3", err2)

	if !errors.Is(err3, err2) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err3, err1) {
		t.Errorf("invalid errors.Is")
	}

	if !errors.Is(err3, err0) {
		t.Errorf("invalid errors.Is")
	}
}

type myErr struct{ msg string }

func (e *myErr) Error() string { return e.msg }

// this function name is used for check stack trace
func func201df1b9f41c6cbf81ee12d34e90e26b(ok bool) error {
	err := errors.Wrap(&myErr{msg: "test err"})
	if ok {
		return err
	}
	panic(err)
}

var zero = 0

func func8efcd880187ff6088559a26f82659290(ok bool, f func() error) error {
	errCh := make(chan error)
	errors.Go(
		func(err error) {
			errCh <- err
		},
		func() error {
			if ok {
				return f()
			}

			// this will panic
			return fmt.Errorf("%d", 10/zero)
		},
	)
	return <-errCh
}

func TestFormat(t *testing.T) {
	funcName := "func201df1b9f41c6cbf81ee12d34e90e26b"
	e1 := func8efcd880187ff6088559a26f82659290(true, func() error {
		return func201df1b9f41c6cbf81ee12d34e90e26b(true)
	})

	haveTrace := false
	for _, t := range errors.StackTrace(e1) {
		if strings.Contains(t.Func(), funcName) {
			haveTrace = true
			break
		}
	}

	if !haveTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	e2 := errors.NewWithCause("test cause", e1)

	if !strings.Contains(errors.Format(e2), funcName) {
		t.Errorf("invalid errors.Format")
	}
}

func TestErrorsAs(t *testing.T) {
	var target *myErr

	e := func201df1b9f41c6cbf81ee12d34e90e26b(true)

	if !errors.As(e, &target) {
		t.Errorf("invalid errors.As")
	}

	if e.Error() != target.msg {
		t.Errorf("invalid errors.As")
	}
}

func TestGo(t *testing.T) {
	err := func8efcd880187ff6088559a26f82659290(true, func() error {
		return func201df1b9f41c6cbf81ee12d34e90e26b(true)
	})

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "func201df1b9f41c6cbf81ee12d34e90e26b") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "func8efcd880187ff6088559a26f82659290") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestGoNil(t *testing.T) {
	err := func8efcd880187ff6088559a26f82659290(true, func() error { return nil })
	if err != nil {
		t.Errorf("invalid errors.Go")
	}
}

func TestGoPanic(t *testing.T) {
	err := func8efcd880187ff6088559a26f82659290(false, nil)

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "func8efcd880187ff6088559a26f82659290") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "func8efcd880187ff6088559a26f82659290") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestGoTracedPanic(t *testing.T) {
	err := func8efcd880187ff6088559a26f82659290(true, func() error {
		return func201df1b9f41c6cbf81ee12d34e90e26b(false)
	})

	goodTrace := false
	for _, l := range errors.StackTrace(err) {
		if strings.Contains(l.Func(), "func8efcd880187ff6088559a26f82659290") {
			goodTrace = true
			break
		}
	}

	if !goodTrace {
		t.Errorf("invalid errors.StackTrace")
	}

	goodParent := false
	for _, l := range errors.ParentStackTrace(err) {
		if strings.Contains(l.Func(), "func8efcd880187ff6088559a26f82659290") {
			goodParent = true
			break
		}
	}

	if !goodParent {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}

func TestNilTrace(t *testing.T) {
	err := fmt.Errorf("test err")
	if len(errors.StackTrace(err)) != 0 {
		t.Errorf("invalid errors.StackTrace")
	}
	if len(errors.ParentStackTrace(err)) != 0 {
		t.Errorf("invalid errors.ParentStackTrace")
	}
}
