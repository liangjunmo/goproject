package codes

import (
	"github.com/liangjunmo/gocode"
)

type Language string

const (
	LanguageZhCn Language = "zh_CN"
)

var i18n = map[Language]map[gocode.Code]string{
	LanguageZhCn: zhCn,
}

func Translate(code gocode.Code, lang Language) string {
	if _, ok := i18n[lang]; !ok {
		lang = LanguageZhCn
	}

	if _, ok := i18n[lang][code]; !ok {
		code = Unknown
	}

	return i18n[lang][code]
}
