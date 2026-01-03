// This program retrieves programming practice questions from online platforms
// based on a user-specified topic or tag. Its purpose is to provide a structured
// and extensible way to discover relevant coding problems for learning,
// practice, or analysis.
//
// Current Scope
//
// At present, the implementation supports fetching questions exclusively from
// Codeforces. The design intentionally keeps the data source abstract so that
// additional platforms can be integrated in the future with minimal changes.

package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

const (
	CODEFORCES_BASE_URL = "https://codeforces.com/api/"
	// To access the data you just send a HTTP-request to address
	// https://codeforces.com/api/{methodName} with method-specific parameters.
	PROBLEMSET_METHOD = "problemset.problems"
)

// available tags on codeforces
var CodeForcesTags = []string{
	"dp", "greedy", "math", "geometry", "string",
	"data structures", "trees", "graphs", "sorting", "binary search",
	"hashing", "bitmasks", "dp", "trees", "graphs", "sorting",
	"binary search", "hashing", "bitmasks",
}

// represents a problem object from codeforces response
type Problem struct {
	ContestID      int      `json:"contestId"`
	ProblemSetName string   `json:"problemsetName"`
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Points         float64  `json:"points"`
	Rating         int      `json:"rating"`
	Tags           []string `json:"tags"`
}

// represents a response from codeforces API
type CodeforcesResponse struct {
	Status string `json:"status"`
	Result struct {
		Problems []Problem `json:"problems"`
	} `json:"result"`
}

// fetchCodeforcesProblemSetWithTag fetches a list of problems based on the provided tag
// this method returns all the problems with the provided tag with other tags too
func fetchCodeforcesProblemSetWithTag(tag string) ([]Problem, error) {
	url := CODEFORCES_BASE_URL + PROBLEMSET_METHOD + "?tags=" + tag
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching problem set, err=%v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading problem set, err=%v", err)
		return nil, err
	}

	var codeforcesResponse CodeforcesResponse
	err = json.Unmarshal(body, &codeforcesResponse)
	if err != nil {
		log.Printf("Error unmarshalling problem set, err=%v", err)
		return nil, err
	}

	// sort the problems by rating
	sort.Slice(codeforcesResponse.Result.Problems, func(i, j int) bool {
		return codeforcesResponse.Result.Problems[i].Rating < codeforcesResponse.Result.Problems[j].Rating
	})

	return codeforcesResponse.Result.Problems, nil
}

// fetchCodeforcesProblemSetWithTagOnly fetches a list of problems based on the provided tag
// this method returns only problems with the provided tag only
func fetchCodeforcesProblemSetWithTagOnly(tag string) ([]Problem, error) {
	allProblems, err := fetchCodeforcesProblemSetWithTag(tag)
	if err != nil {
		log.Printf("Error fetching problem set, err=%v", err)
		return nil, err
	}

	var problems []Problem
	for _, problem := range allProblems {
		if len(problem.Tags) == 1 {
			problems = append(problems, problem)
		}
	}
	return problems, nil
}

// fetchCodeforcesProblemSetWithTags fetches a list of problems based on the provided tags
// this method returns problems with all the provided tags
func fetchCodeforcesProblemSetWithTags(tags []string) ([]Problem, error) {
	url := CODEFORCES_BASE_URL + PROBLEMSET_METHOD + "?tags=" + strings.Join(tags, ";")
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching problem set, err=%v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading problem set, err=%v", err)
		return nil, err
	}

	var codeforcesResponse CodeforcesResponse
	err = json.Unmarshal(body, &codeforcesResponse)
	if err != nil {
		log.Printf("Error unmarshalling problem set, err=%v", err)
		return nil, err
	}

	// sort the problems by rating
	sort.Slice(codeforcesResponse.Result.Problems, func(i, j int) bool {
		return codeforcesResponse.Result.Problems[i].Rating < codeforcesResponse.Result.Problems[j].Rating
	})

	return codeforcesResponse.Result.Problems, nil
}

// fetchCodeforcesProblemSetWithTagsOnly fetches a list of problems based on the provided tags
// this method returns only problems with the provided tags only
func fetchCodeforcesProblemSetWithTagsOnly(tags []string) ([]Problem, error) {
	allProblems, err := fetchCodeforcesProblemSetWithTags(tags)
	if err != nil {
		log.Printf("Error fetching problem set, err=%v", err)
		return nil, err
	}

	var problems []Problem
	for _, problem := range allProblems {
		if len(problem.Tags) == len(tags) {
			problems = append(problems, problem)
		}
	}
	return problems, nil
}

// <--------------------------------- handlers --------------------------------->

// problemsByTagHandler is a handler to get problems by tag
func problemsByTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	problems, err := fetchCodeforcesProblemSetWithTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// problemsByTagsHandler is a handler to get problems by tags
func problemsByTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	problems, err := fetchCodeforcesProblemSetWithTags(strings.Split(tags, ","))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// problemsByTagOnlyHandler is a handler to get problems by tag only
func problemsByTagOnlyHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	problems, err := fetchCodeforcesProblemSetWithTagOnly(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// problemsByTagsOnlyHandler is a handler to get problems by tags only
func problemsByTagsOnlyHandler(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query().Get("tags")
	problems, err := fetchCodeforcesProblemSetWithTagsOnly(strings.Split(tags, ","))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// <--------------------------------- server code --------------------------------->

func main() {
	const PORT = "49160"

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("/problems", problemsByTagHandler)
	mux.HandleFunc("/problems/multi", problemsByTagsHandler)
	mux.HandleFunc("/problems/only", problemsByTagOnlyHandler)
	mux.HandleFunc("/problems/multi/only", problemsByTagsOnlyHandler)

	server := &http.Server{
		Addr:              ":" + PORT,
		Handler:           loggingMiddleware(mux),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Println("HTTP server started on port=" + PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error, err=%v", err)
		}
	}()

	shutdownGracefully(server)
}

// healthHandler is a simple handler to check if the server is running
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// loggingMiddleware is a middleware to log the requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// shutdownGracefully is a function to shutdown the server gracefully
// it waits for 10 seconds for the server to shutdown gracefully
// if the server does not shutdown gracefully, it forcefully shuts down the server
func shutdownGracefully(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed, err=%v", err)
	}
	log.Println("server shutdown gracefully")
}
