package internal

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func PrintInfo(a ...interface{}) {
	silent := viper.GetBool("silent")
	if !silent {
		fmt.Fprintf(os.Stdout, "%s %s\n", promptui.IconGood, fmt.Sprint(a...))
	}
}

func PrintInfof(format string, a ...interface{}) {
	silent := viper.GetBool("silent")
	if !silent {
		fmt.Fprintf(os.Stdout, "%s %s\n", promptui.IconGood, fmt.Sprintf(format, a...))
	}
}

func PrintError(a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", promptui.IconBad, fmt.Sprint(a...))
}

func PrintErrorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", promptui.IconBad, fmt.Sprintf(format, a...))
}
