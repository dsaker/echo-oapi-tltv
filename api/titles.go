package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
)

// FindTitles implements all the handlers in the ServerInterface
func (p *Server) FindTitles(ctx echo.Context, params FindTitlesParams) error {

	titles, err := p.queries.ListTitles(
		ctx.Request().Context(),
		db.ListTitlesParams{
			Similarity: *params.Similarity,
			Limit:      *params.Limit,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, titles)
}

func (p *Server) AddTitle(ctx echo.Context) error {
	// We expect a NewTitle object in the request body.
	var newTitle NewTitle
	err := ctx.Bind(&newTitle)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock()
	defer p.Unlock()

	title, err := p.queries.InsertTitle(
		ctx.Request().Context(),
		db.InsertTitleParams{
			Title:      newTitle.Title,
			NumSubs:    newTitle.NumSubs,
			LanguageID: newTitle.LanguageId,
		})

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, title)
}

func (p *Server) FindTitleByID(ctx echo.Context, id int64) error {

	title, err := p.queries.SelectTitleById(ctx.Request().Context(), id)

	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, title)
}

func (p *Server) DeleteTitle(ctx echo.Context, id int64) error {

	err := p.queries.DeleteTitleById(ctx.Request().Context(), id)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}
