package dataBase

type RegisterTable struct {
	Pid        string `sql:"pid"`
	Username   string `sql:"username"`
	Password   string `sql:"password"`
	UserType   byte   `sql:"user_type"`
	Email      string `sql:"email"`
	Photo      string `sql:"photo"`
	RealName   string `sql:"realName"`
	Profession string `sql:"profession"`
	Age        byte   `sql:"age"`
	Region     string `sql:"region"`
	Hobby      string `sql:"hobby"`
}



