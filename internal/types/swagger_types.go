package types

// swagger public structures, it is recommended to write comments to all of them,
// if you use yapi, import swagger.json into yapi and fill in the notes automatically to avoid repeating the comments.

// Result output data format
type Result struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // return information description
	Data interface{} `json:"data"` // return data
}

// Params query parameters
type Params struct {
	Page int    `form:"page" binding:"gte=0" json:"page"`      // page number, starting from page 0
	Size int    `form:"size" binding:"gt=0" json:"size"`       // lines per page
	Sort string `form:"sort" binding:"" json:"sort,omitempty"` // sorted fields, multi-column sorting separated by commas

	Columns []Column `json:"columns,omitempty"` // query conditions
}

// Column search information
type Column struct {
	Name  string      `json:"name"`  // column name
	Exp   string      `json:"exp"`   // expressions, which default to = when the value is null, have =, ! =, >, >=, <, <=, like
	Value interface{} `json:"value"` // column value
	Logic string      `json:"logic"` // logical type, defaults to and when value is null, only &(and), ||(or)
}
