// +build windows

package common

import "github.com/rs/zerolog/log"

func SetUlimitMax() {
	log.Trace().Msgf("ignore ulimit on 386 os.")
}
