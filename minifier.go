package main

import (
	"bytes"
	"io"

	"github.com/web-assets/go-jsmin"
)

type minifyFunc func(io.Reader) ([]byte, error)

func jsMinify(r io.Reader) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	err := jsmin.Min(r, b)

	return b.Bytes(), err
}

type minifier struct {
	mf      minifyFunc
	sources []io.Reader
	out     io.Writer
}

func newMinifier(mf minifyFunc, sources []io.Reader, out io.Writer) *minifier {
	return &minifier{
		mf:      mf,
		sources: sources,
		out:     out,
	}
}

func (m *minifier) minify() error {
	for _, source := range m.sources {
		b, err := m.mf(source)
		if err != nil {
			return err
		}

		m.out.Write(b)
	}

	return nil
}
