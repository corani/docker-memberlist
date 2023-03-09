package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/memberlist"
)

func main() {
	list, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		panic(err)
	}

	me, err := os.Hostname()
	if err != nil {
		panic("Failed to get hostname: " + err.Error())
	}

	log.Printf("Hostname: %v", me)

	if me != "app1" {
		log.Println("Trying to join app1")

		_, err := list.Join([]string{"app1"})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	}

	log.Println("Known members:")
	for _, member := range list.Members() {
		log.Printf("- %s @ %s", member.Name, member.Addr)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "<h1>%s (%d)</h1>\n", me, list.GetHealthScore())

		fmt.Fprintf(w, "<ul>\n")
		for _, member := range list.Members() {
			fmt.Fprintf(w, "<li>%s @ %s</li>", member.Name, member.Addr)
		}
		fmt.Fprintf(w, "</ul>\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
