package vies

// VatCheckRequest is the request body for the VIES VAT check API.
type VatCheckRequest struct {
	CountryCode string `json:"countryCode"`
	VatNumber   string `json:"vatNumber"`
}

// VatCheckResponse is the response from the VIES VAT check API.
type VatCheckResponse struct {
	CountryCode            string `json:"countryCode"`
	VatNumber              string `json:"vatNumber"`
	RequestDate            string `json:"requestDate"`
	Valid                  bool   `json:"valid"`
	RequestIdentifier      string `json:"requestIdentifier"`
	Name                   string `json:"name"`
	Address                string `json:"address"`
	TraderName             string `json:"traderName"`
	TraderStreet           string `json:"traderStreet"`
	TraderPostalCode       string `json:"traderPostalCode"`
	TraderCity             string `json:"traderCity"`
	TraderCompanyType      string `json:"traderCompanyType"`
	TraderNameMatch        string `json:"traderNameMatch"`
	TraderStreetMatch      string `json:"traderStreetMatch"`
	TraderPostalCodeMatch  string `json:"traderPostalCodeMatch"`
	TraderCityMatch        string `json:"traderCityMatch"`
	TraderCompanyTypeMatch string `json:"traderCompanyTypeMatch"`
}

// ParsedAddress extracts city and postal code from the Address field.
func (r *VatCheckResponse) ParsedAddress() (city, postalCode string) {
	// Address format is typically: "\nSTREET NUMBER\nPOSTALCODE CITY\n"
	// We try to use TraderCity/TraderPostalCode first if available
	if r.TraderCity != "" && r.TraderCity != "---" {
		city = r.TraderCity
	}
	if r.TraderPostalCode != "" && r.TraderPostalCode != "---" {
		postalCode = r.TraderPostalCode
	}
	return city, postalCode
}
