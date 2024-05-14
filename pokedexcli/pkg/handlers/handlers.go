package handlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/abi01shek/pokedexcli/pkg/apiCalls"
	"github.com/abi01shek/pokedexcli/pkg/pokecache"
)

const defatulApiAddress = "https://pokeapi.co/api/v2/location-area/?limit=20&offset=20"
const exploreBaseAddress = "https://pokeapi.co/api/v2/location-area/"
const pokemonBaseAddress = "https://pokeapi.co/api/v2/pokemon/"

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

type locationApiResT struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type exploreAreaT struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type pokemonT struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type config struct {
	locationPrev        string
	locationNext        string
	locationCache       *pokecache.Cache
	exploreCache        *pokecache.Cache
	pokemonInCurrentLoc map[string]bool
	caughtPokemon       map[string]pokemonT
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Printf("Welcome to Pokedex!\nUsage: \n")
	for _, cmd := range getCommand() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)

	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}

// getLocation: Gets location data for a given address either from
// the cache or through API
func getLocation(cfg *config, addr string) ([]byte, error) {
	var body []byte
	var err error
	body, found := cfg.locationCache.Get(addr)
	if !found {
		body, err = apiCalls.GetBodyApiCall(addr)
		if err != nil {
			return nil, err
		}
		cfg.locationCache.Add(addr, body)
	}
	return body, nil
}

// exploreLocation: Gets the explore location data for given address
// either from cache or through API
func exploreLocation(cfg *config, addr string) ([]byte, error) {
	var body []byte
	var err error
	body, found := cfg.exploreCache.Get(addr)
	if !found {
		body, err = apiCalls.GetBodyApiCall(addr)
		if err != nil {
			return nil, err
		}
		cfg.exploreCache.Add(addr, body)
	}
	return body, nil
}

// commandMap: Get the next 20 locations
func commandMap(cfg *config, args ...string) error {
	nextLocAddr := cfg.locationNext
	if nextLocAddr == "" {
		nextLocAddr = defatulApiAddress
	}

	body, err := getLocation(cfg, nextLocAddr)
	if err != nil {
		return err
	}

	locRes := locationApiResT{}
	err = json.Unmarshal(body, &locRes)
	if err != nil {
		return err
	}

	cfg.locationNext = *locRes.Next
	cfg.locationPrev = *locRes.Previous

	for _, myLoc := range locRes.Results {
		fmt.Printf("%s\n", myLoc.Name)
	}
	return nil
}

// commandMapb : get previous 20 locations
func commandMapb(cfg *config, args ...string) error {
	prevLocAddr := cfg.locationPrev
	if prevLocAddr == "" {
		return errors.New("no previous locations found")
	}

	body, err := getLocation(cfg, prevLocAddr)
	if err != nil {
		return err
	}

	locRes := locationApiResT{}
	err = json.Unmarshal(body, &locRes)
	if err != nil {
		return err
	}

	if locRes.Next != nil {
		cfg.locationNext = *locRes.Next
	} else {
		cfg.locationNext = ""
	}

	if locRes.Previous != nil {
		cfg.locationPrev = *locRes.Previous
	} else {
		cfg.locationPrev = ""
	}

	for _, myLoc := range locRes.Results {
		fmt.Printf("%s\n", myLoc.Name)
	}
	return nil

}

func commandExplore(cfg *config, args ...string) error {
	expLoc := strings.Join(args[:], "")
	apiAddr := exploreBaseAddress + expLoc
	fmt.Printf("Exploring %s ...\n", expLoc)

	body, err := exploreLocation(cfg, apiAddr)
	if err != nil {
		return err
	}

	exploreRes := exploreAreaT{}
	err = json.Unmarshal(body, &exploreRes)
	if err != nil {
		return err
	}

	cfg.pokemonInCurrentLoc = make(map[string]bool)
	fmt.Printf("Found Pokemon:\n")
	for _, pe := range exploreRes.PokemonEncounters {
		fmt.Printf("\t- %s\n", pe.Pokemon.Name)
		cfg.pokemonInCurrentLoc[pe.Pokemon.Name] = true
	}

	return nil
}

// commandCatch: try to catch a pokemon
func commandCatch(cfg *config, args ...string) error {
	pokemonName := strings.Join(args[:], "")
	if _, exists := cfg.pokemonInCurrentLoc[pokemonName]; !exists {
		fmt.Printf("Pokemon %s not found in current location\n", pokemonName)
		return nil
	}

	apiAddr := pokemonBaseAddress + pokemonName
	body, err := apiCalls.GetBodyApiCall(apiAddr)
	if err != nil {
		return nil
	}

	pokemonRes := pokemonT{}
	err = json.Unmarshal(body, &pokemonRes)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	rval := rand.Intn(pokemonRes.BaseExperience)
	if rval > 40 {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)
	cfg.caughtPokemon[pokemonName] = pokemonRes
	return nil
}

// commandInspect: inpsect a pokemon if it is in your pokedex
func commandInspect(cfg *config, args ...string) error {
	pokemonName := strings.Join(args[:], "")
	if pe, exists := cfg.caughtPokemon[pokemonName]; exists {
		fmt.Printf("Name: %s\n", pe.Name)
		fmt.Printf("Height: %d\n", pe.Height)
		fmt.Printf("Weight: %d\n", pe.Weight)
		return nil
	}
	fmt.Printf("Pokemon %s does not exist in your pokedex\n", pokemonName)
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Printf("Your Pokedex:\n")
	for pokemonName, _ := range cfg.caughtPokemon {
		fmt.Printf("\t- %s\n", pokemonName)
	}
	return nil
}

// getCommand: list all the commands available
func getCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Get next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore a given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon with its name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon in your pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "lists all pokemons in your pokedex",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(inp string) []string {
	inp = strings.ToLower(inp)
	words := strings.Fields(inp)
	return words
}

func StartRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}
	cfg.locationCache = pokecache.NewCache(5 * time.Minute)
	cfg.exploreCache = pokecache.NewCache(5 * time.Minute)
	go cfg.locationCache.ReadLoop()
	go cfg.exploreCache.ReadLoop()
	cfg.caughtPokemon = make(map[string]pokemonT)
	for {
		fmt.Printf("Pokedex> ") // shell prompt
		scanner.Scan()
		if scanner.Text() == "" {
			continue
		}
		words := cleanInput(scanner.Text())
		commandName := words[0]
		args := words[1:]

		if command, exists := getCommand()[commandName]; !exists {
			fmt.Printf("Unknown command %v", commandName)
			continue
		} else {
			err := command.callback(&cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
