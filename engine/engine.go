package engine

import (
	"fmt"
	"path/filepath"
)

type Engine struct {
	languages []*Language
}

func (eng *Engine) Languages() []*Language {
	return eng.languages
}

func (eng *Engine) Run(languageName string, opts *RunOptions) (*RunResult, error) {
	lang := eng.findLanguage(languageName)
	if lang == nil {
		return nil, fmt.Errorf("Language not found: '%s'", languageName)
	}

	return lang.Run(opts)
}

func New(confPath string, disableApparmor bool) (result *Engine, err error) {
	result = &Engine{}
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
		lang, err = loadLanguage(config)
		if err != nil {
			err = fmt.Errorf("Error loading language '%s': %s", config, err)
			// fail everything if one language fails to load
			return
		}
		if disableApparmor {
			lang.disableAppArmor()
		}
		result.languages = append(result.languages, lang)
	}

	return
}

func (engine *Engine) findLanguage(name string) *Language {
	for _, lang := range engine.Languages() {
		if lang.Name == name {
			return lang
		}
	}
	return nil
}
