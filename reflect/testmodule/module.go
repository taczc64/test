package module

import (
	"fmt"
)

type Module struct {
	User string
	Age  int
	Sex  string
}

func (this *Module) GetUser() string {
	fmt.Println("this is function get user")
	return this.User
}

func (this *Module) GetAge(username string) int {
	fmt.Println("this is function get age")
	fmt.Println("user name:", username)
	return this.Age
}

func (this *Module) GetSex(username string) string {
	fmt.Println("this is function get sex")
	fmt.Println("user name:", username)
	return this.Sex
}
