package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", standupOrderHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func standupOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	team := getTeam(r)
	date := getDate(r)

	standupOrder := getStandupOrder(team, date)

	fmt.Printf("Standup order for %s: %v", date.Format("2006-01-02"), standupOrder)

	response := "Standup order for "
	response += date.Format("2006-01-02")
	response += "\n============================\n\n"
	response += strings.Join(standupOrder, "\n")
	response += "\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func getTeam(r *http.Request) []string {
	q := r.URL.Query()
	team := q["team"]
	sort.Strings(team)
	return team
}

func getDate(r *http.Request) time.Time {
	q := r.URL.Query()
	date := time.Now().UTC()
	if len(q["date"]) == 1 {
		fmt.Println(q["date"][0])
		urldate, err := time.Parse("2006-01-02", q.Get("date"))
		if err == nil {
			date = urldate
		}
	}
	return date
}

func removeElement(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func generateSeed(team []string, date time.Time) int64 {
	dateteam := date.Format("2006-01-02")
	dateteam += strings.Join(team, "-")
	sum := sha256.Sum256([]byte(dateteam))
	return int64(binary.BigEndian.Uint64(sum[:]))
}

func getStandupOrder(team []string, date time.Time) []string {
	rand.Seed(generateSeed(team, date))
	teamsize := len(team)
	teamcopy := make([]string, teamsize)
	copy(teamcopy, team)
	standupOrder := []string{}
	for i := 0; i < teamsize; i++ {
		num := rand.Intn(len(teamcopy))
		standupOrder = append(standupOrder, teamcopy[num])
		teamcopy = removeElement(teamcopy, num)
	}
	return standupOrder
}
