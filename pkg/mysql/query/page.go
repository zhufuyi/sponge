package query

import "strings"

var defaultMaxSize = 1000

// SetMaxSize change the default maximum number of pages per page
// Deprecated: moved to package pkg/gorm/query SetMaxSize
func SetMaxSize(max int) {
	if max < 10 {
		max = 10
	}
	defaultMaxSize = max
}

// Page info
// Deprecated: moved to package pkg/gorm/query Page
type Page struct {
	page int    // page number, starting from page 0
	size int    // number per page
	sort string // sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
}

// Page get page value
// Deprecated: moved to package pkg/gorm/query Page
func (p *Page) Page() int {
	return p.page
}

// Size number per page
// Deprecated: moved to package pkg/gorm/query Size
func (p *Page) Size() int {
	return p.size
}

// Sort get sort field
// Deprecated: moved to package pkg/gorm/query Sort
func (p *Page) Sort() string {
	return p.sort
}

// Offset get offset value
// Deprecated: moved to package pkg/gorm/query Offset
func (p *Page) Offset() int {
	return p.page * p.size
}

// DefaultPage default page, number 20 per page, sorted by id backwards
// Deprecated: moved to package pkg/gorm/query DefaultPage
func DefaultPage(page int) *Page {
	if page < 0 {
		page = 0
	}
	return &Page{
		page: page,
		size: 20,
		sort: "id DESC",
	}
}

// NewPage custom page, starting from page 0.
// the parameter columnNames indicates a sort field, if empty means id descending, if there are multiple column names, separated by a comma,
// a '-' sign in front of each column name indicates descending order, otherwise ascending order.
// Deprecated: moved to package pkg/gorm/query NewPage
func NewPage(page int, size int, columnNames string) *Page {
	if page < 0 {
		page = 0
	}
	if size > defaultMaxSize || size < 1 {
		size = defaultMaxSize
	}

	return &Page{
		page: page,
		size: size,
		sort: getSort(columnNames),
	}
}

// convert to mysql sort, each column name preceded by a '-' sign, indicating descending order, otherwise ascending order, example:
//
//	columnNames="name" means sort by name in ascending order,
//	columnNames="-name" means sort by name descending,
//	columnNames="name,age" means sort by name in ascending order, otherwise sort by age in ascending order,
//	columnNames="-name,-age" means sort by name descending before sorting by age descending.
func getSort(columnNames string) string {
	columnNames = strings.Replace(columnNames, " ", "", -1)
	if columnNames == "" {
		return "id DESC"
	}

	names := strings.Split(columnNames, ",")
	strs := make([]string, 0, len(names))
	for _, name := range names {
		if name[0] == '-' && len(name) > 1 {
			strs = append(strs, name[1:]+" DESC")
		} else {
			strs = append(strs, name+" ASC")
		}
	}

	return strings.Join(strs, ", ")
}
