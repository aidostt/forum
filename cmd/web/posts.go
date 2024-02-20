package main

import (
	"errors"
	"forum.aidostt-buzuk/internal/data"
	"forum.aidostt-buzuk/internal/validator"
	"github.com/jackc/pgtype"
	"net/http"
)

type postCreateForm struct {
	Heading     *string   `form:"heading"`
	Description *string   `form:"description"`
	Tags        *[]string `form:"tags"`
	validator.Validator
}

func (app *application) filtersSearchHandlerGet(w http.ResponseWriter, r *http.Request) {
	d := app.newTemplateData(r)
	app.render(w, http.StatusOK, "filters.tmpl", d)
}
func (app *application) filtersSearchHandlerPost(w http.ResponseWriter, r *http.Request) {
}

// TODO: post should be filtered by time, tags and quantity of likes
func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	//d := app.newTemplateData(r)
	//var input struct {
	//	Heading string
	//	Tags    []string
	//	data.Filters
	//}

}
func (app *application) createPostHandlerGet(w http.ResponseWriter, r *http.Request) {
	d := app.newTemplateData(r)
	d.Form = postCreateForm{}
	//TODO: if user not registered/logged in or didn't pass the activation redirect to log in
	app.render(w, http.StatusOK, "create.tmpl", d)
}
func (app *application) createPostHandlerPost(w http.ResponseWriter, r *http.Request) {
	var form postCreateForm
	d := app.newTemplateData(r)
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var authorID pgtype.UUID
	err = authorID.Set("6f52ada3-0c28-45d7-9b35-e420277bc127")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	//TODO: retrieve authors id
	var post *data.Post
	if form.Heading != nil {
		post.Heading = *form.Heading
	}
	if form.Description != nil {
		post.Description = *form.Description
	}
	if form.Tags != nil {
		post.Tags = *form.Tags
	}
	form.Validator = *(validator.New())
	if data.ValidatePost(&(form.Validator), post); !form.Valid() {
		d.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", d)
		return
	}
	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	d.Post = post
	app.render(w, http.StatusOK, "view.tmpl", d)
}

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	d := app.newTemplateData(r)
	post, err := app.models.Posts.GetById(id)
	if err != nil {
		if errors.Is(err, data.ErrNoRecord) {
		}
		return
	}
	d.Post = post
	app.render(w, http.StatusOK, "view.tmpl", d)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	d := app.newTemplateData(r)
	post, err := app.models.Posts.GetById(id)
	if err != nil {
		if errors.Is(err, data.ErrNoRecord) {
		}
		return
	}

	var input postCreateForm
	err = app.decodePostForm(r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if input.Heading != nil {
		post.Heading = *input.Heading
	}
	if input.Description != nil {
		post.Description = *input.Description
	}
	if input.Tags != nil {
		post.Tags = *input.Tags
	}
	var authorID pgtype.UUID
	err = authorID.Set("6f52ada3-0c28-45d7-9b35-e420277bc127")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//TODO: retrieve authors id
	input.Validator = *(validator.New())
	if data.ValidatePost(&(input.Validator), post); !input.Valid() {
		d.Form = input
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", d)
		return
	}
	err = app.models.Posts.Update(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	d.Post = post
	app.render(w, http.StatusOK, "view.tmpl", d)
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, ErrNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.models.Posts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	//TODO: send "post has been deleted successfully" message
	app.render(w, http.StatusOK, "home.tmpl", nil)
}
