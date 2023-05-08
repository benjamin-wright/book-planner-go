package manager

type demand[T any] struct {
	toAdd    []T
	toRemove []T
}

func getOneForOneDemand[T any, U any](state map[string]T, existing map[string]U, transform func(T) U) demand[U] {
	toAdd := []U{}
	toRemove := []U{}

	for name, obj := range state {
		if _, ok := existing[name]; !ok {
			toAdd = append(toAdd, transform(obj))
		}
	}

	for name, obj := range existing {
		if _, ok := state[name]; !ok {
			toRemove = append(toRemove, obj)
		}
	}

	return demand[U]{
		toAdd:    toAdd,
		toRemove: toRemove,
	}
}

func getOrphanedDemand[T any, U any](state map[string]T, existing map[string]U, equals func(T, U) bool) []U {
	toRemove := []U{}

	for _, obj := range existing {
		missing := true

		for _, ref := range state {
			if equals(ref, obj) {
				missing = false
				break
			}
		}

		if missing {
			toRemove = append(toRemove, obj)
		}
	}

	return toRemove
}

type hasStorage[T any] interface {
	*T
	GetStorage() string
}

func getStorageBoundDemand[
	T any,
	U any,
	TP hasStorage[T],
	UP hasStorage[U],
](
	state map[string]T,
	existing map[string]U,
	transform func(T) U,
) demand[U] {
	toAdd := []U{}
	toRemove := []U{}

	for name, db := range state {
		if ss, ok := existing[name]; !ok {
			toAdd = append(toAdd, transform(db))
		} else {
			dbPtr := TP(&db)
			ssPtr := UP(&ss)

			if dbPtr.GetStorage() != ssPtr.GetStorage() {
				toRemove = append(toRemove, transform(db))
				toAdd = append(toAdd, transform(db))
			}
		}
	}

	for name, db := range existing {
		if _, ok := state[name]; !ok {
			toRemove = append(toRemove, db)
		}
	}

	return demand[U]{
		toAdd:    toAdd,
		toRemove: toRemove,
	}
}

type readyable[T any] interface {
	*T
	IsReady() bool
}

func getServiceBoundDemand[T comparable, U any, V any, W any, PT Nameable[T], PV readyable[V]](
	state map[string]T,
	existing map[string]U,
	servers map[string]V,
	services map[string]W,
	transform func(T) U,
) demand[U] {
	d := demand[U]{
		toAdd:    []U{},
		toRemove: []U{},
	}

	seen := map[string]U{}

	for _, client := range state {
		clientPtr := PT(&client)
		name := clientPtr.GetName()

		ss, hasSS := servers[name]
		_, hasSvc := services[name]

		ssPtr := PV(&ss)

		if !hasSS || !hasSvc || !ssPtr.IsReady() {
			continue
		}

		desired := transform(client)
		seen[name] = desired

		if _, ok := existing[name]; !ok {
			d.toAdd = append(d.toAdd, desired)
		}
	}

	for current, db := range existing {
		if _, ok := seen[current]; !ok {
			d.toRemove = append(d.toRemove, db)
		}
	}

	return d
}
