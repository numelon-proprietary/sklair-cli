local function tokenise(html)
    local tokens = {}
    local pos = 1

    while pos <= #html do
        local start, stop, tag, attributes, selfClose = html:find("<(%/?%w+)(.-)(%/?)>", pos)
        if start then
            -- text between tags / body
            if start > pos then
                local text = html:sub(pos, start - 1):match("^%s*(.-)%s*$")
                if #text > 0 then
                    table.insert(tokens, { type = "text", content = text })
                end
            end

            tag = tag:gsub("/", "")
            local isClosing = html:sub(start + 1, start + 1) == "/"

            if isClosing then
                table.insert(tokens, { type = "close", tag = tag })
            elseif #selfClose > 0 then
                table.insert(tokens, { type = "selfclose", tag = tag, attrs = attributes })
            else
                table.insert(tokens, { type = "open", tag = tag, attrs = attributes })
            end

            pos = stop + 1
        else
            -- No more tags, just trailing text
            local trailing = html:sub(pos):match("^%s*(.-)%s*$")
            if #trailing > 0 then
                table.insert(tokens, { type = "text", content = trailing })
            end
            break
        end
    end

    return tokens
end

local fs = require("coro-fs")
local aa = fs.readFile("../index.html")

p(tokenise(aa))
