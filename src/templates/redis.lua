description = "Show a log of redis commands (get/set/expire/del)"
short_description = "redis requests log"
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
    sysdig.set_filter("(fd.sport=6379 or proc.name=redis) and evt.is_io=true and evt.type=write")
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
  local op = string.sub(data , 1,11)
  if op == "*2..$4..key" then
        method = "key"
        local content = string.sub(data , 18, -1)
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end
   if op == "*2..$3..get" then
        method = "get"
        local content = string.sub(data , 18, -1)
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end
   if op == "*3..$3..set" then
        method = "set"
        local content = string.sub(data , 18, -1)
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end
   if op == "*3..$6..exp" then
        method = "expire"
        local content = string.sub(data , 21, -1)
        print(string.format("%-23.23s %-20.20s %-20.20s %-23.23s %-20s %s",evt.field(datetime),
    evt.field(fcontainerid),evt.field(fcontainername),method,format_time_interval(evt.field(latency)),content))
  end

  return true
end