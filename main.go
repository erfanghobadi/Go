package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type user struct {
	ID                                        int
	Name, Family                              string
	Age                                       int
	IsMale                                    bool
	NationalityId, PhoneNumb, Email, Password string
}

// Declaration variables ===>

var person = user{}
var Users = []user{}
var NumberOfUsers int
var CounterID int = 0
var ResultOfWriting bool = true
var WantToEdit, Loggin string
var UserStruct user
var JsonData []byte
var UserMap map[string]interface{}
var SaveUserMapInArray []map[string]interface{}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func main() {

	var wantToEdit string

	fmt.Print("Number of Users : ")
	fmt.Scanln(&NumberOfUsers)
	Users := InputUserData(NumberOfUsers)

	JsonData := ConverToJson(Users)
	fmt.Println(JsonData)
	ResultOfWriting = SaveJsonFile(JsonData)
	fmt.Println("\n\n", "***** Result of Writing:", ResultOfWriting)

	fmt.Print("\n", "Are You Want To Edit Users Data ? (if you want type 'yes' , if you're not press 'enter') :")
	fmt.Scanln(&wantToEdit)
	if wantToEdit == "yes" {
		Edit(Users)
	}

	fmt.Print("\n", "Are you Want to Loggin ? (if you wat type 'yes' , if you're not press 'enter') :")
	fmt.Scanln(&Loggin)
	if Loggin == "yes" {
		Log(SaveUserMapInArray)
	}
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func InputUserData(NumberOfUsers int) []user {

	for i := 1; i <= NumberOfUsers; i++ {
		CounterID++
		person.ID = CounterID

		fmt.Print("\n", "Enter Name:")
		fmt.Scanln(&person.Name)

		fmt.Print("Enter Family:")
		fmt.Scanln(&person.Family)

		fmt.Print("Enter Age:")
		fmt.Scanln(&person.Age)

		fmt.Print("Enter NationalityId:")
		fmt.Scanln(&person.NationalityId)

		fmt.Print("Enter PhoneNumber:")
		fmt.Scanln(&person.PhoneNumb)

		fmt.Print("Enter IsMale (Please Enter true or false):")
		fmt.Scanln(&person.IsMale)

		fmt.Print("Enter Email:")
		fmt.Scanln(&person.Email)

		fmt.Print("Enter Password:")
		fmt.Scanln(&person.Password)

		Users = append(Users, person)
		if i < NumberOfUsers {
			fmt.Println("\n****Enter Next User's Data****")
		}
	}
	return Users
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func ConverToJson([]user) string {
	JsonData, _ = json.Marshal(Users)
	return string(JsonData)
}

func SaveJsonFile(JsonData string) bool {

	file, _ := os.Create("Users.json")

	defer file.Close()

	_, err := file.Write([]byte(JsonData))
	if err != nil {
		fmt.Println("Error:", err)
		ResultOfWriting = false
	}

	return ResultOfWriting
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func Edit(Users []user) {
	var fieldName string
	var newValue string
	var another string
	var inputID int

	for _, UserStruct = range Users {
		UserMap = StructToMap(UserStruct)
		SaveUserMapInArray = append(SaveUserMapInArray, UserMap)
	}
	fmt.Print("Enter User ID to Change:")
	fmt.Scanln(&inputID)
	for index, UserMap := range SaveUserMapInArray {
		indexMap := index + 1
		for indexMap == inputID {
			fmt.Print("Enter Field Name Correctly (Like : Family , ...) :")
			fmt.Scanln(&fieldName)
			for key, _ := range UserMap {
				if key == fieldName {
					fmt.Print("Enter New Value:")
					fmt.Scanln(&newValue)
					UserMap[key] = newValue
				}
			}
			fmt.Print("Are Want to Change Another Field ? (if you want type 'yes' if you're not type 'no') :")
			fmt.Scanln(&another)
			if another == "yes" {
				continue
			} else if another == "no" {
				break
			}
		}
		fmt.Print("If You Want to Edit Another User Enter it's ID :")
		fmt.Scanln(&inputID)
	}
	JsonData, _ = json.Marshal(SaveUserMapInArray)
	file, _ := os.Create("Users.json")

	defer file.Close()

	_, err := file.Write([]byte(JsonData))
	if err != nil {
		fmt.Println("Error:", err)
		ResultOfWriting = false
	}
	fmt.Println("***** Result of Editing & Writing:", ResultOfWriting)
}

func StructToMap(UserStruct user) map[string]interface{} {
	UserMap := map[string]interface{}{
		"ID":            UserStruct.ID,
		"Name":          UserStruct.Name,
		"Family":        UserStruct.Family,
		"Age":           UserStruct.Age,
		"IsMale":        UserStruct.IsMale,
		"NationalityId": UserStruct.NationalityId,
		"PhoneNumb":     UserStruct.PhoneNumb,
		"Email":         UserStruct.Email,
		"Password":      UserStruct.Password,
	}
	return UserMap
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func Log(SaveUserMapInArray []map[string]interface{}) {
	var emailAddress, passwd string
	var Authentication bool

	for _, UserStruct = range Users {
		UserMap = StructToMap(UserStruct)
		SaveUserMapInArray = append(SaveUserMapInArray, UserMap)
	}

	fmt.Print("Enter Your Email:")
	fmt.Scanln(&emailAddress)

	fmt.Print("Enter Your Password:")
	fmt.Scanln(&passwd)

OuterLoop: // Label for Outerloop to break
	for _, UserMap := range SaveUserMapInArray {
		i := 0
		for _, value := range UserMap {
			if value == emailAddress || value == passwd {
				i++
				if i == 2 {
					Authentication = true
					break OuterLoop
				}
			}
		}
	}
	fmt.Println("Authentication : ", Authentication)
}
