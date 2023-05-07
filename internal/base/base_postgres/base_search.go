package base_postgres

type Searcher interface {
	getQuery() string
	getIsJoin() bool
	getJoinModels() string
	getQueryJoin() string
}

type Search struct {
	query      string
	isJoin     bool
	JoinModels string
	queryJoin  string
}

func NewSearcher(query string, isJoin bool, JoinModels string, queryJoin string) Searcher {
	return &Search{
		query:      query,
		isJoin:     isJoin,
		JoinModels: JoinModels,
		queryJoin:  queryJoin,
	}
}

func (p *Search) getQuery() string {
	return p.query
}

func (p *Search) getIsJoin() bool {
	return p.isJoin
}

func (p *Search) getJoinModels() string {
	return p.JoinModels
}

func (p *Search) getQueryJoin() string {
	return p.queryJoin
}
