package models

import (
	"fmt"
	"strconv"
	"strings"
)

type CommitLog struct {
	PreviousCommitId string
	CommitId         string
	Committer        string
	Message          string
	Time             int64
}

func (c *CommitLog) Serialize() string {
	return fmt.Sprintf("%s %s %s %d %s\n", c.PreviousCommitId, c.CommitId, c.Committer, c.Time, c.Message)
}

func (c *CommitLog) Deserialize(input string) (*CommitLog, error) {
	data := strings.Split(input, " ")
	c.PreviousCommitId = data[0]
	c.CommitId = data[1]
	c.Committer = data[2] + " " + data[3]

	time, err := strconv.ParseInt(data[4], 10, 64)
	if err != nil {
		return nil, err
	}
	c.Time = time

	message := ""
	if len(data) >= 6 {
		for i := 5; i < len(data); i++ {
			message += data[i] + " "
		}
	}

	message = strings.TrimSpace(message)
	c.Message = message

	return c, nil
}
