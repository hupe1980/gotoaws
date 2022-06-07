package exec

import (
	"os"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractiveRun(t *testing.T) {
	t.Run("should enable default stdin, stdout, and stderr before running the command", func(t *testing.T) {
		// GIVEN
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := &Cmd{
			command: func(name string, args []string, opts ...CmdOption) cmdRunner {
				assert.Equal(t, "date", name)

				// Make sure that the options applied are what we expect.
				cmd := &exec.Cmd{}
				for _, opt := range opts {
					opt(cmd)
				}

				assert.Equal(t, os.Stdin, cmd.Stdin)
				assert.Equal(t, os.Stdout, cmd.Stdout)
				assert.Equal(t, os.Stderr, cmd.Stderr)

				m := NewMockcmdRunner(ctrl)
				m.EXPECT().Run().Return(nil)

				return m
			},
		}

		// WHEN
		err := cmd.InteractiveRun("date")

		// THEN
		assert.NoError(t, err)
	})
}
