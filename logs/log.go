package logs

import (
	"fmt"
	"time"

	"github.com/reverbdotcom/sbx/commit"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://app.datadoghq.com/logs?query=version:1.0.0-sha-%v&from_ts=%v&live=true"

var openURL = open.Open

func Run() (string, error) {
	err := openURL(Url())

	if err != nil {
		return "", err
	}

	return "", nil
}

var headSHA = commit.HeadSHA

func Url() string {
	sha, _ := headSHA()
	return fmt.Sprintf(template, sha, unixOneHourAgo())
}

var now = time.Now

func unixOneHourAgo() int64 {
	oneHourAgo := now().Add(-1 * time.Hour)

	return oneHourAgo.UnixMilli()
}
