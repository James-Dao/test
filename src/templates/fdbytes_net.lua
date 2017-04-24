description = "Counts the total bytes read from and written to the network, and prints the result every second";
short_description = "Show total network I/O bytes";
category = "Net";

-- Chisel argument list
args = {}

tot = 0
totin = 0
totout = 0
fdfilename_value = ""
fdl4proto_value = ""

require "common"

-- Initialization callback
function on_init()
	-- Request the fields
	fbytes = chisel.request_field("evt.rawarg.res")
	ftime = chisel.request_field("evt.time.s")
	fisread = chisel.request_field("evt.is_io_read")
	fdfilename  = chisel.request_field("fd.name")
	fdl4proto  = chisel.request_field("fd.l4proto")
	
	-- set the filter
	chisel.set_filter("evt.is_io=true and (fd.type=ipv4 or fd.type=ipv6)")
	
	chisel.set_interval_s(1)
	
	return true
end

-- Event parsing callback
function on_event()
	bytes = evt.field(fbytes)
	isread = evt.field(fisread)

	fdfilename_value = evt.field(fdfilename)
	fdl4proto_value = evt.field(fdl4proto)

	if bytes ~= nil and bytes > 0 then
		tot = tot + bytes
		
		if isread then
			totin = totin + bytes
		else
			totout = totout + bytes
		end
	end

	return true
end

function on_interval(delta)
	etime = evt.field(ftime)

	print(etime .. "     " .. fdfilename_value .. "   "  .. fdl4proto_value .. "       " .. "       in:" .. format_bytes(totin) .. "        out:" .. format_bytes(totout) .. "       tot:" .. format_bytes(tot))
	tot = 0
	totin = 0
	totout = 0
	return true
end