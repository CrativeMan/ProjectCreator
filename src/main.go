package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	lip "github.com/charmbracelet/lipgloss"
)

const (
	C     = 0
	CPP   = 1
	GO    = 2
	JAVA  = 3
	ENVRC = 4
)

var (
	language     int
	path         string
	sty          styles
	Hostname     string
	GoModuleName string
)

type styles struct {
	success lip.Style
	fail    lip.Style
	warning lip.Style
}

func main() {
	sty.success = lip.NewStyle().Bold(true).Foreground(lip.Color("86"))
	sty.fail = lip.NewStyle().Bold(true).Foreground(lip.Color("9"))
	sty.warning = lip.NewStyle().Bold(true).Foreground(lip.Color("#ffb300"))

	fmt.Println("Version 0.1.9")

	initialForm := promptUserWithChoices()
	err := initialForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating project at: %s\n", path)
	path = _makeGlobalPath(path)

	switch language {
	case C:
		createCEnv(path)
	case CPP:
		fmt.Println(sty.success.Render("C++"))
	case GO:
		createGoEnv(path)
	case JAVA:
		fmt.Println(sty.success.Render("Java"))
	default:
		log.Fatal("Failed to create languageEnv.\nUnexpected language detected.")
	}

	fmt.Println(sty.success.Render("Successfully created project"))
}

func promptUserWithChoices() *huh.Form {
	WIP := sty.warning.Render(" (WIP)")
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose programming language: ").
				Options(
					huh.NewOption("C", C),
					huh.NewOption("C++"+WIP, CPP),
					huh.NewOption("Go", GO),
					huh.NewOption("Java"+WIP, JAVA),
				).
				Value(&language),

			huh.NewInput().
				Title("Enter path to project").
				Prompt("Path: ").
				Placeholder("If left empty it's this path.").
				Validate(_isValidPath).
				Value(&path),
		).WithTheme(huh.ThemeDracula()),
	)
}

func createCEnv(path string) {
	projType := cProjectType()
	writeRunFile(path, C)

	_chmodFile(path, "run")
	_chmodFile(path, "build")

	switch projType {
	case NORM:
		dependencies := []string{
			"clang-tools",
			"llvmPackages.clangUseLLVM",
			"gcc",
		}
		writeEnvrc(path)
		_allowDirenv(path)

		writeFlake(path, C, dependencies)
		writeMain(path, C)
	case RAYLIB:
		dependencies := []string{
			"clang-tools",
			"llvmPackages.clangUseLLVM",
			"gcc",
			"raylib",
		}
		writeEnvrc(path)
		_allowDirenv(path)

		writeFlake(path, C, dependencies)
		writeMain(path, C)
	case SUB:
		writeMain(path, C)
	}
}

func createGoEnv(path string) {
	projType := goProjectType()
	GoModuleName = askUserForGoModuleName()
	writeRunFile(path, GO)

	_chmodFile(path, "run")
	_chmodFile(path, "build")

	switch projType {
	case NORM:
		dependencies := []string{
			"go",
			"gofumpt",
			"goimports-reviser",
			"golines",
			"delve",
		}
		writeEnvrc(path)
		_allowDirenv(path)
		writeFlake(path, GO, dependencies)
		writeGoMod(path)
		writeMain(path, GO)
	case RAYLIB:
		dependencies := []string{
			"go",
			"gofumpt",
			"goimports-reviser",
			"golines",
			"delve",
			"wayland",
			"wayland-protocols",
			"glew",
			"glfw",
			"libxkbcommon",
			"xorg.libX11",
			"xorg.libXcursor",
			"xorg.libXi",
			"xorg.libXrandr",
			"xorg.libXineram",
		}
		writeEnvrc(path)
		_allowDirenv(path)
		writeFlake(path, GO, dependencies)
		writeGoMod(path)
		writeMain(path, GO)
	case SUB:
		writeGoMod(path)
		writeMain(path, GO)
	}
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

func _isValidPath(path string) error {
	path = _makeGlobalPath(path)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err2 := _createFolder(path)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func _createFolder(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	fmt.Println(sty.success.Render("Created folder"))
	return nil
}

func _makeGlobalPath(path string) string {
	// check if path ends with a / and if not add it
	pathRunes := []rune(path)
	if pathRunes[len(pathRunes)-1] != '/' {
		path += "/"
	}
	if pathRunes[0] == '~' {
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to get user home dir:"))
			log.Fatal(err)
		}
		path = dir + path[1:]
	}
	return path
}

func _chmodFile(path string, filename string) {
	fileInfo, err := os.Stat(path + filename)
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to get fileinfo from ", filename))
		log.Fatal(err)
	}

	mode := fileInfo.Mode()
	execMode := mode | 0100

	err = os.Chmod(path+filename, execMode)
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to make " + filename + " executable"))
		log.Fatal(err)
	}

	fmt.Println(sty.success.Render("Made " + filename + " executable"))
}
