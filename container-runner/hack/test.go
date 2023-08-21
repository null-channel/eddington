package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/null-channel/eddington/container-runner/internal/templates"
)

func main() {
	fmt.Println("Hello, world!")
	serviceTemplate := templates.ServiceTemplate{
		NullApplicationName: "my-nullapp",
		AppName:             "my-app",
		CustomerID:          "my-customer",
	}

	t, err := template.New("service").Parse(templates.Service)

	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, serviceTemplate)

	if err != nil {
		panic(err)
	}
}
