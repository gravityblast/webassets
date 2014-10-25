package main

import (
	"flag"
	"os"
)

func main() {
	var configPath string
	var sourcePaths stringslice
	var destPath string

	flag.BoolVar(&logger.verbose, "v", false, "verbose")
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Var(&sourcePaths, "src", "source paths")
	flag.StringVar(&destPath, "dest", "", "output path")
	flag.Parse()

	if configPath == "" || destPath == "" || len(sourcePaths) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	config, err := parseConfigFile(configPath)
	if err != nil {
		logger.Fatalf("%v\n", err)
	}

	logger.Printf("using config: %s\n", configPath)
	logger.Printf("sources: %s\n", &sourcePaths)
	logger.Printf("dest: %s\n", destPath)

	b := newBundler(config, sourcePaths, destPath)
	b.bundle()
}
