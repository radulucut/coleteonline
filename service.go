package coleteonline

type ServiceExtraOption struct {
	Id             int64    `json:"id"`
	Name           string   `json:"name"`
	RequiredFields []string `json:"requiredFields"`
	OptionalFields []string `json:"optionalFields"`
}

type ServiceResponse struct {
	Id           int64                `json:"id"`
	CourierName  string               `json:"courierName"`
	Name         string               `json:"name"`
	ExtraOptions []ServiceExtraOption `json:"extraOptions"`
}
