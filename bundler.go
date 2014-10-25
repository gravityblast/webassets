package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
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

func (b *bundler) bundleFiles(outName string, srcNames []string, mf minifyFunc, ext string) {
	var sourceFiles []io.Reader
	buf := bytes.NewBuffer([]byte{})

	for _, srcName := range srcNames {
		srcPath := b.buildSourcePath(srcName, ext)
		sourceFile, err := os.Open(srcPath)
		if err != nil {
			log.Fatal(err)
		}
		defer sourceFile.Close()

		sourceFiles = append(sourceFiles, sourceFile)
	}

	c := newMinifier(mf, sourceFiles, buf)
	c.minify()

	content, err := ioutil.ReadAll(buf)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	h := md5.New()
	h.Write(content)
	outName = fmt.Sprintf("%s-%x", outName, h.Sum(nil))
	outPath := b.buildOuputPath(outName, ext)

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	w, err := outFile.Write(content)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	logger.Printf("writtent %d bytes to %s", w, outPath)
}

func (b *bundler) bundleJavascripts() {
	var wg sync.WaitGroup
	for outName, srcNames := range b.config.Javascripts {
		wg.Add(1)
		go func(outName string, srcNames []string) {
			b.bundleFiles(outName, srcNames, jsMinify, "js")
			wg.Done()
		}(outName, srcNames)
	}
	wg.Wait()
}
