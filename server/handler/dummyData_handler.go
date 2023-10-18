package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/server/db"
	"main/server/model"
	"main/server/utils"
)

func ReadJSONFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AddDummyDataHandler() {
	dataFiles := []struct {
		tableName string
		filePath  string
		dataPtr   interface{}
	}{
		{"cars", "server/dummyData/car.json", &[]model.Car{}},
		{"car_customizations", "server/dummyData/customization.json", &[]model.PartCustomization{}},
		{"default_customizations", "server/dummyData/defaultCustomization.json", &[]model.DefualtCustomisation{}},
		{"rewards", "server/dummyData/rewards.json", &[]model.RaceRewards{}},
		{"race_types", "server/dummyData/raceTypes.json", &[]model.RaceTypes{}},
		{"rating_multis", "server/dummyData/classMultiplier.json", &[]model.RatingMulti{}},
		{"arena_perks", "server/dummyData/arenaPerks.json", &[]model.ArenaLevelPerks{}},
	}

	for _, dataFile := range dataFiles {
		if !utils.TableIsEmpty(dataFile.tableName) {
			addtoDb(dataFile.filePath, dataFile.dataPtr)
		}
	}
}

func addtoDb(filePath string, modelType interface{}) {

	carData, err := ReadJSONFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(carData, &modelType)
	if err != nil {
		log.Fatal(err)
	}

	switch slice := modelType.(type) {
	case *[]model.Car:
		for _, item := range *slice {
			fmt.Println("Car data:", item)
			db.CreateRecord(&item)
		}
	case *[]model.DefualtCustomisation:
		for _, item := range *slice {
			fmt.Println("default customization data:", item)
			db.CreateRecord(&item)
		}
	case *[]model.RaceRewards:
		for _, item := range *slice {
			fmt.Println("win rewards data:", item)
			db.CreateRecord(&item)
		}

	case *[]model.RatingMulti:
		for _, item := range *slice {
			fmt.Println("rating data:", item)
			db.CreateRecord(&item)
		}
	case *[]model.RaceTypes:
		for _, item := range *slice {
			fmt.Println("race types data:", item)
			db.CreateRecord(&item)
		}
	case *[]model.ArenaLevelPerks:
		for _, item := range *slice {
			fmt.Println("Arena perks are:", item)
			db.CreateRecord(&item)
		}
	default:
		log.Fatal("Invalid modelType provided")
	}

}
