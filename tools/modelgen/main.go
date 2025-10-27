package main

import (
	"os"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/Goryudyuma/techblog-go-mysql-tools/internal/infra/mysql/modelcommon"
)

func main() {
	genConfig := gen.Config{
		ModelPkgPath:  "./internal/infra/mysql/model",
		Mode:          0, // generate mode
		FieldNullable: true,
	}
	genConfig.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		"datetime": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "sql.NullTime"
			}
			return "time.Time"
		},
		"text": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "sql.NullString"
			}
			return "string"
		},
		"varchar": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "sql.NullString"
			}
			return "string"
		},
		"int": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "sql.NullInt64"
			}
			return "int64"
		},
		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "sql.NullInt64"
			}
			return "int64"
		},
		"enum": func(columnType gorm.ColumnType) (dataType string) {
			// プロダクトに合わせて特殊な型に変換する
			if columnType.Name() == "gender" {
				return "publicType.Gender"
			}
			return "string"
		},
	})
	genConfig.WithTableNameStrategy(func(tableName string) (targetTableName string) {
		if tableName == "gorp_migrations" { // sql-migrateの管理テーブルなのでmodelとして生成しない
			return ""
		}
		return tableName
	})
	genConfig.WithImportPkgPath( // 特殊な型を使いたい場合は、ここでimportするパスを書いておく
		"github.com/Goryudyuma/techblog-go-mysql-tools/pkg/publicType",
	)
	g := gen.NewGenerator(genConfig)
	g.WithOpts(
		gen.FieldModify(func(f gen.Field) gen.Field {
			switch f.Type { // WithDataTypeMapで指定しても、ポインタがついてしまうので、ここで調整
			case "*sql.NullTime":
				f.Type = "sql.NullTime"
			case "*sql.NullString":
				f.Type = "sql.NullString"
			case "*sql.NullInt64":
				f.Type = "sql.NullInt64"
			}
			f.Tag.Set("db", f.Name)
			f.GORMTag = map[string][]string{}
			f.Tag.Remove("json")
			return f
		}),
		gen.FieldIgnore("created_at", "updated_at"),
		gen.WithMethod(modelcommon.Common{}),
	)

	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "test_service",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	dsn := config.FormatDSN()
	if modelgenDSN := os.Getenv("MODELGEN_DSN"); modelgenDSN != "" {
		dsn = modelgenDSN
	}
	db, err := gorm.Open(gormmysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	g.UseDB(db) // reuse your gorm db
	g.GenerateAllTable()
	g.Execute()
}
