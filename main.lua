local fs = require("fs")

local args = require("./argv")()
p(args)
if not args.build then
    print("Invalid arguments"); os.exit(1)
end

local BUILD_DIR = "build"
fs.mkdirpSync(BUILD_DIR)
local source = [[
    <xyz>
    <abc />
    <foo></foo>

    <bar>hello</bar>

    <SomeTag aaa="sure">Hello</SomeTag>

    <hi>
    <closinggg/>

    <Atsakakaka name="Hello">Bye</Atsakakaka>

    <lua>
        return 234234 * 348923849723
    </lua>

    <p>
        <lua>
            return 234234 * 348923849723
        </lua>
    </p>
]]

local native = {
    html = true, body = true, div = true, span = true, p = true -- TODO: extend later
}

-- for tag, body in source:gmatch("<(%w+)[^/>]*>(.-)</%1>") do
--     if not native[tag] then
--         local component_src = fs.readFile("./components/" .. tag .. ".html")
--         if component_src then
--             seen[tag] = "paired with content"
--             p(body)
--         else
--             print(string.format("Component '%s' could not be found inside the components directory!", tag))
--         end
--     end
-- end

local componentCache = {}
local function getComponent(name)
    if not componentCache[name] then
        local src = fs.readFileSync("./src/components/" .. name .. ".html")
        if not src then
            print(string.format("Component '%s' could not be found inside the components directory!", name))
            return nil, "Component does not exist"
        end

        componentCache[name] = src
    end

    return componentCache[name]
end

--------------------------------------------------
--------------------------------------------------
--- 1. finding positions first -------------------
--------------------------------------------------
--------------------------------------------------

local matches = {}

-- full pair w content
for startPos, tag, body, endPos in source:gmatch("()<(%w+)[^/>]*>(.-)</%2>()") do
    endPos = endPos - 1

    matches[#matches + 1] = {
        type = "full",
        start = startPos,
        stop = endPos,
        body = source:sub(startPos, endPos)
    }
end

-- self closing tags
for startPos, tag, endPos in source:gmatch("()<(%w+)[^/>]*/>()") do
    endPos = endPos - 1

    matches[#matches + 1] = {
        type = "selfclose",
        start = startPos,
        stop = endPos,
        body = source:sub(startPos, endPos)
    }
end

-- opening tags only
for startPos, tag, endPos in source:gmatch("()<(%w+)[^/>]*>()") do
    endPos = endPos - 1

    -- skip if already captured by full / selfclose
    local duplicate = false
    for _, node in ipairs(matches) do
        if node.start == startPos then
            duplicate = true
            break
        end
    end
    if not duplicate then
        matches[#matches + 1] = {
            type = "open",
            start = startPos,
            stop = endPos,
            body = source:sub(startPos, endPos)
        }
    end
end

-- sort by order of nodes in source
-- this is purely for clarity during debugging
table.sort(matches, function(a, b)
    return a.start < b.start
end)

p(matches)

--------------------------------------------------
--------------------------------------------------
--- 2. replacement with offsets etc --------------
--------------------------------------------------
--------------------------------------------------

local function header(raw)
    -- TODO: find is really naive, we need to ensure that these <> havent been escaped or arent inside of another string
    local x = raw:find("<")
    local y = raw:find(">")

    if not x or not y then return nil end

    raw = raw:sub(x + 1, y - 1) -- chop the rest off, irrelevantt
    -- chop tag into raw pieces, preserving spaces in quotes
    local props = {}
    local tagName = nil

    local i = 1
    local len = #raw
    local current = ""
    local quoteChar = nil
    local escape = false
    local tokens = {}

    while i <= len do
        local c = raw:sub(i, i)

        if quoteChar then
            if escape then
                current = current .. c
                escape = false
            elseif c == "\\" then
                escape = true
            elseif c == quoteChar then
                quoteChar = nil
                current = current .. c
            else
                current = current .. c
            end
        elseif c == "'" or c == '"' then
            quoteChar = c
            current = current .. c
        elseif c:match("%s") then
            if #current > 0 then
                tokens[#tokens + 1] = current
                current = ""
            end
        else
            current = current .. c
        end

        i = i + 1
    end

    if #current > 0 then tokens[#tokens + 1] = current end

    tagName = tokens[1]

    for i = 2, #tokens do
        local token = tokens[i]
        local key, val = token:match("^(%w+)%s*=%s*['\"](.-)['\"]$")
        if key then
            props[key] = val
        elseif not token:find("=") then
            props[token] = true
        end
    end

    return tagName, props
end

local function render(node)
    local tag, props = header(node.type == "selfclose" and node.body:gsub(" /", ""):gsub("/", "") or node.body)
    if native[tag] then return node.body end

    local INNER_PATTERN = "<[^>]+>(.-)</[^>]+>"

    -- check this (and potential builtin components) before getComponent so that we dont error
    if tag == "lua" then p(node.type) end
    if tag == "lua" and node.type == "full" then
        local _, _, inner = node.body:find(INNER_PATTERN)
        inner = inner or ""

        local fn, err = load(inner, "lua_block", "t", {
            os = os,
            math = math,
            table = table,
            string = string,
            props = props
            -- TODO: probably make a safe env handler to make this more powerful but still limited
        })

        if not fn then
            return "<!-- lua error: " .. err .. " -->"
        end

        local ok, result = pcall(fn)
        if not ok then
            return "<!-- lua runtime error: " .. result .. " -->"
        end

        return tostring(result or "")
    end

    local component_src = getComponent(tag)
    if not component_src then
        return "<!-- missing component: " .. tag .. " -->"
    end

    if node.type == "full" then
        -- extract just the inner content
        local _, _, inner = node.body:find(INNER_PATTERN)
        inner = inner or ""

        component_src = component_src:gsub("<!%-%-%s*%$%$BODY%s*%-%->", inner)
    end

    -- props injection
    component_src = component_src:gsub("<!%-%-%s*%$%$PROPS%.(%w+)%s*%-%->", function(key)
        return props[key] or "" ---@diagnostic disable-line need-check-nil
    end)

    return component_src
end

local offset = 0
for _, node in ipairs(matches) do
    local start        = node.start + offset
    local stop         = node.stop + offset

    local originalSpan = stop - start + 1
    local replacement  = render(node)

    -- actual replacement
    source             = source:sub(1, start - 1) .. replacement .. source:sub(stop + 1)

    local delta        = #replacement - originalSpan
    offset             = offset + delta
end

print(source)
