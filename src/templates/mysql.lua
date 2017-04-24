description = "Show a log of mysql commands (select/insert/update/delete)"
short_description = "mysql requests log"
category = "Application"

-- Chisel argument list
args ={}

require "common"

-- Helpers --
function split(s, delimiter)
  result = {};
  for match in (s..delimiter):gmatch("(.-)"..delimiter) do
    table.insert(result, match);
  end
  return result;
end

-- Argument notification callback
function on_set_arg(name, val)
   return true
end

-- Initialization callback
function on_init()
    util = {}
    start_time = os.time()
    sysdig.set_filter("(fd.sport=3306 or proc.name=mysql) and evt.is_io=true and evt.type=write")
    sysdig.set_snaplen(4096)
    data = chisel.request_field("evt.arg[1]")
    datetime = chisel.request_field("evt.datetime")
    latency = chisel.request_field("evt.latency")
    fcontainername = chisel.request_field("container.name")
    fcontainerid = chisel.request_field("container.id")
    return true
end

-- Event callback
function on_event()
  local data = evt.field(data)
  local line = split(data, " ")
  local op = string.lower(line[1])
  if line[1] == ".....select" then
        local method = string.sub(op , 6,-1)
        local content = "select" .. string.sub(data , 12,-1)
        query = method .. ":" .. content
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end

  if line[1] == "!....insert" then
        local method = string.sub(op , 6,-1)
        local content = "insert" .. string.sub(data , 12,-1)
        query = method .. ":" .. content
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end

  if line[1] == ".....update" then
        local method = string.sub(op , 6,-1)
        local content = "update" .. string.sub(data , 12,-1)
        query = method .. ":" .. content
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end

  if line[1] == ".....delete" then
        local method = string.sub(op , 6,-1)
        local content = "delete" .. string.sub(data , 12,-1)
        query = method .. ":" .. content
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end

  return true
end