package handler

// swagger公共结构体，每个字段建议都写注释，生成的swagger.json也带有注释，
// 把swagger.json导入yapi后自动填写备注，避免重复填写

// Result 输出数据格式
type Result struct {
	Code int         `json:"code"` // 返回码
	Msg  string      `json:"msg"`  // 返回信息说明
	Data interface{} `json:"data"` // 返回数据
}

// Params 查询原始参数
type Params struct {
	Page int    `form:"page" binding:"gte=0" json:"page"`      // 页码
	Size int    `form:"size" binding:"gt=0" json:"size"`       // 每页行数
	Sort string `form:"sort" binding:"" json:"sort,omitempty"` // 排序字段，默认值为-id，字段前面有-号表示倒序，否则升序，多个字段用逗号分隔

	Columns []Column `json:"columns,omitempty"` // 列查询条件
}

// Column 表的列查询信息
type Column struct {
	Name  string      `json:"name"`  // 列名
	Exp   string      `json:"exp"`   // 表达式，值为空时默认为=，有=、!=、>、>=、<、<=、like七种类型
	Value interface{} `json:"value"` // 列值
	Logic string      `json:"logic"` // 逻辑类型，值为空时默认为and，有&(and)、||(or)两种类型
}
