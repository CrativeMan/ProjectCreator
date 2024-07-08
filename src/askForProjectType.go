package main

import "github.com/charmbracelet/huh"

var projType int = 0

const (
	NORM   = 0
	RAYLIB = 1
	SUB    = 2
)

func cProjectType() int {
	cForm := cProjTypeForm()
	err := cForm.Run()
	if err != nil {
		panic(err)
	}

	return projType
}

func goProjectType() int {
	goForm := goProjTypeForm()
	err := goForm.Run()
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
