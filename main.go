package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

var home = loadNestedTemplates("template/home.htm")
var submit = loadNestedTemplates("template/submit.htm", "template/mosquito.htm")
var results = loadNestedTemplates("template/results.htm")
var assets = loadNestedTemplates("template/assets.htm")

var static = http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

func main() {
	ds := SuccessStore{}
	env := &Env{store: ds}

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/receive", env.receiveHandler)
	http.HandleFunc("/results", env.resultsHandler)
	http.HandleFunc("/assets", assetsHandler)
	http.Handle("/static/", static)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	homeHandler(w, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, home, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, submit, nil)
}

func parseExperimentFormData(r *http.Request) (Experiment, error) {
	var responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2 bool
	if r.FormValue("responsive") == "yes" {
		responsive = true
	} else if r.FormValue("responsive") == "no" {
		responsive = false
	} else {
		return Experiment{}, errors.New("Incorrect responsive value")
	}
	if r.FormValue("head") == "on" {
		head = false
	} else {
		head = true
	}
	if r.FormValue("leg1") == "on" {
		leg1 = false
	} else {
		leg1 = true
	}
	if r.FormValue("leg2") == "on" {
		leg2 = false
	} else {
		leg2 = true
	}
	if r.FormValue("leg3") == "on" {
		leg3 = false
	} else {
		leg3 = true
	}
	if r.FormValue("leg4") == "on" {
		leg4 = false
	} else {
		leg4 = true
	}
	if r.FormValue("leg5") == "on" {
		leg5 = false
	} else {
		leg5 = true
	}
	if r.FormValue("leg6") == "on" {
		leg6 = false
	} else {
		leg6 = true
	}
	if r.FormValue("wing1") == "on" {
		wing1 = false
	} else {
		wing1 = true
	}
	if r.FormValue("wing2") == "on" {
		wing2 = false
	} else {
		wing2 = true
	}

	return Experiment{
		Responsive: responsive,
		Head:       head,
		Leg1:       leg1,
		Leg2:       leg2,
		Leg3:       leg3,
		Leg4:       leg4,
		Leg5:       leg5,
		Leg6:       leg6,
		Wing1:      wing1,
		Wing2:      wing2,
	}, nil
}

func (env *Env) receiveHandler(w http.ResponseWriter, r *http.Request) {
	experiment, err := parseExperimentFormData(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.store.storeExperiment(experiment)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func (env *Env) resultsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := env.store.getResult()

	if err != nil {
		log.Fatal(err)
	}

	render(w, results, res)
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	render(w, assets, nil)
}

func render(w http.ResponseWriter, t *template.Template, data interface{}) {
	err := t.ExecuteTemplate(w, "layout.htm", data)
	if err != nil {
		log.Fatal(err)
	}
}

func loadNestedTemplates(filenames ...string) *template.Template {
	t := append([]string{"template/layout.htm"}, filenames...)
	return template.Must(template.ParseFiles(t...))
}
