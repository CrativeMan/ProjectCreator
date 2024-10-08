package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	lip "github.com/charmbracelet/lipgloss"
)

func parseArgs() Args {
	var args Args
	flag.BoolVar(&args.Help, "h", false, "show help")
	flag.BoolVar(&args.Version, "v", false, "show version")
	flag.BoolVar(&args.Flake, "nf", false, "dont generate a flake.nix file")

	flag.Parse()

	return args
}

func main() {
	exit := initial()

	if !exit {
		initialForm := prompUserWithLanguage()
		err := initialForm.Run()
		if err != nil {
			log.Fatalf("Failed to run init form: %v\n", err)
		}

		if language != CLOSE {
			pathForm := promptUserWithPath()
			err = pathForm.Run()
			if err != nil {
				log.Fatalf("Failed to get path from user: %v\n", err)
			}
			if path == "" {
				path, err = os.Getwd()
				if err != nil {
					panic(err)
				}
			}
			fmt.Printf("Creating project at: %s\n", path)
			path = _makeGlobalPath(path)

			switch language {
			case C:
				createCEnv(path)
			case CPP:
				createCppEnv(path)
			case GO:
				createGoEnv(path)
			case RUST:
				createRustEnv(path)
			case CLOSE:
				break
			default:
				log.Fatal("Failed to create languageEnv.\nUnexpected language detected.")
			}

			fmt.Println(sty.success.Render("Successfully created project"))
		}
	}
}

func initial() bool {
	var exit bool

	sty.success = lip.NewStyle().Bold(true).Foreground(lip.Color("86"))
	sty.fail = lip.NewStyle().Bold(true).Foreground(lip.Color("9"))
	sty.warning = lip.NewStyle().Bold(true).Foreground(lip.Color("#ffb300"))

	arguments := parseArgs()

	if arguments.Version {
		fmt.Printf("createp: %s\n", version)
		exit = true
	}

	if arguments.Help {
		fmt.Printf(helpText, version)
		exit = true
	}

	if arguments.Flake {
		GenerateFlake = false
		fmt.Println("Not generating a flake.nix file")
	}

	return exit
}

func createCEnv(path string) {
	// ask for project type
	projType := cProjectType()

	// make file
	writeRunFile(path, C)

	// standard dependencies
	dependencies := []string{
		"clang-tools",
		"llvmPackages.clangUseLLVM",
		"gcc",
		"clang",
		"cmake",
	}

	switch projType {
	case NORM:
		// direnv
		writeEnvrc(path)
		_allowDirenv(path)

		writeFlake(path, dependencies)

		writeMain(path, C)
	case NCURSES:
		// direnv
		writeEnvrc(path)
		_allowDirenv(path)

		dependencies = append(dependencies, "ncurses")
		writeFlake(path, dependencies)

		writeMain(path, C)
	case RAYLIB:
		writeEnvrc(path)
		_allowDirenv(path)

		// extra dependencies
		dependencies = append(dependencies, "raylib")
		writeFlake(path, dependencies)

		writeMain(path, C)
	case SUB:
		writeMain(path, C)
	default:
		log.Fatalf("Unknown project type detected")
	}
}

func createCppEnv(path string) {
	// ask for project type
	projType := cppProjectType()

	// run files
	writeRunFile(path, CPP)

	// standard dependencies
	dependencies := []string{
		"clang-tools",
		"llvmPackages.clangUseLLVM",
		"gcc",
		"clang",
		"cmake",
	}

	switch projType {
	case NORM:
	case NCURSES:
		// direnv
		writeEnvrc(path)
		_allowDirenv(path)

		// extra dependencies
		if projType == NCURSES {
			dependencies = append(dependencies, "ncurses")
		}

		writeFlake(path, dependencies)

		writeMain(path, CPP)
	case SUB:
		writeMain(path, CPP)
	default:
		log.Fatalf("Unknown project type detected")
	}
}

func createGoEnv(path string) {
	projType := goProjectType()
	GoModuleName = askUserForGoModuleName()
	writeRunFile(path, GO)

	dependencies := []string{
		"go",
		"gofumpt",
		"goimports-reviser",
		"golines",
		"delve",
	}

	switch projType {
	case NORM:

		writeEnvrc(path)
		_allowDirenv(path)
		writeFlake(path, dependencies)
		writeGoMod(path)
		writeMain(path, GO)
	case RAYLIB:
		writeEnvrc(path)
		_allowDirenv(path)

		dependencies = append(dependencies, "wayland")
		dependencies = append(dependencies, "wayland-protocols")
		dependencies = append(dependencies, "glew")
		dependencies = append(dependencies, "glfw")
		dependencies = append(dependencies, "libxkbcommon")
		dependencies = append(dependencies, "xorg.libX11")
		dependencies = append(dependencies, "xorg.libXcursor")
		dependencies = append(dependencies, "xorg.libXi")
		dependencies = append(dependencies, "xorg.libXrandr")
		dependencies = append(dependencies, "xorg.libXineram")

		writeFlake(path, dependencies)

		writeGoMod(path)

		writeMain(path, GO)
	case SUB:
		writeGoMod(path)
		writeMain(path, GO)
	case COBRA:
		writeEnvrc(path)
		_allowDirenv(path)
		dependencies = append(dependencies, "cobra-cli")
		writeFlake(path, dependencies)
		writeGoMod(path)
		writeMain(path, GO)
	default:
		log.Fatalf("Unknown project type detected")
	}
}

func createRustEnv(path string) {
	projType := goProjectType()
	GoModuleName = askUserForGoModuleName()
	writeRunFile(path, GO)

	dependencies := []string{
		"cargo",
		"rustc",
		"rust-analyzer",
		"rustfmt",
		"clippy",
		"bacon",
	}

	switch projType {
	case NORM:
		writeEnvrc(path)
		_allowDirenv(path)
		writeFlake(path, dependencies)
		writeMain(path, RUST)
	default:
		log.Fatalf("Unknown project type detected")
	}
}

func _isValidPath(path string) error {
	if language != CLOSE {
		if path == "\n" || path == "" {
			var err error
			return err
		}

		path = _makeGlobalPath(path)

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				err2 := _createFolder(path)
				if err2 != nil {
					return err2
				}
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
