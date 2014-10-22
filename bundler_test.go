package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	assert "github.com/pilu/miniassert"
)

const (
	testDataFolder     = "__test-data__"
	testFixturesFolder = "__test-fixtures__"
)

func testNewBundler() *bundler {
	c := config{
		Javascripts: map[string][]string{
			"application": []string{"file-1", "file-2"},
			"admin":       []string{"admin-1", "admin-2", "editor/main"},
		},
	}
	sources := []string{
		fmt.Sprintf("%s/public", testFixturesFolder),
		fmt.Sprintf("%s/vendor", testFixturesFolder),
	}
	return newBundler(c, sources, fmt.Sprintf("%s/public/assets", testDataFolder))
}

func testRemoveTestDataFolder() {
	os.RemoveAll(testDataFolder)
}

func TestNewBundler(t *testing.T) {
	b := testNewBundler()
	assert.Equal(t, []string{fmt.Sprintf("%s/public", testFixturesFolder), fmt.Sprintf("%s/vendor", testFixturesFolder)}, b.sourcePaths)
	assert.Equal(t, fmt.Sprintf("%s/public/assets", testDataFolder), b.outputPath)
}

func TestBundler_InitOutputPath(t *testing.T) {
	defer testRemoveTestDataFolder()

	b := testNewBundler()

	_, err := os.Stat(testDataFolder)
	if err == nil {
		t.Errorf("%s followed should must be removed before runnning tests", testDataFolder)
	}

	b.initOutputPath()

	fi, err := os.Stat(b.outputPath)
	assert.Nil(t, err)
	assert.True(t, fi.IsDir())
}

func TestBundler_BuildOutputPath(t *testing.T) {
	b := testNewBundler()
	path := b.buildOuputPath("application", "js")
	expected := fmt.Sprintf("%s/application.js", b.outputPath)
	assert.Equal(t, expected, path)
}

func TestBundler_BuildSourcePath(t *testing.T) {
	defer testRemoveTestDataFolder()

	b := testNewBundler()
	path := b.buildSourcePath("file-1", "js")
	assert.Equal(t, fmt.Sprintf("%s/public/file-1.js", testFixturesFolder), path)

	path = b.buildSourcePath("admin-1", "js")
	assert.Equal(t, fmt.Sprintf("%s/vendor/admin-1.js", testFixturesFolder), path)

	path = b.buildSourcePath("editor/main", "js")
	assert.Equal(t, fmt.Sprintf("%s/vendor/editor/main.js", testFixturesFolder), path)
}

func TestBundler_Bundle(t *testing.T) {
	defer testRemoveTestDataFolder()
	b := testNewBundler()
	b.bundle()

	tests := map[string]string{
		"application.js": `
function file1(){return"file1"}
function file2(){return"file2"}`,
		"admin.js": `
function admin1(){return"admin1"}
function admin2(){return"admin2"}
function editor(){return"editor"}`,
	}

	for filename, expectedContent := range tests {
		f, err := os.Open(fmt.Sprintf("%s/%s", b.outputPath, filename))
		assert.Nil(t, err)

		content, err := ioutil.ReadAll(f)
		assert.Nil(t, err)

		assert.Equal(t, expectedContent, string(content))
	}
}
