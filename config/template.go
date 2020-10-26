package config

type XormTpl struct {
	FileName string
	Ext      string
	Content  string
}

func GetDefaultXormTpl() []XormTpl {
	return []XormTpl{
		XormTpl{
			FileName: "struct.go",
			Ext:      ".go",
			Content:  GetXormGOTpl(),
		},
	}
}

func GetXormGOTpl() string {
	return `package torm

{{$ilen := len .Imports}}
import (
        {{range .Imports}}"{{.}}"{{end}}
        "github.com/tal-tech/torm"
		"context"
       )

{{range .Tables}}
{{$tb := Mapper .Name}}
{{$table := .}}
{{$dao := printf "%sDao" $tb}}

type {{$tb}} struct {
    {{range .ColumnsSeq}}{{$col := $table.GetColumn .}} {{Mapper $col.Name}}    {{Type $col}} {{Tag $table $col}}
    {{end}}
}

type {{$dao}} struct {
    torm.DbBaseDao
	ctx context.Context
}

func New{{$dao}}(ctx context.Context, v ...interface{}) *{{$dao}} {
    this := new({{$dao}})
	this.ctx = ctx
	this.InitSession(ctx)
    if ins := torm.GetDbInstance("default", "writer"); ins != nil {
        this.UpdateEngine(ins.Engine)
    } else {
        return nil
    }
    if len(v) != 0 {
        this.UpdateEngine(v...)
    }
    return this
}

{{$pl := len .PrimaryKeys}}
{{if gt $pl 0}}
func (this *{{$dao}})Get({{genParams $table .PrimaryKeys true}}) (ret []{{$tb}}, err error) {
    ret = make([]{{$tb}},0)
    {{range .PrimaryKeys}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    err = this.Session.Find(&ret)
    return
}
func (this *{{$dao}})GetLimit({{genParams $table .PrimaryKeys true}}, pn, rn int) (ret []{{$tb}}, err error) {
    ret = make([]{{$tb}},0)
    {{range .PrimaryKeys}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    err = this.Session.Limit(rn,pn).Find(&ret)
    return
}
func (this *{{$dao}})GetCount({{genParams $table .PrimaryKeys true}}) (ret int64, err error) {
    {{range .PrimaryKeys}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    ret, err = this.Session.Count(new({{$tb}}))
    return
}
{{end}}

{{range .Indexes}}
func (this *{{$dao}})GetByIdx{{getMethodName .Name}}({{genParams $table .Cols true}}) (ret []{{$tb}}, err error) {
    ret = make([]{{$tb}},0)
    {{range .Cols}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    err = this.Session.Find(&ret)
    return
}
func (this *{{$dao}})GetByIdx{{getMethodName .Name}}Count({{genParams $table .Cols true}}) (ret int64, err error) {
    {{range .Cols}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    ret, err = this.Session.Count(new({{$tb}}))
    return
}
func (this *{{$dao}})GetByIdx{{getMethodName .Name}}Limit({{genParams $table .Cols true}}, pn,rn int) (ret []{{$tb}}, err error) {
    ret = make([]{{$tb}},0)
    {{range .Cols}}
    {{$p := Mapper .}}
    this.BuildQuery(m{{$p}}, "{{.}}")
    {{end}}
    err = this.Session.Limit(rn,pn).Find(&ret)
    return
}
{{end}}

{{end}}

`
}
