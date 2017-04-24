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
grtable_net = {}
grtable_file = {}
grtable_errors = {}
grtable_fdcount = {}
islive = false
fkeys = {}
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
  key_fld = {"fd.name","proc.name","proc.pid","proc.vpid","fd.sport","fd.sproto","evt.count","container.id","container.name"},
  key_desc = {"Connection","Process","Host_pid","Container_pid","Server_Port","PROTO","IOPS","Container.id","Container.name"},
  value_fld = "evt.rawarg.res",
  value_desc = "Net",
  value_units = "bytes",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_file =
{
  key_fld = {"proc.name","proc.pid","proc.vpid","thread.vmsize","thread.vmrss","container.id","container.name"},
  key_desc = {"Process","Host_pid","Container_pid","VIRT","RES","Container.id","Container.name"},
  value_fld = "evt.rawarg.res",
  value_desc = "Bytes",
  value_units = "bytes",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_fdcount =
{
  key_fld = {"proc.name","proc.pid","proc.vpid","proc.fdlimit","proc.fdusage","container.id","container.name","proc.exeline"},
  key_desc = {"Process","Host_pid","Container_pid","Max","PCT","Container.id","Container.name","Command"},
  value_fld = "proc.fdopencount",
  value_desc = "Open",
  value_units = "none",
  top_number = 1000,
  output_format = "normal"
}

vizinfo_errors =
{
  key_fld = {"proc.name","proc.pid","proc.vpid","container.id","container.name"},
  key_desc = {"Process","Host_pid","Container_pid","Container.id","Container.name"},
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
  fvalue_net = chisel.request_field(vizinfo_net.value_fld)
  fvalue_file = chisel.request_field(vizinfo_file.value_fld)
  fvalue_fdcount = chisel.request_field(vizinfo_fdcount.value_fld)
  fvalue_errors = chisel.request_field(vizinfo_errors.value_fld)

  fcpu = chisel.request_field("thread.cpu")
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
  chisel.set_interval_s(1)
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
            grtable_httprequest[key] = format_bytes(total_bytes) .. "             " .. format_time_interval(total_time / #transactions)
    end
end

-- Event parsing callback
function on_event()
  -- CPU
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
    netbytes = "NetBytes" .. "     " .. containerid_net .. "      " .. containername_net  .. "     " .. fdfilename_value .. "                         "  .. fdl4proto_value .. "       " .. "       " .. format_bytes(totin) .. "        " .. format_bytes(totout) .. "       " .. format_bytes(tot)
    
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

  --ERRORS
  if evt.field(evtfailed) == true then
    local key_errors = nil
    local kv_errors = nil

    for i_errors, fld_errors in ipairs(fkeys_errors) do
      kv_errors = evt.field(fld_errors)
      if kv_errors == nil then
        return
      end

      if key_errors == nil then
        key_errors = kv_errors
      else
        key_errors = key_errors .. "\001\001" .. evt.field(fld_errors)
      end
    end

    value_errors = evt.field(fvalue_errors)

    if value_errors ~= nil and value_errors > 0 then
      entryval_errors = grtable_errors[key_errors]

      if entryval_errors == nil then
        grtable_errors[key_errors] = value_errors
      else
        grtable_errors[key_errors] = grtable_errors[key_errors] + value_errors
      end
    end

  end

-- ECHO
  if evt.field(containername) ~="host" and evt.field(evtisio) == true and evt.field(evtdir) == "<"  and evt.field(evtrawres) > 0 then
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

-- Net Slower
  if (evt.field(fdtype) == "ipv4" or evt.field(fdtype) == "ipv6") and evt.field(evtisio) == true then
    lat = evt.field(latency) / 1000000
    fn = evt.field(fdname)
    if evt.field(evtdir) == "<" then
      netslower = "NetSlower    " .. string.format("%-20.20s %-20.20s %-12.12s %-8s %12d %s",
        evt.field(fcontainerid),evt.field(containername),evt.field(fpname),evt.field(eventtype),lat,fn)
    end
  end

-- Http Request
  if evt.field(evtisio) == true and (evt.field(fdsockfamily) == "ip" or evt.field(fdsockfamily) == "unix") then
    local evtbuflen_value = evt.field(evtbuflen)
    if evtbuflen_value ~= nil then
      run_http_parser(evt, on_transaction)
    end
  end

-- Host Mysql
  if (evt.field(fdsport)==3306 or evt.field(fpname)=="mysql") and evt.field(evtisio)==true and evt.field(eventtype)=="write" then
    local data = evt.field(data)
    local line = split(data, " ")
    local op = string.lower(line[1])
    if line[1] == ".....select" then
          local method = string.sub(op , 6,-1)
          local content = "select" .. string.sub(data , 12,-1)
          query = method .. ":" .. content
          print("Mysql   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end

    if line[1] == "!....insert" then
          local method = string.sub(op , 6,-1)
          local content = "insert" .. string.sub(data , 12,-1)
          query = method .. ":" .. content
          print("Mysql   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end

    if line[1] == ".....update" then
          local method = string.sub(op , 6,-1)
          local content = "update" .. string.sub(data , 12,-1)
          query = method .. ":" .. content
          print("Mysql   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end

    if line[1] == ".....delete" then
          local method = string.sub(op , 6,-1)
          local content = "delete" .. string.sub(data , 12,-1)
          query = method .. ":" .. content
          print("Mysql   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end
  end

-- Host Redis
  if (evt.field(fdsport)==6379 or evt.field(fpname)=="redis-server") and evt.field(evtisio)==true and evt.field(eventtype)=="write" then
    local data = evt.field(data)
    local op = string.sub(data , 1,11)
    if op == "*2..$4..key" then
          method = "key"
          local content = string.sub(data , 18, -1)
          print("Redis   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end
     if op == "*2..$3..get" then
          method = "get"
          local content = string.sub(data , 18, -1)
          print("Redis   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end
     if op == "*3..$3..set" then
          method = "set"
          local content = string.sub(data , 18, -1)
          print("Redis   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end
     if op == "*3..$6..exp" then
          method = "expire"
          local content = string.sub(data , 21, -1)
          print("Redis   " .. string.format("%-20.20s %-20.20s %-23.23s %-20s %s",
      evt.field(fcontainerid),evt.field(containername),method,format_time_interval(evt.field(latency)),content))
    end
  end

  return true
end

-- Periodic timeout callback
function on_interval(ts_s, ts_ns, delta)
  
  print_sorted_table("TopContainersCpu", grtable, ts_s, 0, delta, vizinfo)
  
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
