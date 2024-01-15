package main

import (
	"forum.aidostt-buzuk/internal/data"
	"net/http"
)

func (app *application) testModel() error {
	//user := data.User{
	//	Name:      "buzuk",
	//	Nickname:  "buzuk",
	//	Email:     "buzuk@gmail.com",
	//	Password:  []byte("buz"),
	//	Activated: true,
	//}
	user, err := app.models.Users.GetByNickname("buzuk")
	if err != nil {
		return err
	}
	user.Email = "ExampleEmail@yahoo.com"
	err = app.models.Users.Update(user)
	app.infoLog.Println(user)
	return nil
}

func (app *application) CreateUserHandlerPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `form:"name"`
		Nickname string `form:"nickname"`
		Email    string `form:"email"`
		password string `form:"password"`
	}
	err := app.decodePostForm(r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user := &data.User{
		Name:      input.Name,
		Nickname:  input.Nickname,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//TODO: validate input data
	err = app.models.Users.Insert(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.render(w, http.StatusOK, "home.tmpl", nil)
}

func (app *application) CreateUserHandlerGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "register.tmpl", nil)
}

func (app *application) AxtivateUserHandler(w http.ResponseWriter, r *http.Request) {

}
