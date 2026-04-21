package gatusctl

import (
	"fmt"
	"os/exec"
)

func Restart(container string) error {
	out, err := exec.Command("docker", "restart", container).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker restart failed: %v: %s", err, string(out))
	}
	return nil
}
