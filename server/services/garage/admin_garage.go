package garage

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/request"
	"main/server/response"
	"main/server/utils"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddGarageService(ctx *gin.Context, addGarageReq request.AddGarageRequest) {

	// Check if a garage already exists at the specified latitude and longitude.
	var exists bool

	//check that no two same garages are on same locations
	query := "SELECT EXISTS (SELECT * FROM garages WHERE latitude=? AND longitude=?)"
	err := db.QueryExecutor(query, &exists, addGarageReq.Latitude, addGarageReq.Longitude)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	if exists {
		response.ShowResponse(utils.GARAGE_ALREADY_PRESENT, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	//capacity calculation

	capacity := int64(math.Round(float64(addGarageReq.GarageType)/2 + float64(addGarageReq.Rarity)/2))

	//	var newGarage model.Garage

	newGarage := model.Garage{
		GarageName:    addGarageReq.GarageName,
		Latitude:      addGarageReq.Latitude,
		Longitude:     addGarageReq.Longitude,
		Level:         addGarageReq.Level,
		CoinsRequired: addGarageReq.CoinsRequired,
		GarageType:    addGarageReq.GarageType,
		Rarity:        addGarageReq.Rarity,
		Capacity:      capacity,
		Locked:        true,
	}

	err = db.CreateRecord(&newGarage)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	// Return the success response along with the new garage details.
	response.ShowResponse(utils.GARAGE_ADD_SUCCESS, utils.HTTP_OK, utils.SUCCESS, newGarage, ctx)
}

func DeleteGarageService(ctx *gin.Context, deleteReq request.DeletGarageReq) {
	//validate the garage id
	if !db.RecordExist("garages", deleteReq.GarageId, "garage_id") {
		response.ShowResponse(utils.GARAGE_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	query := "DELETE FROM garages WHERE garage_id =?"
	err := db.RawExecutor(query, deleteReq.GarageId)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	response.ShowResponse(utils.GARAGE_DELETE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, nil, ctx)

}

func UpdateGarageService(ctx *gin.Context, updateReq request.UpdateGarageReq) {
	var garageDetails model.Garage

	//check if the garage exists or not
	if !db.RecordExist("garages", updateReq.GarageId, "garage_id") {
		response.ShowResponse(utils.GARAGE_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
		return
	}

	err := db.FindById(&garageDetails, updateReq.GarageId, "garage_id")
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	//Null check on the inputs
	if updateReq.GarageName != "" {
		garageDetails.GarageName = updateReq.GarageName
	}

	if updateReq.Latitude != 0 {
		garageDetails.Latitude = updateReq.Latitude
	}
	if updateReq.Longitude != 0 {
		garageDetails.Longitude = updateReq.Longitude
	}
	if updateReq.Level != 0 {
		garageDetails.Level = updateReq.Level
	}
	if updateReq.CoinsRequired != 0 {
		garageDetails.CoinsRequired = updateReq.CoinsRequired
	}

	if updateReq.GarageType != 0 {
		garageDetails.GarageType = updateReq.GarageType
	}

	if updateReq.Rarity != 0 {
		garageDetails.Rarity = updateReq.Rarity
	}

	err = db.UpdateRecord(&garageDetails, updateReq.GarageId, "garage_id").Error
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	response.ShowResponse(utils.GARAGE_UPDATE_SUCCESS, utils.HTTP_OK, utils.SUCCESS, garageDetails, ctx)

}
func GetAllGarageListService(ctx *gin.Context) {
	var garageList = []model.Garage{}
	var dataresp response.DataResponse
	// Get the query parameters for skip and limit from the request
	skipParam := ctx.DefaultQuery("skip", "0")
	limitParam := ctx.Query("limit")
	var query string
	if limitParam == "" {
		query = "SELECT * FROM garages ORDER BY created_at DESC"
	} else {
		skip, err := strconv.Atoi(skipParam)
		if err != nil {
			response.ShowResponse("Invalid skip value", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}

		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			response.ShowResponse("Invalid limit value", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, ctx)
			return
		}
		query = fmt.Sprintf("SELECT * FROM garages ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, skip)

	}

	// Get the query parameters for skip and limit from the request
	err := db.QueryExecutor(query, &garageList)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM garages"
	err = db.QueryExecutor(countQuery, &totalCount)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		return
	}
	dataresp.TotalCount = totalCount
	dataresp.Data = garageList

	response.ShowResponse(utils.GARAGE_LIST_FETCHED, utils.HTTP_OK, utils.SUCCESS, dataresp, ctx)
}

func GetGarageTypes(ctx *gin.Context) {
	var garageTypeList = []struct {
		TypeName string `json:"label,omitempty" gorm:"unique"`
		TypeId   int    `json:"value"`
	}{
		struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "The Mu",
			TypeId:   int(utils.THE_MU),
		}, struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "Red's Hotspot ",
			TypeId:   int(utils.REDS_HOTSPOT),
		}, struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "The Bear's Hideaway",
			TypeId:   int(utils.THE_BEARS_HIDEAWAY),
		}, struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "Princes Palace",
			TypeId:   int(utils.PRINCES_PALACE),
		}, struct {
			TypeName string "json:\"label,omitempty\" gorm:\"unique\""
			TypeId   int    "json:\"value\""
		}{
			TypeName: "The Great Spot",
			TypeId:   int(utils.THE_GREAT_SPOT),
		},
	}
	var dataresp response.DataResponse
	// Get the query parameters for skip and limit from the request

	// Build the SQL query with skip and limit

	dataresp.TotalCount = len(garageTypeList)
	dataresp.Data = garageTypeList

	response.ShowResponse(utils.GARAGE_LIST_FETCHED, utils.HTTP_OK, utils.SUCCESS, dataresp, ctx)
}

func GetRarity(ctx *gin.Context) {
	var rarityList = []struct {
		RarityId   int    `json:"value"`
		RarityName string `json:"label"`
	}{
		struct {
			RarityId   int    "json:\"value\""
			RarityName string "json:\"label\""
		}{
			RarityId:   utils.D,
			RarityName: "D class",
		},
		struct {
			RarityId   int    "json:\"value\""
			RarityName string "json:\"label\""
		}{
			RarityId:   utils.C,
			RarityName: "C class",
		}, struct {
			RarityId   int    "json:\"value\""
			RarityName string "json:\"label\""
		}{
			RarityId:   utils.B,
			RarityName: "B class",
		}, struct {
			RarityId   int    "json:\"value\""
			RarityName string "json:\"label\""
		}{
			RarityId:   utils.A,
			RarityName: "A class",
		}, struct {
			RarityId   int    "json:\"value\""
			RarityName string "json:\"label\""
		}{
			RarityId:   utils.S,
			RarityName: "S class",
		},
	}
	var dataresp response.DataResponse
	// Get the query parameters for skip and limit from the request

	dataresp.TotalCount = len(rarityList)
	dataresp.Data = rarityList

	response.ShowResponse(utils.GARAGE_LIST_FETCHED, utils.HTTP_OK, utils.SUCCESS, dataresp, ctx)
}
