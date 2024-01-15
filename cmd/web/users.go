package main

import (
	"errors"
	"forum.aidostt-buzuk/internal/data"
	"forum.aidostt-buzuk/internal/validator"
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
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponce(w, r, v.Errors)
	}
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists.")
			app.failedValidationResponce(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateNickname):
			v.AddError("email", "a user with this nickname already exists.")
			app.failedValidationResponce(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	app.render(w, http.StatusOK, "home.tmpl", nil)
}

func (app *application) CreateUserHandlerGet(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "register.tmpl", nil)
}

func (app *application) AxtivateUserHandler(w http.ResponseWriter, r *http.Request) {

}
