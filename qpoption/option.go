package qpoption

// shamelessly yoinked from <https://christine.website/blog/gonads-2022-04-24>

type Option[T any] struct{ v *T }

func (o *Option[T]) Set(val T) Option[T] { o.v = &val; return *o }
func (o Option[T]) IsSome() bool         { return o.v != nil }
func (o Option[T]) IsNone() bool         { return !o.IsSome() }

func (o Option[T]) OrElse(other T) T {
	if o.IsSome() {
		return *o.v
	}
	return other
}

func (o Option[T]) Yank() T {
	if o.IsNone() {
		panic("yank on empty option")
	}
	return *o.v
}

func Value[T any](v T) (t Option[T]) { return t.Set(v) }
