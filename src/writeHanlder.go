package main

import (
	"fmt"
	"os"
	"os/exec"
)

func writeFlake(path string, language int, dependencies []string) {
	depString := parseDependencies(dependencies)

	flakeName := "flake.nix"
	flake, err := os.Create(path + flakeName)
	if err != nil {
		panic(err)
	}
	defer flake.Close()

	switch language {
	case C:
		depAll := combineDeps(depString, CFLAKECONTENT)
		_writeCFlake(flake, depAll)
	case CPP:
		_writeCppFlake(flake)
	case GO:
		depAll := combineDeps(depString, GOFLAKECONTENT)
		_writeGoFlake(flake, depAll)
	case JAVA:
		_writeJavaFlake(flake)
	}

	err = flake.Sync()
	if err != nil {
		panic(err)
	}
	fmt.Println(sty.success.Render("Created flake"))
}

func writeMain(path string, language int) {
	switch language {
	case C:
		_writeCMain(path)
	case CPP:
		_writeCppMain(path)
	case GO:
		_writeGoMain(path)
	case JAVA:
		_writeJavaMain(path)
	}
}

func writeRunFile(path string, language int) {
	switch language {
	case C:
		_writeCRun(path)
	case GO:
		_writeGoRun(path)
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

// WRITE FLAKE FILES

func _writeCFlake(flake *os.File, contents string) {
	_, err := flake.WriteString(contents)
	if err != nil {
		panic(err)
	}
}

// TODO: change this
func _writeCppFlake(flake *os.File) {
	_, err := flake.WriteString(CFLAKECONTENT)
	if err != nil {
		panic(err)
	}
}

func _writeGoFlake(flake *os.File, contents string) {
	_, err := flake.WriteString(contents)
	if err != nil {
		panic(err)
	}
}

// TODO: change this
func _writeJavaFlake(flake *os.File) {
	_, err := flake.WriteString(CFLAKECONTENT)
	if err != nil {
		panic(err)
	}
}

// WRITE MAIN FILES

func _writeCMain(path string) {
	name := "main.c"
	mainC, err := os.Create(path + name)
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
	fmt.Println(path)
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

func _writeJavaMain(path string) {
	fmt.Println(path)
}

// WRITE RUN FILES

func _writeCRun(path string) {
	goBuildName := "build"
	goRunName := "run"

	goBuild, err := os.Create(path + goBuildName)
	if err != nil {
		panic(err)
	}
	goRun, err := os.Create(path + goRunName)
	if err != nil {
		panic(err)
	}
	defer goBuild.Close()
	defer goRun.Close()

	_, err = goBuild.WriteString("gcc main.c -o main")
	if err != nil {
		panic(err)
	}
	_, err = goRun.WriteString("./build\n./main")
	if err != nil {
		panic(err)
	}

	err = goBuild.Sync()
	if err != nil {
		panic(err)
	}
	err = goRun.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created build and run file"))
}

func _writeGoRun(path string) {
	goBuildName := "build"
	goRunName := "run"

	goBuild, err := os.Create(path + goBuildName)
	if err != nil {
		panic(err)
	}
	goRun, err := os.Create(path + goRunName)
	if err != nil {
		panic(err)
	}
	defer goBuild.Close()
	defer goRun.Close()

	_, err = goBuild.WriteString("go build -o main -v")
	if err != nil {
		panic(err)
	}
	_, err = goRun.WriteString("./build\n./main")
	if err != nil {
		panic(err)
	}

	err = goBuild.Sync()
	if err != nil {
		panic(err)
	}
	err = goRun.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Println(sty.success.Render("Created build and run file"))
}

func _allowDirenv(path string) {
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

func parseDependencies(dep []string) string {
	var depen string
	for _, elem := range dep {
		depen += "\t\t\t\t" + elem + "\n"
	}

	return depen
}

func combineDeps(dep string, flake string) string {
	return fmt.Sprintf("%s\n%s%s", flake, dep, FLAKECONTENT_END)
}
