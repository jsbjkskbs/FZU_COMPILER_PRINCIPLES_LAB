package utils

import "app/utils/log"

func Divider() log.Argument {
	return log.Argument{
		Highlight: true,
		Format:    "====================\n",
	}
}
