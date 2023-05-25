package pipe

import (
	"context"
	"io"
	"log"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

type Command struct {
	Name string
	Args []string
}

// Pipe запускает комманды параллельно, передавая STDOUT процесса на STDIN
// следуещего по порядку процесса.
func Pipe(ctx context.Context, r io.Reader, commands ...Command) io.Reader {
	if len(commands) == 0 {
		return r
	}
	errg, ctx := errgroup.WithContext(ctx)

	cmds := make([]*exec.Cmd, 0, len(commands))
	for _, cmd := range commands {
		osCmd := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
		cmds = append(cmds, osCmd)
	}
	cmds[0].Stdin = r
	var pr io.ReadCloser
	for i := 0; i < len(cmds); i++ {
		var pw *io.PipeWriter
		pr, pw = io.Pipe()

		cmd := cmds[i]

		cmd.Stdout = pw
		if i != len(cmds)-1 {
			cmds[i+1].Stdin = pr
		}

		errg.Go(func() error {
			err := cmd.Run()
			log.Printf("Program [%v] exited: err=%v", cmd, err)
			pw.CloseWithError(err)
			return err
		})
	}

	return pr
}
