return {
    name = "numelon-oss/sklair-cli",
    version = "0.0.1",
    description = "A zero-abstraction web framework for building fast, modular sites with pure HTML, CSS, and JS.",
    tags = { "html", "js", "css", "framework", "modular" },
    license = "MIT",
    author = { name = "Richard Ziupsnys", email = "64844585+Richy-Z@users.noreply.github.com" },
    homepage = "https://github.com/numelon-oss/sklair-cli",
    dependencies = {
        "creationix/coro-fs",
        "creationix/coro-http",
        "luvit/secure-socket"
    },
    files = {
        "**.lua",
        "!test*",
        "!deps"
    }
}
