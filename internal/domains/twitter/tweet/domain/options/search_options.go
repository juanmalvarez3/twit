package options

type SearchOptions struct {
	Filters    SearchFilters
	Pagination SearchPagination
}

func NewSearchOptions() SearchOptions {
	return SearchOptions{
		Filters:    NewSearchFilters(),
		Pagination: NewSearchPagination(),
	}
}

func (so SearchOptions) WithFilters(filters SearchFilters) SearchOptions {
	so.Filters = filters
	return so
}

func (so SearchOptions) WithPagination(pagination SearchPagination) SearchOptions {
	so.Pagination = pagination
	return so
}
