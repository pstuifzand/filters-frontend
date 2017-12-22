package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Filters: []Filter{
// 	Filter{
// 		Name:    "Alerts",
// 		Combine: "and",
// 		Rules: []FilterRule{
// 			FilterRule{
// 				Function: "ends",
// 				Field:    "to",
// 				Arg:      "@stuifzand.com",
// 			},
// 			FilterRule{
// 				Function: "begins",
// 				Field:    "subject",
// 				Arg:      "Alert:",
// 			},
// 		},
// 		Actions: []Action{
// 			Action{
// 				Action:      "move_to_folder",
// 				ActionValue: "Spam",
// 			},
// 			Action{
// 				Action:      "remove",
// 				ActionValue: "",
// 			},
// 		},
// 	},
// },

type FilterRule struct {
	Field    string `json:"field"`
	Function string `json:"func"`
	Arg      string `json:"arg"`
}

type Action struct {
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
}

type Filter struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Combine string       `json:"combine"`
	Rules   []FilterRule `json:"rules"`
	Actions []Action     `json:"actions"`
}

type Filters struct {
	Filters []Filter `json:"filters"`
}

type apiHandler struct {
}

func (handler *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f, err := os.Open("store.json")
		defer f.Close()
		dec := json.NewDecoder(f)
		var fg Filters
		err = dec.Decode(&fg)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		n := 1
		for v := range fg.Filters {
			fg.Filters[v].Id = n
			n++
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		enc := json.NewEncoder(w)
		enc.Encode(fg)
		return
	} else if r.Method == http.MethodPost {
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		var fg Filters
		err := dec.Decode(&fg)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		fp, err := os.OpenFile("store.json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			log.Print(err)
			return
		}

		defer fp.Close()
		enc := json.NewEncoder(fp)
		enc.Encode(fg)

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(204)
		return
	} else if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(405)
}

type generateHandler struct {
}

func (handler *generateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		defer r.Body.Close()

		f, err := os.Open("store.json")
		defer f.Close()
		dec := json.NewDecoder(f)
		var fg Filters
		err = dec.Decode(&fg)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		generateFilterFile(w, &fg)

		return

	} else if r.Method == http.MethodPost {
		defer r.Body.Close()

		f, err := os.Open("store.json")
		defer f.Close()
		dec := json.NewDecoder(f)
		var fg Filters
		err = dec.Decode(&fg)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		generateFilterFile(w, &fg)

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteHeader(204)
		return

	} else if r.Method == http.MethodOptions {
		defer r.Body.Close()

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(405)
}

func main() {
	http.Handle("/store", &apiHandler{})
	http.Handle("/generate", &generateHandler{})
	http.Handle("/", http.FileServer(http.Dir("dist")))
	log.Fatal(http.ListenAndServe(":8088", nil))
}
