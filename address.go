package coleteonline

type Contact struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Phone2  string `json:"phone2,omitempty"`
	Company string `json:"company,omitempty"`
	Email   string `json:"email,omitempty"`
}

type Address struct {
	CountryCode    string `json:"countryCode"`
	PostalCode     string `json:"postalCode"`
	City           string `json:"city"`
	County         string `json:"county"`
	CountyCode     string `json:"countyCode"`
	Street         string `json:"street"`
	Number         string `json:"number"`
	Building       string `json:"building,omitempty"`
	Entrance       string `json:"entrance,omitempty"`
	Intercom       string `json:"intercom,omitempty"`
	Floor          string `json:"floor,omitempty"`
	Apartment      string `json:"apartment,omitempty"`
	Landmark       string `json:"landmark,omitempty"`
	AdditionalInfo string `json:"additionalInfo,omitempty"`
}

type ValidationStrategyType string

const (
	ValidationStrategyTypeMinimal      ValidationStrategyType = "minimal"
	ValidationStrategyTypePriceMinimal ValidationStrategyType = "priceMinimal"
)

type OrderAddress struct {
	AddressId          int64                  `json:"addressId"`
	Contact            Contact                `json:"contact"`
	Address            Address                `json:"address"`
	ValidationStrategy ValidationStrategyType `json:"validationStrategy"`
}

type Pagination struct {
	TotalItems  int64 `json:"totalItems"`
	CurrentPage int64 `json:"currentPage"`
	TotalPages  int64 `json:"totalPages"`
}

type AddressListResponse struct {
	Data       []OrderAddress `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
