package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	stringContent := "string to log"

	t.Run("basic", func(t *testing.T) {
		f, err := os.CreateTemp("", "logfile")
		if err != nil {
			return
		}
		defer f.Close()
		defer os.Remove(f.Name())
		logg := New("Info", f, os.Stdout)
		logg.Debug(stringContent)
		logg.Error(stringContent)
		logg.Info(stringContent)

		require.FileExists(t, f.Name())

		nf, err := os.Open(f.Name())
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)
		expectedString := `{"level":"DEBUG","message":"string to log"}` +
			"\n" + `{"level":"ERROR","message":"string to log"}` +
			"\n" + `{"level":"INFO","message":"string to log"}` +
			"\n"

		buf := make([]byte, len(expectedString))
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(expectedString), buf)
		os.Remove(f.Name())
	})

	t.Run("only error logs because level", func(t *testing.T) {
		f, err := os.CreateTemp("", "logfile")
		if err != nil {
			return
		}
		defer f.Close()
		defer os.Remove(f.Name())
		logg := New("error", f, os.Stdout)
		logg.Info(stringContent)
		logg.Debug(stringContent)
		logg.Warn(stringContent)
		logg.Error(stringContent)

		require.FileExists(t, f.Name())

		nf, err := os.Open(f.Name())
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)
		expectedString := "{\"level\":\"ERROR\",\"message\":\"string to log\"}\n"

		buf := make([]byte, len(expectedString))
		_, err = nf.Read(buf)

		require.NoError(t, err)
		require.Equal(t, []byte(expectedString), buf)
		os.Remove(f.Name())
	})

	t.Run("multiwriter", func(t *testing.T) {
		f, err := os.CreateTemp("", "logfile")
		if err != nil {
			return
		}
		defer f.Close()
		defer os.Remove(f.Name())

		f2, err := os.CreateTemp("", "logfile2")
		if err != nil {
			return
		}
		defer f2.Close()
		defer os.Remove(f2.Name())

		logg := New("Info", f, f2)
		logg.Debug(stringContent)
		logg.Error(stringContent)
		logg.Info(stringContent)

		require.FileExists(t, f.Name())
		require.FileExists(t, f2.Name())

		nf, err := os.Open(f.Name())
		if err != nil {
			return
		}
		defer nf.Close()
		require.NoError(t, err)

		nf2, err := os.Open(f2.Name())
		if err != nil {
			return
		}
		defer nf2.Close()
		require.NoError(t, err)

		expectedString := `{"level":"DEBUG","message":"string to log"}` +
			"\n" + `{"level":"ERROR","message":"string to log"}` +
			"\n" + `{"level":"INFO","message":"string to log"}` +
			"\n"
		buf := make([]byte, len(expectedString))
		buf2 := make([]byte, len(expectedString))
		_, err = nf.Read(buf)
		require.NoError(t, err)
		_, err = nf2.Read(buf2)
		require.NoError(t, err)

		require.Equal(t, []byte(expectedString), buf)
		require.Equal(t, []byte(expectedString), buf2)
		os.Remove(f.Name())
		os.Remove(f2.Name())
	})
}
