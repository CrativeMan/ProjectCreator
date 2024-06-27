package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	lip "github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"os/exec"
)

const (
	C     = 0
	CPP   = 1
	GO    = 2
	JAVA  = 3
	ENVRC = 4
)

var (
	language int
	path     string
	sty      styles
	Hostname string
	err      error
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

	initialForm := promptUserWithChoices()
	err := initialForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	// check if path ends with a / and if not add it
	pathRunes := []rune(path)
	lastPathRune := pathRunes[len(pathRunes)-1]
	if lastPathRune != '/' {
		path += "/"
	}

	switch language {
	case C:
		// Run createCEnv
		fmt.Println(sty.success.Render("C"))
	case CPP:
		// Run createCppEnv
		fmt.Println(sty.success.Render("C++"))
	case GO:
		_writeFiles(path, GO)
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


	_writeFiles(path, ENVRC)
	_allowDirenv(err)

	fmt.Println(sty.success.Render("Successfully created project"))
}

func promptUserWithChoices() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose programming language: ").
				Options(
					huh.NewOption("C", C),
					huh.NewOption("C++", CPP),
					huh.NewOption("Go", GO),
					huh.NewOption("Java", JAVA),
				).
				Value(&language),

			huh.NewInput().
				Title("Enter path to project (if left empty this path)").
				Prompt("Path:").
				Validate(_isValidPath).
				Value(&path),
		).WithTheme(huh.ThemeDracula()),
	)
}

func createGoEnv(path string) error {
	// ask for module name
	moduleName := askUserForGoModuleName()
	// create go mod file
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir(path)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(dir)
	fmt.Println(sty.success.Render("Created go mod file"))
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
		break
	case CPP:
		break
	case GO:
		goMainName := "main.go"
		flake, err := os.Create(path + flakeName)
		if err != nil {
			log.Fatal(err)
		}
		mainGo, err := os.Create(path + goMainName)
		if err != nil {
			log.Fatal(err)
		}
		defer flake.Close()
		defer mainGo.Close()

		_, err = flake.WriteString(GOFLAKECONTENT)
		if err != nil {
			log.Fatal(err)
		}
		_, err = mainGo.WriteString(GOMAINCONTENTS)
		if err != nil {
			log.Fatal(err)
		}
		err = flake.Sync()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created flake.nix file"))
		err = mainGo.Sync()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created main.go file"))
	case JAVA:
		break
	case ENVRC:
		filename := ".envrc"
		file, err := os.Create(path + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(ENVRCCONTENT)
		if err != nil {
			log.Fatal(err)
		}
		err = file.Sync()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sty.success.Render("Created .envrc file"))
	}
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
}

func _isValidPath(path string) error {
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

	return nil
}

const (
	GOFLAKECONTENT string = `{
  description = "Basic Go development environment";

  inputs = {
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    flake-utils,
    nixpkgs,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
        ];
      };
    });
}`
	GOMAINCONTENTS string = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}`
	ENVRCCONTENT string = "use flake"
)
