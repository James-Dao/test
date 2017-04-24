package dao

import (
	"bufio"
	"config"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"model"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type CommandService struct {
	conf           *config.Config
	bufferDuration time.Duration
	lastWrite      time.Time
	lock           sync.Mutex
	readyToFlush   func() bool
	points         []*model.Command
}

func NewCommandService(conf *config.Config) *CommandService {
	i := &CommandService{
		conf:           conf,
		bufferDuration: time.Duration(60 * time.Second),
		lastWrite:      time.Now(),
		points:         make([]*model.Command, 0),
	}
	i.readyToFlush = i.defaultReadyToFlush
	return i
}

func (i *CommandService) Run() error {
	cmd := "sysdig -pc -c /gopath/app/bin/containercommand | awk '{for(i=1;i<=NF;i++) printf\"%s \",$i} {print \"\"}'"
	input := exec.Command("/bin/sh", "-c", cmd)
	input_pipe, err := input.StdoutPipe()
	if err != nil {
		fmt.Errorf(fmt.Sprintf("Run  input_pipe Error: [ %s ] ", err))
		return err
	}
	err = input.Start()
	if err != nil {
		fmt.Errorf(fmt.Sprintf("Run  input.Start(): [ %s ] ", err))
		return err
	}
	input_reader := bufio.NewReader(input_pipe)
	scanner := bufio.NewScanner(input_reader)
	for scanner.Scan() {
		line := string(scanner.Bytes())
		columes := strings.Split(line, " ")
		pid := columes[0]
		commandtime := columes[1]
		userandcontainername := columes[2]
		commands := columes[3:]
		command_string := strings.Join(commands, " ")
		command := model.NewCommand(pid, commandtime, userandcontainername, command_string)
		log.Info(fmt.Sprintf("command:%+v", command))
		if strings.Contains(userandcontainername, "@host") {
			command.Level = "host"
		} else {
			command.Level = "container"
			i.points = append(i.points, command)
		}
		if i.readyToFlush() {
			pointsToFlush := i.points
			i.points = make([]*model.Command, 0)
			i.lastWrite = time.Now()
			if len(pointsToFlush) > 0 {
				fmt.Println(fmt.Sprintf("= start to flush point to storage, point size= [ %d ]", len(pointsToFlush)))
				points := make([]model.Command, len(pointsToFlush))
				for i, p := range pointsToFlush {
					points[i] = *p
				}
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		log.Error(fmt.Sprintf("read error"), err)
	}
	log.Info("Exit Scanner")
	return nil
}

func (i *CommandService) defaultReadyToFlush() bool {
	return time.Since(i.lastWrite) >= i.bufferDuration
}
