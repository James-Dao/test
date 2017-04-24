package dao

// import (
// 	"bufio"
// 	"config"
// 	"fmt"
// 	info "github.com/google/cadvisor/info/v1"
// 	"github.com/google/cadvisor/manager"
// 	_ "github.com/google/cadvisor/storage/bigquery"
// 	_ "github.com/google/cadvisor/storage/elasticsearch"
// 	_ "github.com/google/cadvisor/storage/influxdb"
// 	_ "github.com/google/cadvisor/storage/kafka"
// 	_ "github.com/google/cadvisor/storage/redis"
// 	_ "github.com/google/cadvisor/storage/statsd"
// 	_ "github.com/google/cadvisor/storage/stdout"
// 	influxdb "github.com/influxdb/influxdb/client"
// 	"logger"
// 	//"math"
// 	"common"
// 	"math"
// 	"model"
// 	"net/url"
// 	"os"
// 	"os/exec"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"
// )

// type InfluxdbService struct {
// 	conf             *config.Config
// 	logger           logger.Logger
// 	machineName      string
// 	bufferDuration   time.Duration
// 	lastWrite        time.Time
// 	points           []*influxdb.Point
// 	lock             sync.Mutex
// 	readyToFlush     func() bool
// 	containerManager manager.Manager
// 	machineInfo      *info.MachineInfo
// 	hostContainer    *info.ContainerInfo
// 	host_ip          string
// 	host_id          string
// }

// func NewInfluxdbService(conf *config.Config, logger logger.Logger, containerManager manager.Manager) *InfluxdbService {
// 	hostname, _ := os.Hostname()
// 	machineInfo, _ := containerManager.GetMachineInfo()
// 	dockerInfo, _ := containerManager.DockerInfo()
// 	host_ip := strings.Join(strings.Split(dockerInfo.Hostname, "-"), ".")
// 	host_id := common.GenerateUUID()
// 	hostcontainer, _ := containerManager.GetContainerInfo("/", &info.ContainerInfoRequest{NumStats: -1})
// 	i := &InfluxdbService{
// 		conf:             conf,
// 		logger:           logger,
// 		machineName:      hostname,
// 		bufferDuration:   time.Duration(60 * time.Second),
// 		lastWrite:        time.Now(),
// 		points:           make([]*influxdb.Point, 0),
// 		containerManager: containerManager,
// 		machineInfo:      machineInfo,
// 		hostContainer:    hostcontainer,
// 		host_ip:          host_ip,
// 		host_id:          host_id,
// 	}
// 	i.readyToFlush = i.defaultReadyToFlush
// 	return i
// }

// func (i *InfluxdbService) GetConnection() (*influxdb.Client, error) {
// 	url := &url.URL{
// 		Scheme: "http",
// 		Host:   i.conf.Influxdb_Host,
// 	}
// 	config := &influxdb.Config{
// 		URL:       *url,
// 		Username:  i.conf.Influxdb_Username,
// 		Password:  i.conf.Influxdb_Password,
// 		UserAgent: fmt.Sprintf("%v/%v", "sysdig", "1.0"),
// 	}
// 	client, err := influxdb.NewClient(*config)
// 	if err != nil {
// 		fmt.Errorf(fmt.Sprintf("MetricsService.GetConnection Error: [ %s ] ", err))
// 		return nil, err
// 	}
// 	return client, nil
// }

// func (i *InfluxdbService) Run() error {
// 	influxdb_client, err := i.GetConnection()
// 	if err != nil {
// 		fmt.Errorf(fmt.Sprintf("Run  MetricsService.GetConnection Error: [ %s ] ", err))
// 		return err
// 	}
// 	cmd := "sysdig -c all | awk '{for(i=1;i<=NF;i++) printf\"%s \",$i} {print \"\"}'"
// 	input := exec.Command("/bin/sh", "-c", cmd)
// 	input_pipe, err := input.StdoutPipe()
// 	if err != nil {
// 		fmt.Errorf(fmt.Sprintf("Run  input_pipe Error: [ %s ] ", err))
// 		return err
// 	}
// 	err = input.Start()
// 	if err != nil {
// 		fmt.Errorf(fmt.Sprintf("Run  input.Start(): [ %s ] ", err))
// 		return err
// 	}
// 	input_reader := bufio.NewReader(input_pipe)
// 	scanner := bufio.NewScanner(input_reader)
// 	for scanner.Scan() {
// 		line := string(scanner.Bytes())
// 		columes := strings.Split(line, " ")
// 		title := columes[0]
// 		if title == "TopContainersCpu" {
// 			topcontainerscpu := model.NewTopContainersCpu(columes[1], columes[2])
// 			tags_cpu := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "containername": topcontainerscpu.ContainerName}
// 			cpu, _ := strconv.ParseFloat(topcontainerscpu.CPU, 32)
// 			i.points = append(i.points, makePoint("ContainerCpu", tags_cpu, toSignedIfUnsigned(cpu)))
// 		} else if title == "TopContainersMemory" {
// 			processmemory := model.NewTopContainersMemory(columes[1], columes[2], columes[3], columes[4])
// 			tags_memory := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "containername": processmemory.ContainerName}
// 			memory_res, _ := strconv.ParseFloat(processmemory.Res, 32)
// 			i.points = append(i.points, makePoint("ContainerMemory", tags_memory, toSignedIfUnsigned(memory_res)))
// 		} else if title == "TopContainersProcess" {
// 			processcpu := model.NewTopContainersProcess(columes[1], columes[2], columes[3], columes[4], columes[5], columes[6], columes[7], columes[8], columes[9], columes[10])
// 			tags_cpu_memory := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": processcpu.Process, "hostpid": processcpu.Host_Pid, "containerid": processcpu.ContainerId, "containername": processcpu.ContainerName}
// 			cpu, _ := strconv.ParseFloat(processcpu.CPU, 32)
// 			i.points = append(i.points, makePoint("ProcessCpu", tags_cpu_memory, toSignedIfUnsigned(cpu)))
// 		} else if title == "TopProcessMemory" {
// 			processmemory := model.NewTopContainersProcess(columes[1], columes[2], columes[3], columes[4], columes[5], columes[6], columes[7], columes[8], columes[9], columes[10])
// 			tags_cpu_memory := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": processmemory.Process, "hostpid": processmemory.Host_Pid, "containerid": processmemory.ContainerId, "containername": processmemory.ContainerName}
// 			memory_res, _ := strconv.ParseFloat(processmemory.Res, 32)
// 			i.points = append(i.points, makePoint("ProcessMemory", tags_cpu_memory, toSignedIfUnsigned(memory_res)))
// 		} else if title == "TopContainersThread" {
// 			threadcpu := model.NewTopContainersThread(columes[1], columes[2], columes[3], columes[4], columes[5], columes[6], columes[7], columes[8], columes[9], columes[10], columes[11], columes[12])
// 			tags_cpu_memory := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": threadcpu.Process, "hostpid": threadcpu.Host_Pid, "threadid": threadcpu.Thread_Id, "containerpid": threadcpu.Container_Pid, "containerid": threadcpu.ContainerId, "containername": threadcpu.ContainerName}
// 			cpu, _ := strconv.ParseFloat(threadcpu.CPU, 32)
// 			memory_res, _ := strconv.ParseFloat(threadcpu.Res, 32)
// 			i.points = append(i.points, makePoint("ThreadCpu", tags_cpu_memory, toSignedIfUnsigned(cpu)))
// 			i.points = append(i.points, makePoint("ThreadMemory", tags_cpu_memory, toSignedIfUnsigned(memory_res)))
// 		} else if title == "TopContainersNet" {
// 			topcontainersnet := model.NewTopContainersNet(columes[1], columes[2], columes[3], columes[4], columes[5])
// 			tags_net := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": topcontainersnet.Process, "connection": topcontainersnet.Connection, "containername": topcontainersnet.ContainerName}
// 			bytes, _ := strconv.Atoi(topcontainersnet.Bytes)
// 			iops, _ := strconv.Atoi(topcontainersnet.IOPS)
// 			i.points = append(i.points, makePoint("NetBytes", tags_net, toSignedIfUnsigned(bytes)))
// 			i.points = append(i.points, makePoint("NetIops", tags_net, toSignedIfUnsigned(iops)))
// 		} else if title == "TopContainersFile" {
// 			topcontainersfile := model.NewTopContainersFile(columes[1], columes[2], columes[3])
// 			tags_file_bytes := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": topcontainersfile.Process, "containername": topcontainersfile.ContainerName}
// 			bytes, _ := strconv.Atoi(topcontainersfile.Bytes)
// 			i.points = append(i.points, makePoint("FileBytes", tags_file_bytes, toSignedIfUnsigned(bytes)))
// 		} else if title == "TopContainersErrors" {
// 			topcontainerserror := model.NewTopContainersError(columes[1], columes[2], columes[3])
// 			tags_errors := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "process": topcontainerserror.Process, "containername": topcontainerserror.ContainerName}
// 			errors_count, _ := strconv.Atoi(topcontainerserror.ErrorCount)
// 			i.points = append(i.points, makePoint("Errors", tags_errors, toSignedIfUnsigned(errors_count)))
// 		} else if title == "TopFdCount" {
// 			topcontainersfdcount := model.NewTopContainersFDCount(columes[1], columes[2], columes[3], columes[4], columes[5])
// 			tags_fdcount := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "processname": topcontainersfdcount.ProcessName, "containername": topcontainersfdcount.ContainerName}
// 			fd_count_open, _ := strconv.Atoi(topcontainersfdcount.Open)
// 			fd_count_percent, _ := strconv.ParseFloat(topcontainersfdcount.PCT, 32)
// 			fd_count_max, _ := strconv.Atoi(topcontainersfdcount.Max)
// 			i.points = append(i.points, makePoint("FdCountOpen", tags_fdcount, toSignedIfUnsigned(fd_count_open)))
// 			i.points = append(i.points, makePoint("FdCountPercent", tags_fdcount, toSignedIfUnsigned(fd_count_percent)))
// 			i.points = append(i.points, makePoint("FdCountMax", tags_fdcount, toSignedIfUnsigned(fd_count_max)))
// 		} else if title == "HttpRequest" {
// 			httprequest := model.NewTopHttpRequest(columes[1], columes[2], columes[3], columes[4], columes[5], columes[6], columes[7], columes[8], columes[9])
// 			tags_http_request := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "hostpid": httprequest.Host_Pid, "connection": httprequest.Connection, "method": httprequest.HttpMethod, "url": httprequest.Url, "statuscode": httprequest.StatusCode, "containerid": httprequest.ContainerId, "containername": httprequest.ContainerName}
// 			bytes, _ := strconv.Atoi(httprequest.Bytes)
// 			spent_time, _ := strconv.Atoi(httprequest.Spend_Time)
// 			i.points = append(i.points, makePoint("HttpRequestBytes", tags_http_request, toSignedIfUnsigned(bytes)))
// 			i.points = append(i.points, makePoint("HttpRequestSpent", tags_http_request, toSignedIfUnsigned(spent_time)))
// 			i.points = append(i.points, makePoint("HttpRequestCount", tags_http_request, toSignedIfUnsigned(1)))
// 		} else if title == "Mysql" {
// 			content := columes[5:]
// 			mysql := model.NewMysql(columes[1], columes[2], columes[3], columes[4], strings.Join(content, " "))
// 			tags_mysql := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "method": mysql.Method, "content": mysql.Content, "containerid": mysql.ContainerId, "containername": mysql.ContainerName}
// 			latency, _ := strconv.ParseFloat(mysql.Latency, 32)
// 			i.points = append(i.points, makePoint("Mysql", tags_mysql, toSignedIfUnsigned(latency)))
// 		} else if title == "Redis" {
// 			content := columes[5:]
// 			redis := model.NewRedis(columes[1], columes[2], columes[3], columes[4], strings.Join(content, " "))
// 			tags_redis := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "method": redis.Method, "content": redis.Content, "containerid": redis.ContainerId, "containername": redis.ContainerName}
// 			latency, _ := strconv.ParseFloat(redis.Latency, 32)
// 			i.points = append(i.points, makePoint("Redis", tags_redis, toSignedIfUnsigned(latency)))
// 		}

// 		if i.readyToFlush() {
// 			hostContainer, _ := i.containerManager.GetContainerInfo("/", &info.ContainerInfoRequest{NumStats: -1})
// 			//host cpu information
// 			if hostContainer.Spec.HasCpu && len(hostContainer.Stats) >= 2 {
// 				cur := hostContainer.Stats[len(hostContainer.Stats)-1]
// 				cpuUsage := float64(0)
// 				pre := hostContainer.Stats[len(hostContainer.Stats)-2]
// 				rawUsage := cur.Cpu.Usage.Total - pre.Cpu.Usage.Total
// 				intervalInNs := getInterval(cur.Timestamp, pre.Timestamp)
// 				cpuUsage = (float64(rawUsage*100) / float64(intervalInNs))
// 				fmt.Println(fmt.Sprintf("rawUsage : %d, intervalInNs: %d", rawUsage, intervalInNs))
// 				cpuUsage = cpuUsage / float64(i.machineInfo.NumCores)
// 				if cpuUsage > 100 {
// 					cpuUsage = 100
// 				}
// 				cpuusagepercent := math.Floor(cpuUsage)
// 				fmt.Println(fmt.Sprintf("host cpu:  %d ", cpuusagepercent))
// 				tags_host_cpu := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "NumCores": strconv.Itoa(i.machineInfo.NumCores)}
// 				i.points = append(i.points, makePoint("HostCpuPercent", tags_host_cpu, toSignedIfUnsigned(int64(cpuusagepercent))))
// 			}

// 			if hostContainer.Spec.HasMemory {
// 				//host memory information
// 				cur := hostContainer.Stats[len(hostContainer.Stats)-1]
// 				tags_host_memory := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "MemoryCapacity": strconv.Itoa(int(i.machineInfo.MemoryCapacity))}
// 				memoryusage := cur.Memory.Usage
// 				memorytotal := i.machineInfo.MemoryCapacity
// 				memoryusagepercent := math.Floor(float64(memoryusage*100) / float64(memorytotal))
// 				fmt.Println(fmt.Sprintf("host memory: %d", memoryusagepercent))
// 				i.points = append(i.points, makePoint("HostMemoryUsage", tags_host_memory, toSignedIfUnsigned(memoryusage)))
// 				i.points = append(i.points, makePoint("HostMemoryPercent", tags_host_memory, toSignedIfUnsigned(int64(memoryusagepercent))))
// 			}

// 			// wait cadvisor to support loadaverage in container

// 			// if hostContainer.Spec.HasCpu {
// 			// 	//host load average
// 			// 	cur := hostContainer.Stats[len(hostContainer.Stats)-1]
// 			// 	loadaverage := float64(0)
// 			// 	loadaverage = float64(cur.Cpu.LoadAverage) / float64(1000)
// 			// 	fmt.Println(fmt.Sprintf("host load average:  %d ", loadaverage))
// 			// 	tags_host_loadaverage := map[string]string{"hostid": i.host_id, "hostip": i.host_ip}
// 			// 	i.points = append(i.points, makePoint("HostLoadAverage", tags_host_loadaverage, toSignedIfUnsigned(int64(cur.Cpu.LoadAverage))))
// 			// }

// 			loadaverage, err := i.getLoadAverage()
// 			if err == nil {
// 				tags_host_loadaverage := map[string]string{"hostid": i.host_id, "hostip": i.host_ip}
// 				i.points = append(i.points, makePoint("HostLoadAverage", tags_host_loadaverage, toSignedIfUnsigned(loadaverage)))
// 			}

// 			if hostContainer.Spec.HasFilesystem {
// 				cur := hostContainer.Stats[len(hostContainer.Stats)-1]
// 				for _, fs := range cur.Filesystem {
// 					device := fs.Device
// 					fsuseagepercent := (float64(fs.Usage*100) / float64(fs.Limit))
// 					fmt.Println(fmt.Sprintf("host fs device: %s,  %d ", device, fsuseagepercent))
// 					tags_host_fssagepercent := map[string]string{"hostid": i.host_id, "hostip": i.host_ip, "Device": device}
// 					i.points = append(i.points, makePoint("HostFsUsagePercent", tags_host_fssagepercent, toSignedIfUnsigned(int64(fsuseagepercent))))
// 				}
// 			}

// 			if hostContainer.Spec.HasNetwork {
// 				cur := hostContainer.Stats[len(hostContainer.Stats)-1]
// 				tcp_closewait := cur.Network.Tcp.CloseWait
// 				tcp6_closewait := cur.Network.Tcp6.CloseWait
// 				fmt.Println(fmt.Sprintf("host tcp closewait  %d ", tcp_closewait))
// 				fmt.Println(fmt.Sprintf("host tcp6 closewait %d ", tcp6_closewait))
// 				tags_host_tcp := map[string]string{"hostid": i.host_id, "hostip": i.host_ip}
// 				i.points = append(i.points, makePoint("HostTcpCloseWait", tags_host_tcp, toSignedIfUnsigned(int64(tcp_closewait))))
// 				i.points = append(i.points, makePoint("HostTcp6CloseWait", tags_host_tcp, toSignedIfUnsigned(int64(tcp6_closewait))))
// 			}

// 			pointsToFlush := i.points
// 			i.points = make([]*influxdb.Point, 0)
// 			i.lastWrite = time.Now()
// 			if len(pointsToFlush) > 0 {
// 				fmt.Println(fmt.Sprintf("= start to flush point to storage, point size= [ %d ]", len(pointsToFlush)))
// 				points := make([]influxdb.Point, len(pointsToFlush))
// 				for i, p := range pointsToFlush {
// 					points[i] = *p
// 				}
// 				bp := influxdb.BatchPoints{
// 					Points:   points,
// 					Database: i.conf.Influxdb_Database,
// 				}
// 				response, err := influxdb_client.Write(bp)
// 				if err != nil || checkResponseForErrors(response) != nil {
// 					return fmt.Errorf("failed to write stats to influxDb - %s", err)
// 				}
// 			}
// 		}
// 	}
// 	err = scanner.Err()
// 	if err != nil {
// 		i.logger.Error(fmt.Sprintf("read error"), err)
// 	}
// 	i.logger.Info("Exit Scanner")
// 	return nil
// }

// func (i *InfluxdbService) defaultReadyToFlush() bool {
// 	return time.Since(i.lastWrite) >= i.bufferDuration
// }

// func (i *InfluxdbService) getLoadAverage() (int64, error) {
// 	shellcommand := "uptime"
// 	loadaverage_string, err := common.ExecShell(shellcommand)
// 	if err != nil {
// 		i.logger.Error(fmt.Sprintf("Get Load Average have err: ", err), err)
// 		return 0, err
// 	}
// 	loadaverage_string_loadaverage := strings.Split(loadaverage_string, ",")
// 	loadaverage_string_loadaverage1 := strings.Split(loadaverage_string_loadaverage[3], ":")
// 	loadaverage := strings.TrimSpace(loadaverage_string_loadaverage1[1])
// 	loadaverage_float, _ := strconv.ParseFloat(loadaverage, 64)
// 	loadaverage_int := int64(loadaverage_float * 1000)
// 	return loadaverage_int, nil
// }

// func makePoint(name string, tags map[string]string, value interface{}) *influxdb.Point {
// 	fields := map[string]interface{}{
// 		"value": toSignedIfUnsigned(value),
// 	}
// 	point := &influxdb.Point{
// 		Measurement: name,
// 		Fields:      fields,
// 		Time:        time.Now(),
// 	}
// 	addTagsToPoint(point, tags)
// 	return point
// }

// func addTagsToPoint(point *influxdb.Point, tags map[string]string) {
// 	if point.Tags == nil {
// 		point.Tags = tags
// 	} else {
// 		for k, v := range tags {
// 			point.Tags[k] = v
// 		}
// 	}
// }

// func checkResponseForErrors(response *influxdb.Response) error {
// 	const msg = "failed to write stats to influxDb - %s"

// 	if response != nil && response.Err != nil {
// 		return fmt.Errorf(msg, response.Err)
// 	}
// 	if response != nil && response.Results != nil {
// 		for _, result := range response.Results {
// 			if result.Err != nil {
// 				return fmt.Errorf(msg, result.Err)
// 			}
// 			if result.Series != nil {
// 				for _, row := range result.Series {
// 					if row.Err != nil {
// 						return fmt.Errorf(msg, row.Err)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// func toSignedIfUnsigned(value interface{}) interface{} {
// 	switch v := value.(type) {
// 	case uint64:
// 		return int64(v)
// 	case uint32:
// 		return int32(v)
// 	case uint16:
// 		return int16(v)
// 	case uint8:
// 		return int8(v)
// 	case uint:
// 		return int(v)
// 	}
// 	return value
// }

// func getInterval(current, previous time.Time) int64 {
// 	return current.Sub(previous).Nanoseconds()
// 	//return current.UnixNano() - previous.UnixNano()
// }
