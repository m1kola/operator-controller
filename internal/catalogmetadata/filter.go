package catalogmetadata

// Predicate returns true if the object should be kept when filtering
type Predicate[T Schemas] func(entity *T) bool

// Filter filters a slice accordingly to
func Filter[T Schemas](in []*T, test Predicate[T]) []*T {
	out := []*T{}
	for i := range in {
		if test(in[i]) {
			out = append(out, in[i])
		}
	}
	return out
}

func And[T Schemas](predicates ...Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		eval := true
		for _, predicate := range predicates {
			eval = eval && predicate(obj)
			if !eval {
				return false
			}
		}
		return eval
	}
}

func Or[T Schemas](predicates ...Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		eval := false
		for _, predicate := range predicates {
			eval = eval || predicate(obj)
			if eval {
				return true
			}
		}
		return eval
	}
}

func Not[T Schemas](predicate Predicate[T]) Predicate[T] {
	return func(obj *T) bool {
		return !predicate(obj)
	}
}
