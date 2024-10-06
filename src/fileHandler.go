package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func writeFlake(path string, dependencies []string) {
	if GenerateFlake {
		flakeName := "flake.nix"
		flake, err := os.Create(path + flakeName)
		if err != nil {
			panic(err)
		}
		defer flake.Close()

		depAll := parseDependencies(dependencies)
		contents := fmt.Sprintf(FLAKECONTENT, depAll)

		_, err = flake.WriteString(contents)
		if err != nil {
			panic(err)
		}

		err = flake.Sync()
		if err != nil {
			panic(err)
		}
		fmt.Println(sty.success.Render("Created flake"))
	}
}

func writeMain(path string, language int) {
	switch language {
	case C:
		_writeCMain(path)
	case CPP:
		_writeCppMain(path)
	case GO:
		_writeGoMain(path)
	case RUST:
		_writeRustMain(path)
	default:
		log.Fatalf("Unknown language detected")
	}
}

func writeRunFile(path string, language int) {
	switch language {
	case C:
		_writeCRun(path)
	case CPP:
		_writeCppRun(path)
	case GO:
		_writeGoRun(path)
	default:
		log.Fatalf("Unknown language detected")
	}
}

func writeGoMod(path string) {
	goModName := "go.mod"
	goMod, err := os.Create(path + goModName)
	if err != nil {
		panic(err)
	}
	defer goMod.Close()

	_, err = goMod.WriteString(fmt.Sprintf("module %s\n\ngo 1.22.2", GoModuleName))
	if err != nil {
		panic(err)
	}

	err = goMod.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created go.mod file"))
}

func writeEnvrc(path string) {
	if GenerateFlake {
		filename := ".envrc"
		envrc, err := os.Create(path + filename)
		if err != nil {
			panic(err)
		}
		defer envrc.Close()

		_, err = envrc.WriteString(ENVRCCONTENT)
		if err != nil {
			panic(err)
		}

		err = envrc.Sync()
		if err != nil {
			panic(err)
		}

		fmt.Println(sty.success.Render("Created .envrc file"))
	}
}

// WRITE MAIN FILES

func _writeCMain(path string) {
	name := "main.c"
	if err := os.Mkdir(path+"src", os.ModePerm); err != nil {
		panic(err)
	}
	mainC, err := os.Create(path + "/src/" + name)
	if err != nil {
		panic(err)
	}
	defer mainC.Close()

	_, err = mainC.WriteString(CMAINCONTENTS)
	if err != nil {
		panic(err)
	}

	err = mainC.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Creatd main.c file"))
}

func _writeCppMain(path string) {
	name := "main.cpp"
	main, err := os.Create(path + name)
	if err != nil {
		panic(err)
	}
	defer main.Close()

	_, err = main.WriteString(CPPMAINCONTENTS)
	if err != nil {
		panic(err)
	}

	err = main.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Creatd main.cpp file"))
}

func _writeGoMain(path string) {
	name := "main.go"
	mainGo, err := os.Create(path + name)
	if err != nil {
		panic(err)
	}
	defer mainGo.Close()

	_, err = mainGo.WriteString(GOMAINCONTENTS)
	if err != nil {
		panic(err)
	}

	err = mainGo.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created main.go file"))
}

func _writeRustMain(path string) {
	fmt.Println(path)
}

// WRITE RUN FILES

func _writeCRun(path string) {
	makeName := "Makefile"

	make, err := os.Create(path + makeName)
	if err != nil {
		panic(err)
	}
	defer make.Close()

	_, err = make.WriteString(CMAKECONTENTS)
	if err != nil {
		panic(err)
	}

	err = make.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created make file"))
}

func _writeCppRun(path string) {
	makeName := "Makefile"

	make, err := os.Create(path + makeName)
	if err != nil {
		panic(err)
	}
	defer make.Close()

	_, err = make.WriteString(CPPMAKECONTENTS)
	if err != nil {
		panic(err)
	}

	err = make.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created make file"))
}

func _writeGoRun(path string) {
	makeName := "Makefile"

	make, err := os.Create(path + makeName)
	if err != nil {
		panic(err)
	}
	defer make.Close()

	_, err = make.WriteString(GOMAKECONTENTS)
	if err != nil {
		panic(err)
	}

	err = make.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created make file"))
}

// MISC STUFF

func _allowDirenv(path string) {
	if GenerateFlake {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = os.Chdir(path)
		if err != nil {
			panic(err)
		}
		cmd := exec.Command("direnv", "allow")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		err = os.Chdir(dir)
		if err != nil {
			panic(err)
		}
		fmt.Println(sty.success.Render("Allowed direnv"))
	}
}

func parseDependencies(dep []string) string {
	var depen string
	for _, elem := range dep {
		depen += elem + "\n"
	}

	return depen
}
