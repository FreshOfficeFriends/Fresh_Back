package domain

//func init() {
//	validate = validator.New()
//}

//var validate *validator.Validate

//type User struct {
//	Id           int64     `json:"id"`
//	FirstName    string    `json:"first_name"`
//	SecondName   string    `json:"second_name"`
//	Email        string    `json:"email"`
//	Password     string    `json:"password"`
//	RegisteredAt time.Time `json:"registered_at"`
//}

// на фронте тоже валидируется
type SignUp struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Birthday   string `json:"birthday"`
	Password   string `json:"password"`
}

//func (s SignUp) Validate() error {
//	return validate.Struct(SignUp{})
//}
