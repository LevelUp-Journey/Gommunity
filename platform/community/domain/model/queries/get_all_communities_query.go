package queries

type GetAllCommunitiesQuery struct {
	limit  *int
	offset *int
}

func NewGetAllCommunitiesQuery() GetAllCommunitiesQuery {
	return GetAllCommunitiesQuery{}
}

func (q GetAllCommunitiesQuery) WithPagination(limit, offset int) GetAllCommunitiesQuery {
	q.limit = &limit
	q.offset = &offset
	return q
}

func (q GetAllCommunitiesQuery) Limit() *int {
	return q.limit
}

func (q GetAllCommunitiesQuery) Offset() *int {
	return q.offset
}
