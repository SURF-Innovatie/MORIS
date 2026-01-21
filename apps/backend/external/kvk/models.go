package kvk

// SearchResponse represents the response from the Search API
// Based on https://developers.kvk.nl/nl/documentation/zoeken-api
type SearchResponse struct {
	Pagina              int          `json:"pagina"`
	ResultatenPerPagina int          `json:"resultatenPerPagina"`
	TotaalAantal        int          `json:"totaalAantal"`
	Resultaten          []ResultItem `json:"resultaten"`
}

// ResultItem represents a single item in the search results
type ResultItem struct {
	KvkNummer        string `json:"kvkNummer"`
	Handelsnaam      string `json:"handelsnaam"`
	Vestigingsnummer string `json:"vestigingsnummer,omitempty"`
	Type             string `json:"type"`
	Actief           string `json:"actief,omitempty"`
	VervallenNaam    string `json:"vervallenNaam,omitempty"`
	Adres            *Adres `json:"adres,omitempty"`
}

type Adres struct {
	BinnenlandsAdres *BinnenlandsAdres `json:"binnenlandsAdres,omitempty"`
}

type BinnenlandsAdres struct {
	Plaats     string `json:"plaats,omitempty"`
	Straatnaam string `json:"straatnaam,omitempty"`
	Type       string `json:"type,omitempty"`
}

// BasicProfile represents the response from the Basic Profile API
// Based on https://developers.kvk.nl/nl/documentation/open-dataset-basis-bedrijfsgegevens-api
type BasicProfile struct {
	KvkNummer               string                `json:"kvkNummer"`
	Naam                    string                `json:"naam,omitempty"` // Statutaire naam usually
	Handelsnamen            []Handelsnaam         `json:"handelsnamen,omitempty"`
	Rechtsvorm              string                `json:"rechtsvorm,omitempty"`
	FormeleRegistratiedatum string                `json:"formeleRegistratiedatum,omitempty"`
	MaterieleRegistratie    *MaterieleRegistratie `json:"materieleRegistratie,omitempty"`
	Hoofdvestiging          *Hoofdvestiging       `json:"hoofdvestiging,omitempty"`
}

type MaterieleRegistratie struct {
	DatumAanvang string `json:"datumAanvang,omitempty"`
	DatumEinde   string `json:"datumEinde,omitempty"`
}

type Hoofdvestiging struct {
	Vestigingsnummer  string `json:"vestigingsnummer,omitempty"`
	EersteHandelsnaam string `json:"eersteHandelsnaam,omitempty"`
}

type Handelsnaam struct {
	Naam     string `json:"naam,omitempty"`
	Volgorde int    `json:"volgorde,omitempty"`
}
