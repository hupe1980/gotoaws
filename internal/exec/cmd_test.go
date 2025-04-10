package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	t.Run("should delegate to exec and call Run", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cmd := &Cmd{
			command: func(name string, _ []string, _ ...CmdOption) cmdRunner {
				assert.Equal(t, "date", name)

				m := NewMockcmdRunner(ctrl)
				m.EXPECT().Run().Return(nil)

				return m
			},
		}

		err := cmd.Run("date", nil)

		assert.NoError(t, err)
	})
}
