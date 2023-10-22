package utils

const (
	ARENA_LEVEL_NOT_VALID  string = "Arena level not valid"
	ARENA_NOT_FOUND        string = "Arena not found"
	CAR_NOT_OWNED          string = "Car not owned"
	CAR_NOT_FOUND          string = "Car not found"
	ARENA_NOT_OWNED        string = "Arena not owned"
	USER_NOT_FOUND         string = "User not found"
	EMAIL_NOT_FOUND        string = "Email not found"
	COLOR_NOT_FOUND        string = "Color not found"
	GARAGE_TYPE_NOT_VALILD string = "Garage type is not valid"
	GARAGE_NOT_FOUND       string = "Garage not found"
	GARAGE_NOT_OWNED       string = "Garage is not owned"
)
const (
	EASY_PERK   string = "@every 30m"
	MEDIUM_PERK string = "@every 3h"
	HARD_PERK   string = "@every 7h"

	EASY_ARENA_SLOT   int64 = 3
	MEDIUM_ARENA_SLOT int64 = 5
	HARD_ARENA_SLOT   int64 = 7

	ARENA_ADD_SUCCESS     string = "Arena Added successfully"
	ARENA_DELETE_SUCCESS  string = "Arena Deleted successfully"
	ARENA_UPDATE_SUCCESS  string = "Arena Updated successfully"
	ARENA_ALREADY_PRESENT string = "Arena already present at that location"
	ARENA_ALREADY_OWNED   string = "Player already owns the arena"
)

// Car constants
const (
	CAR_ADDED_SUCCESS string = "Car added  sucessfully"
	BUY_CAR_ERROR     string = "Car need to be bought first"

	EQUIP_CORRECT_CAR                string = "Car need to be selected first"
	CAR_ALREADY_ALLOTTED             string = "Car already alotted to others arena"
	CAR_LIMIT_REACHED                string = "Car Limit reached upgarde the garage to increse the limit"
	CAR_SELECTED_SUCCESS             string = "Current car selected successfully"
	CAR_ALREADY_BOUGHT               string = "Car already bought"
	CAR_REPLACED_SUCCESS             string = "Car replaced sucessfully"
	CAR_SOLD_SUCCESS                 string = "Car sold sucessfully"
	CAR_BOUGHT_SUCESS                string = "Car bought successfully"
	UPGRADE_SUCCESS                  string = "Part upgraged successfully"
	LICENSE_PLATE_CUSTOMIZED_SUCCESS string = "License Plate updated sucessfuly"
	INTERIOR_CUSTOMIZED_SUCCESS      string = "Interior updated sucessfuly"
	COLOR_ALREADY_BOUGHT             string = "Color already bought"
	WHEELS_CUSTOMIZED_SUCCESS        string = "Wheels updated succesfully"
	COLOR_CUSTOMIZED_SUCCESS         string = "Color updated succesfully"
	CAR_REPAIR_SUCCESS               string = "Car repaired successfully"
	NO_CARS_ADDED                    string = "No more cars can be added"
	UPGRADE_REACHED_MAX_LEVEL        string = "Part reached to its max level"
	PURCHASE_SUCCESS                 string = "Purchase succesfull"
	CUSTID_REQUIRED                  string = "Car customisation id is required"
)

// Garage constants
const (
	GARAGE_BOUGHT_SUCESS     string = "Garage bought successfully"
	GARAGE_LIST_FETCHED      string = "Garage list fetched successfully"
	GARAGE_UPGRADED          string = "Garage upgrade successfully"
	ADD_CAR_TO_GARAGE_FAILED string = "Unable to add car to garage"

	GARAGE_ADD_SUCCESS      string = "Garage Added successfully"
	GARAGE_DELETE_SUCCESS   string = "Garage Deleted successfully"
	GARAGE_UPDATE_SUCCESS   string = "Garage Updated successfully"
	GARAGE_ALREADY_PRESENT  string = "Garage already present at that location"
	UPGRADE_LEVEL           string = "Upgrade your level to unlock the car"
	NOT_ENOUGH_REPAIR_PARTS string = "Not enough repair parts"
	NOT_ENOUGH_COINS        string = "Not enough coins"
	NOT_ENOUGH_CASH         string = "Not enough cash"
)

const (
	FORBIDDEN_REQUEST  string = "Forbidden Request"
	LOGIN_SUCCESS      string = "Login Successfull"
	LOGIN_FAILED       string = "Login Failed"
	EMAIL_EXISTS       string = "Email is already attached to another player"
	NOT_FOUND          string = "Record not found"
	DATA_FETCH_SUCCESS string = "Data Fetch Successfully"

	EMAIL_UPDATED_SUCCESS  string = "Email updated successfully"
	WON                    string = "You WON"
	LOSE                   string = "You LOST"
	ACCOUNT_DELETE_SUCCESS string = "Account deleted successfully"
	LOGOUT_SUCCESS         string = "Logout successfully"

	LINK_GENERATED_SUCCESS  string = "Link generated successfully"
	PASSWORD_UPDATE_SUCCESS string = "Password updated successfully"
)
