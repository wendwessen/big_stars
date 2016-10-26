package api

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type apiResponse struct  {
	ErrorMessage string `json:"error,omitempty"`
	ID			 string `json:"error,omitempty"`
	Result		 string	`json:"result,omitempty"`
}

type apiHandlerFunc func(Storage,http.ResponseWriter,map[string]string,*json.Encoder,*json.Decoder)

func Mount(router *mux.Router, storage Storage,auth func(http.HandlerFunc) http.HandlerFunc)  {
	if storage == nil {
		panic("storage is nil")
	}

	if auth == nil {
		auth = func(f http.HandlerFunc) http.HandlerFunc { return f }
	}

	collectionHandlers := map[string]apiHandlerFunc{
		"GET": getAll,
	}

	resourceHandlers := map[string]apiHandlerFunc{
		"GET": get,
	}

	router.HandleFunc("/{collection}",auth(chooseAndInitialize(collectionHandlers,storage)))
	router.HandleFunc("/collection/{id}",auth(chooseAndInitialize(resourceHandlers,storage)))

}

func chooseAndInitialize(handlersByMethod map[string]apiHandlerFunc,storage Storage) http.HandlerFunc  {
	return func(resp http.ResponseWriter, req *http.Request) {
		handler,ok := handlersByMethod[req.Method]
		if !ok {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
		}
		vars := mux.Vars(req)
		enc := json.NewEncoder(resp)
		dec := json.NewDecoder(req.Body)
		handler(storage,resp,vars,enc,dec)
	}

}

func getAll(storage Storage,resp http.ResponseWriter,vars map[string]string, enc *json.Encoder, dec *json.Decoder)  {
	resources,stoResp := storage.GetAll(vars["collection"])
	resp.WriteHeader(stoResp.StatusCode())

	err := enc.Encode(resources)
	if err != nil {
		log.Println(err)
	}
}

func get(storage Storage,resp http.ResponseWriter,vars map[string]string, enc *json.Encoder, dec *json.Decoder)   {
	resources, stoResp := storage.Get(vars["collection"],vars["id"])
	resp.WriteHeader(stoResp.StatusCode())
	err := enc.Encode(resources)
	if err != nil {
		log.Println(err)
	}
}





