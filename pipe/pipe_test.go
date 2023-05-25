package pipe

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	in := strings.NewReader("1\n2\n3\n")
	out := Pipe(context.Background(), in, Command{
		Name: "wc",
		Args: []string{
			"-l",
		},
	})

	b, err := io.ReadAll(out)
	require.NoError(t, err)

	assert.Equal(t, "       3\n", string(b))
}

func TestPipe_unknown_command(t *testing.T) {
	in := strings.NewReader("1\n2\n3\n")
	out := Pipe(context.Background(), in, Command{
		Name: "ccccccc",
		Args: []string{
			"-l",
		},
	})

	_, err := io.ReadAll(out)
	require.Error(t, err)
	t.Logf("err: %v", err)
}

func TestPipe_cancel(t *testing.T) {
	in := strings.NewReader("1\n2\n3\n")
	out := Pipe(context.Background(), in,
		Command{
			Name: "sleep",
			Args: []string{
				"100",
			},
		},
		Command{
			Name: "ccccccc",
			Args: []string{
				"-l",
			},
		},
		Command{
			Name: "wc",
			Args: []string{
				"-l",
			},
		},
	)

	_, err := io.ReadAll(out)
	require.Error(t, err)
	t.Logf("err: %v", err)
}
