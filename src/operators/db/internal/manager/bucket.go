package manager

import "ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"

type bucket[T any, PT k8s_generic.Resource[T]] struct {
	state map[string]T
}

func newBucket[T any, PT k8s_generic.Resource[T]]() bucket[T, PT] {
	return bucket[T, PT]{
		state: map[string]T{},
	}
}

func (b *bucket[T, PT]) apply(update k8s_generic.Update[T]) {
	for _, toRemove := range update.ToRemove {
		b.remove(toRemove)
	}

	for _, toAdd := range update.ToAdd {
		b.add(toAdd)
	}
}

func (b *bucket[T, PT]) add(obj T) {
	ptr := PT(&obj)
	key := ptr.GetName()

	b.state[key] = obj
}

func (b *bucket[T, PT]) remove(obj T) {
	ptr := PT(&obj)
	key := ptr.GetName()

	delete(b.state, key)
}
