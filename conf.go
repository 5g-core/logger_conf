package logger_conf

import (
	"log"
	"os"
	"strconv"

	"github.com/5g-core/path_util"
)

var N5GCLogDir string = path_util.N5GCPath("n5gc/log") + "/"
var LibLogDir string = N5GCLogDir + "lib/"
var NfLogDir string = N5GCLogDir + "nf/"

var N5GCLogfle string = N5GCLogDir + "n5gc.log"

func init() {
	if err := os.MkdirAll(LibLogDir, 0775); err != nil {
		log.Printf("Mkdir %s failed: %+v", LibLogDir, err)
	}
	if err := os.MkdirAll(NfLogDir, 0775); err != nil {
		log.Printf("Mkdir %s failed: %+v", NfLogDir, err)
	}

	// Create log file or if it already exist, check if user can access it
	f, fileOpenErr := os.OpenFile(N5GCLogfle, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if fileOpenErr != nil {
		// user cannot access it.
		log.Printf("Cannot Open %s\n", N5GCLogfle)
	} else {
		// user can access it
		if err := f.Close(); err != nil {
			log.Printf("File %s cannot been closed\n", N5GCLogfle)
		}
	}

	sudoUID, errUID := strconv.Atoi(os.Getenv("SUDO_UID"))
	sudoGID, errGID := strconv.Atoi(os.Getenv("SUDO_GID"))

	if errUID == nil && errGID == nil {
		// if using sudo to run the program, errUID will be nil and sudoUID will get the uid who run sudo
		// else errUID will not be nil and sudoUID will be nil
		// If user using sudo to run the program and create log file, log will own by root,
		// here we change own to user so user can view and reuse the file
		if err := os.Chown(N5GCLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", N5GCLogDir, sudoUID, sudoGID, err)
		}
		if err := os.Chown(LibLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", LibLogDir, sudoUID, sudoGID, err)
		}
		if err := os.Chown(NfLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", NfLogDir, sudoUID, sudoGID, err)
		}

		if fileOpenErr == nil {
			if err := os.Chown(N5GCLogfle, sudoUID, sudoGID); err != nil {
				log.Printf("File %s chown to %d:%d error: %v\n", N5GCLogfle, sudoUID, sudoGID, err)
			}
		}
	}
}
