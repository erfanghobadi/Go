package main

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type user struct {
	ID            int
	Name          string
	Family        string
	Age           int
	IsMale        bool
	NationalityId string
	PhoneNumb     string
	Email         string
	Password      string
}

// Declaration variables ===>

var person = user{}
var Users = []user{}
var NumberOfUsers int
var CounterID int = 0
var Loggin string

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func main() {

	var wantToEdit string

	fmt.Print("Number of Users : ")
	fmt.Scanln(&NumberOfUsers)
	Users := InputUserData(NumberOfUsers)

	SaveToMongodb(Users)

	fmt.Print("\n", "Are You Want To Edit Users Data ? (if you want type 'yes' , if you're not press 'enter'): ")
	fmt.Scanln(&wantToEdit)
	if wantToEdit == "yes" {
		fmt.Print("Update or Delete Document ? (to update type 'up' and to delete type 'de'): ")
		var UpOrDe string
		fmt.Scanln(&UpOrDe)
		if UpOrDe=="up"{
			UpdateMongodb()
		}else if UpOrDe=="de"{
			DeleteMongo()
		}else{
			os.Exit(1)
		}
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

		/////////////////////////////////////////////			 Hash User's Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		person.Password = string(hashedPassword)
		/////////////////////////////////////////////

		Users = append(Users, person)
		if i < NumberOfUsers {
			fmt.Println("\n****Enter Next User's Data****")
		}
	}
	return Users
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func SaveToMongodb(Users []user) {

	/////////
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connecting Done !")
	/////////

	for _, UserStruct := range Users {
		UserMap := StructToMap(UserStruct)

		bsonData, err := bson.Marshal(UserMap)
		if err != nil {
			log.Fatal(err)
		}
		collection := client.Database("userDB").Collection("userscollection")
		insertResult, err := collection.InsertOne(context.TODO(), bsonData)
		fmt.Println("insertResult :", insertResult)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Save in Mongo Successful")

	err = client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("Save in Mongo Failed !")
		log.Fatal(err)
	}
}

///////////////////////***********//////////////////////////////////
///////////////////////***********//////////////////////////////////

func UpdateMongodb() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connecting Done !")

	var filters []FieldFilter
	var updates []FieldUpdate
	var set string
	for {
		var filter FieldFilter
		fmt.Println("\n", "<<< FILTER >>>")
		fmt.Print("\n", "Field Name: ")
		fmt.Scanln(&filter.Name)
		if filter.Name == "done" {
			break
		}
		fmt.Print("Field Value: ")
		reader := bufio.NewReader(os.Stdin)
		filter.Value, _ = reader.ReadString('\n')
		filters = append(filters, filter)

		var update FieldUpdate
		fmt.Println("\n", "<<< UPDATE >>>")
		fmt.Print("\n", "Field Name: ")
		fmt.Scanln(&update.Name)
		fmt.Print("Please Enter '$set' or '$unset' ( '$unset' for Delete Fields ): ")
		fmt.Scanln(&set)
		fmt.Print("New Value: ")
		reader = bufio.NewReader(os.Stdin)
		update.Value, _ = reader.ReadString('\n')
		updates = append(updates, update)
	}

	filter := bson.D{}
	for _, fieldFilter := range filters {
		filter = append(filter, bson.E{Key: fieldFilter.Name, Value: fieldFilter.Value})
	}

	collection := client.Database("userDB").Collection("userscollection")

	updateFields := bson.D{}
	for _, fieldUpdate := range updates {
		updateFields = append(updateFields, bson.E{Key: fieldUpdate.Name, Value: fieldUpdate.Value})
	}
	update := bson.D{{set, updateFields}}
	updateResult, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	if updateResult.ModifiedCount > 0 {
		fmt.Println("Successful")
	} else {
		fmt.Println("Not Changed !")
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println("Save in Mongo Failed !")
		log.Fatal(err)
	}
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

type FieldFilter struct {
	Name  string
	Value interface{}
}

type FieldUpdate struct {
	Name  string
	Value interface{}
}

func DeleteMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connecting Done !")

	var filters []FieldFilter
	for {
		var filter FieldFilter
		fmt.Println("\n", "<<< FILTER >>>")
		fmt.Print("\n", "Field Name: ")
		fmt.Scanln(&filter.Name)
		if filter.Name == "done" {
			break
		}
		fmt.Print("Field Value: ")
		reader := bufio.NewReader(os.Stdin)
		filter.Value, _ = reader.ReadString('\n')
		filters = append(filters, filter)
	}
	
	filter := bson.D{}
	for _, fieldFilter := range filters {
		filter = append(filter, bson.E{Key: fieldFilter.Name, Value: fieldFilter.Value})
	}

	collection := client.Database("userDB").Collection("userscollection")

	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	if deleteResult.DeletedCount > 0 {
		fmt.Println("Successful")
	} else {
		fmt.Println("Not Changed !")
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("Save in Mongo Failed !")
		log.Fatal(err)
	}
}
