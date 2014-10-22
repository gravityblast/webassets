package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type bundler struct {
	config      config
	sourcePaths []string
	outputPath  string
}

func newBundler(config config, sourcePaths []string, outputPath string) *bundler {
	return &bundler{
		config:      config,
		sourcePaths: sourcePaths,
		outputPath:  outputPath,
	}
}

func (b *bundler) bundle() {
	b.initOutputPath()
	b.bundleJavascripts()
}

func (b *bundler) initOutputPath() {
	err := os.MkdirAll(b.outputPath, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *bundler) buildOuputPath(path, extension string) string {
	return filepath.Join(b.outputPath, fmt.Sprintf("%s.%s", path, extension))
}

func (b *bundler) buildSourcePath(name, extension string) string {
	filename := fmt.Sprintf("%s.%s", name, extension)
	for _, path := range b.sourcePaths {
		fullpath := filepath.Join(path, filename)
		_, err := os.Stat(fullpath)
		if err == nil {
			return fullpath
		}
	}

	log.Fatalf("can't find asset %s", filename)
	return ""
}

func (b *bundler) bundleJavascripts() {
	var wg sync.WaitGroup
	for outName, srcNames := range b.config.Javascripts {
		wg.Add(1)
		go func(outName string, srcNames []string) {
			outPath := b.buildOuputPath(outName, "js")

			outFile, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()

			var sourceFiles []io.Reader

			for _, srcName := range srcNames {
				srcPath := b.buildSourcePath(srcName, "js")
				sourceFile, err := os.Open(srcPath)
				if err != nil {
					log.Fatal(err)
				}
				defer sourceFile.Close()

				sourceFiles = append(sourceFiles, sourceFile)
			}

			c := newMinifier(jsMinify, sourceFiles, outFile)
			c.minify()

			wg.Done()
		}(outName, srcNames)
	}
	wg.Wait()
}
