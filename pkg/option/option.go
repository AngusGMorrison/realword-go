package option

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Option represents an optional type.
type Option[T any] struct {
	some  bool
	value T
}

func (o Option[T]) GoString() string {
	return fmt.Sprintf("Option[%[1]T]{some: %[2]t, value: %#[1]v}", o.value, o.some)
}

func (o Option[T]) String() string {
	if o.some {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

func (o *Option[T]) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(bytes, &o.value); err != nil {
		return err // nolint:wrapcheck
	}

	o.some = true

	return nil
}

// Some returns an Option[T] populated with the given type.
func Some[T any](value T) Option[T] {
	return Option[T]{
		some:  true,
		value: value,
	}
}

// None returns an empty Option[T] whose zero-value is considered valid by
// [Option.NonZero] and [NonZero].
func None[S any]() Option[S] {
	return Option[S]{}
}

// Some returns true if the Option is populated with a non-zero value.
func (o Option[T]) Some() bool {
	return o.some
}

// ErrEmptyOption is returned by [Option.Value] when attempting to retrieve a
// value from an empty Option.
var ErrEmptyOption = errors.New("expected Option value was empty")

// Value returns the value of the [Option], or [ErrEmptyOption] error if the
// Option is empty.
func (o Option[T]) Value() (T, error) {
	if !o.Some() {
		return *new(T), ErrEmptyOption
	}
	return o.value, nil
}

// ValueOrZero returns the value of the [Option], which is the zero-value of T
// if the Option is None.
func (o Option[T]) ValueOrZero() T {
	return o.value
}

// Conversion is a function that converts a T into a U.
type Conversion[T any, U any] func(T) (U, error)

// Map transforms an Option[T] into an Option[U] by applying the given conversion to T.
// None[T] is returned as None[U].
//
// # Errors
//   - Any error returned by the conversion.
func Map[T any, U any](opt Option[T], convert Conversion[T, U]) (Option[U], error) {
	if !opt.Some() {
		return None[U](), nil
	}

	u, err := convert(opt.ValueOrZero())
	if err != nil {
		return None[U](), err
	}

	return Some(u), nil
}
