package coleteonline

import "time"

type Sender struct {
	AddressId          int64                  `json:"addressId,omitempty"`
	Contact            *Contact               `json:"contact,omitempty"`
	Address            *Address               `json:"address,omitempty"`
	ValidationStrategy ValidationStrategyType `json:"validationStrategy,omitempty"`
}

type Recipient struct {
	AddressId          int64                  `json:"addressId,omitempty"`
	Contact            *Contact               `json:"contact,omitempty"`
	Address            *Address               `json:"address,omitempty"`
	ValidationStrategy ValidationStrategyType `json:"validationStrategy,omitempty"`
}

type PackageType byte

const (
	PackageTypeEnvelope PackageType = 1
	PackageTypePackage  PackageType = 2
)

type Package struct {
	Weight float64 `json:"weight"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Length float64 `json:"length"`
}

type ServiceType string

const (
	ServiceTypeDirectId  ServiceType = "directId"
	ServiceTypeBestPrice ServiceType = "bestPrice"
	ServiceTypeGrade     ServiceType = "grade"
)

type ServiceGrade string

const (
	ServiceGradeDelivery  ServiceGrade = "delivery"
	ServiceGradePickup    ServiceGrade = "pickUp"
	ServiceGradeRepayment ServiceGrade = "repayment"
)

type OrderService struct {
	SelectionType ServiceType    `json:"selectionType"`
	ServiceIds    []int64        `json:"serviceIds,omitempty"`
	Grades        []ServiceGrade `json:"grades,omitempty"`
}

type Packages struct {
	Type    PackageType `json:"type"`
	Content string      `json:"content"`
	List    []Package   `json:"list"`
}

type Order struct {
	Sender       Sender        `json:"sender"`
	Recipient    Recipient     `json:"recipient"`
	Packages     Packages      `json:"packages"`
	Service      OrderService  `json:"service"`
	ExtraOptions []interface{} `json:"extraOptions,omitempty"`
}

type ServicePrice struct {
	Total float64 `json:"total"`
	NoVat float64 `json:"noVat"`
}

type ServiceDetails struct {
	Id          int64  `json:"id"`
	CourierName string `json:"courierName"`
	Name        string `json:"name"`
}

type OrderResponseService struct {
	Price   ServicePrice   `json:"price"`
	Service ServiceDetails `json:"service"`
}

type ExtraOptionId byte

const (
	ExtraOptionIdStatusChange     ExtraOptionId = 1
	ExtraOptionIdOpenAtDelivery   ExtraOptionId = 2
	ExtraOptionIdSaturdayDelivery ExtraOptionId = 3
	ExtraOptionIdInsurance        ExtraOptionId = 4
	ExtraOptionIdAccountRepayment ExtraOptionId = 5
	ExtraOptionIdCashRepayment    ExtraOptionId = 6
	ExtraOptionIdDeclaredValue    ExtraOptionId = 7
	ExtraOptionIdScheduledPickup  ExtraOptionId = 8
	ExtraOptionIdClientReference  ExtraOptionId = 9
	ExtraOptionIdBaseCurrency     ExtraOptionId = 10
)

type OrderResponse struct {
	Service             OrderResponseService `json:"service"`
	AWB                 string               `json:"awb"`
	UniqueId            string               `json:"uniqueId"`
	EstimatedPickupDate string               `json:"estimatedPickupDate"`
}

type OrderPriceResponse struct {
	Selected OrderResponseService   `json:"selected"`
	List     []OrderResponseService `json:"list"`
}

type StatusSummary struct {
	UniqueId string `json:"uniqueId"`
	AWB      string `json:"awb"`
}

type StatusTextPart struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type StatusTextParts struct {
	Ro StatusTextPart `json:"ro"`
}

type StatusComment struct {
	Ro string `json:"ro"`
}

type StatusHistory struct {
	DateTime        time.Time       `json:"dateTime"`
	UnixDateTime    int64           `json:"unixDateTime"`
	StatusTextParts StatusTextParts `json:"statusTextParts"`
	StatusComment   StatusComment   `json:"comment"`
	Code            int64           `json:"code"`
}

type OrderStatusResponse struct {
	Summary StatusSummary   `json:"summary"`
	History []StatusHistory `json:"history"`
}
