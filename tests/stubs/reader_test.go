package stubs

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderForTextByBlocks(t *testing.T) {
	text := "Hello, world!"

	reader := ReaderForText(text)

	firstEight := make([]byte, 8)
	n, err := reader.Read(firstEight)
	require.NoError(t, err)
	require.Equal(t, 8, n)
	require.Equal(t, []byte(`Hello, w`), firstEight)

	secondEight := make([]byte, 8)
	n, err = reader.Read(secondEight)
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 5, n)
	require.Equal(t, []byte("orld!\x00\x00\x00"), secondEight)

	thirdEight := make([]byte, 8)
	n, err = reader.Read(thirdEight)
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)
}

func TestReaderForTextByReadAll(t *testing.T) {
	text := "Hello, world!"

	reader := ReaderForText(text)

	read, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, []byte(text), read)
}

func TestReaderForTextByReadAllBig(t *testing.T) {
	text := strings.Repeat("Hello, world!", 128)

	reader := ReaderForText(text)

	read, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, []byte(text), read)
}
