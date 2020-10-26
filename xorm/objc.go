// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	//"fmt"
	"strings"
	"text/template"

	"xorm.io/xorm/schemas"
)

var (
	ObjcTmpl LangTmpl = LangTmpl{
		template.FuncMap{"Mapper": mapper.Table2Obj,
			"Type":    objcTypeStr,
			"UnTitle": unTitle,
		},
		nil,
		genCPlusImports,
	}
)

func objcTypeStr(col *schemas.Column) string {
	tp := col.SQLType
	name := strings.ToUpper(tp.Name)
	switch name {
	case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.Serial:
		return "int"
	case schemas.BigInt, schemas.BigSerial:
		return "long"
	case schemas.Char, schemas.Varchar, schemas.TinyText, schemas.Text, schemas.MediumText, schemas.LongText:
		return "NSString*"
	case schemas.Date, schemas.DateTime, schemas.Time, schemas.TimeStamp:
		return "NSString*"
	case schemas.Decimal, schemas.Numeric:
		return "NSString*"
	case schemas.Real, schemas.Float:
		return "float"
	case schemas.Double:
		return "double"
	case schemas.TinyBlob, schemas.Blob, schemas.MediumBlob, schemas.LongBlob, schemas.Bytea:
		return "NSString*"
	case schemas.Bool:
		return "BOOL"
	default:
		return "NSString*"
	}
	return ""
}

func genObjcImports(tables []*schemas.Table) map[string]string {
	imports := make(map[string]string)

	for _, table := range tables {
		for _, col := range table.Columns() {
			switch objcTypeStr(col) {
			case "time_t":
				imports[`<time.h>`] = `<time.h>`
			case "tstring":
				imports["<string>"] = "<string>"
				//case "__int64":
				//    imports[""] = ""
			}
		}
	}
	return imports
}
