package query

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestPage(t *testing.T) {
	page := DefaultPage(-1)
	t.Log(page.Page(), page.Size(), page.Sort(), page.Offset())

	SetMaxSize(1)

	page = NewPage(-1, 100, "id")
	t.Log(page.Page(), page.Size(), page.Sort(), page.Offset())
}

func TestParams_ConvertToPage(t *testing.T) {
	p := &Params{
		Page: 1,
		Size: 50,
		Sort: "age,-name",
	}
	order, limit, offset := p.ConvertToPage()
	t.Logf("order=%s, limit=%d, offset=%d", order, limit, offset)

}

func TestParams_ConvertToGormConditions(t *testing.T) {
	type args struct {
		columns []Column
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		// --------------------------- only 1 column query ------------------------------
		{
			name: "1 column eq",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "ZhangSan",
					},
				},
			},
			want:    "name = ?",
			want1:   []interface{}{"ZhangSan"},
			wantErr: false,
		},
		{
			name: "1 column neq",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "ZhangSan",
						//Exp:   "neq",
						Exp: "!=",
					},
				},
			},
			want:    "name <> ?",
			want1:   []interface{}{"ZhangSan"},
			wantErr: false,
		},
		{
			name: "1 column gt",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Value: 20,
						//Exp:   Gt,
						Exp: ">",
					},
				},
			},
			want:    "age > ?",
			want1:   []interface{}{20},
			wantErr: false,
		},
		{
			name: "1 column gte",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Value: 20,
						//Exp:   Gte,
						Exp: ">=",
					},
				},
			},
			want:    "age >= ?",
			want1:   []interface{}{20},
			wantErr: false,
		},
		{
			name: "1 column lt",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Value: 20,
						//Exp:   Lt,
						Exp: "<",
					},
				},
			},
			want:    "age < ?",
			want1:   []interface{}{20},
			wantErr: false,
		},
		{
			name: "1 column lte",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Value: 20,
						//Exp:   Lte,
						Exp: "<=",
					},
				},
			},
			want:    "age <= ?",
			want1:   []interface{}{20},
			wantErr: false,
		},
		{
			name: "1 column like",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "Li",
						Exp:   Like,
					},
				},
			},
			want:    "name LIKE ?",
			want1:   []interface{}{"%Li%"},
			wantErr: false,
		},

		// --------------------------- query 2 columns  ------------------------------
		{
			name: "2 columns eq and",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "ZhangSan",
					},
					{
						Name:  "gender",
						Value: "male",
					},
				},
			},
			want:    "name = ? AND gender = ?",
			want1:   []interface{}{"ZhangSan", "male"},
			wantErr: false,
		},
		{
			name: "2 columns neq and",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "ZhangSan",
						//Exp:   Neq,
						Exp: "!=",
					},
					{
						Name:  "name",
						Value: "LiSi",
						//Exp:   Neq,
						Exp: "!=",
					},
				},
			},
			want:    "name <> ? AND name <> ?",
			want1:   []interface{}{"ZhangSan", "LiSi"},
			wantErr: false,
		},
		{
			name: "2 columns gt and",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
					},
					{
						Name:  "age",
						Value: 20,
						//Exp:   Gt,
						Exp: ">",
					},
				},
			},
			want:    "gender = ? AND age > ?",
			want1:   []interface{}{"male", 20},
			wantErr: false,
		},
		{
			name: "2 columns gte and",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
					},
					{
						Name:  "age",
						Value: 20,
						//Exp:   Gte,
						Exp: ">=",
					},
				},
			},
			want:    "gender = ? AND age >= ?",
			want1:   []interface{}{"male", 20},
			wantErr: false,
		},
		{
			name: "2 columns lt and",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "female",
					},
					{
						Name:  "age",
						Value: 20,
						//Exp:   Lt,
						Exp: "<",
					},
				},
			},
			want:    "gender = ? AND age < ?",
			want1:   []interface{}{"female", 20},
			wantErr: false,
		},
		{
			name: "2 columns lte and",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "female",
					},
					{
						Name:  "age",
						Value: 20,
						//Exp:   Lte,
						Exp: "<=",
					},
				},
			},
			want:    "gender = ? AND age <= ?",
			want1:   []interface{}{"female", 20},
			wantErr: false,
		},
		{
			name: "2 columns range and",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Value: 10,
						//Exp:   Gte,
						Exp: ">=",
					},
					{
						Name:  "age",
						Value: 20,
						//Exp:   Lte,
						Exp: "<=",
					},
				},
			},
			want:    "age >= ? AND age <= ?",
			want1:   []interface{}{10, 20},
			wantErr: false,
		},
		{
			name: "2 columns eq or",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "LiSi",
						//Logic: OR,
						Logic: "||",
					},
					{
						Name:  "gender",
						Value: "female",
					},
				},
			},
			want:    "name = ? OR gender = ?",
			want1:   []interface{}{"LiSi", "female"},
			wantErr: false,
		},
		{
			name: "2 columns neq or",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "LiSi",
						//Logic: OR,
						Logic: "||",
					},
					{
						Name:  "gender",
						Value: "male",
						//Exp:   Neq,
						Exp: "!=",
					},
				},
			},
			want:    "name = ? OR gender <> ?",
			want1:   []interface{}{"LiSi", "male"},
			wantErr: false,
		},

		// ------------------------------ IN -------------------------------------------------
		{
			name: "3 columns eq and",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "LiSi",
					},
					{
						Name:  "name",
						Value: "ZhangSan",
					},
					{
						Name:  "name",
						Value: "WangWu",
					},
				},
			},
			want:    "name IN (?)",
			want1:   []interface{}{[]interface{}{"LiSi", "ZhangSan", "WangWu"}},
			wantErr: false,
		},

		// ---------------------------- error ----------------------------------------------
		{
			name: "exp type err",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
						Exp:   "xxxxxx",
					},
				},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "logic type err",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
						Logic: "xxxxxx",
					},
				},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				columns: nil,
			},
			want:    "",
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		params := &Params{
			Columns: tt.args.columns,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := params.ConvertToGormConditions()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToGormConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToGormConditions() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ConvertToGormConditions() got1 = %v, want %v", got1, tt.want1)
			}

			got = strings.Replace(got, "?", "%v", -1)
			t.Logf(got, got1...)
		})
	}
}

func Test_getExpsAndLogics(t *testing.T) {
	type args struct {
		keyLen   int
		paramSrc string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []string
	}{
		{
			name: "0 columns",
			args: args{
				keyLen:   0,
				paramSrc: "page=0&size=10",
			},
			want:  []string{},
			want1: []string{},
		},
		{
			name: "1 columns",
			args: args{
				keyLen:   1,
				paramSrc: "k=name&v=LiSi&page=0&size=10",
			},
			want:  []string{""},
			want1: []string{""},
		},
		{
			name: "1 columns gt",
			args: args{
				keyLen:   1,
				paramSrc: "k=age&exp=gt&v=20&page=0&size=10",
			},
			want:  []string{"gt"},
			want1: []string{""},
		},
		{
			name: "1 columns gt and",
			args: args{
				keyLen:   1,
				paramSrc: "k=age&exp=gt&v=20&logic=or&page=0&size=10",
			},
			want:  []string{"gt"},
			want1: []string{"or"},
		},
		{
			name: "2 columns",
			args: args{
				keyLen:   2,
				paramSrc: "k=name&v=LiSi&k=gender&v=male&page=0&size=10",
			},
			want:  []string{"", ""},
			want1: []string{"", ""},
		},
		{
			name: "2 columns gt",
			args: args{
				keyLen:   2,
				paramSrc: "k=name&v=LiSi&k=age&v=20&exp=gt&page=0&size=10",
			},
			want:  []string{"", "gt"},
			want1: []string{"", ""},
		},
		{
			name: "3 columns gt  or",
			args: args{
				keyLen:   3,
				paramSrc: "k=name&v=LiSi&exp=neq&k=age&v=20&exp=gt&k=gender&v=male&logic=or&page=0&size=10",
			},
			want:  []string{"neq", "gt", ""},
			want1: []string{"", "", "or"},
		},
		{
			name: "error",
			args: args{
				keyLen:   1,
				paramSrc: "k=name&exp=gt&v=LiSi&page=0&size=10",
			},
			want:  []string{"gt"},
			want1: []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getExpsAndLogics(tt.args.keyLen, tt.args.paramSrc)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getExpsAndLogics() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getExpsAndLogics() got1 = %v, want %v", got1, tt.want1)
			}
			t.Logf("%+v  %+v", got, got1)
		})
	}
}

func TestConditions_ConvertToGorm(t *testing.T) {
	c := Conditions{
		Columns: []Column{
			{
				Name:  "name",
				Value: "ZhangSan",
			},
			{
				Name:  "gender",
				Value: "male",
			},
		}}
	str, values, err := c.ConvertToGorm()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "name = ? AND gender = ?", str)
	assert.Equal(t, len(values), 2)
}

func TestConditions_checkValid(t *testing.T) {
	// empty error
	c := Conditions{}
	err := c.CheckValid()
	assert.Error(t, err)

	// value is empty error
	c = Conditions{
		Columns: []Column{
			{
				Name:  "foo",
				Value: nil,
			},
		},
	}
	err = c.CheckValid()
	assert.Error(t, err)

	// exp error
	c = Conditions{
		Columns: []Column{
			{
				Name:  "foo",
				Value: "bar",
				Exp:   "unknown-exp",
			},
		},
	}
	err = c.CheckValid()
	assert.Error(t, err)

	// logic error
	c = Conditions{
		Columns: []Column{
			{
				Name:  "foo",
				Value: "bar",
				Logic: "unknown-logic",
			},
		},
	}
	err = c.CheckValid()
	assert.Error(t, err)

	// success
	c = Conditions{
		Columns: []Column{
			{
				Name:  "name",
				Value: "ZhangSan",
			},
			{
				Name:  "gender",
				Value: "male",
			},
		}}
	err = c.CheckValid()
	assert.NoError(t, err)
}
