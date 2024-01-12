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
	user, err := app.models.Users.GetByNickName("buzuk")
	if err != nil {
		return err
	}
	user.Email = "ExampleEmail@yahoo.com"
	err = app.models.Users.Update(user)
	app.infoLog.Println(user)
	return nil
}

func (app *application) CreateUserHandlerPost(w http.ResponseWriter, r *http.Request) {
	var user data.User
	//decode form
	err := app.decodePostForm(r, &user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//TODO: validate input data
	err = app.models.Users.Insert(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.render(w, http.StatusOK, "home.tmpl", nil)
}

func (app *application) CreateUserHandlerGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "register.tmpl", nil)
}
