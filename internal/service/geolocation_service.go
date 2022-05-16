package service

import "fmt"

var countries = map[string]string{
	"NG": "Nigeria",
}

type GeoLocationService struct {
}

func NewGeoLocationService() *GeoLocationService {
	return &GeoLocationService{}
}

func (g *GeoLocationService) GetCountryNameByCode(code string) (string, error) {
	if name, ok := countries[code]; ok {
		return name, nil
	}
	return "", fmt.Errorf("country code %s not found", code)
}
