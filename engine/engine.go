package engine

import (
	"fmt"
	"path/filepath"
)

// good times global
var TheEngine Engine

type Engine struct {
	Languages []*Language
}

func LoadConfig(confPath string) (result Engine, err error) {
	configs, err := filepath.Glob(confPath + "/lang-*.yml")
	if err != nil {
		return
	}

	if len(configs) < 1 {
		err = fmt.Errorf("no languages found at path '%s'", confPath)
		return
	}

	for _, config := range configs {
		var lang *Language
		lang, err = LoadLanguage(config)
		if err != nil {
			err = fmt.Errorf("Error loading language '%s': %s", config, err)
			// fail everything if one language fails to load
			return
		}
		result.Languages = append(result.Languages, lang)
	}

	return
}

func (engine *Engine) FindLanguage(name string) (*Language, error) {
	for _, lang := range engine.Languages {
		if lang.Name == name {
			return lang, nil
		}
	}
	return nil, fmt.Errorf("Language '%s' not found", name)
}
