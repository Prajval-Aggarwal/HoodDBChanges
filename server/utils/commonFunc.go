package utils

import (
	"bytes"
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"math"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"text/template"
	"time"

	"github.com/go-mail/mail"
	"golang.org/x/crypto/bcrypt"
)

func IsPassValid(password string) error {

	if len(password) < 8 {
		return errors.New("password is too short")

	}
	hasUpperCase := false
	hasLowerCase := false
	hasNumbers := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		} else if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		} else if char >= '0' && char <= '9' {
			hasNumbers = true
		} else if char >= '!' && char <= '/' {
			hasSpecial = true
		} else if char >= ':' && char <= '@' {
			hasSpecial = true
		}
	}

	if !hasUpperCase {
		return errors.New("password do not contain upperCase Character")
	}

	if !hasLowerCase {
		return errors.New("password do not contain LowerCase Character")
	}

	if !hasNumbers {
		return errors.New("password do not contain any numbers")
	}

	if !hasSpecial {
		return errors.New("password do not contain any special character")
	}
	return nil
}

func HashPassword(password string) (*string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bs)
	return &hashedPassword, nil
}

func SetCarData(carId string, playerId string) error {
	var carDefaults model.DefaultCustomisation
	query := "SELECT * FROM default_customisations WHERE car_id=? "
	err := db.QueryExecutor(query, &carDefaults, carId)
	if err != nil {
		return err
	}

	playerCarCustomisations := &model.PlayerCarCustomisation{
		PlayerId:          playerId,
		CarId:             carId,
		Power:             carDefaults.Power,
		CarLevel:          1,
		Grip:              carDefaults.Grip,
		ShiftTime:         carDefaults.ShiftTime,
		Weight:            carDefaults.Weight,
		OVR:               carDefaults.OVR,
		Durability:        carDefaults.Durability,
		NitrousTime:       carDefaults.NitrousTime,
		ColorCategory:     carDefaults.ColorCategory,
		ColorType:         carDefaults.ColorType,
		ColorName:         carDefaults.ColorName,
		WheelCategory:     carDefaults.WheelCategory,
		WheelColorName:    carDefaults.WheelColorName,
		InteriorColorName: carDefaults.InteriorColorName,
		LPValue:           carDefaults.LPValue,
	}
	err = db.CreateRecord(&playerCarCustomisations)
	if err != nil {
		return err
	}

	newCarRecord := model.OwnedCars{
		PlayerId: playerId,
		CustId:   playerCarCustomisations.CustId,
		Selected: true,
	}

	err = db.CreateRecord(&newCarRecord)
	if err != nil {

		return err
	}
	return nil
}

func IsCarBought(playerId string, custId string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM owned_cars WHERE player_id =? AND cust_id=?)"
	err := db.QueryExecutor(query, &exists, playerId, custId)
	if err != nil {
		return false
	}
	if !exists {
		return false
	}

	return true
}

func SendEmaillService(adminDetails model.Admin, link string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "hoodRacing@gmail.com")
	m.SetHeader("Subject", "Reset Password!")

	var body bytes.Buffer
	tmp, err := template.ParseFiles("server/emailTemplate/forgot-password.html")
	if err != nil {
		fmt.Println("sajdsajdvsja", err.Error())
	}

	data := struct {
		Username string
		Link     string
	}{
		Username: adminDetails.Username,
		Link:     link,
	}

	tmp.Execute(&body, data)

	m.SetBody("text/html", body.String())
	if err != nil {
		fmt.Println("basdbbkjsad", err)
	}
	m.SetHeader("To", adminDetails.Email)
	d := mail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func IsExisting(query string, values ...interface{}) bool {
	var exists bool
	err := db.QueryExecutor(query, &exists, values...)
	if err != nil {
		return false
	}
	return exists
}

func TimeConversion(stringFormat string) (*time.Time, error) {
	timeFormat, err := time.Parse("00:00:05.0000", stringFormat)
	if err != nil {
		fmt.Println("")
		return nil, errors.New("error in parsing the string format of time")
	}
	return &timeFormat, nil
}

func TableIsEmpty(tablename string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM " + tablename + ");"
	db.QueryExecutor(query, &exists)

	return exists
}
func IsEmail(e string) bool {
	//e = strings.ToLower(e)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func GenerateRandomTime(n int, min, max float64) []string {

	var slice []string
	slice1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Generate the first random number as a starting point
	prev := min + rand.Float64()*(max-min)

	if prev < 10 {
		slice = append(slice, "00:00:0"+fmt.Sprintf("%.4f", prev))
	} else {
		slice = append(slice, "00:00:"+fmt.Sprintf("%.4f", prev))
	}
	fmt.Printf("Random Numbers:\n%.4f\n", prev)

	for i := 0; i < n-1; i++ {
		// Generate a new random number that is less than the previous one
		rangeSpan := prev - min
		newNum := min + rand.Float64()*rangeSpan
		newNum = newNum - (float64(slice1[rand.Intn(len(slice1))]) / 10.0)

		if newNum < 10 {
			slice = append(slice, "00:00:0"+fmt.Sprintf("%.4f", newNum))

		} else {

			slice = append(slice, "00:00:"+fmt.Sprintf("%.4f", newNum))
		}
		prev = newNum
	}
	fmt.Println("Before sorting time slice is:", slice)
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i] > slice[j]
	})
	fmt.Println("After sorting time slice is:", slice)

	return slice

}

func GetPlayerDetails(playerId string) (*model.Player, error) {
	var playerDetails model.Player
	query := "SELECT * FROM players WHERE player_id=?"
	err := db.QueryExecutor(query, &playerDetails, playerId)
	if err != nil {
		return nil, err
	}
	return &playerDetails, nil
}

func CustomiseMapping(id string, tableName string) (*response.Customization, error) {
	var carCustomise struct {
		ColorCategory     string `json:"colorCategory,omitempty"`
		ColorType         string `json:"colorType,omitempty"`
		ColorName         string `json:"colorName,omitempty"`
		WheelCategory     string `json:"wheelCategory,omitempty"`
		WheelColorName    string `json:"wheelColorName,omitempty"`
		InteriorColorName string `json:"interiorColorName,omitempty"`
		LPValue           string `json:"lp_value,omitempty"`
	}

	var query string
	if tableName == "default_customisations" {
		query = "SELECT color_category,color_type,color_name,wheel_category,wheel_color_name,interior_color_name,lp_value FROM " + tableName + " WHERE car_id=?"
	} else if tableName == "player_car_customisations" {
		query = "SELECT color_category,color_type,color_name,wheel_category,wheel_color_name,interior_color_name,lp_value FROM " + tableName + " WHERE cust_id=?"
	}

	err := db.QueryExecutor(query, &carCustomise, id)
	if err != nil {
		return nil, err
	}
	var respCustomise response.Customization

	respCustomise.LPValue = carCustomise.LPValue
	respCustomise.ColorCategory = carCustomise.ColorCategory
	respCustomise.WheelCategory = carCustomise.WheelCategory
	switch carCustomise.ColorType {
	case "default":
		respCustomise.ColorType = 1
	case "fluorescent":
		respCustomise.ColorType = 2
	case "pastel":
		respCustomise.ColorType = 3
	case "gun_Metal":
		respCustomise.ColorType = 4
	case "satin":
		respCustomise.ColorType = 5
	case "metal":
		respCustomise.ColorType = 6
	case "military":
		respCustomise.ColorType = 7
	}

	switch carCustomise.ColorName {
	case "red":
		respCustomise.ColorName = 1
	case "green":
		respCustomise.ColorName = 2
	case "pink":
		respCustomise.ColorName = 3
	case "yellow":
		respCustomise.ColorName = 4
	case "blue":
		respCustomise.ColorName = 5
	}

	if carCustomise.ColorType == "military" {
		switch carCustomise.ColorName {
		case "basic":
			respCustomise.ColorName = 1
		case "black":
			respCustomise.ColorName = 2
		case "desert":
			respCustomise.ColorName = 3
		case "tram":
			respCustomise.ColorName = 4
		case "ucp":
			respCustomise.ColorName = 5
		}
	}

	switch carCustomise.WheelColorName {
	case "black":
		respCustomise.WheelColorName = 1
	case "blue":
		respCustomise.WheelColorName = 2
	case "green":
		respCustomise.WheelColorName = 3
	case "pink":
		respCustomise.WheelColorName = 4
	case "red":
		respCustomise.WheelColorName = 5
	case "yellow":
		respCustomise.WheelColorName = 6
	}

	switch carCustomise.InteriorColorName {
	case "white":
		respCustomise.InteriorColorName = 1
	case "pink":
		respCustomise.InteriorColorName = 2
	case "green":
		respCustomise.InteriorColorName = 3
	case "red":
		respCustomise.InteriorColorName = 4
	case "blue":
		respCustomise.InteriorColorName = 5
	case "yellow":
		respCustomise.InteriorColorName = 6
	}

	return &respCustomise, nil
}

type carDetails struct {
	CustId string
	OVR    float64
}

func GiveArenaToAi(arenaId string, arenaLevel int64) error {
	var AIId string
	query := "SELECT player_id FROM players WHERE role='ai' order by RANDOM() LIMIT 1;"
	err := db.QueryExecutor(query, &AIId)
	if err != nil {
		return err
	}

	aiOwnedArena := model.PlayerRaceStats{
		PlayerId: AIId,
		ArenaId:  &arenaId,
		WinTime:  time.Now(),
		ArenaWon: true,
	}

	var randomTimeSlice []string
	var carSlice []carDetails

	if arenaLevel == int64(EASY) {
		aiOwnedArena.WinStreak = EASY_ARENA_SERIES
		randomTimeSlice = GenerateRandomTime(int(EASY_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, arenaId, 1, 2, int(EASY_ARENA_SLOT))
		if err != nil {
			return err
		}
	} else if arenaLevel == int64(MEDIUM) {
		aiOwnedArena.WinStreak = MEDIUM_ARENA_SERIES
		randomTimeSlice = GenerateRandomTime(int(MEDIUM_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, arenaId, 2, 4, int(MEDIUM_ARENA_SLOT))
		if err != nil {
			return err
		}

	} else if arenaLevel == int64(HARD) {
		aiOwnedArena.WinStreak = HARD_ARENA_SERIES
		randomTimeSlice = GenerateRandomTime(int(HARD_ARENA_SLOT), 22.0, 25.0)
		carSlice, err = GiveRandomCar(aiOwnedArena.PlayerId, arenaId, 4, 5, int(HARD_ARENA_SLOT))
		if err != nil {
			return err
		}
	}

	err = db.CreateRecord(&aiOwnedArena)
	if err != nil {
		return err
	}

	for i, val := range randomTimeSlice {
		newRecord := model.ArenaRaceRecord{
			PlayerId: aiOwnedArena.PlayerId,
			ArenaId:  arenaId,
			TimeWin:  fmt.Sprintf("%v", val),
			CustId:   carSlice[i].CustId,
			Result:   "win",
		}

		err := db.CreateRecord(&newRecord)
		if err != nil {
			return err
		}
	}
	return nil
}
func GiveRandomCar(playerId string, arenaId string, min int64, max int64, slots int) ([]carDetails, error) {

	var carSlice []carDetails
	for i := 0; i < slots; i++ {
		var carId string
		query := ` SELECT car_id FROM cars
					WHERE class >= ? AND class <= ?
					ORDER BY RANDOM() LIMIT 1 ; `
		err := db.QueryExecutor(query, &carId, min, max)
		if err != nil {
			return nil, errors.New("error in selecting the random car from the db for ai")
		}
		var carDefaults model.DefaultCustomisation
		query = "SELECT * FROM default_customisations WHERE car_id=? "
		err = db.QueryExecutor(query, &carDefaults, carId)
		if err != nil {

			return nil, err
		}

		playerCarCustomisations := model.PlayerCarCustomisation{
			PlayerId:          playerId,
			CarId:             carId,
			CarLevel:          1,
			Power:             carDefaults.Power,
			Grip:              carDefaults.Grip,
			ShiftTime:         carDefaults.ShiftTime,
			Weight:            carDefaults.Weight,
			OVR:               carDefaults.OVR,
			Durability:        carDefaults.Durability,
			NitrousTime:       carDefaults.NitrousTime,
			ColorCategory:     carDefaults.ColorCategory,
			ColorType:         carDefaults.ColorType,
			ColorName:         carDefaults.ColorName,
			WheelCategory:     carDefaults.WheelCategory,
			WheelColorName:    carDefaults.WheelColorName,
			InteriorColorName: carDefaults.InteriorColorName,
			LPValue:           carDefaults.LPValue,
		}

		err = db.CreateRecord(&playerCarCustomisations)
		if err != nil {
			return nil, err
		}

		newCarRecord := model.OwnedCars{
			PlayerId: playerId,
			CustId:   playerCarCustomisations.CustId,
			Selected: true,
		}

		err = db.CreateRecord(&newCarRecord)
		if err != nil {
			return nil, err
		}
		//Get the customisation id

		carSlice = append(carSlice, carDetails{
			CustId: playerCarCustomisations.CustId,
			OVR:    carDefaults.OVR,
		})
	}

	fmt.Println("Car slice before sorting is:", carSlice)

	sort.SliceStable(carSlice, func(i, j int) bool {
		return carSlice[i].OVR > carSlice[j].OVR
	})

	fmt.Println("Car slice is after sorting:", carSlice)

	return carSlice, nil
}
