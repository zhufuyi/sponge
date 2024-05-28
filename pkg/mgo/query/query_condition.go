// Package query is a library of custom condition queries, support for complex conditional paging queries.
package query

import (
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Eq equal
	Eq       = "eq"
	eqSymbol = "="
	// Neq not equal
	Neq       = "neq"
	neqSymbol = "!="
	// Gt greater than
	Gt       = "gt"
	gtSymbol = ">"
	// Gte greater than or equal
	Gte       = "gte"
	gteSymbol = ">="
	// Lt less than
	Lt       = "lt"
	ltSymbol = "<"
	// Lte less than or equal
	Lte       = "lte"
	lteSymbol = "<="
	// Like fuzzy lookup
	Like = "like"
	// In include
	In = "in"

	// AND logic and
	AND        string = "and" //nolint
	andSymbol1        = "&"
	andSymbol2        = "&&"
	// OR logic or
	OR        string = "or" //nolint
	orSymbol1        = "|"
	orSymbol2        = "||"

	allLogicAnd = 1
	allLogicOr  = 2
)

var expMap = map[string]string{
	Eq:        eqSymbol,
	eqSymbol:  eqSymbol,
	Neq:       neqSymbol,
	neqSymbol: neqSymbol,
	Gt:        gtSymbol,
	gtSymbol:  gtSymbol,
	Gte:       gteSymbol,
	gteSymbol: gteSymbol,
	Lt:        ltSymbol,
	ltSymbol:  ltSymbol,
	Lte:       lteSymbol,
	lteSymbol: lteSymbol,
	Like:      Like,
	In:        In,
}

var logicMap = map[string]string{
	AND:        andSymbol1,
	andSymbol1: andSymbol1,
	andSymbol2: andSymbol1,
	OR:         orSymbol1,
	orSymbol1:  orSymbol1,
	orSymbol2:  orSymbol1,
}

// Params query parameters
type Params struct {
	Page int    `json:"page" form:"page" binding:"gte=0"`
	Size int    `json:"size" form:"size" binding:"gt=0"`
	Sort string `json:"sort,omitempty" form:"sort" binding:""`

	Columns []Column `json:"columns,omitempty" form:"columns"` // not required
}

// Column query info
type Column struct {
	Name  string      `json:"name" form:"name"`   // column name
	Exp   string      `json:"exp" form:"exp"`     // expressions, which default to = when the value is null, have =, !=, >, >=, <, <=, like, in
	Value interface{} `json:"value" form:"value"` // column value
	Logic string      `json:"logic" form:"logic"` // logical type, defaults to and when the value is null, with &(and), ||(or)
}

func (c *Column) checkValid() error {
	if c.Name == "" {
		return fmt.Errorf("field 'name' cannot be empty")
	}
	if c.Value == nil {
		return fmt.Errorf("field 'value' cannot be nil")
	}
	return nil
}

func (c *Column) convertLogic() error {
	if c.Logic == "" {
		c.Logic = AND
	}
	if v, ok := logicMap[strings.ToLower(c.Logic)]; ok { //nolint
		c.Logic = v
		return nil
	}
	return fmt.Errorf("unknown logic type '%s'", c.Logic)
}

// converting ExpType to sql expressions and LogicType to sql using characters
func (c *Column) convert() error {
	if err := c.checkValid(); err != nil {
		return err
	}

	if c.Name == "id" || c.Name == "_id" {
		if str, ok := c.Value.(string); ok {
			c.Name = "_id"
			c.Value, _ = primitive.ObjectIDFromHex(str)
		}
	} else if strings.Contains(c.Name, ":oid") {
		if str, ok := c.Value.(string); ok {
			c.Name = strings.Replace(c.Name, ":oid", "", 1)
			c.Value, _ = primitive.ObjectIDFromHex(str)
		}
	}

	if c.Exp == "" {
		c.Exp = Eq
	}
	if v, ok := expMap[strings.ToLower(c.Exp)]; ok { //nolint
		c.Exp = v
		switch c.Exp {
		//case eqSymbol:
		case neqSymbol:
			c.Value = bson.M{"$neq": c.Value}
		case gtSymbol:
			c.Value = bson.M{"$gt": c.Value}
		case gteSymbol:
			c.Value = bson.M{"$gte": c.Value}
		case ltSymbol:
			c.Value = bson.M{"$lt": c.Value}
		case lteSymbol:
			c.Value = bson.M{"$lte": c.Value}
		case Like:
			escapedValue := regexp.QuoteMeta(fmt.Sprintf("%v", c.Value))
			c.Value = bson.M{"$regex": escapedValue, "$options": "i"}
		case In:
			val, ok := c.Value.(string)
			if !ok {
				return fmt.Errorf("invalid value type '%s'", c.Value)
			}
			values := []interface{}{}
			ss := strings.Split(val, ",")
			for _, s := range ss {
				values = append(values, s)
			}
			c.Value = bson.M{"$in": values}
		}
	} else {
		return fmt.Errorf("unknown exp type '%s'", c.Exp)
	}

	return c.convertLogic()
}

// ConvertToPage converted to conform to mongo rules based on the page size sort parameter
func (p *Params) ConvertToPage() (sort bson.D, limit int, skip int) { //nolint
	page := NewPage(p.Page, p.Size, p.Sort)
	sort = page.sort
	limit = page.size
	skip = page.page * page.size
	return //nolint
}

// ConvertToMongoFilter conversion to mongo-compliant parameters based on the Columns parameter
// ignore the logical type of the last column, whether it is a one-column or multi-column query
func (p *Params) ConvertToMongoFilter() (bson.M, error) {
	filter := bson.M{}
	l := len(p.Columns)
	switch l {
	case 0:
		return bson.M{}, nil

	case 1: // l == 1
		err := p.Columns[0].convert()
		if err != nil {
			return nil, err
		}
		filter[p.Columns[0].Name] = p.Columns[0].Value
		return filter, nil

	case 2: // l == 2
		err := p.Columns[0].convert()
		if err != nil {
			return nil, err
		}
		err = p.Columns[1].convert()
		if err != nil {
			return nil, err
		}
		if p.Columns[0].Logic == andSymbol1 {
			filter = bson.M{"$and": []bson.M{
				{p.Columns[0].Name: p.Columns[0].Value},
				{p.Columns[1].Name: p.Columns[1].Value}}}
		} else {
			filter = bson.M{"$or": []bson.M{
				{p.Columns[0].Name: p.Columns[0].Value},
				{p.Columns[1].Name: p.Columns[1].Value}}}
		}
		return filter, nil

	default: // l >=3
		return p.convertMultiColumns()
	}
}

func (p *Params) convertMultiColumns() (bson.M, error) {
	filter := bson.M{}
	logicType, groupIndexes, err := checkSameLogic(p.Columns)
	if err != nil {
		return nil, err
	}
	if logicType == allLogicAnd {
		for _, column := range p.Columns {
			err := column.convert()
			if err != nil {
				return nil, err
			}
			if v, ok := filter["$and"]; !ok {
				filter["$and"] = []bson.M{{column.Name: column.Value}}
			} else {
				if cols, ok1 := v.([]bson.M); ok1 {
					cols = append(cols, bson.M{column.Name: column.Value})
					filter["$and"] = cols
				}
			}
		}
		return filter, nil
	} else if logicType == allLogicOr {
		for _, column := range p.Columns {
			err := column.convert()
			if err != nil {
				return nil, err
			}
			if v, ok := filter["$or"]; !ok {
				filter["$or"] = []bson.M{{column.Name: column.Value}}
			} else {
				if cols, ok1 := v.([]bson.M); ok1 {
					cols = append(cols, bson.M{column.Name: column.Value})
					filter["$or"] = cols
				}
			}
		}
		return filter, nil
	}
	orConditions := []bson.M{}
	for _, indexes := range groupIndexes {
		if len(indexes) == 1 {
			column := p.Columns[indexes[0]]
			err := column.convert()
			if err != nil {
				return nil, err
			}
			orConditions = append(orConditions, bson.M{column.Name: column.Value})
		} else {
			andConditions := []bson.M{}
			for _, index := range indexes {
				column := p.Columns[index]
				err := column.convert()
				if err != nil {
					return nil, err
				}
				andConditions = append(andConditions, bson.M{column.Name: column.Value})
			}
			orConditions = append(orConditions, bson.M{"$and": andConditions})
		}
	}
	filter["$or"] = orConditions

	return filter, nil
}

func checkSameLogic(columns []Column) (int, [][]int, error) {
	orIndexes := []int{}
	l := len(columns)
	for i, column := range columns {
		if i == l-1 { // ignore the logical type of the last column
			break
		}
		err := column.convertLogic()
		if err != nil {
			return 0, nil, err
		}
		if column.Logic == orSymbol1 {
			orIndexes = append(orIndexes, i)
		}
	}

	if len(orIndexes) == 0 {
		return allLogicAnd, nil, nil
	} else if len(orIndexes) == l-1 {
		return allLogicOr, nil, nil
	}
	// mix and or
	groupIndexes := groupingIndex(l, orIndexes)

	return 0, groupIndexes, nil
}

func groupingIndex(l int, orIndexes []int) [][]int {
	groupIndexes := [][]int{}
	lastIndex := 0
	for _, index := range orIndexes {
		group := []int{}
		for i := lastIndex; i <= index; i++ {
			group = append(group, i)
		}
		groupIndexes = append(groupIndexes, group)
		if lastIndex == index {
			lastIndex++
		} else {
			lastIndex = index
		}
	}
	group := []int{}
	for i := lastIndex + 1; i < l; i++ {
		group = append(group, i)
	}
	groupIndexes = append(groupIndexes, group)
	return groupIndexes
}

// Conditions query conditions
type Conditions struct {
	Columns []Column `json:"columns" form:"columns" binding:"min=1"` // columns info
}

// CheckValid check valid
func (c *Conditions) CheckValid() error {
	if len(c.Columns) == 0 {
		return fmt.Errorf("field 'columns' cannot be empty")
	}

	for _, column := range c.Columns {
		err := column.checkValid()
		if err != nil {
			return err
		}
		if column.Exp != "" {
			if _, ok := expMap[column.Exp]; !ok {
				return fmt.Errorf("unknown exp type '%s'", column.Exp)
			}
		}
		if column.Logic != "" {
			if _, ok := logicMap[column.Logic]; !ok {
				return fmt.Errorf("unknown logic type '%s'", column.Logic)
			}
		}
	}

	return nil
}

// ConvertToMongo conversion to mongo-compliant parameters based on the Columns parameter
// ignore the logical type of the last column, whether it is a one-column or multi-column query
func (c *Conditions) ConvertToMongo() (bson.M, error) {
	p := &Params{Columns: c.Columns}
	return p.ConvertToMongoFilter()
}
