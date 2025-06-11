package options

type SearchPagination struct {
	Limit  int
	Offset int
}

func NewSearchPagination() SearchPagination {
	return SearchPagination{
		Limit:  0,
		Offset: 0,
	}
}

func (sp SearchPagination) WithLimit(limit int) SearchPagination {
	sp.Limit = limit
	return sp
}

func (sp SearchPagination) WithOffset(offset int) SearchPagination {
	sp.Offset = offset
	return sp
}
