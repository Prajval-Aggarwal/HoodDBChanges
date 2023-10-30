package db

import (
	"fmt"
	"main/server/model"

	"gorm.io/gorm"
)

func AutoMigrateDatabase(db *gorm.DB) {

	var dbVersion model.DbVersion
	err := db.First(&dbVersion).Error
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("db version is:", dbVersion.Version)
	if dbVersion.Version < 1 {
		err := db.AutoMigrate(&model.Player{}, &model.PlayerCarCustomisation{}, &model.OwnedCars{}, &model.Car{}, &model.Garage{}, &model.OwnedGarage{}, &model.GarageCars{})
		if err != nil {
			panic(err)
		}
		db.Create(&model.DbVersion{
			Version: 1,
		})
		dbVersion.Version = 1
	}
	if dbVersion.Version < 2 {
		err := db.AutoMigrate(&model.PartCustomization{}, &model.DefaultCustomisation{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 2,
		})
		dbVersion.Version = 2
	}
	if dbVersion.Version < 4 {
		err := db.AutoMigrate(&model.PlayerRaceStats{}, &model.Admin{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 4,
		})
		dbVersion.Version = 4
	}
	if dbVersion.Version < 5 {
		err := db.AutoMigrate(&model.ResetSession{}, &model.Arena{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 5,
		})
		dbVersion.Version = 5
	}
	if dbVersion.Version < 6 {
		err := db.AutoMigrate(&model.Arena{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 6,
		})
		dbVersion.Version = 6
	}
	if dbVersion.Version < 7 {
		err := db.AutoMigrate(&model.RaceTypes{}, &model.RaceRewards{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 7,
		})
		dbVersion.Version = 7
	}
	if dbVersion.Version < 8 {
		err := db.AutoMigrate(&model.ArenaCars{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 8,
		})
		dbVersion.Version = 8
	}
	if dbVersion.Version < 9 {
		err := db.AutoMigrate(&model.PlayerLevel{}, &model.ArenaRaceRecord{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 9,
		})
		dbVersion.Version = 9
	}

	if dbVersion.Version < 10 {
		err := db.AutoMigrate(&model.Session{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 10,
		})
		dbVersion.Version = 10
	}

	if dbVersion.Version < 11 {
		err := db.AutoMigrate(&model.TempRaceRecords{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 11,
		})
		dbVersion.Version = 11
	}
	if dbVersion.Version < 12 {
		err := db.AutoMigrate(&model.ArenaLevelPerks{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 12,
		})
		dbVersion.Version = 12
	}
	if dbVersion.Version < 13 {
		err := db.AutoMigrate(&model.Shop{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 13,
		})
		dbVersion.Version = 13
	}
	if dbVersion.Version < 14 {
		err := db.AutoMigrate(&model.ArenaReward{})
		if err != nil {
			panic(err)
		}
		db.Where("version=?", dbVersion.Version).Updates(&model.DbVersion{
			Version: 14,
		})
		dbVersion.Version = 14
	}

}
