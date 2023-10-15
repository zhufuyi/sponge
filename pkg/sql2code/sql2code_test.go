package sql2code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sqlData = `
create table user
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment 'username',
    password   char(100)       not null comment 'password',
    email      char(50)        not null comment 'email',
    phone      bigint unsigned not null comment 'phone number',
    age        tinyint         not null comment 'age',
    gender     tinyint         not null comment 'gender, 1:male, 2:female, 3:unknown',
    constraint user_email_uindex
        unique (email)
);
`

func TestGenerateOne(t *testing.T) {
	type args struct {
		args *Args
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sql form param",
			args: args{args: &Args{
				SQL: sqlData,
			}},
			wantErr: false,
		},
		{
			name: "sql from file",
			args: args{args: &Args{
				DDLFile: "test.sql",
			}},
			wantErr: false,
		},
		//{
		//	name: "sql from db",
		//	args: args{args: &Args{
		//		DBDsn:   "root:123456@(192.168.3.37:3306)/test",
		//		DBTable: "user",
		//	}},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOne(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		args *Args
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sql form param",
			args: args{args: &Args{
				SQL: sqlData,
			}},
			wantErr: false,
		},
		//{
		//	name: "sql from db",
		//	args: args{args: &Args{
		//		DBDsn:   "root:123456@(127.0.0.1:3306)/test",
		//		DBTable: "user",
		//	}},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestGenerateError(t *testing.T) {
	a := &Args{}
	_, err := Generate(a)
	assert.Error(t, err)
	_, err = GenerateOne(a)
	assert.Error(t, err)

	a = &Args{DDLFile: "notfound.sql"}
	_, err = Generate(a)
	assert.Error(t, err)

	a = &Args{DBDsn: "root:123456@(127.0.0.1:3306)/test"}
	_, err = Generate(a)
	assert.Error(t, err)

	a = &Args{DBDsn: "root:123456@(127.0.0.1:3306)/test", DBTable: "user"}
	_, err = Generate(a)
	assert.Error(t, err)

	a = &Args{DDLFile: "test.sql", CodeType: "unknown"}
	_, err = GenerateOne(a)
	t.Log(err)
	assert.Error(t, err)
}

func Test_getOptions(t *testing.T) {
	a := &Args{
		Package:        "Package",
		GormType:       true,
		JSONTag:        true,
		ForceTableName: true,
		Charset:        "Charset",
		Collation:      "Collation",
		TablePrefix:    "TablePrefix",
		ColumnPrefix:   "ColumnPrefix",
		NoNullType:     true,
		NullStyle:      "sql",
	}

	o := getOptions(a)
	assert.NotNil(t, o)
	a.NullStyle = "ptr"
	assert.NotNil(t, o)
	a.NullStyle = "default"
	assert.NotNil(t, o)
}
