package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

	fmt.Println("Version 0.1.8")

	initialForm := promptUserWithChoices()
	err := initialForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating project at: %s\n", path)
	path = _makeGlobalPath(path)

	_writeFiles(path, ENVRC)
	_allowDirenv(err)

	switch language {
	case C:
		_writeFiles(path, C)
	case CPP:
		// Run createCppEnv
		fmt.Println(sty.success.Render("C++"))
	case GO:
		err := createGoEnv(path)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create go env:"))
			log.Fatal(err)
		}
	case JAVA:
		// Run createJavaEnv
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

func createGoEnv(path string) error {
	// ask for module name
	GoModuleName = askUserForGoModuleName()
	_writeFiles(path, GO)
	createGoModule(path)
	return nil
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

func _writeFiles(path string, langType int) {
	flakeName := "flake.nix"

	switch langType {
	case C:
		cMainName := "main.c"
		flake, err := os.Create(path + flakeName)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create flake.nix file:"))
			log.Fatal(err)
		}
		defer flake.Close()
		mainC, err := os.Create(path + cMainName)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create main.c file:"))
			log.Fatal(err)
		}
		defer mainC.Close()

		_, err = flake.WriteString(CFLAKECONTENT)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to write c flake:"))
			log.Fatal(err)
		}
		_, err = mainC.WriteString(CMAINCONTENTS)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to write c main:"))
			log.Fatal(err)
		}
		err = mainC.Sync()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to sync c main:"))
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created main.c file"))

		// close flake file buffer
		err = flake.Sync()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to sync go flake:"))
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created flake.nix file"))
	case CPP:
		break
	case GO:
		goMainName := "main.go"
		flake, err := os.Create(path + flakeName)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create flake.nix file:"))
			log.Fatal(err)
		}
		defer flake.Close()
		mainGo, err := os.Create(path + goMainName)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create main.go file:"))
			log.Fatal(err)
		}
		defer mainGo.Close()

		_, err = flake.WriteString(GOFLAKECONTENT)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to write go flake:"))
			log.Fatal(err)
		}
		_, err = mainGo.WriteString(GOMAINCONTENTS)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to write go main:"))
			log.Fatal(err)
		}
		err = mainGo.Sync()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to sync go main:"))
			log.Fatal(err)
		}

		// close flake file buffer
		err = flake.Sync()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to sync go flake:"))
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created flake.nix file"))
	case JAVA:
		break
	case ENVRC:
		filename := ".envrc"
		file, err := os.Create(path + filename)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to create .envrc:"))
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(ENVRCCONTENT)
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to write .envrc:"))
			log.Fatal(err)
		}
		err = file.Sync()
		if err != nil {
			fmt.Println(sty.fail.Render("Failed to sync .envrc:"))
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created .envrc file"))
	}
}

func createGoModule(path string) {
	// create go.mod file
	goModName := "go.mod"
	goMod, err := os.Create(path + goModName)
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to create go.mod file:"))
		log.Fatal(err)
	}
	defer goMod.Close()

	_, err = goMod.WriteString(fmt.Sprintf("module %s\n\ngo 1.22.2", GoModuleName))
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to write go.mod:"))
		log.Fatal(err)
	}
	err = goMod.Sync()
	if err != nil {
		fmt.Println(sty.fail.Render("Failed to sync go.mod:"))
		log.Fatal(err)
	}

	fmt.Println(sty.success.Render("Created go.mod file"))
}

func _allowDirenv(err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(path)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("direnv", "allow")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sty.success.Render("Allowed direnv for " + path))
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
