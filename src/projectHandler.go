package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
)

var projType int = 0

const (
	NORM   = 0
	RAYLIB = 1
	SUB    = 2
)

func prompUserWithLanguage() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose programming language: ").
				Options(
					huh.NewOption("C", C),
					huh.NewOption("C++", CPP),
					huh.NewOption("Go", GO),
					// huh.NewOption("Java"+WIP, JAVA),
					huh.NewOption("Exit", CLOSE),
				).
				Value(&language),
		).WithTheme(huh.ThemeDracula()),
	)
}

func promptUserWithPath() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter path to project").
				Prompt("Path: ").
				Placeholder("If left empty it's this path.").
				Validate(_isValidPath).
				Value(&path),
		).WithTheme(huh.ThemeDracula()),
	)	
}

func askUserForGoModuleName() string {
	var moduleName string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Module name:").
				Placeholder("If left empty its main").
				Value(&moduleName).
				CharLimit(20),
		).WithTheme(huh.ThemeDracula()),
	)
	err := form.Run()
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to run ask user for go module name:"))
		log.Fatal(err)
	}

	if moduleName == "" {
		moduleName = "main"
	}

	return moduleName
}

// C project
func cProjectType() int {
	cForm := cProjTypeForm()
	err := cForm.Run()
	if err != nil {
		panic(err)
	}

	return projType
}

func cProjTypeForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select Project Type").
				Options(
					huh.NewOption("Normal", NORM),
					huh.NewOption("Raylib", RAYLIB),
					huh.NewOption("SubProject", SUB),
				).
				Value(&projType),
		).WithTheme(huh.ThemeDracula()),
	)
}

// CPP project
func cppProjectType() int {
	cForm := cppProjTypeForm()
	err := cForm.Run()
	if err != nil {
		panic(err)
	}

	return projType
}

func cppProjTypeForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select Project Type").
				Options(
					huh.NewOption("Normal", NORM),
					huh.NewOption("SubProject", SUB),
				).
				Value(&projType),
		).WithTheme(huh.ThemeDracula()),
	)
}

// GO project
func goProjectType() int {
	goForm := goProjTypeForm()
	err := goForm.Run()
	if err != nil {
		panic(err)
	}

	return projType
}

func goProjTypeForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Select Project Type").
				Options(
					huh.NewOption("Normal", NORM),
					huh.NewOption("Raylib", RAYLIB),
					huh.NewOption("SubProject", SUB),
				).
				Value(&projType),
		).WithTheme(huh.ThemeDracula()),
	)
}
