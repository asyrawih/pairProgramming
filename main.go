package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponsePokemon struct {
	Count    int       `json:"count"`
	Next     int       `json:"next"`
	Previous int       `json:"previous"`
	Results  []Pokemon `json:"results"`
}

type Pokemon struct {
	Name string
	Url  string
}

var baseUrl = "https://pokeapi.co/api/v2/"

func main() {
	p := ResponsePokemon{}
	pokemon, err := p.getPokemon()

	if err != nil {
		fmt.Print(err)
	}

	// Iterate The Pokemon
	for _, poke := range pokemon.Results {
		fmt.Println(poke.Name)
		fmt.Println(poke.Url)
	}
}

// without Struct
func fetcPokemon() (responsePokemon *ResponsePokemon, err error) {

	response, err := http.Get("https://pokeapi.co/api/v2/berry-flavor/")

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var data ResponsePokemon

	json.Unmarshal(result, &data)

	return &data, nil
}

// Via Structs
func (p *ResponsePokemon) getPokemon() (responsePokemon *ResponsePokemon, err error) {

	response, err := http.Get("https://pokeapi.co/api/v2/berry-flavor/")

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var data ResponsePokemon

	json.Unmarshal(result, &data)

	return &data, nil
}
