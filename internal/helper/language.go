package helper

func GetLanguageList() []string {
	return []string{"ru", "en"}
}

func GetDefaultLg() string {
	return "en"
}

func GetLang(lang string) string {
	for _, v := range GetLanguageList() {
		if v == lang {
			return v
		}
	}

	return GetDefaultLg()
}
