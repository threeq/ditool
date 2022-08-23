package templatex

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type Scanner struct {
	group *template.Template
}

func NewScanner(funcMap template.FuncMap) *Scanner {
	scanner := &Scanner{group: template.New("sql-template-root")}
	scanner.group.Funcs(funcMap)
	return scanner
}

func (scanner *Scanner) ScanPaths(c context.Context, dirs ...string) {
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println("sql template", path, "read error", err.Error())
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if ".gosql" != path[len(path)-6:] {
				return nil
			}

			tplName := strings.TrimPrefix(path, dir+"/")
			tplName = tplName[0 : len(tplName)-6]
			tplRaw := tplRawFilter(c, tplName, getFileRaw(path))
			_, err = scanner.group.New(tplName).Parse(tplRaw)
			if err != nil {
				log.Println("["+tplName+"]", "sql template load", path, "error：", err.Error())
			} else {
				log.Println("["+tplName+"]", "sql template load", path, "success.")
			}
			return nil
		})
	}

}

func (scanner *Scanner) Parse(c context.Context, tpl string, data any) (string, error) {
	buf := &strings.Builder{}
	err := scanner.group.ExecuteTemplate(buf, tpl, data)
	return buf.String(), err
}

func getFileRaw(path string) string {
	data, _ := os.ReadFile(path)
	return strings.TrimSpace(string(data))
}

func tplRawFilter(c context.Context, name string, raw string) string {

	return defineRename(c, name, raw)
}

func defineRename(c context.Context, name string, raw string) string {
	reg, err := regexp.Compile(`{{\s*define\s+"([^"]+)"\s*}}`)
	if err != nil {
		log.Println("define rename error. ", err.Error())
	}

	return reg.ReplaceAllString(raw, `{{define "`+name+`/$1"}}`)
}
