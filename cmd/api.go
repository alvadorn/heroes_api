package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alvadorn/heroes_api/pkg/models"
	"github.com/gorilla/mux"
)

type Hero struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Superpower string `json:"superpower"`
}

var heroes []Hero
var db *models.DB

func init() {
	db = models.Connect()
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.
		NewRoute().
		Methods("post").
		Path("/").
		Name("create").
		HandlerFunc(SaveHeroesHandler)

	router.Methods("delete").Path("/{id}").Name("Delete").HandlerFunc(DeleteAHero)
	router.Methods("patch").Path("/{id}").Name("Update").HandlerFunc(UpdateHeroesHandler)
	router.Methods("get").Path("/{id}").Name("identifier").HandlerFunc(GetAHeroHandler)
	router.Methods("get").Path("/").Name("index").HandlerFunc(GetHeroesHandler)
	fmt.Println("ta rodano")
	http.ListenAndServe(":8080", router)
}

func GetHeroesHandler(write http.ResponseWriter, request *http.Request) {
	//ToDo: Verificar por que a array vazia esta retornando nulo
	write.WriteHeader(http.StatusOK)
	heroes1, _ := db.GetAllHeroes()
	fmt.Println("%v", heroes1)
	json.NewEncoder(write).Encode(heroes1)
}

func SaveHeroesHandler(write http.ResponseWriter, request *http.Request) {
	var hero Hero
	json.NewDecoder(request.Body).Decode(&hero)
	heroDynamo := &models.Hero{hero.ID, hero.Name, hero.Superpower}
	err := db.Insert(heroDynamo)

	if err != nil {
		fmt.Println(err)
	}
	write.WriteHeader(http.StatusOK)
}

func GetAHeroHandler(write http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	heroId, _ := strconv.ParseInt(params["ID"], 10, 64)
	heroDb, err := db.GetAHero(heroId)
	if err != nil {
		var heroJson Hero
		heroJson.ID = heroDb.ID
		write.WriteHeader(http.StatusOK)
		json.NewEncoder(write).Encode(heroJson)
	} else {
		write.WriteHeader(http.StatusNotFound)
	}
}

func UpdateHeroesHandler(write http.ResponseWriter, request *http.Request) {
	foundHeroIndex := FoundHeroIndex(request)
	if foundHeroIndex != -1 {
		var hero Hero
		err := json.NewDecoder(request.Body).Decode(&hero)

		if err != nil {
			fmt.Println(err)
		}

		foundHero := &heroes[foundHeroIndex]
		foundHero.Superpower = hero.Superpower
		foundHero.Name = hero.Name
		write.WriteHeader(http.StatusOK)
		json.NewEncoder(write).Encode(heroes[foundHeroIndex])
	} else {
		write.WriteHeader(http.StatusNotFound)
	}
}

func DeleteAHero(write http.ResponseWriter, request *http.Request) {
	foundHeroIndex := FoundHeroIndex(request)
	if foundHeroIndex != -1 {
		write.WriteHeader(http.StatusOK)
		heroes = append(heroes[:foundHeroIndex], heroes[foundHeroIndex+1:len(heroes)]...)
	} else {
		write.WriteHeader(http.StatusNotFound)
	}
}

func FoundHeroIndex(request *http.Request) int {
	params := mux.Vars(request)
	heroId, _ := strconv.ParseInt(params["id"], 10, 64)
	foundHeroIndex := -1

	for i, hero := range heroes {
		if hero.ID == int(heroId) {
			foundHeroIndex = i
		}
	}
	return foundHeroIndex
}
