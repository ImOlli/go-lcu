package lcu

import (
	"errors"
	"github.com/shirou/gopsutil/process"
	"strings"
)

const processName = "LeagueClientUx.exe"
const processNotFoundErrorMessage = "no process with " + processName + " found"

// ConnectInfo represents information to connect to the lcu.
// Please keep in mind that these information can change when the LCU is restarted or closed.
type ConnectInfo struct {
	// Indicates the port of the LCU
	Port string
	// Indicates the authentication token of the LCU
	AuthToken string
}

// IsProcessNotFoundError checks if the specified error is occurred because the LeagueClient process
// could not be found.
func IsProcessNotFoundError(err error) bool {
	return err.Error() == processNotFoundErrorMessage
}

// FindLCUConnectInfo searches over all processes to find the LeagueClient process (LeagueClientUx.exe)
// Then it checks the full commandline of the process to find the program arguments for the port and the auth-token
//
// If the process could not be found a dedicated error is thrown. You can check this error with the function IsProcessNotFoundError
func FindLCUConnectInfo() (*ConnectInfo, error) {
	var connectInfo = &ConnectInfo{}
	var processFound = false

	processList, err := process.Processes()

	if err != nil {
		return nil, err
	}

	for _, p := range processList {
		exe, err := p.Exe()

		if err != nil {
			continue
		}

		if strings.Contains(exe, processName) {
			processFound = true
			cmdLines, err := p.CmdlineSlice()
			if err != nil {
				continue
			}

			for _, line := range cmdLines {
				if strings.HasPrefix(line, "\"--remoting-auth-token=") {
					token := strings.TrimPrefix(line, "\"--remoting-auth-token=")
					token = strings.TrimSuffix(token, "\"")

					connectInfo.AuthToken = token
				} else if strings.HasPrefix(line, "\"--app-port=") {
					port := strings.TrimPrefix(line, "\"--app-port=")
					port = strings.TrimSuffix(port, "\"")

					connectInfo.Port = port
				}
			}

			// We found the right process, so we can end here
			break
		}
	}

	// Create error if we haven't found the right process
	if !processFound {
		return nil, errors.New(processNotFoundErrorMessage)
	}

	return connectInfo, nil
}
