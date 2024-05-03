package infrastructure

import (
	"auth-api/internal/domain"
	"encoding/json"
	"os"
)

type CityQueryRepository struct {
	cities map[uint]*domain.City
}

func NewCityQueryRepository() *CityQueryRepository {
	file, err := os.Open("config/city.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	cities := make(map[uint]*domain.City)
	err = json.NewDecoder(file).Decode(&cities)
	if err != nil {
		panic(err)
	}

	return &CityQueryRepository{
		cities: cities,
	}
}

func (r *CityQueryRepository) FindByID(id uint) (*domain.City, error) {
	return r.cities[id], nil
}
