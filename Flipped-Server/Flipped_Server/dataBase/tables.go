package dataBase

type RegisterTable struct {
	Pid        string `sql:"pid"`
	Username   string `sql:"username"`
	Password   string `sql:"password"`
	UserType   byte   `sql:"user_type"`
	Email      string `sql:"email"`
	Photo      []byte `sql:"photo"`
	RealName   string `sql:"realName"`
	Profession string `sql:"profession"`
	Age        byte   `sql:"age"`
	Region     string `sql:"region"`
	Hobby      string `sql:"hobby"`
}

func (register *RegisterTable) Insert(table RegisterTable) error {
	//imageData :=
	//var sql = "INSERT INTO im.userinfo VALUES(" +
	return nil
}

