package query

import "strings"

var defaultMaxSize = 1000

// SetMaxSize 修改默认每页数量最大值
func SetMaxSize(max int) {
	if max < 10 {
		max = 10
	}
	defaultMaxSize = max
}

// Page 页
type Page struct {
	page int    // 页码，从第0页开始
	size int    // 每一页数量
	sort string // 字段排序
}

// Page 页码
func (p *Page) Page() int {
	return p.page
}

// Size 每一页数量
func (p *Page) Size() int {
	return p.size
}

// Sort 排序
func (p *Page) Sort() string {
	return p.sort
}

// Offset 偏移量
func (p *Page) Offset() int {
	return p.page * p.size
}

// DefaultPage 默认page，每页数量20，按id倒叙排序
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

// NewPage 自定义page，从第0页开始，
// 参数columnNames表示排序字段，如果为空表示id降序，如果有多个列名，用逗号分隔，
// 每一个列名称前面有'-'号，表示降序，否则升序
func NewPage(page int, size int, columnNames string) *Page {
	if page < 0 {
		page = 0
	}
	if size > defaultMaxSize {
		size = defaultMaxSize
	}

	return &Page{
		page: page,
		size: size,
		sort: getSort(columnNames),
	}
}

// 转换为mysql 排序，每一个列名称前面有'-'号，表示降序，否则升序，示例：
// columnNames="name"表示按name升序排序
// columnNames="-name"表示按name降排序
// columnNames="name,age"表示按name升序排序前提下，按age升序排序
// columnNames="-name,-age"表示按name降排序前提下，按age降序排序
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
