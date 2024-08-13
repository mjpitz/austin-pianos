package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

//go:embed readme.md.gotmpl
var ReadmeTemplate string

type Application struct {
	Config struct {
		Database string
		Readme   string
		JSON     string
	}

	Template *template.Template
	DB       DB

	err error
}

func (app *Application) loadTemplate() {
	if app.err != nil {
		return
	}

	log.Println("[TEMPLATE] parsing contents")
	app.Template, app.err = template.New("gen").
		Funcs(sprig.FuncMap()).
		Funcs(map[string]any{
			"urlenc": func(in string) string {
				return strings.ReplaceAll(in, "-", "·")
			},
		}).
		Parse(ReadmeTemplate)
}

func (app *Application) loadDB() {
	if app.err != nil {
		return
	}

	log.Println("[DB] loading contents")
	var contents []byte
	contents, app.err = os.ReadFile(app.Config.Database)
	if app.err != nil {
		return
	}

	log.Println("[DB] unmarshalling data")
	app.err = yaml.Unmarshal(contents, &app.DB)
}

func (app *Application) renderJSON() {
	if app.err != nil {
		return
	}

	log.Println("[JSON] rendering database")
	db, err := json.MarshalIndent(app.DB, "", "  ")
	if err != nil {
		app.err = err
		return
	}

	log.Println("[JSON] writing file")
	err = os.WriteFile(app.Config.JSON, db, 0644)
	if err != nil {
		app.err = err
	}
}

func (app *Application) renderREADME() {
	if app.err != nil {
		return
	}

	log.Println("[README] rendering contents")
	buf := bytes.NewBuffer(nil)
	app.err = app.Template.Execute(buf, app.DB)
	if app.err != nil {
		return
	}

	log.Println("[README] writing file")
	err := os.WriteFile(app.Config.Readme, buf.Bytes(), 0644)
	if err != nil {
		app.err = err
	}
}

// Run will execute the underlying business logic of the application. This includes reading the database file and
// rendering the various output files. This could definitely be improved to handle partial write situations, but
// for now, it's more than enough.
func (app *Application) Run() error {
	log.Println("[INPUT] -database=", app.Config.Database)
	log.Println("[INPUT] -json=", app.Config.JSON)
	log.Println("[INPUT] -readme=", app.Config.Readme)

	app.loadTemplate()
	app.loadDB()
	app.renderJSON()
	app.renderREADME()

	return app.err
}

func main() {
	app := &Application{}

	cli := flag.NewFlagSet("gen", flag.ExitOnError)

	cli.StringVar(&app.Config.Database, "database", "db.yaml", "Override the input db.yaml location.")
	cli.StringVar(&app.Config.JSON, "json", "db.json", "Override the resulting db.json location.")
	cli.StringVar(&app.Config.Readme, "readme", "README.md", "Override the resulting README.md location.")

	err := cli.Parse(os.Args[1:])
	if err != nil {
		return
	}

	err = app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
