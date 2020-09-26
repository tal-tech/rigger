package config

import "strings"

const (
	MacOS = "darwin"
)

type GetReplaceNameFn func(string) string

type ReplaceContentItem struct {
	Key string
	Fn  GetReplaceNameFn
}

type NewProjectInfo struct {
	TplRepo        string
	ReplaceContent []ReplaceContentItem
	ReplaceFile    map[string]string
	ReplaceDir     map[string]string
}

var NewOdinInfo NewProjectInfo = NewProjectInfo{
	TplRepo:        "git@github.com:tal-tech/odin.git",
	ReplaceContent: []ReplaceContentItem{ReplaceContentItem{"encoding", skipTemplateName}, ReplaceContentItem{"odin", DefaultReplaceName}, ReplaceContentItem{"Odin", TitleReplaceName}, ReplaceContentItem{"#TemplateName#", recoverEncoding}},
	ReplaceFile:    map[string]string{"odin": ""},
	ReplaceDir:     map[string]string{"odin": ""},
}

var NewGaeaInfo NewProjectInfo = NewProjectInfo{
	TplRepo:        "git@github.com:tal-tech/gaea.git",
	ReplaceContent: []ReplaceContentItem{ReplaceContentItem{"gaea", DefaultReplaceName}, ReplaceContentItem{"Gaea", TitleReplaceName}},
	ReplaceFile:    map[string]string{"gaea": ""},
	ReplaceDir:     map[string]string{"gaea": ""},
}

var NewTritonInfo NewProjectInfo = NewProjectInfo{
	TplRepo:        "git@github.com:tal-tech/triton.git",
	ReplaceContent: []ReplaceContentItem{ReplaceContentItem{"triton", DefaultReplaceName}},
	ReplaceFile:    map[string]string{},
	ReplaceDir:     map[string]string{},
}

var NewPanInfo NewProjectInfo = NewProjectInfo{
	TplRepo:        "git@github.com:tal-tech/pan.git",
	ReplaceContent: []ReplaceContentItem{{"panic", skipTemplateName}, {"pan", DefaultReplaceName}, {"#TemplateName#", recoverPanic}},
	ReplaceFile:    map[string]string{},
	ReplaceDir:     map[string]string{},
}

var NewJobInfo NewProjectInfo = NewProjectInfo{
	TplRepo:        "git@github.com:tal-tech/hera.git",
	ReplaceContent: []ReplaceContentItem{},
	ReplaceFile:    map[string]string{},
	ReplaceDir:     map[string]string{},
}

func DefaultReplaceName(in string) string {
	return in
}

func TitleReplaceName(in string) string {
	return strings.Title(in)
}

func skipTemplateName(in string) string {
	return "#TemplateName#"
}

func recoverEncoding(in string) string {
	return "encoding"
}

func recoverPanic(in string) string {
	return "panic"
}

func recoverOdin(in string) string {
	return "gaea"
}

func DefaultXormConfig() map[string]string {
	return map[string]string{
		"lang":    "go",
		"genJson": "0",
		"prefix":  "cos_",
	}
}
