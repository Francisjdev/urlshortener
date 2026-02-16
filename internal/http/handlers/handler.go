package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/francisjdev/urlshortener/internal/model"
	"github.com/francisjdev/urlshortener/internal/service"
)

type URLHandler struct {
	Service *service.URLService
}

func (handler URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {

	type createURLRequest struct {
		URL string `json:"url"`
	}

	type createURLResponse struct {
		Code     string `json:"code"`
		ShortURL string `json:"short_url"`
	}
	type errorParam struct {
		Ret_error   string `json:"error"`
		Err_message string `json:"message"`
	}

	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	parameters := createURLRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&parameters)

	if err != nil {
		response := errorParam{
			Ret_error:   "Bad Request",
			Err_message: "Missing or invalid required parameter 'URL'"}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)

		return
	}
	_, error := url.ParseRequestURI(parameters.URL)
	if error != nil {
		response := errorParam{
			Ret_error:   "Bad Request",
			Err_message: "Missing or invalid required parameter 'URL'"}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}
	url := model.URL{}
	url.LongURL = parameters.URL

	err = handler.Service.CreateShortURL(r.Context(), &url)
	if err != nil {

		response := errorParam{
			Ret_error:   "Cant Create Code",
			Err_message: "Short Url could not be created"}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return

	}
	successResp := createURLResponse{
		Code:     url.Code,
		ShortURL: "https://urlshortener-zdcd.onrender.com/" + url.Code,
	}

	data, err := json.Marshal(successResp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

}

func (h URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:] // strip the leading "/"
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	urlValue, err := h.Service.GetCode(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_ = h.Service.IncrementHitCount(r.Context(), code)

	http.Redirect(w, r, urlValue.LongURL, http.StatusFound)
}
