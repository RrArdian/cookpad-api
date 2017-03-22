package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type Resep struct {
	Nama        string `json:"nama"`
	Ingredients string `json:"ingredients"`
	Image       string `json:"image"`
}

type ResepData struct {
	Results []Resep `json:"results"`
}

func ResepHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	keyword := param.Get("q")
	page := param.Get("page")
	if page == "" {
		page = "1"
	}

	doc, err := goquery.NewDocument("https://cookpad.com/id/cari/" + keyword + "?page=" + page)
	if err != nil {
		log.Fatal(err)
	}

	var resepData ResepData
	doc.Find(".recipe").Each(func(index int, item *goquery.Selection) {
		title := item.Find(".media__body-overflow header span").Text()
		ingredients := item.Find(".media__body-overflow .recipe__ingredients").Text()
		img := item.Find("img")
		image, _ := img.Attr("src")
		resepData.Results = append(resepData.Results, Resep{Nama: title, Ingredients: ingredients, Image: image})
	})

	js, err := json.Marshal(resepData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/resep", ResepHandler)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
