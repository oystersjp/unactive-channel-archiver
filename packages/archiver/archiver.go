package archiver

import (
	"fmt"
)

func Say() {
	fmt.Println("hello hoge")
}

type Archiver struct {
	c Config
}

type Config struct {
	Token            string
	ExcludeChannels  []string
	ArchiveLimitDays int
}

type Channel struct {
	ID       int
	Name     string
	IsMember bool
	Created  string
}

func NewArchiver(t string, chs []string, d int) *Archiver {
	return &Archiver{
		c: Config{
			Token:            t,
			ExcludeChannels:  chs,
			ArchiveLimitDays: d,
		},
	}
}

func (a Archiver) Exec() {
	// todo main logic
}

func (a Archiver) GetChannels() *[]Channel {
	// todo
	return &[]Channel{
		{ID: 0, Name: "_zoe", IsMember: true, Created: "2016/01/01 00:00:00"},
		{ID: 1, Name: "general", IsMember: false, Created: "2016/01/01 00:00:00"},
	}
}

func (c *Channel) Join() {
	// todo
}

func (c *Channel) DecideArchive() bool {
	// todo
	return true
}

func (c *Channel) Archive() {
	// todo
}
