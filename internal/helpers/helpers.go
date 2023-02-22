package helpers

import (
	"bufio"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"os"
	"strings"
)

func Read() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "failed to read from console")
	}

	text = strings.ReplaceAll(text, "\n", "")
	return text, nil
}
