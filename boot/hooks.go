package boot

import (
	"fmt"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

func (b *Boot) DispatchBootNode(genesisJSON, publicKey, privateKey string) error {
	return b.dispatch("boot", []string{
		genesisJSON,
		publicKey,
		privateKey,
	}, nil)
}

// dispatch to both exec calls, and remote web hooks.
func (b *Boot) dispatch(hookName string, args []string, f func() error) error {
	zlog.Info("BEGIN", zap.String("hook", hookName))

	executable := fmt.Sprintf("./%s.sh", hookName)

	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	//fmt.Printf("  Executing hook: %q\n", cmd.Args)

	err := cmd.Run()
	if err != nil {
		return err
	}

	zlog.Info("ENDED", zap.String("hook", hookName))

	return nil
}
