package pagination

type Pagination struct {
	limit  *uint64 `json:"limit,omitempty"`
	offset *uint64 `json:"offset,omitempty"`
}

func (p *Pagination) Limit() (uint64, bool) {
	if p.limit == nil {
		return 0, false
	}

	return *p.limit, true
}

func (p *Pagination) SetLimit(limit uint64) {
	p.limit = &limit
}

func (p *Pagination) Offset() (uint64, bool) {
	if p.offset == nil {
		return 0, false
	}

	return *p.offset, true
}

func (p *Pagination) SetOffset(offset uint64) {
	p.offset = &offset
}
