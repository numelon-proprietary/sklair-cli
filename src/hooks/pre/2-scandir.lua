for _, v in pairs(fs.scanDir("cache:")) do
    print(v.name, v.isDir)
end