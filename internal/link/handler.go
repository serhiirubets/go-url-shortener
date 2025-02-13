package link

import (
	"api/configs"
	"api/pkg/event"
	"api/pkg/middleware"
	"api/pkg/req"
	"api/pkg/res"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
	Config         *configs.Config
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.Handle("GET /{hash}", middleware.Auth(handler.GoTo(), deps.Config))
	router.Handle("GET /link", middleware.Auth(handler.GetAll(), deps.Config))
	router.Handle("POST /link", middleware.Auth(handler.Create(), deps.Config))
	router.Handle("PUT /link/{id}", middleware.Auth(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.Auth(handler.Delete(), deps.Config))

}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[CreateLinkRequest](&w, r)

		if err != nil {
			return
		}

		link := NewLink(body.Url)
		for {
			existingLink, _ := handler.LinkRepository.GetByHash(link.Hash)

			if existingLink == nil {
				break
			}

			link.GenerateHash()
		}

		createdLink, errLink := handler.LinkRepository.Create(link)
		if errLink != nil {
			http.Error(w, errLink.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)

		if ok {
			fmt.Println(email)
		}
		body, err := req.HandleBody[UpdateLinkRequest](&w, r)

		if err != nil {
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, link, http.StatusOK)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseInt(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, deleteError := handler.LinkRepository.GetById(uint(id))
		if deleteError != nil {
			http.Error(w, deleteError.Error(), http.StatusBadRequest)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, nil, http.StatusOK)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()

		res.Json(w, GetAllLinksResponse{
			Links: links,
			Count: count,
		}, http.StatusOK)
	}
}
