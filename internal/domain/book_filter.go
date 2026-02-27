package domain

type BookFilter struct {
	Author    string
	Title     string
	Genre     string
	YearFrom  int
	YearTo    int
	PagesFrom int
	PagesTo   int
	Limit     uint64
	Offset    uint64
}
