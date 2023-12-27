package win_process

import (
	"github.com/pkg/errors"
	"os"
)

func KillProcessByID(id int) error {
	if p, pErr := os.FindProcess(id); pErr != nil {
		return errors.Wrapf(pErr, "fail to find pid %d", id)
	} else {
		if err := p.Kill(); err != nil {
			return errors.Wrapf(err, "fail to kill pid %d", id)
		}
	}
	return nil
}
