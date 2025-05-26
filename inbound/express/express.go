package inbound

import "fmt"

type GetPreImportManifestModel struct {
	UUID              string
	UploadLoggingUUID string
	DischargePort     string
	VasselName        string
	ArrivalDate       string
	CustomerName      string
	Details           []*GetPreImportManifestDetilModel
}

type GetPreImportManifestDetilModel struct {
	UUID                     string
	MasterAirWaybill         string
	HouseAirWaybill          string
	Category                 string
	ConsigneeTax             string
	ConsigneeBranch          string
	ConsigneeName            string
	ConsigneeAddress         string
	ConsigneeDistrict        string
	ConsigneeSubprovince     string
	ConsigneeProvince        string
	ConsigneePostcode        string
	ConsigneeCountryCode     string
	ConsigneeEmail           string
	ConsigneePhoneNumber     string
	ShipperName              string
	ShipperAddress           string
	ShipperDistrict          string
	ShipperSubprovince       string
	ShipperProvince          string
	ShipperPostcode          string
	ShipperCountryCode       string
	ShipperEmail             string
	ShipperPhoneNumber       string
	TariffCode               string
	TariffSequence           string
	StatisticalCode          string
	EnglishDescriptionOfGood string
	ThaiDescriptionOfGood    string
	Quantity                 int64
	QuantityUnitCode         string
	NetWeight                float64
	NetWeightUnitCode        string
	GrossWeight              float64
	GrossWeightUnitCode      string
	Package                  string
	PackageUnitCode          string
	CifValueForeign          float64
	FobValueForeign          float64
	ExchangeRate             float64
	CurrencyCode             string
	ShippingMark             string
	ConsignmentCountry       string
	FreightValueForeign      float64
	FreightCurrencyCode      string
	InsuranceValueForeign    float64
	InsuranceCurrencyCode    string
	OtherChargeValueForeign  string
	OtherChargeCurrencyCode  string
	InvoiceNo                string
	InvoiceDate              string
}

type UpdatePreImportManifestDetailModel struct {
	UUID                     string
	MasterAirWaybill         string
	HouseAirWaybill          string
	Category                 string
	ConsigneeTax             string
	ConsigneeBranch          string
	ConsigneeName            string
	ConsigneeAddress         string
	ConsigneeDistrict        string
	ConsigneeSubprovince     string
	ConsigneeProvince        string
	ConsigneePostcode        string
	ConsigneeCountryCode     string
	ConsigneeEmail           string
	ConsigneePhoneNumber     string
	ShipperName              string
	ShipperAddress           string
	ShipperDistrict          string
	ShipperSubprovince       string
	ShipperProvince          string
	ShipperPostcode          string
	ShipperCountryCode       string
	ShipperEmail             string
	ShipperPhoneNumber       string
	TariffCode               string
	TariffSequence           string
	StatisticalCode          string
	EnglishDescriptionOfGood string
	ThaiDescriptionOfGood    string
	Quantity                 int64
	QuantityUnitCode         string
	NetWeight                float64
	NetWeightUnitCode        string
	GrossWeight              float64
	GrossWeightUnitCode      string
	Package                  string
	PackageUnitCode          string
	CifValueForeign          float64
	FobValueForeign          float64
	ExchangeRate             float64
	CurrencyCode             string
	ShippingMark             string
	ConsignmentCountry       string
	FreightValueForeign      float64
	FreightCurrencyCode      string
	InsuranceValueForeign    float64
	InsuranceCurrencyCode    string
	OtherChargeValueForeign  string
	OtherChargeCurrencyCode  string
	InvoiceNo                string
	InvoiceDate              string
}

var expectedHeadersUploadUpdateRawPreImport = []string{
	"//1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52",
}

func validateHeadersUploadUpdateRawPreImport(headers []string) error {
	for i, expected := range expectedHeadersUploadUpdateRawPreImport {
		if i >= len(headers) {
			return fmt.Errorf("missing header at column %d: expected '%s'", i+1, expected)
		}
		if headers[i] != expected {
			return fmt.Errorf("header mismatch at column %d: expected '%s', got '%s'", i+1, expected, headers[i])
		}
	}
	return nil
}
