-- Chisel description
description = "Show the top containers defined by the highest CPU utilization."
short_description = "Top containers by CPU usage"
category = "CPU Usage"

-- Chisel argument list
args = {}

require "commonpackage"
terminal = require "ansiterminal"
require "httppackage"

grtable = {}
grtable_memory = {}
grtable_process = {}
grtable_process_memory = {}
grtable_thread = {}
grtable_net = {}
grtable_file = {}
grtable_errors = {}
grtable_fdcount = {}
islive = false
fkeys = {}
fkeys_memory = {}
fkeys_process = {}
fkeys_process_memory = {}
fkeys_thread = {}
fkeys_net = {}
fkeys_file = {}
fkeys_fdcount = {}
fkeys_errors = {}
infostr = ""
netslower = ""

fdfilename_value = ""
fdl4proto_value = ""
netbytes = ""
tot = 0
totin = 0
totout = 0

grtable_httprequest = {}
partial_transactions = {}

vizinfo =
{
  key_fld = {"container.name"},
  key_desc = {"Container.name"},
  value_fld = "thread.exectime",
  value_desc = "CPU",
  value_units = "timepct",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_memory =
{
  key_fld = {"thread.vmsize","thread.vmrss","container.name"},
  key_desc = {"VIRT","RES","Container.name"},
  value_fld = "thread.exectime",
  value_desc = "CPU",
  value_units = "timepct",
  top_number = 1000,
  output_format = "normal"
}


vizinfo_process =
{
  key_fld = {"proc.name","user.name","proc.nthreads","proc.pid","thread.vmsize","thread.vmrss","evt.cpu","container.id","container.name"},
  key_desc = {"Process","User","ThreadCount","Host_pid","VIRT","RES","Cpu_No","Container.id","Container.name"},
  value_fld = "thread.exectime",
  value_desc = "CPU",
  value_units = "timepct",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_process_memory =
{
  key_fld = {"proc.name","user.name","proc.nthreads","proc.pid","thread.vmsize","thread.vmrss","evt.cpu","container.id","container.name"},
  key_desc = {"Process","User","ThreadCount","Host_pid","VIRT","RES","Cpu_No","Container.id","Container.name"},
  value_fld = "thread.exectime",
  value_desc = "CPU",
  value_units = "timepct",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_thread =
{
  key_fld = {"proc.name","user.name","proc.nthreads","proc.pid","thread.tid","proc.vpid","thread.vmsize","thread.vmrss","evt.cpu","container.id","container.name"},
  key_desc = {"Process","User","ThreadCount","Host_pid","ThreadId","Container_pid","VIRT","RES","Cpu_No","Container.id","Container.name"},
  value_fld = "thread.exectime",
  value_desc = "CPU",
  value_units = "timepct",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_net =
{
  key_fld = {"fd.name","proc.name","evt.count","container.name"},
  key_desc = {"Connection","Process","IOPS","Container.name"},
  value_fld = "evt.rawarg.res",
  value_desc = "Net",
  value_units = "bytes",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_file =
{
  key_fld = {"proc.name","container.name"},
  key_desc = {"Process","Container.name"},
  value_fld = "evt.rawarg.res",
  value_desc = "Bytes",
  value_units = "bytes",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_fdcount =
{
  key_fld = {"proc.name","proc.fdlimit","proc.fdusage","container.name"},
  key_desc = {"Process","Max","PCT","Container.name"},
  value_fld = "proc.fdopencount",
  value_desc = "Open",
  value_units = "none",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_errors =
{
  key_fld = {"proc.name","container.name"},
  key_desc = {"Process","Container.name"},
  value_fld = "evt.count",
  value_desc = "#Errors",
  value_units = "none",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_httprequest =
{
    key_fld = {"method", "url", "containerid", "containername"},
    key_desc = {"method", "url", "containerid", "containername"},
    value_fld = "ncalls",
    value_desc = "ncalls",
    value_units = "none",
    top_number = 1000,
    output_format = "normal"
}

-- Initialization callback
function on_init()

  sysdig.set_snaplen(2000)

  -- Request the fields we need
  for i, name in ipairs(vizinfo.key_fld) do
    fkeys[i] = chisel.request_field(name)
  end

   for i_memory, name_memory in ipairs(vizinfo_memory.key_fld) do
    fkeys_memory[i_memory] = chisel.request_field(name_memory)
  end

  for i_process, name_process in ipairs(vizinfo_process.key_fld) do
    fkeys_process[i_process] = chisel.request_field(name_process)
  end

  for i_process_memory, name_process_memory in ipairs(vizinfo_process_memory.key_fld) do
    fkeys_process_memory[i_process_memory] = chisel.request_field(name_process_memory)
  end

  for i_thread, name_thread in ipairs(vizinfo_thread.key_fld) do
    fkeys_thread[i_thread] = chisel.request_field(name_thread)
  end

  for i_net, name_net in ipairs(vizinfo_net.key_fld) do
    fkeys_net[i_net] = chisel.request_field(name_net)
  end

  for i_file, name_file in ipairs(vizinfo_file.key_fld) do
    fkeys_file[i_file] = chisel.request_field(name_file)
  end

  for i_fdcount, name_fdcount in ipairs(vizinfo_fdcount.key_fld) do
    fkeys_fdcount[i_fdcount] = chisel.request_field(name_fdcount)
  end

  for i_errors, name_errors in ipairs(vizinfo_errors.key_fld) do
    fkeys_errors[i_errors] = chisel.request_field(name_errors)
  end

  -- Request the fields we need
  fvalue = chisel.request_field(vizinfo.value_fld)
  fvalue_process = chisel.request_field(vizinfo_process.value_fld)
  fvalue_thread = chisel.request_field(vizinfo_thread.value_fld)
  fvalue_net = chisel.request_field(vizinfo_net.value_fld)
  fvalue_file = chisel.request_field(vizinfo_file.value_fld)
  fvalue_fdcount = chisel.request_field(vizinfo_fdcount.value_fld)
  fvalue_errors = chisel.request_field(vizinfo_errors.value_fld)

  fcpu = chisel.request_field("thread.cpu")
  fcpu_memory = chisel.request_field("thread.cpu")
  fcpu_process = chisel.request_field("thread.cpu")
  fcpu_thread = chisel.request_field("thread.cpu") 
  eventtype = chisel.request_field("evt.type")
  fdtype = chisel.request_field("fd.type")
  evtisio = chisel.request_field("evt.is_io")
  evtdir = chisel.request_field("evt.dir")
  evtrawres = chisel.request_field("evt.rawres")
  containername = chisel.request_field("container.name")
  evtfailed = chisel.request_field("evt.failed")

  fdnum = chisel.request_field("fd.num")
  fdname = chisel.request_field("fd.name")

  fbuf = chisel.request_field("evt.rawarg.data")
  fisread = chisel.request_field("evt.is_io_read")
  fres = chisel.request_field("evt.rawarg.res")
  fpname = chisel.request_field("proc.name")
  fcontainerid = chisel.request_field("container.id")

  datetime = chisel.request_field("evt.datetime")
  latency = chisel.request_field("evt.latency")
  evtbuflen = chisel.request_field("evt.buflen")
  fdsockfamily = chisel.request_field("fd.sockfamily")

  fdl4proto  = chisel.request_field("fd.l4proto")

  data = chisel.request_field("evt.arg[1]")

  fdsport = chisel.request_field("fd.sport")  
  mem_virt= chisel.request_field("thread.vmsize")
  mem_res= chisel.request_field("thread.vmrss")

  process_mem_virt= chisel.request_field("thread.vmsize")
  process_mem_res= chisel.request_field("thread.vmrss")

  http_init()

  return true
end

-- Helpers --
function split(s, delimiter)
  result = {};
  for match in (s..delimiter):gmatch("(.-)"..delimiter) do
    table.insert(result, match);
  end
  return result;
end

function build_grtable_key(transaction)
    request = transaction["request"]
    response = transaction["response"]
    statuscode = response["code"]
    ret = ""
    ret = transaction["containerid"] .. "\001\001" .. transaction["containername"] .. "\001\001" .. transaction["pid_value"] .. "\001\001" .. transaction["fd_value"] .. "\001\001" .. statuscode .. "\001\001"
    ret = ret .. string.format("%s\001\001%s", request["method"], request["url"])
    return ret
end

function on_transaction(transaction)
    grtable_key = build_grtable_key(transaction)
    if not grtable_httprequest[grtable_key] then
        grtable_httprequest[grtable_key] = {}
    end
    table.insert(grtable_httprequest[grtable_key], transaction)
end

-- Final chisel initialization
function on_capture_start()
  chisel.set_interval_s(10)
  return true
end

function aggregate_grtable()
    for key, transactions in pairs(grtable_httprequest) do
            total_bytes = 0
            for _, tr in ipairs(transactions) do
                total_bytes = total_bytes + tr["response"]["length"]
            end
            total_time = 0
            for _, tr in ipairs(transactions) do
                total_time = total_time + tr["response"]["ts"] - tr["request"]["ts"]
            end
            grtable_httprequest[key] = total_bytes .. "             " .. total_time / #transactions
    end
end

-- Event parsing callback
function on_event()
  -- Container CPU
  if evt.field(eventtype) == "procinfo" then
    local key = nil
    local kv = nil

    for i, fld in ipairs(fkeys) do
      kv = evt.field(fld)
      if kv == nil then
        return
      end

      if key == nil then
        key = kv
      else
        key = key .. "\001\001" .. evt.field(fld)
      end
    end

    local cpu = evt.field(fcpu)

    if grtable[key] == nil then
      grtable[key] = cpu * 10000000
    else
      grtable[key] = grtable[key] + (cpu * 10000000)
    end
  end

  -- Container Memory
  if evt.field(mem_virt) ~= nil and evt.field(mem_virt) > 0 and evt.field(mem_res) ~= nil and evt.field(mem_res) > 0 and  evt.field(eventtype) == "procinfo" then
    local key_memory = nil
    local kv_memory = nil

    for i_memory, fld_memory in ipairs(fkeys_memory) do
      kv_memory = evt.field(fld_memory)
      if kv_memory == nil then
        return
      end

      if key_memory == nil then
        key_memory = kv_memory
      else
        key_memory = key_memory .. "\001\001" .. evt.field(fld_memory)
      end
    end

    local cpu_memory = evt.field(fcpu_memory)

    if grtable_memory[key_memory] == nil then
      grtable_memory[key_memory] = cpu_memory * 10000000
    else
      grtable_memory[key_memory] = grtable_memory[key_memory] + (cpu_memory * 10000000)
    end
  end

  -- Process CPU
  if evt.field(eventtype) == "procinfo" then
    local key_process = nil
    local kv_process = nil

    for i_process, fld_process in ipairs(fkeys_process) do
      kv_process = evt.field(fld_process)
      if kv_process == nil then
        return
      end

      if key_process == nil then
        key_process = kv_process
      else
        key_process = key_process .. "\001\001" .. evt.field(fld_process)
      end
    end

    local cpu_process = evt.field(fcpu_process)

    if grtable_process[key_process] == nil then
      grtable_process[key_process] = cpu_process * 10000000
    else
      grtable_process[key_process] = grtable_process[key_process] + (cpu_process * 10000000)
    end
  end

  -- Process Mempry
  if evt.field(process_mem_virt) ~= nil and evt.field(process_mem_virt) > 0 and evt.field(process_mem_res) ~= nil and evt.field(process_mem_res) > 0 and  evt.field(eventtype) == "procinfo" then
    local key_process_memory = nil
    local kv_process_memory = nil

    for i_process_memory, fld_process_memory in ipairs(fkeys_process_memory) do
      kv_process_memory = evt.field(fld_process_memory)
      if kv_process_memory == nil then
        return
      end

      if key_process_memory == nil then
        key_process_memory = kv_process_memory
      else
        key_process_memory = key_process_memory .. "\001\001" .. evt.field(fld_process_memory)
      end
    end

    local cpu_process_memory = evt.field(fcpu_process)

    if grtable_process_memory[key_process_memory] == nil then
      grtable_process_memory[key_process_memory] = cpu_process_memory * 10000000
    else
      grtable_process_memory[key_process_memory] = grtable_process[key_process_memory] + (cpu_process_memory * 10000000)
    end
  end


  -- Thread CPU/Mempry
  if evt.field(eventtype) == "procinfo" then
    local key_thread = nil
    local kv_thread = nil

    for i_thread, fld_thread in ipairs(fkeys_thread) do
      kv_thread = evt.field(fld_thread)
      if kv_thread == nil then
        return
      end

      if key_thread == nil then
        key_thread = kv_thread
      else
        key_thread = key_thread .. "\001\001" .. evt.field(fld_thread)
      end
    end

    local cpu_thread = evt.field(fcpu_thread)

    if grtable_thread[key_thread] == nil then
      grtable_thread[key_thread] = cpu_thread * 10000000
    else
      grtable_thread[key_thread] = grtable_thread[key_thread] + (cpu_thread * 10000000)
    end
  end
  
  --NET
  if (evt.field(fdtype) == "ipv4" or evt.field(fdtype) == "ipv6") and evt.field(evtisio) == true then
    local bytes = evt.field(fres)
    local isread = evt.field(fisread)
    local containerid_net = evt.field(fcontainerid)
    local containername_net = evt.field(containername)
    fdfilename_value = evt.field(fdname)
    fdl4proto_value = evt.field(fdl4proto)
    if bytes ~= nil and bytes > 0 then
      tot = tot + bytes
      if isread then
        totin = totin + bytes
      else
        totout = totout + bytes
      end
    end
    netbytes = "NetBytes" .. "     " .. containerid_net .. "      " .. containername_net  .. "     " .. fdfilename_value .. "                         "  .. fdl4proto_value .. "       " .. "       " .. totin .. "        " .. totout .. "       " .. tot
    
    local key_net = nil
    local kv_net = nil

    for i_net, fld_net in ipairs(fkeys_net) do
      kv_net = evt.field(fld_net)
      if kv_net == nil then
        return
      end

      if key_net == nil then
        key_net = kv_net
      else
        key_net = key_net .. "\001\001" .. evt.field(fld_net)
      end
    end

    value_net = evt.field(fvalue_net)

    if value_net ~= nil and value_net > 0 then
      entryval_net = grtable_net[key_net]

      if entryval_net == nil then
        grtable_net[key_net] = value_net
      else
        grtable_net[key_net] = grtable_net[key_net] + value_net
      end
    end

  end

  --FILE
  if evt.field(fdtype) == "file" and evt.field(evtisio) == true then
    local key_file = nil
    local kv_file = nil

    for i_file, fld_file in ipairs(fkeys_file) do
      kv_file = evt.field(fld_file)
      if kv_file == nil then
        return
      end

      if key_file == nil then
        key_file = kv_file
      else
        key_file = key_file .. "\001\001" .. evt.field(fld_file)
      end
    end

    value_file = evt.field(fvalue_file)

    if value_file ~= nil and value_file > 0 then
      entryval_file = grtable_file[key_file]

      if entryval_file == nil then
        grtable_file[key_file] = value_file
      else
        grtable_file[key_file] = grtable_file[key_file] + value_file
      end
    end

  end

-- FD COUNT
  if evt.field(eventtype) ~= "switch" then
    local key_fdcount = nil 
    local kv_fdcount = nil

    for i_fdcount, fld_fdcount in ipairs(fkeys_fdcount) do
      kv_fdcount = evt.field(fld_fdcount)
      if kv_fdcount == nil then
        return
      end

      if key_fdcount == nil then
        key_fdcount = kv_fdcount
      else
        key_fdcount = key_fdcount .. "\001\001" .. evt.field(fld_fdcount)
      end
    end

    value_fdcount = evt.field(fvalue_fdcount)

    if value_fdcount ~= nil and value_fdcount > 0 then
      entryval_fdcount = grtable_fdcount[key_fdcount]

      if entryval_fdcount == nil then
        grtable_fdcount[key_fdcount] = value_fdcount
      else
        grtable_fdcount[key_fdcount] = grtable_fdcount[key_fdcount] + value_fdcount
      end
    end
  end

-- ECHO
  if evt.field(evtisio) == true and evt.field(evtdir) == "<"  and evt.field(evtrawres) > 0 then
    local buf_echo = evt.field(fbuf)
    local isread_echo = evt.field(fisread)
    local res_echo = evt.field(fres)
    local name_echo = evt.field(fdname)
    local pname_echo = evt.field(fpname)
    local containername_echo = evt.field(containername)
    local containerid_echo = evt.field(fcontainerid)

    if name_echo == nil then
      name_echo = "<NA>"
    end
    if res_echo <= 0 then
      return true
    end
    local container = ""  
    container = string.format("[%s] [%s]", containername_echo, containerid_echo );
    if isread_echo then
      name_pname = string.format("%s (%s)", name_echo, pname_echo );
      infostr = string.format("------ Read %s from %s %s", format_bytes(res_echo), container, name_pname)
    else
      name_pname = string.format("%s (%s)", name_echo, pname_echo );
      infostr = string.format("------ Write %s to %s %s", format_bytes(res_echo), container, name_pname)
    end
    --print(infostr)
  end

-- Http Request
  if evt.field(evtisio) == true and (evt.field(fdsockfamily) == "ip" or evt.field(fdsockfamily) == "unix") then
    local evtbuflen_value = evt.field(evtbuflen)
    if evtbuflen_value ~= nil then
      run_http_parser(evt, on_transaction)
    end
  end


  return true
end

-- Periodic timeout callback
function on_interval(ts_s, ts_ns, delta)
  
  print_sorted_table("TopContainersCpu", grtable, ts_s, 0, delta, vizinfo)

  print_sorted_table("TopContainersMemory", grtable_memory, ts_s, 0, delta, vizinfo_memory)

  print_sorted_table("TopContainersProcess", grtable_process, ts_s, 0, delta, vizinfo_process)

  print_sorted_table("TopProcessMemory", grtable_process_memory, ts_s, 0, delta, vizinfo_process_memory)

  print_sorted_table("TopContainersThread", grtable_thread, ts_s, 0, delta, vizinfo_thread)
  
  print_sorted_table("TopContainersNet", grtable_net, ts_s, 0, delta, vizinfo_net)

  print_sorted_table("TopContainersFile", grtable_file, ts_s, 0, delta, vizinfo_file)

  print_sorted_table("TopFdCount", grtable_fdcount, ts_s, 0, delta, vizinfo_fdcount)

  print_sorted_table("TopContainersErrors", grtable_errors, ts_s, 0, delta, vizinfo_errors)
  
  etime = evt.field(ftime)

  print(netslower)

  aggregate_grtable()
  print_sorted_table("HttpRequest ", grtable_httprequest, ts_s, 0, delta, vizinfo_httprequest)

  print(netbytes)

  -- Clear the table
  grtable = {}
  grtable_memory = {}
  grtable_process = {}
  grtable_process_memory = {}
  grtable_thread = {}
  grtable_net = {}
  grtable_file = {}
  grtable_errors = {}
  grtable_fdcount = {}
  grtable_httprequest = {}

  fdfilename_value = ""
  fdl4proto_value = ""
  netbytes = ""
  tot = 0
  totin = 0
  totout = 0

  infostr = ""
  netslower = ""
  return true
end

-- Called by the engine at the end of the capture (Ctrl-C)
function on_capture_end(ts_s, ts_ns, delta) 
  return true
end
