package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	InitTemp "pokemon/temp"
	"strconv"
	"strings"
	"time"
)

var AllPokemonTypes []string

type Favorites struct {
	IDs []int `json:"ids"`
}

const FavoritesFilePath = "favorites.json"

type PokemonTypeResponse struct {
	Name    string `json:"name"`
	Pokemon []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon"`
}

type ViewData struct {
	AllGenerations interface{}
	AllTypes       []string  // Ajoutez ce champ pour stocker les types de Pokémon
	Pokemons       []Pokemon // Ajoutez ou modifiez ce champ en fonction de vos besoins
}

type ViewDatas struct {
	Pokemons []PokemonSpecies
}

type Database struct {
	Generations []Generation
}

type PokemonSpecies struct {
	Name string `json:"name"`
}

type Generations struct {
	PokemonSpecies []struct {
		Name string `json:"name"`
	} `json:"pokemon_species"`
}

type Generation struct {
	ID             int              `json:"id"`
	NameGene       string           `json:"name"`
	PokemonSpecies []PokemonSpecies `json:"pokemon_species"`
}

type GenerationsList struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Pokemon struct {
	ID                     int             `json:"id"`
	Name                   string          `json:"name"`
	Height                 int             `json:"height"`
	Weight                 int             `json:"weight"`
	Forms                  []PokemonForm   `json:"forms"`
	Species                PokemonSpecies  `json:"species"`
	Type                   []string        `json:"types"`
	Image                  string          `json:"image"`
	Abilities              []Ability       `json:"abilities"` // Nouveau champ pour les capacités
	DamageRelations        DamageRelations `json:"damage_relations"`
	LocationAreaEncounters string          `json:"location_area_encounters"`
}

type PokemonForm struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonAbility struct {
	Ability struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"ability"`
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
}

type SimplePokemon struct {
	Name  string
	Type  string
	Image string
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ApiResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Pokemon `json:"results"`
}

type DamageRelations struct {
	DoubleDamageFrom []TypeRelation `json:"double_damage_from"`
	DoubleDamageTo   []TypeRelation `json:"double_damage_to"`
	HalfDamageFrom   []TypeRelation `json:"half_damage_from"`
	HalfDamageTo     []TypeRelation `json:"half_damage_to"`
	NoDamageFrom     []TypeRelation `json:"no_damage_from"`
	NoDamageTo       []TypeRelation `json:"no_damage_to"`
}

type TypeRelation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type TypesResponse struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func ToLower(str string) string {
	return strings.ToLower(str)
}

type TestTempResult struct {
	Data        ViewData
	PokemonTest Pokemon
}

func GetRandomPokemons() ([]Pokemon, error) {
	rand.Seed(time.Now().UnixNano()) // Initialise le générateur de nombres aléatoires

	var pokemons []Pokemon

	for i := 0; i < 20; i++ {
		// Générer un ID aléatoire pour un Pokémon. Ajustez le max selon le nombre total de Pokémon disponibles.
		pokemonID := rand.Intn(898) + 1
		pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokemonID)

		// La fonction FetchPokemonDetails est modifiée pour retourner le nom, types, abilities, et l'image du Pokémon.
		id, locationAreaEncounters, height, weight, name, types, abilities, image, err := FetchPokemonDetails(pokemonURL)
		if err != nil {
			log.Printf("Failed to fetch Pokemon details for ID %d: %v", pokemonID, err)
			continue // Continue to the next iteration if an error occurs
		}

		// Crée un nouvel objet Pokémon avec les détails récupérés et l'ajoute à la liste
		pokemon := Pokemon{
			ID:                     id,
			LocationAreaEncounters: locationAreaEncounters,
			Height:                 height,
			Weight:                 weight,
			Name:                   name,
			Type:                   types,
			Abilities:              abilities, // Assurez-vous d'ajouter les capacités ici
			Image:                  image,
		}
		pokemons = append(pokemons, pokemon)
	}

	return pokemons, nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	pokemons, err := GetRandomPokemons()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des Pokémon", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Nombre de Pokémons récupérés : %d\n", len(pokemons))
	// Passez les Pokémon au template
	InitTemp.Temp.ExecuteTemplate(w, "index", pokemons)
}

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	InitTemp.Temp.ExecuteTemplate(w, "error", nil)
}

func FetchPokemonDetails(pokemonURL string) (id int, locationAreaEncounters string, height int, weight int, name string, types []string, abilities []Ability, image string, err error) {
	resp, err := http.Get(pokemonURL)
	if err != nil {
		return 0, "", 0, 0, "", nil, nil, "", err
	}
	defer resp.Body.Close()

	var detailResp struct {
		ID                     int    `json:"id"`
		Name                   string `json:"name"`
		Image                  string `json:"image"`
		Height                 int    `json:"height"`
		Weight                 int    `json:"weight"`
		LocationAreaEncounters string `json:"location_area_encounters"`
		Forms                  []struct {
			Name string `json:"name"`
		} `json:"forms"`
		Species struct {
			Name string `json:"name"`
		} `json:"species"`
		Sprites struct {
			Other struct {
				OfficialArtwork struct {
					FrontDefault string `json:"front_default"`
				} `json:"official-artwork"`
			} `json:"other"`
		} `json:"sprites"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
		Abilities []struct {
			Ability struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"ability"`
		} `json:"abilities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
		return 0, "", 0, 0, "", nil, nil, "", err
	}

	id = detailResp.ID
	name = detailResp.Name
	image = detailResp.Sprites.Other.OfficialArtwork.FrontDefault
	height = detailResp.Height
	weight = detailResp.Weight
	locationAreaEncounters = detailResp.LocationAreaEncounters

	for _, t := range detailResp.Types {
		types = append(types, t.Type.Name)
	}

	for _, a := range detailResp.Abilities {
		abilities = append(abilities, Ability{
			Name: a.Ability.Name,
			URL:  a.Ability.URL,
		})
	}

	return id, locationAreaEncounters, height, weight, name, types, abilities, image, nil
}

func SearchPokemon(w http.ResponseWriter, r *http.Request) {
	// Extrait le terme de recherche de la requête
	searchQuery := r.URL.Query().Get("query")
	searchQuery = ToLower(searchQuery)
	if searchQuery == "" {
		http.Error(w, "Vous devez fournir un terme de recherche", http.StatusBadRequest)
		return
	}

	// Utilisez searchQuery pour faire une requête à l'API et obtenir des données sur le Pokémon
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", searchQuery)
	id, locationAreaEncounters, height, weight, name, types, abilities, image, err := FetchPokemonDetails(pokemonURL) // Modifié pour inclure les abilities
	if err != nil {
		log.Printf("Failed to fetch details for %s: %v", searchQuery, err) // Add logging here
		if strings.Contains(err.Error(), "404") {
			InitTemp.Temp.ExecuteTemplate(w, "search", map[string]string{
				"ErrorMessage": "Aucun Pokémon trouvé",
			})
		} else {
			http.Error(w, fmt.Sprintf("Erreur lors de la récupération des détails de Pokémon: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Si un Pokémon est trouvé, continuez comme d'habitude
	pokemon := Pokemon{
		ID:                     id,
		LocationAreaEncounters: locationAreaEncounters,
		Height:                 height,
		Weight:                 weight,
		Name:                   name,
		Type:                   types,
		Abilities:              abilities, // Assurez-vous d'ajouter les capacités ici
		Image:                  image,
	}

	log.Printf("Fetched details for %s: %+v", searchQuery, pokemon)
	InitTemp.Temp.ExecuteTemplate(w, "search", pokemon)
}

func PokemonDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extrayez le nom du Pokémon de l'URL
	name := strings.TrimPrefix(r.URL.Path, "/pokemon/")

	// Utilisez `FetchPokemonDetails` pour obtenir les détails du Pokémon
	id, locationAreaEncounters, height, weight, name, types, abilities, image, err := FetchPokemonDetails("https://pokeapi.co/api/v2/pokemon/" + name) // Ajouté abilities dans la récupération
	if err != nil {
		http.Error(w, "Pokémon non trouvé", http.StatusNotFound)
		return
	}
	if err != nil {
		// Gérez l'erreur, peut-être en continuant sans les informations d'évolution
		log.Printf("Erreur lors de la récupération des détails d'évolution: %v", err)
	}

	damageRelations, err := FetchTypeDamageRelations(types[0])
	if err != nil {
		log.Printf("Erreur lors de la récupération des relations de dommages pour le type %s: %v", types[0], err)
		// Vous pouvez soit ignorer cette erreur soit retourner une erreur HTTP
	}

	// Créez une instance de Pokémon avec les détails obtenus
	pokemon := Pokemon{
		ID:                     id,
		LocationAreaEncounters: locationAreaEncounters,
		Height:                 height,
		Weight:                 weight,
		Name:                   name,
		Type:                   types,
		Abilities:              abilities, // Incluez les capacités ici
		Image:                  image,
		DamageRelations:        damageRelations,
	}

	// Passez le Pokémon au template de détail
	InitTemp.Temp.ExecuteTemplate(w, "pokemon", pokemon)
}

func FetchEvolutionDetails(evolutionChainID int) (evolutions []string, err error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/evolution-chain/%d/", evolutionChainID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Chain struct {
			EvolvesTo []struct {
				Species struct {
					Name string `json:"name"`
				} `json:"species"`
				EvolvesTo []struct {
					Species struct {
						Name string `json:"name"`
					} `json:"species"`
				} `json:"evolves_to"`
			} `json:"evolves_to"`
		} `json:"chain"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Ajouter le premier Pokémon (base) à la liste des évolutions
	evolutions = append(evolutions, data.Chain.EvolvesTo[0].Species.Name)

	// Parcourir la chaîne d'évolution et ajouter chaque évolution
	for _, evolution := range data.Chain.EvolvesTo {
		if len(evolution.EvolvesTo) > 0 {
			evolutions = append(evolutions, evolution.EvolvesTo[0].Species.Name)
		}
	}

	return evolutions, nil
}

func FetchTypeDamageRelations(typeName string) (DamageRelations, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/type/%s", typeName)
	resp, err := http.Get(url)
	if err != nil {
		return DamageRelations{}, err
	}
	defer resp.Body.Close()

	var damageRelations DamageRelations
	if err := json.NewDecoder(resp.Body).Decode(&damageRelations); err != nil {
		return DamageRelations{}, err
	}

	return damageRelations, nil
}

func FiltrerTypeHandler(w http.ResponseWriter, r *http.Request) {
	// Ici, tu récupères les types pour les passer au template.
	types, err := FetchPokemonTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Type     []string
		Pokemons []Pokemon // Change ici pour utiliser []Pokemon au lieu de []string
	}{
		Type: types,
	}

	if r.Method == "POST" {
		typeName := r.FormValue("type")
		pokemons, err := FetchPokemonsByType(typeName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data.Pokemons = pokemons
	}

	InitTemp.Temp.ExecuteTemplate(w, "filtrerType", data)
}

func FetchPokemonsByType(typeName string) ([]Pokemon, error) {
	// Appelle l'API pour obtenir tous les Pokémon d'un certain type
	url := "https://pokeapi.co/api/v2/type/" + typeName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var typeResp struct {
		Pokemon []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
		} `json:"pokemon"`
	}
	if err := json.Unmarshal(body, &typeResp); err != nil {
		return nil, err
	}

	var pokemons []Pokemon
	for i, p := range typeResp.Pokemon {
		if i >= 20 { // Limite à 20 Pokémon
			break
		}
		// Utilise l'URL de chaque Pokémon pour obtenir des détails supplémentaires
		id, locationAreaEncounters, height, weight, name, types, abilities, image, err := FetchPokemonDetails(p.Pokemon.URL)
		if err != nil {
			// Gérer l'erreur ou continuer avec le prochain Pokémon
			continue
		}
		pokemons = append(pokemons, Pokemon{
			ID:                     id,
			LocationAreaEncounters: locationAreaEncounters,
			Height:                 height,
			Weight:                 weight,
			Name:                   name,
			Type:                   types,
			Abilities:              abilities,
			Image:                  image,
		})
	}

	return pokemons, nil
}

func FetchPokemonTypes() ([]string, error) {
	url := "https://pokeapi.co/api/v2/type/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response TypesResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var types []string
	for _, t := range response.Results {
		types = append(types, t.Name)
	}

	return types, nil
}

func ServePokemonsHandlers(w http.ResponseWriter, r *http.Request) {

	starterIDs := map[int][]int{
		1: {1, 4, 7},       // Génération 1: Bulbizarre, Salamèche, Carapuce
		2: {152, 155, 158}, // Génération 2: Germignon, Héricendre, Kaiminus
		3: {252, 255, 258}, // Génération 3: Arcko, Poussifeu, Gobou
		4: {387, 390, 393}, // Génération 4: Tortipouss, Ouisticram, Tiplouf
		5: {495, 498, 501}, // Génération 5: Vipélierre, Gruikui, Moustillon
		6: {650, 653, 656}, // Génération 6: Marisson, Feunnec, Grenousse
		7: {722, 725, 728}, // Génération 7: Brindibou, Flamiaou, Otaquin
		8: {810, 813, 816}, // Génération 8: Ouistempo, Flambino, Larméléon
		// Ajouter d'autres générations si nécessaire
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Assume deux générations par page pour simplifier l'exemple
	genStart := (page-1)*2 + 1
	genEnd := genStart + 1

	var pokemons []SimplePokemon
	for gen := genStart; gen <= genEnd; gen++ {
		if ids, ok := starterIDs[gen]; ok {
			for _, id := range ids {
				pokemon, err := GetPokemonInfoByID(id)
				if err != nil {
					fmt.Printf("Failed to get pokemon info for ID %d: %v\n", id, err)
					continue
				}
				pokemons = append(pokemons, pokemon)
			}
		}
	}

	// Calcul de la pagination
	isLast := genEnd >= len(starterIDs)

	tmplData := struct {
		Pokemons []SimplePokemon
		PrevPage int
		NextPage int
		IsFirst  bool
		IsLast   bool
	}{
		Pokemons: pokemons,
		PrevPage: max(1, page-1),
		NextPage: page + 1,
		IsFirst:  page == 1,
		IsLast:   isLast,
	}

	if err := InitTemp.Temp.ExecuteTemplate(w, "starter", tmplData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetPokemonInfoByID(id int) (SimplePokemon, error) {
	var sp SimplePokemon
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return sp, err
	}
	defer resp.Body.Close()

	var data struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Sprites struct {
			Other struct {
				OfficialArtwork struct {
					FrontDefault string `json:"front_default"`
				} `json:"official-artwork"`
			} `json:"other"`
		} `json:"sprites"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return sp, err
	}

	sp.Name = data.Name
	sp.Type = data.Type
	sp.Image = data.Sprites.Other.OfficialArtwork.FrontDefault

	return sp, nil
}

func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	pokemonID, err := strconv.Atoi(r.FormValue("pokemonId"))
	if err != nil {
		http.Error(w, "Invalid Pokémon ID", http.StatusBadRequest)
		return
	}

	// Lire le fichier JSON
	var favs Favorites
	file, err := ioutil.ReadFile("favorites.json")
	if err != nil {
		// Gérer l'erreur ou créer le fichier s'il n'existe pas
	}
	json.Unmarshal(file, &favs)

	// Ajouter l'ID aux favoris s'il n'est pas déjà présent
	for _, id := range favs.IDs {
		if id == pokemonID {
			// Pokémon déjà en favori, redirection ou message approprié
			return
		}
	}
	favs.IDs = append(favs.IDs, pokemonID)

	// Écrire le fichier JSON mis à jour
	updatedFavs, err := json.Marshal(favs)
	if err != nil {
		// Gérer l'erreur
	}

	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Erreur lors de la lecture des favoris", http.StatusInternalServerError)
		return
	} else if os.IsNotExist(err) {
		favs = Favorites{IDs: []int{}}
	}

	if err := ioutil.WriteFile("favorites.json", updatedFavs, 0644); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde des favoris", http.StatusInternalServerError)
		return
	}

	// Redirection vers la page des favoris
	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}

func ReadFavorites() (*Favorites, error) {
	var favs Favorites

	// Tentez de lire le contenu du fichier JSON des favoris
	bytes, err := ioutil.ReadFile("favorites.json")
	if err != nil {
		// Si le fichier n'existe pas, retournez une nouvelle instance de Favorites
		if os.IsNotExist(err) {
			return &Favorites{IDs: []int{}}, nil
		}
		// Pour toute autre erreur de lecture, la retourner
		return nil, err
	}

	// Décodez le contenu JSON en objet Favorites
	err = json.Unmarshal(bytes, &favs)
	if err != nil {
		return nil, err
	}

	return &favs, nil
}

// Handler pour afficher les Pokémon favoris
func ShowFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	favs, err := ReadFavorites()
	if err != nil {
		http.Error(w, "Failed to read favorites", http.StatusInternalServerError)
		return
	}

	var pokemons []Pokemon
	for _, id := range favs.IDs {
		pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)
		id, locationAreaEncounters, height, weight, name, types, abilities, image, err := FetchPokemonDetails(pokemonURL)
		if err != nil {
			log.Printf("Failed to fetch details for pokemon with ID %d: %v\n", id, err)
			continue // Passer au prochain ID en cas d'erreur
		}

		pokemon := Pokemon{
			ID:                     id,
			LocationAreaEncounters: locationAreaEncounters,
			Height:                 height,
			Weight:                 weight,
			Name:                   name,
			Type:                   types,
			Abilities:              abilities,
			Image:                  image,
		}

		pokemons = append(pokemons, pokemon)
	}
	// Assurez-vous que vous avez le template HTML approprié pour afficher la liste des pokémons

	InitTemp.Temp.ExecuteTemplate(w, "favorites", pokemons) // Assurez-vous que le template est correctement défini pour traiter []Pokemon
}

func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	pokemonID, err := strconv.Atoi(r.FormValue("pokemonId"))
	if err != nil {
		http.Error(w, "Invalid Pokémon ID", http.StatusBadRequest)
		return
	}

	// Lire les favoris actuels
	favs, err := ReadFavorites()
	if err != nil {
		http.Error(w, "Failed to read favorites", http.StatusInternalServerError)
		return
	}

	// Trouver et supprimer l'ID du slice
	for i, id := range favs.IDs {
		if id == pokemonID {
			favs.IDs = append(favs.IDs[:i], favs.IDs[i+1:]...)
			break
		}
	}

	// Réécrire le fichier JSON avec l'ID supprimé
	updatedFavs, err := json.Marshal(favs)
	if err != nil {
		http.Error(w, "Failed to update favorites", http.StatusInternalServerError)
		return
	}
	if err := ioutil.WriteFile("favorites.json", updatedFavs, 0644); err != nil {
		http.Error(w, "Failed to save favorites", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page des favoris ou confirmer la suppression
	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}
