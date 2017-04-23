package main

type Options struct {
	MonthName       map[string]string
	link            string
	Source          string
	CloudTags       []string
	DescriptionTags []string
	Translations    map[string]string
	ShowAllLimit    int
	ShowRandomLimit int
}

var options = Options{
	MonthName: map[string]string{"01": "Leden", "02": "Únor", "03": "Březen", "04": "Duben", "05": "Květen", "06": "Červen", "07": "Červenec", "08": "Srpen", "09": "Září", "10": "Říjen", "11": "Listopad", "12": "Prosinec"},
	link:      "/gallery",
	Source:    "/Users/rdvorak/Pictures",
	CloudTags: []string{"Rating", "Keyword", "Month", "Year", "State", "Country", "Location", "Sublocation", "Geoname"},
	// options.CloudTags = []string{"Rating", "FocalLength", "Keyword", "Month", "Year", "State", "Country", "Location", "Sublocation", "Geoname"}
	DescriptionTags: []string{"Rating", "Model", "Lens", "Keyword", "Month", "Year", "State", "Country", "Location", "Sublocation", "Geoname"},
	Translations: map[string]string{
		"Pozdechov": "Pozděchov",
		"Mala":      "Malá",
		"Velka":     "Velká",
		"Velke":     "Velké",
		"Vratna":    "Vrátna",
		"Vysoke":    "Vysoké",
		"Nizke":     "Nízke",
		"Javorniky": "Javorníky",
		"Krkonose":  "Krkonoše",
		"Jeseniky":  "Jeseníky",
		"Sulov":     "Súlov",
		"Kutna":     "Kutná",
		"Vsetin":    "Vsetín"},
	ShowAllLimit:    200,
	ShowRandomLimit: 100,
}
