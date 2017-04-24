package model

type TopContainersCpu struct {
	CPU           string
	ContainerName string
}

func NewTopContainersCpu(cpu, containername string) *TopContainersCpu {
	return &TopContainersCpu{
		CPU:           cpu,
		ContainerName: containername,
	}
}

type TopContainersMemory struct {
	CPU           string
	Virt          string
	Res           string
	ContainerName string
}

func NewTopContainersMemory(cpu, virt, res, containername string) *TopContainersMemory {
	return &TopContainersMemory{
		CPU:           cpu,
		Virt:          virt,
		Res:           res,
		ContainerName: containername,
	}
}

type TopContainersProcess struct {
	CPU           string
	Process       string
	User          string
	ThreadCount   string
	Host_Pid      string
	Virt          string
	Res           string
	CpuNo         string
	ContainerId   string
	ContainerName string
}

func NewTopContainersProcess(cpu, process, user, threadcount, hostid, virt, res, cpuno, containerid, containername string) *TopContainersProcess {
	return &TopContainersProcess{
		CPU:           cpu,
		Process:       process,
		User:          user,
		ThreadCount:   threadcount,
		Host_Pid:      hostid,
		Virt:          virt,
		Res:           res,
		CpuNo:         cpuno,
		ContainerId:   containerid,
		ContainerName: containername,
	}
}

type TopContainersThread struct {
	CPU           string
	Process       string
	User          string
	ThreadCount   string
	Host_Pid      string
	Thread_Id     string
	Container_Pid string
	Virt          string
	Res           string
	CpuNo         string
	ContainerId   string
	ContainerName string
}

func NewTopContainersThread(cpu, process, user, threadcount, hostid, threadid, containerpid, virt, res, cpuno, containerid, containername string) *TopContainersThread {
	return &TopContainersThread{
		CPU:           cpu,
		Process:       process,
		User:          user,
		ThreadCount:   threadcount,
		Host_Pid:      hostid,
		Thread_Id:     threadid,
		Container_Pid: containerpid,
		Virt:          virt,
		Res:           res,
		CpuNo:         cpuno,
		ContainerId:   containerid,
		ContainerName: containername,
	}
}

type TopContainersNet struct {
	Bytes         string
	Connection    string
	Process       string
	IOPS          string
	ContainerName string
}

func NewTopContainersNet(bytes, connection, process, iops, containername string) *TopContainersNet {
	return &TopContainersNet{
		Bytes:         bytes,
		Connection:    connection,
		Process:       process,
		IOPS:          iops,
		ContainerName: containername,
	}
}

type TopContainersFile struct {
	Bytes         string
	Process       string
	ContainerName string
}

func NewTopContainersFile(bytes, process, containername string) *TopContainersFile {
	return &TopContainersFile{
		Bytes:         bytes,
		Process:       process,
		ContainerName: containername,
	}
}

type TopContainersError struct {
	ErrorCount    string
	Process       string
	ContainerName string
}

func NewTopContainersError(errorcount, process, containername string) *TopContainersError {
	return &TopContainersError{
		ErrorCount:    errorcount,
		Process:       process,
		ContainerName: containername,
	}
}

type TopContainersFDCount struct {
	Open          string
	ProcessName   string
	Max           string
	PCT           string
	ContainerName string
}

func NewTopContainersFDCount(open, processname, max, pct, containername string) *TopContainersFDCount {
	return &TopContainersFDCount{
		Open:          open,
		ProcessName:   processname,
		Max:           max,
		PCT:           pct,
		ContainerName: containername,
	}
}

type TopHttpRequest struct {
	Bytes         string
	Spend_Time    string
	Host_Pid      string
	Connection    string
	HttpMethod    string
	Url           string
	StatusCode    string
	ContainerId   string
	ContainerName string
}

func NewTopHttpRequest(bytes, spendtime, containerid, containername, hostpid, connection, statuscode, httpmethod, url string) *TopHttpRequest {
	return &TopHttpRequest{
		Bytes:         bytes,
		Spend_Time:    spendtime,
		Host_Pid:      hostpid,
		Connection:    connection,
		HttpMethod:    httpmethod,
		Url:           url,
		StatusCode:    statuscode,
		ContainerId:   containerid,
		ContainerName: containername,
	}
}

type NetSlower struct {
	ContainerId   string
	ContainerName string
	Process       string
	EventType     string
	Latency       string
	Connection    string
}

func NewNetSlower(containerid, containername, process, eventtype, latency, connection string) *NetSlower {
	return &NetSlower{
		ContainerId:   containerid,
		ContainerName: containername,
		Process:       process,
		EventType:     eventtype,
		Latency:       latency,
		Connection:    connection,
	}
}

type Mysql struct {
	ContainerId   string
	ContainerName string
	Method        string
	Latency       string
	Content       string
}

func NewMysql(containerid, containername, method, latency, content string) *Mysql {
	return &Mysql{
		ContainerId:   containerid,
		ContainerName: containername,
		Method:        method,
		Latency:       latency,
		Content:       content,
	}
}

type Redis struct {
	ContainerId   string
	ContainerName string
	Method        string
	Latency       string
	Content       string
}

func NewRedis(containerid, containername, method, latency, content string) *Redis {
	return &Redis{
		ContainerId:   containerid,
		ContainerName: containername,
		Method:        method,
		Latency:       latency,
		Content:       content,
	}
}

type Command struct {
	Pid                  string
	CommandTime          string
	UserAndContainername string
	Command              string
	Level                string
}

func NewCommand(pid, commandTime, userAndContainername, command string) *Command {
	return &Command{
		Pid:                  pid,
		CommandTime:          commandTime,
		UserAndContainername: userAndContainername,
		Command:              command,
	}
}
