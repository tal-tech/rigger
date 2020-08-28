// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/lunny/log"
	"github.com/spf13/cobra"
	"github.com/tal-tech/rigger/config"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"
)

type Reverse struct {
	single   bool
	tmplPath string
}

type ReverseOption func(r *Reverse) error

func NewReverse(option ...ReverseOption) *Reverse {
	r := new(Reverse)
	for _, f := range option {
		if err := f(r); err != nil {
			return nil
		}
	}
	return r
}

func SetReverseSingle(s bool) ReverseOption {
	return func(r *Reverse) error {
		r.single = s
		return nil
	}
}

func SetReverseTmplPath(path string) ReverseOption {
	return func(r *Reverse) error {
		r.tmplPath = path
		return nil
	}
}

var (
	genJson bool = false
)

func printReversePrompt(flag string) {
}

type Tmpl struct {
	Tables  []*core.Table
	Imports map[string]string
	Models  string
}

func dirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}

	return true
}

func (r *Reverse) Run(cmd *cobra.Command, args []string) {
	isMultiFile := !r.single

	curPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	var genDir string
	var model string
	var filterPat *regexp.Regexp
	var dir string
	if len(args) >= 3 {
		genDir, err = filepath.Abs(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		//[SWH|+] 经测试，path.Base不能解析windows下的“\”，需要替换为“/”
		genDir = strings.Replace(genDir, "\\", "/", -1)
		model = path.Base(genDir)

		if len(args) >= 4 {
			filterPat, err = regexp.Compile(args[3])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		model = "models"
		genDir = path.Join(curPath, model)
	}

	var configs map[string]string
	if r.tmplPath == "" {
		configs = config.DefaultXormConfig()
	} else {
		dir, err = filepath.Abs(r.tmplPath)
		if err != nil {
			log.Errorf("%v", err)
			return
		}

		if !dirExists(dir) {
			log.Errorf("Template %v path is not exist", dir)
			return
		}
		cfgPath := path.Join(dir, "config")
		info, err := os.Stat(cfgPath)
		if err == nil && !info.IsDir() {
			configs = loadConfig(cfgPath)
		}
	}

	var langTmpl LangTmpl
	var ok bool
	var lang string = "go"
	var prefix string = "" //[SWH|+]

	if l, ok := configs["lang"]; ok {
		lang = l
	}
	if j, ok := configs["genJson"]; ok {
		genJson, err = strconv.ParseBool(j)
	}

	//[SWH|+]
	if j, ok := configs["prefix"]; ok {
		prefix = j
	}

	if langTmpl, ok = langTmpls[lang]; !ok {
		fmt.Println("Unsupported programing language", lang)
		return
	}

	os.MkdirAll(genDir, os.ModePerm)

	supportComment = (args[0] == "mysql" || args[0] == "mymysql")

	Orm, err := xorm.NewEngine(args[0], args[1])
	if err != nil {
		log.Errorf("%v", err)
		return
	}

	tables, err := Orm.DBMetas()
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	if filterPat != nil && len(tables) > 0 {
		size := 0
		for _, t := range tables {
			if filterPat.MatchString(t.Name) {
				tables[size] = t
				size++
			}
		}
		tables = tables[:size]
	}

	commonRender := &Render{
		isMultiFile: isMultiFile,
		genDir:      genDir,
		prefix:      prefix,
		model:       model,
		tables:      tables,
		langTmpl:    langTmpl,
	}

	if r.tmplPath == "" {
		defaultTpl := config.GetDefaultXormTpl()
		for _, cfg := range defaultTpl {
			t := template.New(cfg.FileName)
			t.Funcs(langTmpl.Funcs)

			tmpl, err := t.Parse(cfg.Content)
			if err != nil {
				log.Errorf("%v", err)
				return
			}
			commonRender.Do(tmpl, cfg.FileName, cfg.Ext)
		}
	} else {
		filepath.Walk(dir, func(f string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if info.Name() == "config" {
				return nil
			}

			bs, err := ioutil.ReadFile(f)
			if err != nil {
				log.Errorf("%v", err)
				return err
			}
			t := template.New(f)
			t.Funcs(langTmpl.Funcs)

			tmpl, err := t.Parse(string(bs))
			if err != nil {
				log.Errorf("%v", err)
				return err
			}

			fileName := info.Name()
			newFileName := fileName[:len(fileName)-4]
			ext := path.Ext(newFileName)

			err = commonRender.Do(tmpl, newFileName, ext)
			if err != nil {
				return err
			}
			return nil
		})
	}

}

func table2Obj(name string) string {
	newstr := make([]rune, 0)
	upNextChar := false

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}
