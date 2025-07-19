--local process = require("process")

return function(argv)
    if not argv then argv = process.argv end
    local opts = {}
    local args = {}
    for i, arg in ipairs(argv) do
        -- option?
        local opt = arg:match("^%-%-(.*)")
        if opt then
            -- extract option name and value
            local key, value = opt:match("([a-z_%-]*)=(.*)")
            --p('OPT', opt, key, value)
            -- value provided?
            if value then
                -- option seen once?
                if type(opts[key]) == 'string' then
                    -- transform option to array of values
                    opts[key] = { opts[key], value }
                    -- options seen many times?
                elseif type(opts[key]) == 'table' then
                    -- append value
                    table.insert(opts[key], value)
                    -- options was not seen
                else
                    -- assign value
                    opts[key] = value
                end
                -- no value provided. just set option to true
            elseif opt ~= '' then
                opts[opt] = true
                -- options stop
            else
                -- copy left arguments
                for i = i + 1, #argv do
                    table.insert(args, argv[i])
                end
                break
            end
            -- argument
        else
            table.insert(args, arg)
        end
    end
    -- report options and arguments
    return opts, args
end
