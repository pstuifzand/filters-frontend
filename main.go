package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func generateField(r FilterRule) string {
	switch r.Field {
	case "to":
		return `$mail->to()`
	case "from":
		return `$mail->from()`
	case "cc":
		return `$mail->cc()`
	case "bcc":
		return `$mail->bcc()`
	case "subject":
		return `$mail->subject()`
	case "body":
		return `$mail->body()`
	}
	return "ERROR"
}

func generateRule(w io.Writer, r FilterRule) {
	switch r.Function {
	case "begins":
		fmt.Fprintf(w, "%s =~ m{^\\Q%s\\E}", generateField(r), r.Arg)
	case "ends":
		fmt.Fprintf(w, "%s =~ m{\\Q%s\\E$}", generateField(r), r.Arg)
	case "equal":
		fmt.Fprintf(w, "%s eq %q", generateField(r), r.Arg)
	case "not_equal":
		fmt.Fprintf(w, "%s ne %q", generateField(r), r.Arg)
	case "contains":
		fmt.Fprintf(w, "%s =~ m{\\Q%s\\E}", generateField(r), r.Arg)
	default:
		fmt.Fprintf(w, "unknown function %s\n", r.Function)
	}
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

		fmt.Fprint(w, "use Email::Filter;\n\n")
		fmt.Fprint(w, `my $email = Email::Filter->new(emergency => "~/emergency_mbox");`)
		fmt.Fprintln(w)

		for _, f := range fg.Filters {
			fmt.Fprintf(w, "\n// Filter - %s\n", f.Name)

			op := "||"
			startValue := 0
			if f.Combine == "and" {
				startValue = 1
				op = "&&"
			}
			fmt.Fprintf(w, "my $result_%d = %d;\n", f.Id, startValue)
			for _, r := range f.Rules {
				fmt.Fprintf(w, "$result_%d %s= ", f.Id, op)
				generateRule(w, r)
				fmt.Fprint(w, ";\n")
			}
			fmt.Fprintf(w, "if ($result_%d) {\n", f.Id)
			for _, action := range f.Actions {
				fmt.Fprintf(w, "    // actions %s %s\n", action.Action, action.ActionValue)
			}
			fmt.Fprint(w, "    return;\n}\n")
		}

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
