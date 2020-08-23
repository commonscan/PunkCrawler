package common

import (
	"github.com/rs/zerolog/log"
	"syscall"
)

func SetUlimitMax() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Warn().Msg("Error Getting Rlimit " + err.Error())
		return
	}
	rLimit.Cur = rLimit.Max
	if rLimit.Cur != rLimit.Max {
		log.Trace().Msgf("current rlimit: %d , max %d", rLimit.Cur, rLimit.Max)
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			log.Warn().Msg("Error Setting Rlimit" + err.Error())
		}
		log.Trace().Msgf("after update rlimit: %d , max %d", rLimit.Cur, rLimit.Max)
	}

	//err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	//log.Info().Msgf("set limit to %d", rLimit.Cur)
}
