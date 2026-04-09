package google

import (
	"github.com/example/location-demo/internal/domain"
)

type GoogleType string

const (
	// --- HIERARCHY & POLITICAL (Cấu trúc hành chính) ---
	GoogleTypeCountry                  GoogleType = "country"
	GoogleTypeAdministrativeAreaLevel1 GoogleType = "administrative_area_level_1" // Tỉnh/Thành phố trực thuộc TW
	GoogleTypeAdministrativeAreaLevel2 GoogleType = "administrative_area_level_2" // Quận/Huyện
	GoogleTypeAdministrativeAreaLevel3 GoogleType = "administrative_area_level_3"
	GoogleTypeAdministrativeAreaLevel4 GoogleType = "administrative_area_level_4"
	GoogleTypeAdministrativeAreaLevel5 GoogleType = "administrative_area_level_5"
	GoogleTypeLocality                 GoogleType = "locality" // Thành phố/Thị xã

	GoogleTypeSublocality       GoogleType = "sublocality"
	GoogleTypeSublocalityLevel1 GoogleType = "sublocality_level_1" // Phường/Xã
	GoogleTypeSublocalityLevel2 GoogleType = "sublocality_level_2"

	GoogleTypeNeighborhood GoogleType = "neighborhood" // Khu phố
	GoogleTypePolitical    GoogleType = "political"

	// --- INFRASTRUCTURE & ADDRESS (Hạ tầng & Địa chỉ) ---
	GoogleTypePostalCode    GoogleType = "postal_code"
	GoogleTypeRoute         GoogleType = "route" // Tên đường
	GoogleTypeStreetAddress GoogleType = "street_address"
	GoogleTypePremise       GoogleType = "premise"      // Tòa nhà/Cơ sở
	GoogleTypeSubpremise    GoogleType = "subpremise"   // Căn hộ/Số phòng
	GoogleTypeIntersection  GoogleType = "intersection" // Ngã tư
	GoogleTypePlusCode      GoogleType = "plus_code"

	// --- LANDMARKS & RECREATION (Danh lam thắng cảnh & Giải trí) ---
	GoogleTypeTouristAttraction GoogleType = "tourist_attraction"
	GoogleTypePointOfInterest   GoogleType = "point_of_interest"
	GoogleTypeNaturalFeature    GoogleType = "natural_feature"
	GoogleTypePark              GoogleType = "park"
	GoogleTypeMuseum            GoogleType = "museum"
	GoogleTypeAmusementPark     GoogleType = "amusement_park"
	GoogleTypeAquarium          GoogleType = "aquarium"
	GoogleTypeArtGallery        GoogleType = "art_gallery"
	GoogleTypeStadium           GoogleType = "stadium"
	GoogleTypeZoo               GoogleType = "zoo"

	// --- VENUES & BUSINESSES (Địa điểm kinh doanh) ---
	GoogleTypeEstablishment   GoogleType = "establishment" // Phổ biến nhất (Fallback)
	GoogleTypeRestaurant      GoogleType = "restaurant"
	GoogleTypeCafe            GoogleType = "cafe"
	GoogleTypeBakery          GoogleType = "bakery"
	GoogleTypeBar             GoogleType = "bar"
	GoogleTypeNightClub       GoogleType = "night_club"
	GoogleTypeLodging         GoogleType = "lodging" // Khách sạn/Nhà nghỉ
	GoogleTypeHotel           GoogleType = "hotel"
	GoogleTypeShoppingMall    GoogleType = "shopping_mall"
	GoogleTypeDepartmentStore GoogleType = "department_store"
	GoogleTypeSupermarket     GoogleType = "supermarket"
	GoogleTypeGym             GoogleType = "gym"
	GoogleTypeSpa             GoogleType = "spa"
	GoogleTypeHospital        GoogleType = "hospital"
	GoogleTypePharmacy        GoogleType = "pharmacy"
	GoogleTypeBank            GoogleType = "bank"
	GoogleTypeAtm             GoogleType = "atm"

	// --- TRANSPORTATION (Giao thông) ---
	GoogleTypeAirport        GoogleType = "airport"
	GoogleTypeBusStation     GoogleType = "bus_station"
	GoogleTypeTrainStation   GoogleType = "train_station"
	GoogleTypeSubwayStation  GoogleType = "subway_station"
	GoogleTypeTransitStation GoogleType = "transit_station"
	GoogleTypeTaxiStand      GoogleType = "taxi_stand"
	GoogleTypeGasStation     GoogleType = "gas_station"
	GoogleTypeParking        GoogleType = "parking"
)

var AllowSearchTextTypeMap = map[GoogleType]bool{
	GoogleTypeCountry:                  true,
	GoogleTypeAdministrativeAreaLevel1: true,
	GoogleTypeAdministrativeAreaLevel2: true,
	GoogleTypeLocality:                 true,
	GoogleTypePostalCode:               true,
	GoogleTypeTouristAttraction:        true,
	GoogleTypePark:                     true,
	GoogleTypeMuseum:                   true,
	GoogleTypeAmusementPark:            true,
	GoogleTypeAquarium:                 true,
	GoogleTypeArtGallery:               true,
	GoogleTypeStadium:                  true,
	GoogleTypeZoo:                      true,
	GoogleTypeRestaurant:               true,
	GoogleTypeCafe:                     true,
	GoogleTypeBakery:                   true,
	GoogleTypeBar:                      true,
	GoogleTypeNightClub:                true,
	GoogleTypeLodging:                  true,
	GoogleTypeHotel:                    true,
	GoogleTypeShoppingMall:             true,
	GoogleTypeDepartmentStore:          true,
	GoogleTypeSupermarket:              true,
	GoogleTypeGym:                      true,
	GoogleTypeSpa:                      true,
	GoogleTypeHospital:                 true,
	GoogleTypePharmacy:                 true,
	GoogleTypeBank:                     true,
	GoogleTypeAtm:                      true,
	GoogleTypeAirport:                  true,
	GoogleTypeBusStation:               true,
	GoogleTypeTrainStation:             true,
	GoogleTypeSubwayStation:            true,
	GoogleTypeTransitStation:           true,
	GoogleTypeTaxiStand:                true,
	GoogleTypeGasStation:               true,
	GoogleTypeParking:                  true,
}

func IsAllowSearchTextType(t GoogleType) bool {
	return AllowSearchTextTypeMap[t]
}

// mapGoogleTypes matches Google's type response to our internal domain models accurately.
func mapGoogleTypes(types []string) domain.LocationType {
	// 1. Chuyển slice thành map để lookup chính xác (Exact Match)
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	// 2. Kiểm tra VENUE (Ưu tiên các cơ sở kinh doanh/dịch vụ)
	// Bao gồm các hằng số cụ thể bạn đã định nghĩa
	if typeMap[string(GoogleTypeRestaurant)] ||
		typeMap[string(GoogleTypeCafe)] ||
		typeMap[string(GoogleTypeHotel)] ||
		typeMap[string(GoogleTypeBank)] ||
		typeMap[string(GoogleTypeHospital)] ||
		typeMap[string(GoogleTypeSupermarket)] ||
		typeMap[string(GoogleTypeEstablishment)] {
		return domain.LocationTypeVenue
	}

	// 3. Kiểm tra LANDMARK (Các điểm đến du lịch, văn hóa)
	if typeMap[string(GoogleTypePointOfInterest)] ||
		typeMap[string(GoogleTypeTouristAttraction)] ||
		typeMap[string(GoogleTypeNaturalFeature)] ||
		typeMap[string(GoogleTypePark)] ||
		typeMap[string(GoogleTypeMuseum)] ||
		typeMap[string(GoogleTypeStadium)] {
		return domain.LocationTypeLandmark
	}

	// 4. Kiểm tra ADDRESS (Hạ tầng & Địa chỉ cụ thể)
	if typeMap[string(GoogleTypeStreetAddress)] ||
		typeMap[string(GoogleTypeRoute)] ||
		typeMap[string(GoogleTypePremise)] ||
		typeMap[string(GoogleTypeIntersection)] ||
		typeMap[string(GoogleTypeAirport)] ||
		typeMap[string(GoogleTypeBusStation)] {
		return domain.LocationTypeAddress
	}

	// 5. Kiểm tra HIERARCHY (Từ nhỏ đến lớn)

	// Cấp Phường/Xã/Khu phố
	if typeMap[string(GoogleTypeSublocalityLevel1)] ||
		typeMap[string(GoogleTypeSublocality)] ||
		typeMap[string(GoogleTypeNeighborhood)] {
		return domain.LocationTypeWard
	}

	// Cấp Quận/Huyện/Thành phố thuộc tỉnh (Locality tại VN thường là TP thuộc tỉnh)
	if typeMap[string(GoogleTypeAdministrativeAreaLevel2)] ||
		typeMap[string(GoogleTypeLocality)] {
		return domain.LocationTypeDistrict
	}

	// Cấp Tỉnh/Thành phố trực thuộc TW
	if typeMap[string(GoogleTypeAdministrativeAreaLevel1)] {
		return domain.LocationTypeCity
	}

	// Cấp Quốc gia
	if typeMap[string(GoogleTypeCountry)] {
		return domain.LocationTypeCountry
	}

	// 6. Fallback cuối cùng nếu không khớp bất kỳ loại hành chính nào
	return domain.LocationTypeAddress
}
