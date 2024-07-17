package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const (
	FLAKE_START = "/flake_start"
	FLAKE_END   = "/flake_end"
	C_          = "/c_"
	GO_         = "/go_"
	ENVRC_MAIN  = "/envrc_content"

	MAIN  = "main"
	BUILD = "build"
	RUN   = "run"
)

func writeFlake(path string, language int, dependencies []string) {
	depString := parseDependencies(dependencies)

	flakeName := "flake.nix"
	flake, err := os.Create(path + flakeName)
	if err != nil {
		panic(err)
	}
	defer flake.Close()

	if GetFilesLocaly {
		depAll := combineDeps(depString, FLAKECONTENT_START, FLAKECONTENT_END)

		switch language {
		case C:
			_writeCFlake(flake, depAll)
		case GO:
			_writeGoFlake(flake, depAll)
		}
	} else {
		flakeStart, err := SFTPCLIENT.Open(SFTPPATH + FLAKE_START)
		if err != nil {
			log.Fatalf("Failed to open flakestart: %v", err)
		}
		defer flakeStart.Close()

		flakeEnd, err := SFTPCLIENT.Open(SFTPPATH + FLAKE_END)
		if err != nil {
			log.Fatalf("Failed to open flakeend: %v", err)
		}
		defer flakeEnd.Close()

		flakeStartContents, err := io.ReadAll(flakeStart)
		if err != nil {
			log.Fatalf("Failed to read flakestart: %v", err)
		}

		flakeEndContents, err := io.ReadAll(flakeEnd)
		if err != nil {
			log.Fatalf("Failed to read flakestart: %v", err)
		}

		depAll := combineDeps(depString, string(flakeStartContents), string(flakeEndContents))

		switch language {
		case C:
			_readCFlake(flake, depAll)
		case GO:
			_readGoFlake(flake, depAll)
		}
	}

	err = flake.Sync()
	if err != nil {
		panic(err)
	}
	fmt.Println(sty.success.Render("Created flake"))
}

func writeMain(path string, language int) {
	if GetFilesLocaly {
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
	} else {
		switch language {
		case C:
			_readCMain(path)
		case GO:
			_readGoMain(path)
		}
	}
}

func writeRunFile(path string, language int) {
	if GetFilesLocaly {
		switch language {
		case C:
			_writeCRun(path)
		case GO:
			_writeGoRun(path)
		}
	} else {
		switch language {
		case C:
			_readCRun(path)
		case GO:
			_readGoRun(path)
		}
	}
}

// TODO: adapt this to network
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

// TODO: adapt this to network
func writeEnvrc(path string) {
	filename := ".envrc"
	envrc, err := os.Create(path + filename)
	if err != nil {
		panic(err)
	}
	defer envrc.Close()

	if GetFilesLocaly {
		_, err = envrc.WriteString(ENVRCCONTENT)
		if err != nil {
			panic(err)
		}
	} else {
		remEnv, err := SFTPCLIENT.Open(SFTPPATH + ENVRC_MAIN)
		if err != nil {
			panic(err)
		}
		defer remEnv.Close()

		_, err = envrc.ReadFrom(remEnv)
		if err != nil {
			panic(err)
		}
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

func _writeGoFlake(flake *os.File, contents string) {
	_, err := flake.WriteString(contents)
	if err != nil {
		panic(err)
	}
}

// READ FLAKE FILS
// TODO: clean up

func _readCFlake(flake *os.File, deps string) {
	_, err := flake.WriteString(deps)
	if err != nil {
		panic(err)
	}
}

func _readGoFlake(flake *os.File, deps string) {
	_, err := flake.WriteString(deps)
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

// READ MAIN FILES
// TODO: clean up

func _readCMain(path string) {
	name := "main.c"
	main, err := os.Create(path + name)
	if err != nil {
		panic(err)
	}
	defer main.Close()

	file, err := SFTPCLIENT.Open(SFTPPATH + C_ + MAIN)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = main.ReadFrom(file)
	if err != nil {
		panic(err)
	}
}

func _readGoMain(path string) {
	name := "main.go"
	main, err := os.Create(path + name)
	if err != nil {
		panic(err)
	}
	defer main.Close()

	file, err := SFTPCLIENT.Open(SFTPPATH + GO_ + MAIN)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = main.ReadFrom(file)
	if err != nil {
		panic(err)
	}
}

// WRITE RUN FILES

func _writeCRun(path string) {
	cBuildName := "build"
	cRunName := "run"

	cBuild, err := os.Create(path + cBuildName)
	if err != nil {
		panic(err)
	}
	cRun, err := os.Create(path + cRunName)
	if err != nil {
		panic(err)
	}
	defer cBuild.Close()
	defer cRun.Close()

	_, err = cBuild.WriteString("gcc main.c -o main")
	if err != nil {
		panic(err)
	}
	_, err = cRun.WriteString("./build\n./main")
	if err != nil {
		panic(err)
	}

	err = cBuild.Sync()
	if err != nil {
		panic(err)
	}
	err = cRun.Sync()
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

// READ RUN FILES
// TODO: clean up

func _readCRun(path string) {
	build, err := os.Create(path + "build")
	if err != nil {
		panic(err)
	}
	run, err := os.Create(path + "run")
	if err != nil {
		panic(err)
	}
	defer build.Close()
	defer run.Close()

	rb, err := SFTPCLIENT.Open(SFTPPATH + C_ + BUILD)
	if err != nil {
		panic(err)
	}
	rr, err := SFTPCLIENT.Open(SFTPPATH + C_ + RUN)
	if err != nil {
		panic(err)
	}
	defer rb.Close()
	defer rr.Close()

	_, err = build.ReadFrom(rb)
	if err != nil {
		panic(err)
	}
	_, err = run.ReadFrom(rr)
	if err != nil {
		panic(err)
	}
}

func _readGoRun(path string) {
	build, err := os.Create(path + "build")
	if err != nil {
		panic(err)
	}
	run, err := os.Create(path + "run")
	if err != nil {
		panic(err)
	}
	defer build.Close()
	defer run.Close()

	rb, err := SFTPCLIENT.Open(SFTPPATH + GO_ + BUILD)
	if err != nil {
		panic(err)
	}
	rr, err := SFTPCLIENT.Open(SFTPPATH + GO_ + RUN)
	if err != nil {
		panic(err)
	}
	defer rb.Close()
	defer rr.Close()

	_, err = build.ReadFrom(rb)
	if err != nil {
		panic(err)
	}
	_, err = run.ReadFrom(rr)
	if err != nil {
		panic(err)
	}
}

// ASD

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

func combineDeps(dep string, flake string, end string) string {
	return fmt.Sprintf("%s\n%s%s", flake, dep, end)
}
