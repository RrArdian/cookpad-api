package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/julienschmidt/httprouter"
)

type Resep struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Ingredients string `json:"ingredients"`
	Image       string `json:"image"`
}

type ResepData struct {
	Results []Resep `json:"results"`
}

func ResepHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		rid, _ := item.Attr("data-id")
		id, _ := strconv.Atoi(rid)
		title := item.Find(".media__body-overflow header span").Text()
		ingredients := item.Find(".media__body-overflow .recipe__ingredients").Text()
		img := item.Find("img")
		image, _ := img.Attr("src")
		resepData.Results = append(resepData.Results, Resep{ID: id, Nama: title, Ingredients: ingredients, Image: image})
	})

	js, err := json.Marshal(resepData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type ResepDetail struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Author      string   `json:"author"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

func ResepDetailHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	doc, err := goquery.NewDocument("https://cookpad.com/id/resep/" + id)
	if err != nil {
		log.Fatal(err)
	}

	var resepDetail ResepDetail
	doc.Find(".editor").Each(func(index int, item *goquery.Selection) {
		rid, _ := strconv.Atoi(id)
		title := item.Find(".intro-container h1").Text()
		description := item.Find(".recipe-show__story p").Text()
		img := item.Find(".tofu_image img")
		image, _ := img.Attr("src")
		author := item.Find(".media__img span").Text()
		var ingredients []string
		item.Find(".ingredient").Each(func(index int, q *goquery.Selection) {
			quantity := q.Find(".ingredient__details").Text()
			ingredients = append(ingredients, quantity)
		})
		var step []string
		item.Find(".step").Each(func(index int, q *goquery.Selection) {
			st := q.Find(".step__text").Text()
			step = append(step, st)
		})
		resepDetail = ResepDetail{ID: rid, Title: title, Description: description, Image: image, Author: author, Ingredients: ingredients, Steps: step}
	})

	js, err := json.Marshal(resepDetail)
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

	route := httprouter.New()
	route.GET("/resep", ResepHandler)
	route.GET("/resep/:id", ResepDetailHandler)
	log.Fatal(http.ListenAndServe(":"+port, route))
}
