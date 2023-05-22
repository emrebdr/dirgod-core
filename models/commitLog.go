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
	return fmt.Sprintf("%s %s %s %s %d\n", c.PreviousCommitId, c.CommitId, c.Committer, c.Message, c.Time)
}

func (c *CommitLog) Deserialize(input string) (*CommitLog, error) {
	data := strings.Split(input, " ")
	c.PreviousCommitId = data[0]
	c.CommitId = data[1]
	c.Committer = data[2] + " " + data[3]
	c.Message = data[4]

	time, err := strconv.ParseInt(data[5], 10, 64)
	if err != nil {
		return nil, err
	}
	c.Time = time

	return c, nil
}
