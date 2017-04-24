partial_transactions = {}

function http_init()
    buffer_field = chisel.request_field("evt.buffer")
    fd_field = chisel.request_field("fd.num")
    fdname_field = chisel.request_field("fd.name")
    pid_field = chisel.request_field("proc.pid")
    rawtime_field = chisel.request_field("evt.rawtime")
    datetime_field = chisel.request_field("evt.datetime")
    dir_field = chisel.request_field("evt.io_dir")
    containerid_field = chisel.request_field("container.id")
    containername_field = chisel.request_field("container.name")
    sysdig.set_snaplen(1024)
end

function parse_request(req_buffer)
    method, url = string.match(req_buffer, "^(%u+) (%g+)")
    if method and url then
        host = string.match(req_buffer, "Host: (%g+)%.%.")
        if host then
            url = host .. url
        end
        return {
            method=method,
            url=url
        }
    end

    return nil
end

function parse_response(resp_buffer)
    resp_code = string.match(resp_buffer, "HTTP/[%g]+ (%d+)")
    if resp_code then
        content_length = string.match(resp_buffer, "Content%-Length: (%d+)%.%.")
        if not content_length then
            content_length = 0
        end
        return {
          code = tonumber(resp_code),
          length = tonumber(content_length)
        }
    else
        return nil
    end
end

function run_http_parser(evt, on_transaction)
    buf = evt.field(buffer_field)
    fd = evt.field(fd_field)
    pid = evt.field(pid_field)
    evt_dir = evt.field(dir_field)
    key = string.format("%d\001\001%d", pid, fd)
    timestamp = evt.field(rawtime_field)

    transaction = partial_transactions[key]
    if not transaction then
        request = parse_request(buf)
        if request then
            transaction_dir = "<NA>"
            if evt_dir == "read" then
                transaction_dir = "<"
            elseif evt_dir == "write" then
                transaction_dir = ">"
            end
            request["ts"] = timestamp
            partial_transactions[key] = {
                request= request,
                dir=transaction_dir,
                containerid=evt.field(containerid_field),
                containername=evt.field(containername_field),
                pid_value=evt.field(pid_field),
                fd_value=evt.field(fdname_field),
            }
        end
    else
        response = parse_response(buf)
        if response then
            transaction["response"] = response
            transaction["response"]["ts"] = timestamp
            on_transaction(transaction)
            partial_transactions[key] = nil
        end
    end
end