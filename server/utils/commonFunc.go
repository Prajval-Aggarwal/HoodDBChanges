package utils

import (
	"bytes"
	"errors"
	"fmt"
	"main/server/db"
	"main/server/model"
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

func AlreadyAtMax(val int) bool {
	return val == 5

}

// func SetPlayerCarDefaults(playerId string, carId string) error {

// 	//set default car upgrades
// 	// playerCarUpgrades := model.PlayerCarUpgrades{
// 	// 	PlayerId:     playerId,
// 	// 	CarId:        carId,
// 	// 	Engine:       0,
// 	// 	Turbo:        0,
// 	// 	Intake:       0,
// 	// 	Nitrous:      1,
// 	// 	Body:         0,
// 	// 	Tires:        0,
// 	// 	Transmission: 0,
// 	// }
// 	// err := db.CreateRecord(playerCarUpgrades)
// 	// if err != nil {
// 	// 	//	response.ShowResponse(err.Error(), HTTP_INTERNAL_SERVER_ERROR, FAILURE, nil, ctx)
// 	// 	return err
// 	// }

// 	//set default car customizations
// 	var carDefualtCutomizations []model.DefualtCustomisation
// 	query := "SELECT * FROM default_customizations WHERE car_id=?"
// 	err := db.QueryExecutor(query, &carDefualtCutomizations, carId)
// 	if err != nil {
// 		return err
// 	}
// 	for _, customise := range carDefualtCutomizations {
// 		playerCarCustomizations := model.PlayerCarCustomization{
// 			PlayerId:      playerId,
// 			CarId:         carId,
// 			Part:          customise.Part,
// 			ColorCategory: customise.ColorCategory,
// 			ColorType:     customise.ColorType,
// 			ColorName:     customise.ColorName,
// 			Value:         customise.Value,
// 		}
// 		err = db.CreateRecord(playerCarCustomizations)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	//set default car stats
// 	var currentCarStat model.CarStats
// 	err = db.FindById(&currentCarStat, carId, "car_id")
// 	if err != nil {

// 		return err
// 	}
// 	playerCarStats := model.PlayerCarsStats{
// 		PlayerId:    playerId,
// 		CarId:       carId,
// 		Power:       currentCarStat.Power,
// 		Grip:        currentCarStat.Grip,
// 		ShiftTime:   currentCarStat.ShiftTime,
// 		Weight:      currentCarStat.Weight,
// 		OVR:         currentCarStat.OVR,
// 		Durability:  currentCarStat.Durability,
// 		NitrousTime: float64(currentCarStat.NitrousTime),
// 	}
// 	err = db.CreateRecord(&playerCarStats)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func DeleteCarDetails(tableName string, playerId string, carId string) error {
	query := "DELETE FROM " + tableName + " WHERE car_id =? AND player_id =?"
	err := db.RawExecutor(query, carId, playerId)
	if err != nil {
		return err
	}
	return nil
}

func IsCarBought(playerId string, carId string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM owned_cars WHERE player_id =? AND car_id=?)"
	err := db.QueryExecutor(query, &exists, playerId, carId)
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

// func UpgradeData(playerId string, carId string) (*model.Player, *model.PlayerCarsStats, *model.PlayerCarUpgrades, string, int64, *model.RatingMulti, error) {
// 	var playerDetails model.Player
// 	var playerCarStats model.PlayerCarsStats
// 	var carClassDetails string

// 	var playerCarUpgrades model.PlayerCarUpgrades
// 	var maxUpgradeLevel int64
// 	var classRating model.RatingMulti
// 	//check if the car is owned or not
// 	var exists bool
// 	query := "SELECT EXISTS(SELECT * FROM owned_cars WHERE player_id =? AND car_id=?)"
// 	err := db.QueryExecutor(query, &exists, playerId, carId)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}
// 	if !exists {
// 		return nil, nil, nil, "", 0, nil, errors.New(NOT_FOUND)
// 	}
// 	err = db.FindById(&playerDetails, playerId, PLAYER_ID)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	query = "SELECT * FROM player_cars_stats WHERE player_id=? AND car_id=?"
// 	err = db.QueryExecutor(query, &playerCarStats, playerId, carId)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	query = "SELECT class FROM cars WHERE car_id=?"
// 	err = db.QueryExecutor(query, &carClassDetails, carId)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	query = "SELECT * FROM rating_multis WHERE class=?"
// 	err = db.QueryExecutor(query, &classRating, carClassDetails)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	query = "SELECT * FROM player_car_upgrades WHERE player_id=? AND car_id=?"
// 	err = db.QueryExecutor(query, &playerCarUpgrades, playerId, carId)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	query = "SELECT upgrade_level FROM upgrades WHERE class =? ORDER BY upgrade_level DESC LIMIT 1;"
// 	err = db.QueryExecutor(query, &maxUpgradeLevel, carClassDetails)
// 	if err != nil {
// 		return nil, nil, nil, "", 0, nil, err
// 	}

// 	return &playerDetails, &playerCarStats, &playerCarUpgrades, carClassDetails, maxUpgradeLevel, &classRating, nil
// }

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

func TimeConversion(stringFormat string) *time.Time {
	timeFormat, err := time.Parse("00:00:05.0000", stringFormat)
	if err != nil {
		fmt.Println("error in parsing the string format of time")
		return nil
	}
	return &timeFormat
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
