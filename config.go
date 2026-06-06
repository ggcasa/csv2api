package csv2api

type Config struct {
	Generator struct {
		NumarLinii    int64  `yaml:"numar_linii"`
		LungimeId     int64  `yaml:"lungime_id"`
		CaleIesireCsv string `yaml:"cale_iesire_csv"`
	} `yaml:"generator"`

	Migrare struct {
		CaleIntrareCsv string `yaml:"cale_intrare_csv"`
		CaleBazaDate   string `yaml:"cale_baza_date"`
		NumeTabel      string `yaml:"nume_tabel"`
	} `yaml:"migrare"`

	Server struct {
		Port        string `yaml:"port"`
		ApiEndpoint string `yaml:"api_endpoint"`
	} `yaml:"server"`

	Atribute struct {
		Marci            []string `yaml:"marci"`
		Culori           []string `yaml:"culori"`
		LocatieBlocuri   []string `yaml:"locatie_blocuri"`
		GreutateInterval struct {
			Min int `yaml:"min"`
			Max int `yaml:"max"`
		} `yaml:"greutate_interval"`
		MarimeInterval struct {
			Min int `yaml:"min"`
			Max int `yaml:"max"`
		} `yaml:"marime_interval"`
	} `yaml:"atribute"`
}
