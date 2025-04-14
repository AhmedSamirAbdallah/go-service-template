package database

type SortOption struct {
	Field     string
	Direction int // 1 for asc , -1 for desc
}

type QueryOptions struct {
	Sort       []SortOption
	Projection interface{} // Fields to include or exclude in the result
	Limit      int         // Maximum number of results to return
	Offset     int         // Number of results to skip
}
