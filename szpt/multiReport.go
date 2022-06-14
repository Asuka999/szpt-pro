package szpt

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
)

type UserFile struct {
	Users []UserList
}

type UserList struct {
	Name     string
	Password string
}

func GetUsersJson() []UserList {
	file, err := os.Open("./user.json")
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(file)

	var userJson UserFile
	jsoniter.Unmarshal(body, &userJson)
	return userJson.Users

}

func MultiReport() {
	userList := GetUsersJson()
	for i, list := range userList {
		fmt.Println(i)
		Report(User{
			Account:  list.Name,
			Password: list.Password,
		})
	}

}
