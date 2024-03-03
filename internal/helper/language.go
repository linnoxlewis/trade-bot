package helper

const (
	RuLang = "ru"
	EnLang = "en"
)

func GetLanguageList() []string {
	return []string{RuLang, EnLang}
}

func GetDefaultLg() string {
	return EnLang
}

func GetLang(lang string) string {
	for _, v := range GetLanguageList() {
		if v == lang {
			return v
		}
	}

	return GetDefaultLg()
}
