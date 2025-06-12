Fun little command line pokemon game written in GO

## To Build
go build -o pokedexcli ./cmd/web/

## Play
./pokedexcli

### Get Help
```
Pokedex> help
Welcome to Pokedex!
Usage:
inspect: Inspect a pokemon in your pokedex
pokedex: lists all pokemons in your pokedex
help: Displays a help message
exit: Exit the Pokedex
map: Get next 20 locations
mapb: Get previous 20 locations
explore: explore a given location
catch: Catch a pokemon with its name
```

### Check the map for different regions
```
Pokedex> map
mt-coronet-1f-route-216
mt-coronet-1f-route-211
mt-coronet-b1f
great-marsh-area-1
great-marsh-area-2
...
```

### Explore a region
```
Pokedex> explore mt-coronet-1f-route-216
Exploring mt-coronet-1f-route-216 ...
Found Pokemon:
	- clefairy
	- golbat
	- machoke
	- graveler
	- nosepass
	- meditite
	- chingling
	- bronzor
```

### Catch Pokemons!
```
Pokedex> catch golbat
Throwing a Pokeball at golbat...
golbat was caught!
```

### Inspect Pokemons
```
Pokedex> inspect golbat
Name: golbat
Height: 16
Weight: 550
```

### Check your Pokedex
```
Pokedex> pokedex
Your Pokedex:
	- golbat
```
