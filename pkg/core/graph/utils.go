package graph

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"text/template"
)

//go:embed "graph.tmpl"
var htmlTemplate string

// GenerateHTML generates a new HTML file with the loaded data.
func GenerateHTML(elements Elements, w io.Writer) error {
	// Take the info from gg struct and json indent it to a string

	jsonString, err := json.MarshalIndent(elements, "", "    ")
	if err != nil {
		return err
	}

	vars := struct {
		Elements string
	}{
		Elements: string(jsonString),
	}

	t, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(w, vars)
	if err != nil {
		return err
	}

	return nil
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
