package route

import (
	"fmt"
	"net/http"
	"os"
	controller "pokemon/controller"
)

func InitRoute() {
	http.HandleFunc("/", controller.ErrorPage)
	http.HandleFunc("/accueil", controller.Index)
	http.HandleFunc("/add-favorite", controller.AddFavoriteHandler)
	http.HandleFunc("/favorites", controller.ShowFavoritesHandler)
	http.HandleFunc("/search", controller.SearchPokemon)
	http.HandleFunc("/pokemon/", controller.PokemonDetailHandler)
	http.HandleFunc("/filtrerType", controller.FiltrerTypeHandler)
	http.HandleFunc("/starter", controller.ServePokemonsHandlers)
	http.HandleFunc("/starters", controller.ServePokemonsHandlers)
	http.HandleFunc("/remove-favorite", controller.RemoveFavoriteHandler)

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	fmt.Println("(http://localhost:8080/) - Server started on port:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server closed")
}
