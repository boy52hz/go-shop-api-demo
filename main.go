package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/boy52hz/go-demo-shop-api/db"
	"github.com/boy52hz/go-demo-shop-api/middlewares"
	itemsService "github.com/boy52hz/go-demo-shop-api/services"
)

func itemHandleFunc(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "items/")
	id, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := itemsService.FindOne(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serializedItem, err := json.Marshal(item)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(serializedItem)
	case http.MethodPut:
		updatedItem := &itemsService.Item{}
		err := json.NewDecoder(r.Body).Decode(updatedItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		affectedRows, err := itemsService.Update(id, updatedItem)
		if affectedRows == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		serializedUpdatedItem, err := json.Marshal(updatedItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(serializedUpdatedItem)
	case http.MethodDelete:
		itemsService.Delete(id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func itemsHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := itemsService.FindAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		serializedItems, err := json.Marshal(items)
		if err != nil {
			log.Println(err)
		}
		w.Write(serializedItems)
	case http.MethodPost:
		newItem := &itemsService.Item{}
		err := json.NewDecoder(r.Body).Decode(newItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = itemsService.Create(newItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		serializedNewItem, err := json.Marshal(newItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(serializedNewItem)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func applyMiddlewares(handler http.Handler) http.Handler {
	return middlewares.ContentTypeMiddleware((middlewares.EnableCorsMiddleware(handler)))
}

func main() {
	db.InitializeDatabase()

	itemHandler := http.HandlerFunc(itemHandleFunc)
	itemListHandler := http.HandlerFunc(itemsHandleFunc)
	http.Handle("/items/", applyMiddlewares(itemHandler))
	http.Handle("/items", applyMiddlewares(itemListHandler))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
