package query

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPage(t *testing.T) {
	page := DefaultPage(-1)
	t.Log(page.Page(), page.Size(), page.Sort(), page.Skip())
	page = NewPage(0, 20, "")
	t.Log(page.Page(), page.Size(), page.Sort(), page.Skip())

	SetMaxSize(1)
	page = NewPage(0, 20, "_id")
	t.Log(page.Page(), page.Size(), page.Sort(), page.Skip())
}

func TestParams_ConvertToPage(t *testing.T) {
	p := &Params{
		Page: 0,
		Size: 20,
		Sort: "age,-name",
	}
	order, limit, offset := p.ConvertToPage()
	t.Logf("order=%v, limit=%d, skip=%d", order, limit, offset)
}

func TestParams_ConvertToMongoFilter(t *testing.T) {
	type args struct {
		columns []Column
	}
	tests := []struct {
		name    string
		args    args
		want    bson.M
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
			want:    bson.M{"name": "ZhangSan"},
			wantErr: false,
		},
		{
			name: "1 column neq",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Exp:   "!=",
						Value: "ZhangSan",
					},
				},
			},
			want:    bson.M{"name": bson.M{"$neq": "ZhangSan"}},
			wantErr: false,
		},
		{
			name: "1 column gt",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Exp:   ">",
						Value: 20,
					},
				},
			},
			want:    bson.M{"age": bson.M{"$gt": 20}},
			wantErr: false,
		},
		{
			name: "1 column gte",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Exp:   ">=",
						Value: 20,
					},
				},
			},
			want:    bson.M{"age": bson.M{"$gte": 20}},
			wantErr: false,
		},
		{
			name: "1 column lt",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Exp:   "<",
						Value: 20,
					},
				},
			},
			want:    bson.M{"age": bson.M{"$lt": 20}},
			wantErr: false,
		},
		{
			name: "1 column lte",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Exp:   "<=",
						Value: 20,
					},
				},
			},
			want:    bson.M{"age": bson.M{"$lte": 20}},
			wantErr: false,
		},
		{
			name: "1 column like",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Exp:   Like,
						Value: "Li",
					},
				},
			},
			want:    bson.M{"name": bson.M{"$regex": "Li"}},
			wantErr: false,
		},
		{
			name: "1 column IN",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Exp:   In,
						Value: "ab,cd,ef",
					},
				},
			},
			want:    bson.M{"name": bson.M{"$in": []interface{}{"ab", "cd", "ef"}}},
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
			want:    bson.M{"$and": []bson.M{{"name": "ZhangSan"}, {"gender": "male"}}},
			wantErr: false,
		},
		{
			name: "2 columns neq and",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Exp:   "!=",
						Value: "ZhangSan",
					},
					{
						Name:  "name",
						Exp:   "!=",
						Value: "LiSi",
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"name": bson.M{"$neq": "ZhangSan"}}, {"name": bson.M{"$neq": "LiSi"}}}},
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
						Exp:   ">",
						Value: 20,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "male"}, {"age": bson.M{"$gt": 20}}}},
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
						Exp:   ">=",
						Value: 20,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "male"}, {"age": bson.M{"$gte": 20}}}},
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
						Exp:   "<",
						Value: 20,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "female"}, {"age": bson.M{"$lt": 20}}}},
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
						Exp:   "<=",
						Value: 20,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "female"}, {"age": bson.M{"$lte": 20}}}},
			wantErr: false,
		},
		{
			name: "2 columns range and",
			args: args{
				columns: []Column{
					{
						Name:  "age",
						Exp:   ">=",
						Value: 10,
					},
					{
						Name:  "age",
						Exp:   "<=",
						Value: 20,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"age": bson.M{"$gte": 10}}, {"age": bson.M{"$lte": 20}}}},
			wantErr: false,
		},
		{
			name: "2 columns eq or",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "LiSi",
						Logic: "||",
					},
					{
						Name:  "gender",
						Value: "female",
					},
				},
			},
			want:    bson.M{"$or": []bson.M{{"name": "LiSi"}, {"gender": "female"}}},
			wantErr: false,
		},
		{
			name: "2 columns neq or",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: "LiSi",
						Logic: "||",
					},
					{
						Name:  "gender",
						Exp:   "!=",
						Value: "male",
					},
				},
			},
			want:    bson.M{"$or": []bson.M{{"name": "LiSi"}, {"gender": bson.M{"$neq": "male"}}}},
			wantErr: false,
		},
		{
			name: "2 columns eq and in",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
					},
					{
						Name:  "name",
						Exp:   In,
						Value: "LiSi,ZhangSan,WangWu",
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "male"}, {"name": bson.M{"$in": []interface{}{"LiSi", "ZhangSan", "WangWu"}}}}},
			wantErr: false,
		},

		// --------------------------- query 3 columns  ------------------------------
		{
			name: "3 columns and",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
					},
					{
						Name:  "name",
						Value: "ZhangSan",
					},
					{
						Name:  "age",
						Exp:   "<",
						Value: 12,
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"gender": "male"}, {"name": "ZhangSan"}, {"age": bson.M{"$lt": 12}}}},
			wantErr: false,
		},
		{
			name: "3 columns or",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
						Logic: "||",
					},
					{
						Name:  "name",
						Value: "ZhangSan",
						Logic: "||",
					},
					{
						Name:  "age",
						Exp:   "<",
						Value: 12,
					},
				},
			},
			want:    bson.M{"$or": []bson.M{{"gender": "male"}, {"name": "ZhangSan"}, {"age": bson.M{"$lt": 12}}}},
			wantErr: false,
		},
		{
			name: "3 columns mix (or and)",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
						Logic: "and",
					},
					{
						Name:  "name",
						Value: "ZhangSan",
						Logic: "||",
					},
					{
						Name:  "age",
						Exp:   "<",
						Value: 12,
					},
				},
			},
			want:    bson.M{"$or": []bson.M{{"$and": []bson.M{{"gender": "male"}, {"name": "ZhangSan"}}}, {"age": bson.M{"$lt": 12}}}},
			wantErr: false,
		},

		// --------------------------- query 4 columns  ------------------------------
		{
			name: "4 columns mix (or and)",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Value: "male",
						Logic: "||",
					},
					{
						Name:  "name",
						Value: "ZhangSan",
					},
					{
						Name:  "age",
						Exp:   "<",
						Value: 12,
						Logic: "||",
					},
					{
						Name:  "city",
						Value: "canton",
					},
				},
			},
			want:    bson.M{"$or": []bson.M{{"gender": "male"}, {"$and": []bson.M{{"name": "ZhangSan"}, {"age": bson.M{"$lt": 12}}}}, {"city": "canton"}}},
			wantErr: false,
		},

		{
			name: "convert to object id",
			args: args{
				columns: []Column{
					{
						Name:  "id",
						Value: "65ce48483f11aff697e30d6d",
					},
					{
						Name:  "order_id:oid",
						Value: "65ce48483f11aff697e30d6d",
					},
				},
			},
			want:    bson.M{"$and": []bson.M{{"_id": primitive.ObjectID{0x65, 0xce, 0x48, 0x48, 0x3f, 0x11, 0xaf, 0xf6, 0x97, 0xe3, 0xd, 0x6d}}, {"order_id": primitive.ObjectID{0x65, 0xce, 0x48, 0x48, 0x3f, 0x11, 0xaf, 0xf6, 0x97, 0xe3, 0xd, 0x6d}}}},
			wantErr: false,
		},

		// ---------------------------- error ----------------------------------------------
		{
			name: "exp type err",
			args: args{
				columns: []Column{
					{
						Name:  "gender",
						Exp:   "xxxxxx",
						Value: "male",
					},
				},
			},
			want:    nil,
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "name empty",
			args: args{
				columns: []Column{
					{
						Name:  "",
						Value: "male",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value empty",
			args: args{
				columns: []Column{
					{
						Name:  "name",
						Value: nil,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				columns: nil,
			},
			want:    primitive.M{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{
				Columns: tt.args.columns,
			}
			got, err := params.ConvertToMongoFilter()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToMongoFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToMongoFilter() got = %#v, want = %#v", got, tt.want)
			}
		})
	}
}

func TestConditions_ConvertToMongo(t *testing.T) {
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
	got, err := c.ConvertToMongo()
	if err != nil {
		t.Error(err)
	}
	want := bson.M{"$and": []bson.M{{"name": "ZhangSan"}, {"gender": "male"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ConvertToMongo() got = %+v, want %+v", got, want)
	}
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

func Test_groupingIndex(t *testing.T) {
	type args struct {
		l         int
		orIndexes []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "4 index 1",
			args: args{
				l:         4,
				orIndexes: []int{0, 2},
			},
			want: [][]int{{0}, {1, 2}, {3}},
		},
		{
			name: "4 index 2",
			args: args{
				l:         4,
				orIndexes: []int{1},
			},
			want: [][]int{{0, 1}, {2, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupingIndex(tt.args.l, tt.args.orIndexes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("groupingIndex got = %#v, want = %#v", got, tt.want)
			}
			t.Log(got)
		})
	}
}

func Test_getSort(t *testing.T) {
	names := []string{
		"", "id", "-id", "gender", "gender,id", "-gender,-id",
	}
	for _, name := range names {
		d := getSort(name)
		t.Log(d)
	}
}
