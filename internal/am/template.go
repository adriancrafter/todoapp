package am

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	rootDir    = "assets/templates/web/"
	layoutPath = rootDir + "layout/layout.tmpl"
)

type (
	TemplateManager struct {
		*SimpleCore
		fs          embed.FS
		haltOnError bool
		cache       TemplateCache
	}
)

func NewTemplateManager(fs embed.FS, haltOnError bool, opts ...Option) *TemplateManager {
	return &TemplateManager{
		SimpleCore:  NewCore("template-manager", opts...),
		fs:          fs,
		haltOnError: haltOnError,
		cache:       TemplateCache{},
	}
}

func (tm *TemplateManager) Setup(ctx context.Context) (err error) {
	err = tm.load(rootDir)
	if err != nil {
		return errors.Wrap(err, "template manager load error")
	}

	tm.debugCache()

	return tm.Process()
}

func (tm *TemplateManager) Get(controllerName, handlerName string) (t *template.Template, err error) {
	bundle, ok := tm.cache[controllerName][handlerName]
	if !ok {
		return nil, errors.Wrapf(err, "template not found for %s#%s", controllerName, handlerName)
	}
	return bundle.Template()
}

type (
	TemplateBundle struct {
		files    *TemplateFiles
		template *template.Template
	}

	TemplateFiles struct {
		Layout   string
		Main     string
		Partials []string
	}

	// TemplateCache is a map of template bundles
	// It can be assumed that this is the general format of the map:
	// map[{controller-name}]map[{handler-name}]TemplateBundle
	TemplateCache map[string]map[string]TemplateBundle
)

func (tm *TemplateManager) load(dirPath string) error {
	dirPath = filepath.Clean(dirPath)
	entries, err := tm.fs.ReadDir(dirPath)
	if err != nil {
		if _, notExist := err.(*fs.PathError); notExist {
			tm.Log().Errorf("directory %s does not exist\n", dirPath)
			return nil
		}
		tm.Log().Errorf("error reading dir %s: %v\n", dirPath, err)
		return err
	}

	var partials []string
	var mainTemplates []string

	isLayoutDir := filepath.Base(dirPath) == "layout"
	parentDirName := filepath.Base(filepath.Dir(dirPath))

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			// Skip 'layout' directory inside controller directories
			if parentDirName != "web" && entry.Name() == "layout" {
				continue
			}
			if entry.Name() == "partial" {
				partials, err = tm.processPartials(fullPath)
				if err != nil {
					tm.Log().Errorf("error processing partial in %s: %v\n", fullPath, err)
				}
				continue
			}
			if !isLayoutDir {
				// Recursively load subdirectories, excluding 'layout'
				err := tm.load(fullPath)
				if err != nil {
					tm.Log().Errorf("error loading %s: %v\n", fullPath, err)
				}
			}
			continue
		}

		// Check if it's a main template
		isInLayoutsDir := strings.Contains(filepath.Dir(fullPath), filepath.Join(rootDir, parentDirName, "layout"))
		if !isLayoutDir && !strings.HasPrefix(entry.Name(), "_") && !isInLayoutsDir {
			mainTemplates = append(mainTemplates, fullPath)
		}
	}

	for _, mainTemplate := range mainTemplates {
		err := tm.updateCache(mainTemplate, partials)
		if err != nil {
			tm.Log().Errorf("error processing %s: %v\n", mainTemplate, err)
		}
	}

	return nil
}

func (tm *TemplateManager) processPartials(dirPath string) ([]string, error) {
	var partials []string

	entries, err := tm.fs.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "_") {
			partials = append(partials, fullPath)
		}
	}

	return partials, nil
}

func (tm *TemplateManager) updateCache(filePath string, partials []string) error {
	tm.Log().Debugf("processing file %s", filePath)
	fileName := filepath.Base(filePath)
	dirName := filepath.Base(filepath.Dir(filePath))

	controllerName := dirName
	handlerName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Initialize the bundle if it doesn't exist
	if _, exists := tm.cache[controllerName]; !exists {
		tm.cache[controllerName] = make(map[string]TemplateBundle)
	}

	bundle, exists := tm.cache[controllerName][handlerName]
	if !exists {
		bundle = TemplateBundle{
			files: &TemplateFiles{
				Layout:   layoutPath, // Default layout
				Partials: []string{},
			},
		}
	}

	layoutOverridePath := tm.determineLayoutPath(controllerName, handlerName)
	if layoutOverridePath != "" {
		bundle.files.Layout = layoutOverridePath
	}

	// If it's a main template, set Main
	if !strings.HasPrefix(fileName, "_") {
		bundle.files.Main = filePath
	}

	bundle.files.Partials = append(bundle.files.Partials, partials...)

	tm.cache[controllerName][handlerName] = bundle

	return nil
}

func (tm *TemplateManager) determineLayoutPath(controllerName, handlerName string) string {
	controllerLayoutDir := filepath.Join(rootDir, controllerName, "layout")
	entries, err := tm.fs.ReadDir(controllerLayoutDir)
	if err != nil {
		// ErrorMsg reading the directory, maybe it doesn't exist
		return ""
	}

	controllerLayout := "layout.tmpl"
	handlerLayout := handlerName + "_layout.tmpl"

	for _, entry := range entries {
		if entry.Name() == handlerLayout {
			return filepath.Join(controllerLayoutDir, handlerLayout)
		} else if entry.Name() == controllerLayout {
			controllerLayout = filepath.Join(controllerLayoutDir, controllerLayout) // Store path but keep looking
		}
	}

	if controllerLayout != "layout.tmpl" {
		// Return controller specific layout if handler specific layout was not found
		return controllerLayout
	}

	// Fallback to default layout if no specific layout were found
	return layoutPath
}

func (tm *TemplateManager) Process() error {
	for controllerName, cache := range tm.cache {
		for handlerName, bundle := range cache {
			all := append([]string{bundle.files.Layout, bundle.files.Main}, bundle.files.Partials...)
			t, err := template.New("base").Funcs(tm.TemplateFx()).ParseFS(tm.fs, all...)
			if err != nil {
				err := errors.Wrap(err, "template parse error for %s#%s", controllerName, handlerName)
				if tm.haltOnError {
					return err
				}
				tm.Log().Errorf(err.Error())
			}

			tm.cache[controllerName][handlerName] = TemplateBundle{
				template: t, // We don't need to store the files anymore
			}
		}
	}

	return nil
}

func (tb *TemplateBundle) Template() (t *template.Template, err error) {
	if tb.template == nil {
		return nil, errors.NewError("template not available")
	}
	return tb.template, nil
}

func (tm *TemplateManager) TemplateFx() template.FuncMap {
	return template.FuncMap{
		"toTitle": ToTitle,
		"concat":  Concat,
		"hasRole": HasRole,
	}
}

func ToTitle(str string) string {
	words := strings.Fields(str)
	dont := " a an on the to "

	for i, w := range words {
		if strings.Contains(dont, " "+w+" ") {
			words[i] = w
		} else {
			words[i] = strings.Title(w)
		}
	}
	return strings.Join(words, " ")
}

func Concat(strs ...string) string {
	return strings.Trim(strings.Join(strs, " "), " ")
}

func HasRole(role string, in ...string) bool {
	if len(in) == 0 {
		return false
	}

	for _, r := range in {
		if role == r {
			return true
		}
	}

	return false
}

func (tm *TemplateManager) debugCache() {
	tm.Log().Info("Templates cache:")
	for ctrlDir, v := range tm.cache {
		for handlerDir, bundle := range v {
			tm.Log().Infof("%s#%s:", ctrlDir, handlerDir)
			tm.Log().Info("===================")
			if bundle.files != nil {
				tm.Log().Info("* Files:")
				tm.Log().Infof("layout: %s", bundle.files.Layout)
				tm.Log().Infof("main: %s", bundle.files.Main)
				for _, p := range bundle.files.Partials {
					tm.Log().Infof("partial: %s", p)
				}
			}
			if bundle.template != nil {
				tm.Log().Info("* Template:")
				tm.Log().Infof("%+v", bundle.template)
			}
			tm.Log().Info("===================")
		}
	}
}
