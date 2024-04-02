document.getElementById('addToFavBtn').addEventListener('click', function() {
    const pokemonId = getPokemonId(); // Implémente cette fonction pour récupérer l'ID du Pokémon actuel
    addPokemonToFavorites(pokemonId);
});

function addPokemonToFavorites(pokemonId) {
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    if (!favorites.includes(pokemonId)) {
        favorites.push(pokemonId);
        localStorage.setItem('favorites', JSON.stringify(favorites));
    }
}

function removePokemonFromFavorites(pokemonId) {
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    favorites = favorites.filter(id => id !== pokemonId);
    localStorage.setItem('favorites', JSON.stringify(favorites));
}

function getPokemonId() {
    // Cette fonction doit récupérer l'ID du Pokémon de la page. Cela dépend de ta structure HTML/JS.
    // Par exemple, si l'ID est stocké dans un attribut data-id sur le bouton :
    return document.getElementById('addToFavBtn').getAttribute('data-id');
}

// Appel cette fonction pour afficher les Pokémon favoris
function displayFavorites() {
    let favorites = JSON.parse(localStorage.getItem('favorites')) || [];
    favorites.forEach(pokemonId => {
        // Utilise pokemonId pour récupérer et afficher les détails des Pokémon favoris
        // Cela pourrait impliquer de faire une requête à l'API PokeAPI et de mettre à jour le DOM avec les résultats
    });
}
