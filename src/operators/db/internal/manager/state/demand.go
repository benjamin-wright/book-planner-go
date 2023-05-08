package state

type Demand[T any] struct {
	ToAdd    []T
	ToRemove []T
}

func getOneForOneDemand[T any, U any](state map[string]T, existing map[string]U, transform func(T) U) Demand[U] {
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

	return Demand[U]{
		ToAdd:    toAdd,
		ToRemove: toRemove,
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
) Demand[U] {
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

	return Demand[U]{
		ToAdd:    toAdd,
		ToRemove: toRemove,
	}
}

type readyable[T any] interface {
	*T
	IsReady() bool
}

type targetable[T comparable] interface {
	nameable[T]
	GetTarget() string
}

func getServiceBoundDemand[T comparable, U comparable, V any, W any, PT targetable[T], PU nameable[U], PV readyable[V]](
	state map[string]T,
	existing map[string]U,
	servers map[string]V,
	services map[string]W,
	transform func(T) U,
) Demand[U] {
	d := Demand[U]{
		ToAdd:    []U{},
		ToRemove: []U{},
	}

	seen := map[string]U{}

	for _, client := range state {
		clientPtr := PT(&client)
		target := clientPtr.GetTarget()

		ss, hasSS := servers[target]
		_, hasSvc := services[target]

		ssPtr := PV(&ss)

		if !hasSS || !hasSvc || !ssPtr.IsReady() {
			continue
		}

		desired := transform(client)
		name := PU(&desired).GetName()
		seen[name] = desired

		if _, ok := existing[name]; !ok {
			d.ToAdd = append(d.ToAdd, desired)
		}
	}

	for current, db := range existing {
		if _, ok := seen[current]; !ok {
			d.ToRemove = append(d.ToRemove, db)
		}
	}

	return d
}
