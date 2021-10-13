package main

import (
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

func parseExperimentFormData(f ExperimentForm) Experiment {
	var responsive, head, leg1, leg2, leg3, leg4, leg5, leg6, wing1, wing2 bool
	if f.Responsive == "yes" {
		responsive = true
	} else {
		responsive = false
	}
	if f.Head == "on" {
		head = false
	} else {
		head = true
	}
	if f.Leg1 == "on" {
		leg1 = false
	} else {
		leg1 = true
	}
	if f.Leg2 == "on" {
		leg2 = false
	} else {
		leg2 = true
	}
	if f.Leg3 == "on" {
		leg3 = false
	} else {
		leg3 = true
	}
	if f.Leg4 == "on" {
		leg4 = false
	} else {
		leg4 = true
	}
	if f.Leg5 == "on" {
		leg5 = false
	} else {
		leg5 = true
	}
	if f.Leg6 == "on" {
		leg6 = false
	} else {
		leg6 = true
	}
	if f.Wing1 == "on" {
		wing1 = false
	} else {
		wing1 = true
	}
	if f.Wing2 == "on" {
		wing2 = false
	} else {
		wing2 = true
	}

	experiment := Experiment{
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
	}

	return experiment
}

func (env *Env) receiveHandler(w http.ResponseWriter, r *http.Request) {
	formData := ExperimentForm{
		Responsive: r.FormValue("responsive"),
		Head:       r.FormValue("head"),
		Leg1:       r.FormValue("leg1"),
		Leg2:       r.FormValue("leg2"),
		Leg3:       r.FormValue("leg3"),
		Leg4:       r.FormValue("leg4"),
		Leg5:       r.FormValue("leg5"),
		Leg6:       r.FormValue("leg6"),
		Wing1:      r.FormValue("wing1"),
		Wing2:      r.FormValue("wing2"),
	}

	experiment := parseExperimentFormData(formData)

	_ = env.store.storeExperiment(experiment)
	log.Print(experiment)
	http.Redirect(w, r, "/results", http.StatusSeeOther)
}

func (env *Env) resultsHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := env.store.getResult()

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
