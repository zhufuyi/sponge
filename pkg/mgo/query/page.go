package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var defaultMaxSize = 1000

const oidName = "_id"

// SetMaxSize change the default maximum number of pages per page
func SetMaxSize(max int) {
	if max < 10 {
		max = 10
	}
	defaultMaxSize = max
}

// Page info
type Page struct {
	page  int // page number, starting from page 0
	limit int // number per page

	// sort fields, default is id backwards, you can add - sign before the field to indicate
	// reverse order, no - sign to indicate ascending order, multiple fields separated by comma
	sort bson.D
}

// Page get page value
func (p *Page) Page() int {
	return p.page
}

// Limit number per page
func (p *Page) Limit() int {
	return p.limit
}

// Size number per page
// Deprecated: use Limit instead, will delete it in the future
func (p *Page) Size() int {
	return p.limit
}

// Sort get sort field
func (p *Page) Sort() bson.D {
	return p.sort
}

// Skip get offset value
func (p *Page) Skip() int {
	return p.page * p.limit
}

// DefaultPage default page, number 20 per page, sorted by id backwards
func DefaultPage(page int) *Page {
	if page < 0 {
		page = 0
	}
	return &Page{
		page:  page,
		limit: 10,
		sort:  bson.D{{oidName, -1}}, //nolint
	}
}

// NewPage custom page, starting from page 0.
// the parameter columnNames indicates a sort field, if empty means id descending, if there are multiple column names, separated by a comma,
// a '-' sign in front of each column name indicates descending order, otherwise ascending order.
func NewPage(page int, limit int, columnNames string) *Page {
	if page < 0 {
		page = 0
	}
	if limit > defaultMaxSize {
		limit = defaultMaxSize
	} else if limit < 1 {
		limit = 10
	}

	return &Page{
		page:  page,
		limit: limit,
		sort:  getSort(columnNames),
	}
}

// convert to mysql sort, each column name preceded by a '-' sign, indicating descending order, otherwise ascending order, example:
//
//	columnNames="name" means sort by name in ascending order,
//	columnNames="-name" means sort by name descending,
//	columnNames="name,age" means sort by name in ascending order, otherwise sort by age in ascending order,
//	columnNames="-name,-age" means sort by name descending before sorting by age descending.
func getSort(columnNames string) bson.D {
	d := bson.D{}
	columnNames = strings.Replace(columnNames, " ", "", -1)
	if columnNames == "" {
		d = bson.D{{oidName, -1}} //nolint
		return d
	}

	names := strings.Split(columnNames, ",")
	for _, name := range names {
		if name[0] == '-' && len(name) > 1 {
			col := name[1:]
			if col == "id" {
				col = oidName
			}
			d = append(d, bson.E{col, -1}) //nolint
		} else {
			if name == "id" {
				name = oidName
			}
			d = append(d, bson.E{name, 1}) //nolint
		}
	}

	return d
}
