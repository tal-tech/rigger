// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"strings"
	"text/template"

	"xorm.io/xorm/schemas"
)

var (
	CPlusTmpl LangTmpl = LangTmpl{
		template.FuncMap{"Mapper": mapper.Table2Obj,
			"Type":    cPlusTypeStr,
			"UnTitle": unTitle,
		},
		nil,
		genCPlusImports,
	}
)

func cPlusTypeStr(col *schemas.Column) string {
	tp := col.SQLType
	name := strings.ToUpper(tp.Name)
	switch name {
	case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.Serial:
		return "int"
	case schemas.BigInt, schemas.BigSerial:
		return "__int64"
	case schemas.Char, schemas.Varchar, schemas.TinyText, schemas.Text, schemas.MediumText, schemas.LongText:
		return "tstring"
	case schemas.Date, schemas.DateTime, schemas.Time, schemas.TimeStamp:
		return "time_t"
	case schemas.Decimal, schemas.Numeric:
		return "tstring"
	case schemas.Real, schemas.Float:
		return "float"
	case schemas.Double:
		return "double"
	case schemas.TinyBlob, schemas.Blob, schemas.MediumBlob, schemas.LongBlob, schemas.Bytea:
		return "tstring"
	case schemas.Bool:
		return "bool"
	default:
		return "tstring"
	}
	return ""
}

func genCPlusImports(tables []*schemas.Table) map[string]string {
	imports := make(map[string]string)

	for _, table := range tables {
		for _, col := range table.Columns() {
			switch cPlusTypeStr(col) {
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
