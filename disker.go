package pegger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Disker has the configuration for disker and it's Run method.
type Disker struct {
	Dir             string `help:"Directory in which one file per goroutine will be created."`
	Concurrency     int    `help:"Number of concurrent goroutines writing files."`
	FileSizeBytes   int    `help:"Size of files to write."`
	FileSizeBytesP2 uint   `help:"Power of two to set the file size to. Ignored if FileSizeBytes is set."`
	Iterations      int    `help:"Number of times each goroutine will overwrite its file."`
}

// NewDisker gets a new Disker with default values.
func NewDisker() *Disker {
	return &Disker{
		Iterations:      100,
		Concurrency:     runtime.NumCPU(),
		FileSizeBytesP2: 20,
	}
}

// Run runs
func (m *Disker) Run() (err error) {
	if m.Dir == "" {
		m.Dir, err = ioutil.TempDir("", "disker")
		if err != nil {
			return errors.Wrap(err, "getting temp dir")
		}
	}
	if m.FileSizeBytesP2 > 63 {
		return errors.Errorf("FileSizeBytesP2 cannot be greater than 63, but is set to %d", m.FileSizeBytesP2)
	}
	if m.FileSizeBytes == 0 {
		m.FileSizeBytes = 1 << m.FileSizeBytesP2
	}
	fmt.Printf("%#v\n", m)
	eg := &errgroup.Group{}
	start := time.Now()
	for i := 0; i < m.Concurrency; i++ {
		i := i
		eg.Go(func() error {
			for j := 0; j < m.Iterations; j++ {
				f, err := os.Create(filepath.Join(m.Dir, strconv.Itoa(i)))
				if err != nil {
					return errors.Wrap(err, "creating file")
				}
				nr := &nopReader{length: m.FileSizeBytes}
				_, err = io.Copy(f, nr)
				if err != nil {
					return errors.Wrap(err, "copying data")
				}
				err = f.Sync()
				if err != nil {
					return errors.Wrap(err, "syncing file")
				}
				err = f.Close()
				if err != nil {
					return errors.Wrap(err, "closing file")
				}
			}
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "err in writing routines")
	}
	duration := time.Since(start)
	totalSizeMB := float64(m.Concurrency*m.FileSizeBytes*m.Iterations) / 1024 / 1024
	rateMB := totalSizeMB / duration.Seconds()
	fmt.Printf("%#v\n", m)
	fmt.Printf("Wrote: %.0f MB in %v. %.0f MB/s\n", totalSizeMB, duration, rateMB)
	return m.cleanup()
}

func (m *Disker) cleanup() error {
	err := os.RemoveAll(m.Dir)
	return errors.Wrap(err, "removing directory")
}

// nopReader is an io.Reader which returns the slice given to its read method
// unchanged up to a certain number of bytes.
type nopReader struct {
	length int
}

func (nr *nopReader) Read(b []byte) (int, error) {
	if len(b) < nr.length {
		nr.length -= len(b)
		return len(b), nil
	}
	remaining := nr.length
	nr.length = 0
	return remaining, io.EOF
}
