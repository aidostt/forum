package main

import (
	"forum.aidostt-buzuk/internal/data"
	"forum.aidostt-buzuk/internal/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type postCreateForm struct {
	Heading     string   `form:"heading"`
	Description string   `form:"description"`
	Tags        []string `form:"tags"`
	validator.Validator
}

//TODO: post should be filtered by time, tags and quantity of likes

func (app *application) CreatePostHandlerGet(w http.ResponseWriter, r *http.Request) {
	d := app.newTemplateData(r)
	d.Form = postCreateForm{}
	//TODO: if user not registered/logged in or didn't pass the activation redirect to log in
	app.render(w, http.StatusOK, "create.tmpl", d)
}
func (app *application) CreatePostHandlerPost(w http.ResponseWriter, r *http.Request) {
	var form postCreateForm
	d := app.newTemplateData(r)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var authorID pgtype.UUID
	//TODO: retrieve authors id
	post := &data.Post{
		AuthorID:    authorID,
		Heading:     form.Heading,
		Description: form.Description,
		Tags:        form.Tags,
	}
	form.Validator = *(validator.New())
	if data.ValidatePost(&(form.Validator), post); !form.Valid() {
		d.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", d)
	}
	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.render(w, http.StatusOK, "home.tmpl", nil)
}
