package ApiHandler

import (
	"encoding/json"
	"net/http"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/bhopaleneha/tp/Converter"	
	"gorm.io/gorm"
	
	
)


type Env struct{
	Db *gorm.DB
}
func PostFilePath(env *Env) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
	var file Converter.File
	json.NewDecoder(r.Body).Decode(&file)
	csv_file, err := os.Open(file.PathOfFile)
	if err != nil {
			logrus.Fatal("unable to open csv file",err)
			return 
	}else{
		logrus.Info("File with posted path is available")
	}
	
	defer csv_file.Close()
	rec:=ReadCsvFile(csv_file)
	users:=ListValidUsers(rec)
	env.createUsersInDatabase(users)
	Converter.ConToJson(users)
}
}


func GetUserInfo(env *Env)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

	var user []Converter.User
	env.Db.Find(&user)
	json.NewEncoder(w).Encode(user)
}

}
