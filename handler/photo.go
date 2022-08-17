package handler

import (
	"encoding/json"
	"go-redis/cache"
	"go-redis/photo"
	"io/ioutil"
	"log"
	"net/http"
)

type photoHandler struct {
	logger     *log.Logger
	redisCache cache.RedisCache
}

func NewPhotoHandler(logger *log.Logger, redisCache cache.RedisCache) *photoHandler {
	return &photoHandler{logger, redisCache}
}

func (h *photoHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.getAllPhotos(rw, r)
		return
	}
}

func (h *photoHandler) getAllPhotos(rw http.ResponseWriter, r *http.Request) {
	// get photos from cache
	val, err := h.redisCache.Get("photos")
	photos := photo.Photos{}

	// if photos not exist in cache, fetch
	if err != nil {
		res, err := http.Get("https://jsonplaceholder.typicode.com/photos")
		if err != nil {
			h.logger.Fatalln(err)
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			h.logger.Fatalln(err)
		}

		err = json.Unmarshal(body, &photos)

		if err != nil {
			h.logger.Fatalln(err)
		}

		jdata, err := json.Marshal(photos)

		if err != nil {
			h.logger.Fatalln(err)
		}

		h.redisCache.Set("photos", jdata)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(jdata)

	} else {
		rw.Header().Set("Content-Type", "application/json")
		jraw := val.(string)
		err := json.Unmarshal([]byte(jraw), &photos)
		if err != nil {
			h.logger.Fatalln(err)
		}

		jdata, err := json.Marshal(photos)

		if err != nil {
			h.logger.Fatalln(err)
		}
		rw.Write(jdata)
	}

}
