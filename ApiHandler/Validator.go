package ApiHandler

import (
	"encoding/csv"
	"fmt"
	"github.com/bhopaleneha/tp/Converter"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func ReadCsvFile(csv_file *os.File) [][]string {
	r := csv.NewReader(csv_file)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	return records
}
func (env *Env) createUsersInDatabase(users []Converter.User) {
	for _, user := range users {
		if env.Db.Model(&user).Where("id=?", user.Id).Updates(&user).RowsAffected == 0 {
			env.Db.Create(&user)
		}
	}

}

var MapUserId = make(map[string]struct{})

type StandardLogger struct {
	logging *logrus.Logger
}

func ListValidUsers(records [][]string) []Converter.User {

	var ListOfValidUser []Converter.User
	logger := logrus.New()
	baselogger := &StandardLogger{logging: logger}

	for _, line := range records {
		Active, _ := strconv.ParseBool(line[4])
		user := Converter.User{
			Id:       line[0],
			Name:     line[1],
			Email:    line[2],
			Phone:    line[3],
			IsActive: Active,
		}

		isValidUser := baselogger.Validate(&user)
		if isValidUser {
			ListOfValidUser = append(ListOfValidUser, user)
			MapUserId[user.Id] = struct{}{}
		}

	}
	return ListOfValidUser
}

func (l *StandardLogger) Validate(user *Converter.User) bool {
	Validations := true
	if user.Id == "" {
		id := uuid.New()
		user.Id = id.String()
	}

	if user.Name == "" {
		l.logging.Error(" name as blank field for id", " ", user.Id)
		Validations = false
	}
	if user.Email == "" {
		l.logging.Error("email as a blank field for id", " ", user.Id)
		Validations = false
	}
	if len(user.Phone) != 10 {
		l.logging.Error("Phone Number is wrong for id ", " ", user.Id)
		Validations = false
	}
	_, isUserAlreadyExist := MapUserId[user.Id]
	if isUserAlreadyExist {
		l.logging.Error("existing user with id ", user.Id)
		Validations = false

	}
	return Validations
}
