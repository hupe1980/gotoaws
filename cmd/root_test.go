package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmdHelp(t *testing.T) {
	var b bytes.Buffer
	cmd := newRootCmd("")
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-h"})
	require.NoError(t, cmd.Execute())
}

func TestRootCmdVersion(t *testing.T) {
	var b bytes.Buffer
	cmd := newRootCmd("1.2.3")
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-v"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, "awsconnect version 1.2.3\n", b.String())
}
