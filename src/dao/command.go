package dao

import (
	"bufio"
	"config"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"model"
	"os/exec"
	"strings"
)

type CommandService struct {
	conf   *config.Config
	points []*model.Command
}

func NewCommandService(conf *config.Config) *CommandService {
	i := &CommandService{
		conf:   conf,
		points: make([]*model.Command, 0),
	}
	return i
}

func (i *CommandService) Run() error {
	log.Infof("%s", "CommandService Run")
	//cmd := "sysdig -pc -c /gopath/app/bin/containercommand | awk '{for(i=1;i<=NF;i++) printf\"%s \",$i} {print \"\"}'"
	cmd := "sysdig -pc -c /gopath/app/bin/all"
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
	// input_reader := bufio.NewReaderSize(input_pipe, 10)
	scanner := bufio.NewScanner(input_reader)
	for scanner.Scan() {
		line := string(scanner.Bytes())
		log.Infof("line %s", line)
		columes := strings.Split(line, " ")
		pid := columes[0]
		commandtime := columes[1]
		userandcontainername := columes[2]
		commands := columes[3:]
		command_string := strings.Join(commands, " ")
		command := model.NewCommand(pid, commandtime, userandcontainername, command_string)
		// log.Info(fmt.Sprintf("command:%+v", command))
		if strings.Contains(userandcontainername, "@host") {
			command.Level = "host"
		} else {
			command.Level = "container"
			i.points = append(i.points, command)
		}
	}
	err = scanner.Err()
	if err != nil {
		log.Error(fmt.Sprintf("read error"), err)
	}
	log.Info("Exit Scanner")
	return nil
}
