package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"
)

// templates is a global var that will hold all of the parsed module templates
var templates = &template.Template{}

var (
	// moduleRoot points to the directory that contains the module templates
	moduleRoot = "./data/modules/*"
	// themesRoot points to the directory that contains the themes in the format of JSON files
	themesRoot = "./data/themes/"
)

// themes contain the json globals that can be referenced by a Modules
var themes = map[string]map[string]interface{}{}

// Parse modules into a template
// NOTE: init functions get run automatically at runtime before main is executed
func init() {
	// Use '[[  ]]' instead of the default '{{  }}'
	templates = templates.Delims("[[", "]]")

	// Parse all templates in the modules directory
	var err error
	templates, err = templates.ParseGlob(moduleRoot)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Parse all themes in the themes directory
	//
	// Get the list of files in the themes directory
	files, err := ioutil.ReadDir(themesRoot)
	if err != nil {
		log.Fatalf("Error reading themes: %v", err)
	}

	// Loop over the files we found in the themes
	// dir and parse thier contents
	for _, f := range files {
		// Read the contents of the file
		b, err := ioutil.ReadFile(themesRoot + f.Name())
		if err != nil {
			log.Fatalf("Error reading theme '%s': %v", f.Name(), err)
		}

		// Parse contents assuming each file is a JSON file
		theme := map[string]interface{}{}
		if err := json.Unmarshal(b, &theme); err != nil {
			log.Fatalf("Error decoding json theme '%s': %v", f.Name(), err)
		}

		// add the parsed contents by the name of the file to the global
		// var named theme
		themes[f.Name()] = theme
	}
}

// Files represents a collection of Modules.
type Files struct {
	OutputDir string                 `json:"outputDir"`
	Files     []File                 `json:"files"`
	Globals   map[string]interface{} `json:"globals"`
	Themes    []string               `json:"themes"`
}

func (fs *Files) Build() (string, error) {
	fs.mergeVars()
	output := ""
	sep := ""
	for _, f := range fs.Files {
		o, err := f.Build()
		if err != nil {
			return "", err
		}
		output += sep
		output += o
		sep = "\n\nNEW FILE\n\n"
	}
	return output, nil
}

// Add all of the globals and theme vars to each module, this allows
// all of the modules to have access to all of the variables
func (fs *Files) mergeVars() {
	for _, f := range fs.Files {
		for _, m := range f.Modules {
			updateMap(m.Keys, f.Globals)
			updateMap(m.Keys, fs.Globals)
			addThemesToMap(m.Keys, f.Themes)
			addThemesToMap(m.Keys, fs.Themes)
		}
		f.outputDir = fs.OutputDir
	}
}

// addThemesToMap adds the given themes to the passed in map m, if a
// theme is not found an error is returned
func addThemesToMap(m map[string]interface{}, thms []string) error {
	for _, t := range thms {
		thm, ok := themes[t]
		if !ok {
			return fmt.Errorf("Theme '%s' not found", t)
		}
		updateMap(m, thm)
	}
	return nil
}

// updateMap adds m2 keys into m1, if m1 already contains the key
// m1 is not updated
func updateMap(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		if _, ok := m1[k]; !ok {
			m1[k] = v
		}
	}
}

// File represents a collection of modules
type File struct {
	Modules  []Module               `json:"modules"` // The modules that make up the finished template
	Globals  map[string]interface{} `json:"globals"` // Globals can be applied to all of the modules
	Themes   []string               `json:"themes"`
	Filename string                 `json:"filename"` // The name of the file the finished template will be written

	outputDir string
}

// Module represents a single template module
type Module struct {
	Name string                 `json:"name"`
	Keys map[string]interface{} `json:"values"`
}

func (f *File) BuildAndWrite() error {
	output, err := f.Build()
	if err != nil {
		return err
	}
	if len(f.Filename) == 0 {
		return fmt.Errorf("cannot write to file, missing filename")
	}
	if len(f.outputDir) == 0 {
		return fmt.Errorf("cannot write to file, missing OutputDir")
	}

	fp := filepath.Join(f.outputDir, f.Filename)
	if err := ioutil.WriteFile(fp, []byte(output), 0664); err != nil {
		return fmt.Errorf("Error writing file: %v for filename: %v", err, fp)
	}

	return nil
}

// Build
func (f *File) Build() (string, error) {
	output := ""

	for _, mod := range f.Modules {
		buf := &bytes.Buffer{}
		if err := templates.ExecuteTemplate(buf, mod.Name, mod.Keys); err != nil {
			return "", fmt.Errorf("Error executing module '%s': %v", mod.Name, err)
		}

		output += buf.String()
	}

	return output, nil
}
