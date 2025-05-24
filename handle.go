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

func (m *HandleMap[T, PT]) Add(value T) (*Handle, bool) {
	if m.nextUnused != 0 {
		idx := m.nextUnused
		item := PT(m.items[idx])
		m.nextUnused = m.unused[idx]
		m.unused[idx] = 0
		handle := item.Handle()
		if handle == nil {
			handle = &Handle{}
			item.SetHandle(handle)
		}
		gen := item.Handle().gen
		*item = value
		item.SetHandle(handle)
		item.Handle().idx = idx
		item.Handle().gen = gen + 1
		m.totalUnused -= 1
		return item.Handle(), true
	}

	// 0 is a dummy element for 'no item'
	if m.usedItems == 0 {
		m.items[0] = new(T)
		m.usedItems += 1
	}

	if m.usedItems == uint32(len(m.items)) {
		return nil, false
	}

	item := PT(m.items[m.usedItems])
	handle := item.Handle()
	if handle == nil {
		handle = new(Handle)
	}
	*item = value
	item.SetHandle(handle)
	item.Handle().idx = uint32(m.usedItems)
	item.Handle().gen = 1
	m.usedItems += 1
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

func (m *HandleMap[T, PT]) Remove(handle *Handle) {
	if handle.idx <= 0 || handle.idx >= m.usedItems {
		return
	}

	if item := PT(m.items[handle.idx]); item.Handle() == handle {
		m.unused[handle.idx] = m.nextUnused
		m.nextUnused = handle.idx
		m.totalUnused += 1
		handle.idx = 0
	}
}
