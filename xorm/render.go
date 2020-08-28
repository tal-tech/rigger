package xorm

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-xorm/core"
	"github.com/lunny/log"
)

type Render struct {
	isMultiFile bool
	genDir      string
	prefix      string
	model       string
	tables      []*core.Table
	langTmpl    LangTmpl
}

func (r *Render) Do(tmpl *template.Template, newFileName, ext string) error {
	var (
		w   *os.File
		err error
	)
	if !r.isMultiFile {
		w, err = os.Create(path.Join(r.genDir, newFileName))
		if err != nil {
			log.Errorf("%v", err)
			return err
		}

		imports := r.langTmpl.GenImports(r.tables)

		tbls := make([]*core.Table, 0)
		for _, table := range r.tables {
			//[SWH|+]
			if r.prefix != "" {
				table.Name = strings.TrimPrefix(table.Name, r.prefix)
			}
			tbls = append(tbls, table)
		}

		newbytes := bytes.NewBufferString("")

		t := &Tmpl{Tables: tbls, Imports: imports, Models: r.model}
		err = tmpl.Execute(newbytes, t)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}

		tplcontent, err := ioutil.ReadAll(newbytes)
		if err != nil {
			log.Errorf("%v", err)
			return err
		}
		var source string
		if r.langTmpl.Formater != nil {
			source, err = r.langTmpl.Formater(string(tplcontent))
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
		} else {
			source = string(tplcontent)
		}

		w.WriteString(source)
		w.Close()
	} else {
		for _, table := range r.tables {
			//[SWH|+]
			if r.prefix != "" {
				table.Name = strings.TrimPrefix(table.Name, r.prefix)
			}
			filename := table2Obj(table.Name)
			// imports
			tbs := []*core.Table{table}
			imports := r.langTmpl.GenImports(tbs)

			w, err := os.Create(path.Join(r.genDir, filename+ext))
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			defer w.Close()

			newbytes := bytes.NewBufferString("")

			t := &Tmpl{Tables: tbs, Imports: imports, Models: r.model}
			err = tmpl.Execute(newbytes, t)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}

			tplcontent, err := ioutil.ReadAll(newbytes)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			var source string
			if r.langTmpl.Formater != nil {
				source, err = r.langTmpl.Formater(string(tplcontent))
				if err != nil {
					log.Errorf("%v-%v", err, string(tplcontent))
					//return err
				}
			} else {
				source = string(tplcontent)
			}

			w.WriteString(source)
			w.Close()
		}
	}
	return nil
}
