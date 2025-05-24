package handles

func NewMap[T any, PT interface {
	Handle() *Handle
	SetHandle(*Handle)
	*T
}](len int) HandleMap[T, PT] {
	items := make([]*T, 0, len)
	for i := 0; i < len; i++ {
		items = append(items, new(T))
	}

	return HandleMap[T, PT]{
		items:  items,
		unused: make([]uint32, len),
	}
}

type HandleMap[T any, PT interface {
	Handle() *Handle
	SetHandle(*Handle)
	*T
}] struct {
	items       []*T
	usedItems   uint32
	nextUnused  uint32
	unused      []uint32
	totalUnused uint32
}

type Handle struct {
	idx uint32
	gen uint32
}

func (handles *HandleMap[T, PT]) Add(value T) (*Handle, bool) {
	if handles.nextUnused != 0 {
		idx := handles.nextUnused
		item := PT(handles.items[idx])
		handles.nextUnused = handles.unused[idx]
		handles.unused[idx] = 0
		gen := item.Handle().gen
		*item = value
		item.Handle().idx = idx
		item.Handle().gen = gen + 1
		handles.totalUnused -= 1
		return item.Handle(), true
	}

	// 0 is a dummy element for 'no item'
	if handles.usedItems == 0 {
		handles.items[0] = new(T)
		handles.usedItems += 1
	}

	if handles.usedItems == uint32(len(handles.items)) {
		return nil, false
	}

	item := PT(handles.items[handles.usedItems])
	handle := item.Handle()
	if handle == nil {
		handle = new(Handle)
	}
	*item = value
	item.SetHandle(handle)
	item.Handle().idx = uint32(handles.usedItems)
	item.Handle().gen = 1
	handles.usedItems += 1
	return item.Handle(), true
}

func (handles *HandleMap[T, PT]) Get(handle *Handle) *T {
	if handle.idx <= 0 || handle.idx >= handles.usedItems {
		return nil
	}

	if item := PT(handles.items[handle.idx]); item.Handle() == handle {
		return item
	}

	return nil
}
