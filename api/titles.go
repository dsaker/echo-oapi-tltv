package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	db "talkliketv.click/tltv/db/sqlc"
	"talkliketv.click/tltv/internal/oapi"
)

// FindTitles implements all the handlers in the ServerInterface
func (p *Api) FindTitles(w http.ResponseWriter, r *http.Request, params oapi.FindTitlesParams) {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	titles, err := p.queries.ListTitles(
		ctx,
		db.ListTitlesParams{
			Similarity: *params.Similarity,
			Limit:      *params.Limit,
		})

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error finding titles: %s", err))
		return
	}

	if err = json.NewEncoder(w).Encode(titles); err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error encoding titles: %s", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *Api) AddTitle(w http.ResponseWriter, r *http.Request) {
	// We expect a NewTitle object in the request body.
	var newTitle oapi.NewTitle
	if err := json.NewDecoder(r.Body).Decode(&newTitle); err != nil {
		sendApiError(w, http.StatusBadRequest, "Invalid format for NewTitle")
		return
	}

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	title, err := p.queries.InsertTitle(
		ctx, db.InsertTitleParams{
			Title:      newTitle.Title,
			NumSubs:    newTitle.NumSubs,
			LanguageID: newTitle.LanguageId,
		})

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error creating title: %s", err))
		return
	}

	// Now, we have to return the NewTitle
	if err = json.NewEncoder(w).Encode(title); err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error encoding title: %s", err))
		return
	}
}

func (p *Api) FindTitleByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	title, err := p.queries.SelectTitleById(ctx, id)

	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error selecting title by id: %s", err))
		return
	}

	if err = json.NewEncoder(w).Encode(title); err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error encoding title: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Api) DeleteTitle(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), p.config.CtxTimeout)
	defer cancel()

	w.WriteHeader(http.StatusNoContent)
	err := p.queries.DeleteTitleById(ctx, id)
	if err != nil {
		sendApiError(w, http.StatusBadRequest, fmt.Sprintf("Error selecting title by id: %s", err))
		return
	}
}
