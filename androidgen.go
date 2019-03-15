package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

const ktTemplate = `
{{- if ne .Config.Package "" -}} package {{.Config.Package}} {{- end}}

import com.tencent.mmkv.MMKV

/**
 * {{.Group.Name}}
 * Created by {{.Config.Author}} on {{.Date}}.
 */
class {{.FileName}} {

    companion object {
        val instance: {{.FileName}} by lazy { {{.FileName}}() }
    }

    private val mmkv: MMKV = MMKV.defaultMMKV()

    private constructor()

	{{with .Group -}} // === Total {{ len .Items }} ===
	{{range .Items}}
	/**
	 * {{.Comment}}
	 */
    var {{.FuncName}}: {{ typeSymbol .Type }}
        get() {
            return mmkv.decode{{ typeSymbol .Type }}("{{ .Name }}")
        }
        set(value) {
            mmkv.encode("{{ .Name }}", value)
		}
	{{end}}

	{{end}}
}
`

type KtGenerator struct {
	RootDir string
}

func (k *KtGenerator) BuildTemplate() *template.Template {
	funcMap := template.FuncMap{
		"typeSymbol": k.MapTypeSymbol,
	}
	temp := template.New("kotlin").Funcs(funcMap)
	t, err := temp.Parse(ktTemplate)
	if err != nil {
		return nil
	}
	return t
}

func (k *KtGenerator) GenSp(spConfig *SpConfig, groups []*SpGroup) {
	// package to dir struct
	packageName := strings.TrimSpace(spConfig.Package)
	rootDir := strings.TrimSpace(spConfig.ExportDir)
	if rootDir == "" {
		rootDir = "."
	}
	if !filepath.IsAbs(rootDir) {
		absDir, err := filepath.Abs(rootDir)
		if err != nil {
			// ignore
		} else {
			rootDir = absDir
		}

	}
	if packageName != "" {
		dirs := strings.Split(packageName, ".")
		packageDir := filepath.Join(rootDir, filepath.Join(dirs...))
		stat, err := os.Stat(packageDir)
		if err != nil && os.IsNotExist(err) {
			// not exist
			err := os.MkdirAll(packageDir, os.FileMode(0755))
			if err != nil {
				fmt.Println("make package dir error.", err)
			}
		} else if !stat.IsDir() {
			// error , package is a file
		}

		rootDir = packageDir
	}
	// template to file
	temp := k.BuildTemplate()
	if temp == nil {
		// error
		return
	}

	k.RootDir = rootDir

	wg := new(sync.WaitGroup)

	wg.Add(len(groups))
	for _, group := range groups {
		// create entity
		clazzName := fmt.Sprintf("%sSp", group.Name)
		fileName := clazzName + ".kt"
		entity := &GenEntity{
			Config: spConfig,
			Group:  group,
		}
		groupFile := filepath.Join(k.RootDir, fileName)

		entity.FileName = clazzName
		entity.FilePath = groupFile

		entity.Date = time.Now().Format("15:04:05 2006-01-02")

		go func() {
			k.ExecTemplate(temp, groupFile, entity)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (k *KtGenerator) ExecTemplate(temp *template.Template, file string, entity *GenEntity) error {
	exportFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer exportFile.Close()
	return temp.Execute(exportFile, entity)
}

func (k *KtGenerator) MapTypeSymbol(itemType ItemType) string {
	switch itemType {
	case TypeNone:
		return "Unit"
	case TypeInt:
		return "Int"
	case TypeFloat:
		return "Float"
	case TypeLong:
		return "Long"
	case TypeString:
		return "String"
	default:
		return "Unit"
	}
}
