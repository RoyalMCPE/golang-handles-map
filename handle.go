package handles

func NewMap[T any](len int) HandleMap[T] {
	items := make([]wrapper[T], 0, len)
	for i := 0; i < len; i++ {
		items = append(items, *new(wrapper[T]))
	}

	return HandleMap[T]{
		items:  items,
		unused: make([]uint32, len),
	}
}

type wrapper[T any] struct {
	data   T
	handle Handle
}

type HandleMap[T any] struct {
	items       []wrapper[T]
	usedItems   uint32
	nextUnused  uint32
	unused      []uint32
	totalUnused uint32
}

type Handle struct {
	idx uint32
	gen uint32
}

func (handles *HandleMap[T]) Add(value T) (Handle, bool) {
	if handles.nextUnused != 0 {
		idx := handles.nextUnused
		inst := &handles.items[idx]
		item := &inst.data
		handles.nextUnused = handles.unused[idx]
		handles.unused[idx] = 0
		gen := inst.handle.gen
		*item = value
		inst.handle.idx = idx
		inst.handle.gen = gen + 1
		handles.totalUnused -= 1
		return inst.handle, true
	}

	// 0 is a dummy element for 'no item'
	if handles.usedItems == 0 {
		handles.items[0] = *new(wrapper[T])
		handles.usedItems += 1
	}

	if handles.usedItems == uint32(len(handles.items)) {
		return Handle{}, false
	}

	item := &handles.items[handles.usedItems].data
	inst := &handles.items[handles.usedItems]
	*item = value
	inst.handle.idx = handles.usedItems
	inst.handle.gen = 1
	handles.usedItems += 1
	return inst.handle, true
}

func (handles *HandleMap[T]) Get(handle Handle) *T {
	if handle.idx <= 0 || handle.idx >= handles.usedItems {
		return nil
	}

	inst := &handles.items[handle.idx]
	if inst.handle.idx == handle.idx && inst.handle.gen == handle.gen {
		return &inst.data
	}

	return nil
}
