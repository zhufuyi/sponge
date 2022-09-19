package query

import (
	"fmt"
	"strings"
)

const (
	// Eq equal
	Eq = "eq"
	// Neq not equal
	Neq = "neq"
	// Gt greater than
	Gt = "gt"
	// Gte greater than or equal
	Gte = "gte"
	// Lt less than
	Lt = "lt"
	// Lte less than or equal
	Lte = "lte"
	// Like like
	Like = "like"

	// AND logic and
	AND string = "and"
	// OR logic or
	OR string = "or"
)

var expMap = map[string]string{
	Eq:   " = ",
	Neq:  " <> ",
	Gt:   " > ",
	Gte:  " >= ",
	Lt:   " < ",
	Lte:  " <= ",
	Like: " LIKE ",

	"=":  " = ",
	"!=": " <> ",
	">":  " > ",
	">=": " >= ",
	"<":  " < ",
	"<=": " <= ",
}

var logicMap = map[string]string{
	AND: " AND ",
	OR:  " OR ",

	"&":  " AND ",
	"&&": " AND ",
	"|":  " OR ",
	"||": " OR ",
}

// Params 查询原始参数
type Params struct {
	Page int    `form:"page" binding:"gte=0" json:"page"`
	Size int    `form:"size" binding:"gt=0" json:"size"`
	Sort string `form:"sort" binding:"" json:"sort,omitempty"`

	Columns []Column `json:"columns,omitempty"` // 非必须
}

// Column 表的列查询信息
type Column struct {
	Name  string      `json:"name"`  // 列名
	Exp   string      `json:"exp"`   // 表达式，值为空时默认为=，有=、!=、>、>=、<、<=、like七种类型
	Value interface{} `json:"value"` // 列值
	Logic string      `json:"logic"` // 逻辑类型，值为空时默认为and，有&(and)、||(or)两种类型
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

// 把ExpType转换为sql表达式，把LogicType转换为sql使用字符
func (c *Column) convert() error {
	if c.Exp == "" {
		c.Exp = Eq
	}
	if v, ok := expMap[strings.ToLower(c.Exp)]; ok {
		c.Exp = v
		if c.Exp == " LIKE " {
			c.Value = fmt.Sprintf("%%%v%%", c.Value)
		}
	} else {
		return fmt.Errorf("unknown c expression type '%s'", c.Exp)
	}

	if c.Logic == "" {
		c.Logic = AND
	}
	if v, ok := logicMap[strings.ToLower(c.Logic)]; ok {
		c.Logic = v
	} else {
		return fmt.Errorf("unknown logic type '%s'", c.Logic)
	}

	return nil
}

// ConvertToPage 根据参数page size sort转换成符合gorm规则参数
func (p *Params) ConvertToPage() (order string, limit int, offset int) {
	page := NewPage(p.Page, p.Size, p.Sort)
	order = page.sort
	limit = page.size
	offset = page.page * page.size
	return
}

// ConvertToGormConditions 根据参数Columns转换成符合gorm规则参数
// 无论是一列还是多列查询，忽略最后一列的逻辑类型
func (p *Params) ConvertToGormConditions() (string, []interface{}, error) {
	str := ""
	args := []interface{}{}
	l := len(p.Columns)
	if l == 0 {
		return "", nil, nil
	}

	isUseIN := true
	if l == 1 {
		isUseIN = false
	}
	field := p.Columns[0].Name

	for i, column := range p.Columns {
		if err := column.checkValid(); err != nil {
			return "", nil, err
		}

		err := column.convert()
		if err != nil {
			return "", nil, err
		}

		if i == l-1 { // 忽略最后一列的逻辑类型
			str += column.Name + column.Exp + "?"
		} else {
			str += column.Name + column.Exp + "?" + column.Logic
		}
		args = append(args, column.Value)

		if isUseIN {
			if field != column.Name {
				isUseIN = false
				continue
			}
			if column.Exp != expMap[Eq] {
				isUseIN = false
			}
		}
	}

	if isUseIN {
		str = field + " IN (?)"
		args = []interface{}{args}
	}

	return str, args, nil
}

func getExpsAndLogics(keyLen int, paramSrc string) ([]string, []string) { //nolint
	exps, logics := []string{}, []string{}
	param := strings.Replace(paramSrc, " ", "", -1)
	sps := strings.SplitN(param, "?", 2)
	if len(sps) == 2 {
		param = sps[1]
	}

	num := keyLen
	if num == 0 {
		return exps, logics
	}

	fields := []string{}
	kvs := strings.Split(param, "&")
	for _, kv := range kvs {
		if strings.Contains(kv, "page=") || strings.Contains(kv, "size=") || strings.Contains(kv, "sort=") {
			continue
		}
		fields = append(fields, kv)
	}

	// 根据不重复的key分为num组，在每组中判断exp和logic是否存在
	group := map[string]string{}
	for _, field := range fields {
		split := strings.SplitN(field, "=", 2)
		if len(split) != 2 {
			continue
		}

		if _, ok := group[split[0]]; ok {
			// 在一组中，如果exp不存在则填充默认值空，logic不存在则填充充默认值空
			exps = append(exps, group["exp"])
			logics = append(logics, group["logic"])

			group = map[string]string{}
			continue
		} else {
			group[split[0]] = split[1]
		}
	}

	// 处理最后一组
	exps = append(exps, group["exp"])
	logics = append(logics, group["logic"])

	return exps, logics
}
