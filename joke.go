package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sort"
	"strconv"
)

var jokes = map[int]Joke{
	1: {
		Id:   1,
		Joke: "What did one ocean say to the other ocean? Nothing, it just waved.",
	},
	2: {
		Id:   2,
		Joke: "What's the best thing about Switzerland? I don't know, but the flag is a big plus.",
	},
	3: {
		Id:   3,
		Joke: "How many tickles does it take to get an octopus to laugh? Ten tickles.",
	},
}
var newJokeId = 4

type Joke struct {
	Id   int    `json:"id"`
	Joke string `json:"joke"`
}

func searchJokesHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	jokeSlice := make([]Joke, 0, len(jokes))
	for _, joke := range jokes {
		jokeSlice = append(jokeSlice, joke)
	}
	sort.Slice(jokeSlice, func(i, j int) bool {
		return jokeSlice[i].Id < jokeSlice[j].Id
	})

	response := struct {
		Jokes []Joke `json:"jokes"`
	}{
		jokeSlice,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func createJokeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := struct {
		Joke string `json:"joke"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Joke == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	joke := Joke{
		Id:   newJokeId,
		Joke: request.Joke,
	}
	jokes[joke.Id] = joke
	newJokeId++

	response := struct {
		Joke Joke `json:"joke"`
	}{
		joke,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func deleteJokeHandler(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	delete(jokes, id)
	w.WriteHeader(http.StatusNoContent)
}
