package car

import (
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/socket"
	"main/server/utils"

	socketio "github.com/googollee/go-socket.io"
)

func ColorCustomization(s socketio.Conn, req map[string]interface{}) {

	fmt.Println("Color customization socket hit")
	//fmt.Println("Request body of color customisation socket is:", req)
	playerId := s.Context().(string)

	//validation
	custId, ok := req["custId"].(string)
	if !ok {
		response.SocketResponse("Car id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	colorCategory, ok := req["colorCategory"].(string)
	if !ok {
		response.SocketResponse("color category is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	colorTypeId, ok := req["colorTypeId"].(float64)
	if !ok {
		response.SocketResponse("color type id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	var colorType string
	switch int64(colorTypeId) {
	case int64(utils.FLUORESCENT):
		colorType = "fluorescent"
	case int64(utils.PASTEL):
		colorType = "pastel"
	case int64(utils.GUN_METAL):
		colorType = "gun_Metal"
	case int64(utils.SATIN):
		colorType = "satin"
	case int64(utils.METAL):
		colorType = "metal"
	case int64(utils.MILITARY):
		colorType = "military"
	}

	colorId, ok := req["colorId"].(float64)
	if !ok {
		response.SocketResponse("color id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	var colorName string

	switch colorId {
	case float64(utils.CRED):
		colorName = "red"
	case float64(utils.CGREEN):
		colorName = "green"
	case float64(utils.CPINK):
		colorName = "pink"
	case float64(utils.CYELLOW):
		colorName = "yellow"
	case float64(utils.CBLUE):
		colorName = "blue"
	}

	if colorType == "military" {
		switch colorId {
		case float64(utils.MCBLACK):
			colorName = "black"
		case float64(utils.MCDESERT):
			colorName = "desert"
		case float64(utils.MCTRAM):
			colorName = "tram"
		case float64(utils.MCUCP):
			colorName = "ucp"
		}

	}

	//check if that player has bought taht color or not for a specifiic car
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM player_car_customisations WHERE color_category=? AND color_type=? AND color_name=? AND cust_id=? and player_id=?)"
	err := db.QueryExecutor(query, &exists, colorCategory, colorType, colorName, custId, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}
	if exists {
		response.SocketResponse(utils.COLOR_ALREADY_BOUGHT, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//check that the color is present in tht database or not
	query = "SELECT EXISTS(SELECT * FROM part_customizations WHERE color_category=? AND color_type=? AND color_name=?)"
	err = db.QueryExecutor(query, &exists, colorCategory, colorType, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	if !exists {
		response.SocketResponse(utils.COLOR_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	var customise model.PartCustomization
	query = "SELECT * FROM part_customizations WHERE color_category=? AND color_type=? AND color_name=?"
	err = db.QueryExecutor(query, &customise, colorCategory, colorType, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	var amount int64
	if customise.CurrType == "coins" {
		amount = playerDetails.Coins
	} else {
		amount = playerDetails.Cash
	}
	if amount < int64(customise.CurrAmount) {
		response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	// Start a database transaction to handle the purchase.
	tx := db.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Deduct the currency from the player's account.
	if customise.CurrType == "coins" {
		playerDetails.Coins -= int64(customise.CurrAmount)
	} else {
		playerDetails.Cash -= int64(customise.CurrAmount)
	}

	//update player details
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//update the color

	query = "UPDATE player_car_customizations SET color_name=? ,color_category=?, color_type=? WHERE player_id=? AND cust_id=? "
	err = tx.Exec(query, colorName, colorCategory, colorType, playerId, custId).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//bradcasting the player details
	playerResponse, err := socket.GetPlayerDetailsCopy(playerId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//braodcasting the updated player details to the front end

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", *playerResponse)
	response.SocketResponse(utils.COLOR_CUSTOMIZED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, req, "colorCustomise", s)
}

func WheelCustomize(s socketio.Conn, req map[string]interface{}) {
	fmt.Println("Wheel customize socket handler called.")
	playerId := s.Context().(string)

	//validation
	custId, ok := req["custId"].(string)
	if !ok {
		response.SocketResponse(utils.CUSTID_REQUIRED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	wheelCategory, ok := req["colorCategory"].(string)
	if !ok {
		response.SocketResponse("color category is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	colorId, ok := req["colorId"].(float64)
	if !ok {
		response.SocketResponse("color id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	if !utils.IsCarBought(playerId, custId) {
		response.SocketResponse(utils.BUY_CAR_ERROR, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	//mappping of wheel colors
	var colorName string
	switch colorId {
	case float64(utils.CCBLACK):
		colorName = "black"
	case float64(utils.CCBLUE):
		colorName = "blue"
	case float64(utils.CCGREEN):
		colorName = "green"
	case float64(utils.CCPINK):
		colorName = "pink"
	case float64(utils.CCRED):
		colorName = "red"
	case float64(utils.CCYELLOW):
		colorName = "yellow"
	}

	var exists bool
	query := "SELECT EXISTS(SELECT * FROM player_car_customizations WHERE wheel_category=? AND wheel_color_name=? AND cust_id=? and player_id=?)"
	err := db.QueryExecutor(query, &exists, wheelCategory, colorName, custId, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}
	if exists {
		response.SocketResponse(utils.COLOR_ALREADY_BOUGHT, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//check that the color is present in tht database or not
	query = "SELECT EXISTS(SELECT * FROM part_customizations WHERE  wheel_category=? AND wheel_color_name=?)"
	err = db.QueryExecutor(query, &exists, wheelCategory, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "wheelCustomise", s)

		return
	}

	if !exists {
		response.SocketResponse(utils.COLOR_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	var customise model.PartCustomization
	query = "SELECT * FROM part_customizations WHERE  wheel_category=? AND wheel_color_name=?"
	err = db.QueryExecutor(query, &customise, wheelCategory, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	var amount int64
	if customise.CurrType == "coins" {
		amount = playerDetails.Coins
	} else {
		amount = playerDetails.Cash
	}
	if amount < int64(customise.CurrAmount) {
		response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	// Start a database transaction to handle the purchase.
	tx := db.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Deduct the currency from the player's account.
	if customise.CurrType == "coins" {
		playerDetails.Coins -= int64(customise.CurrAmount)
	} else {
		playerDetails.Cash -= int64(customise.CurrAmount)
	}

	//update player details
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	//update the wheel color with proper subcateory

	query = "UPDATE player_car_customizations SET wheel_color_name=?, wheel_category=? WHERE player_id=? AND cust_id=?"
	err = tx.Exec(query, colorName, wheelCategory, playerId, custId).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	//bradcasting the player details
	playerResponse, err := socket.GetPlayerDetailsCopy(playerDetails.PlayerId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "wheelCustomise", s)
		return
	}

	//braodcasting the updated player details to the front end

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", *playerResponse)
	response.SocketResponse(utils.WHEELS_CUSTOMIZED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, req, "wheelCustomise", s)
}

func InteriorCustomize(s socketio.Conn, req map[string]interface{}) {

	fmt.Println("Interior customize socket handler called")
	playerId := s.Context().(string)
	//check the custId s of selected car only
	custId, ok := req["custId"].(string)
	if !ok {
		response.SocketResponse(utils.CUSTID_REQUIRED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	colorId, ok := req["colorId"].(float64)
	if !ok {
		response.SocketResponse("color id is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	var colorName string
	switch colorId {
	case float64(utils.ITWHITE):
		colorName = "white"
	case float64(utils.ITPINK):
		colorName = "pink"
	case float64(utils.ITGREEN):
		colorName = "green"
	case float64(utils.ITRED):
		colorName = "red"
	case float64(utils.ITBLUE):
		colorName = "blue"
	case float64(utils.ITYELLOW):
		colorName = "yellow"
	}

	var exists bool
	query := "SELECT EXISTS(SELECT * FROM player_car_customisations WHERE interior_color_name=? AND cust_id=? and player_id=?)"
	err := db.QueryExecutor(query, &exists, colorName, custId, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "colorCustomise", s)
		return
	}
	if exists {
		response.SocketResponse(utils.COLOR_ALREADY_BOUGHT, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "colorCustomise", s)
		return
	}

	//check that the color is present in tht database or not
	query = "SELECT EXISTS(SELECT * FROM part_customizations WHERE interior_color_name=?)"
	err = db.QueryExecutor(query, &exists, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	if !exists {
		response.SocketResponse(utils.COLOR_NOT_FOUND, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	var customise model.PartCustomization
	query = "SELECT * FROM part_customizations WHERE interior_color_name=?"
	err = db.QueryExecutor(query, &customise, colorName)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	var amount int64
	if customise.CurrType == "coins" {
		amount = playerDetails.Coins
	} else {
		amount = playerDetails.Cash
	}
	if amount < int64(customise.CurrAmount) {
		response.SocketResponse(utils.NOT_ENOUGH_COINS, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	// Start a database transaction to handle the purchase.
	tx := db.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Deduct the currency from the player's account.
	if customise.CurrType == "coins" {
		playerDetails.Coins -= int64(customise.CurrAmount)
	} else {
		playerDetails.Cash -= int64(customise.CurrAmount)
	}

	//update player details
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	//query to update interior color
	query = "UPDATE player_car_customizations SET interior_color_name=? WHERE player_id=? AND cust_id=?"

	err = db.RawExecutor(query, colorName, playerId, custId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "interiorCustomise", s)

		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	//bradcasting the player details
	playerResponse, err := socket.GetPlayerDetailsCopy(playerDetails.PlayerId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "interiorCustomise", s)
		return
	}

	//braodcasting the updated player details to the front end

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", *playerResponse)
	response.SocketResponse(utils.INTERIOR_CUSTOMIZED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, req, "interiorCustomise", s)
}

func LicenseCustomize(s socketio.Conn, req map[string]interface{}) {

	fmt.Println("Licensce customise socket handler called.")
	//check if the car id is equiped or not
	playerId := s.Context().(string)

	//validation
	custId, ok := req["custId"].(string)
	if !ok {
		response.SocketResponse(utils.CUSTID_REQUIRED, utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}
	value, ok := req["value"].(string)
	if !ok {
		response.SocketResponse("value is required", utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	playerDetails, err := utils.GetPlayerDetails(playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}
	var exists bool

	query := "SELECT EXISTS(SELECT * FROM player_car_customisations where cust_id=? and lp_value=? and player_id=?)"
	err = db.QueryExecutor(query, &exists, custId, value, playerId)
	if err != nil {
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}
	if exists {
		response.SocketResponse("Same value not allowed", utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	// Start a database transaction to handle the purchase.
	tx := db.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Deduct the currency from the player's account.
	playerDetails.Coins -= int64(500)

	//update player details
	err = tx.Where(utils.PLAYER_ID, playerId).Updates(&playerDetails).Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	query = "UPDATE player_car_customisations SET lp_value=? WHERE player_id=? AND cust_id=?"
	err = db.RawExecutor(query, value, playerId, custId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	// broadcasting the player details
	playerResponse, err := socket.GetPlayerDetailsCopy(playerDetails.PlayerId)
	if err != nil {
		tx.Rollback()
		response.SocketResponse(err.Error(), utils.HTTP_BAD_REQUEST, utils.FAILURE, nil, "licenseCustomise", s)
		return
	}

	// braodcasting the updated player details to the front end

	utils.SocketServerInstance.BroadcastToRoom("/", playerId, "playerDetails", *playerResponse)
	response.SocketResponse(utils.LICENSE_PLATE_CUSTOMIZED_SUCCESS, utils.HTTP_OK, utils.SUCCESS, req, "licenseCustomise", s)

}
