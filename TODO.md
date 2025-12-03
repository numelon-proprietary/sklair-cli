# Todo list

- define whatever this is, is this a framework or a templating engine?
- fs.statSync(componentsDir): if folder then expect body.html, optional head.html, body.start.html, body.end.html idk. if file then treat as normal

## todo list transferred from `numelon-proprietary/website` repo

- create numelon web packing tool for creating websites
  - this web packing tool should have hot reload
  - allow for components
    - e.g. `/components/someComponent.html` and `/components/someComponent.css`
    - numelon web packing tool will search for components dir
      - if html found then it will replace the tag inside the html with the actual content of the html file
      - if css found for component in component dir then it will be automatically copied to the build directory into css dir and also the link href stylesheet injected into the head of the html file where the component was used
      - if component not found then just assume its a component registered from within js and leave it as-is, but issue a warning in the logs.
      - if js file is found then assume regular html component registered via js and add the js to head tag. ofc if css with same name is found then import that too
      - therefore the order from most important to least is like this:
        1. componentName.js -> import js (& optional css, if html is present then ignored) in head
        2. componentName.html -> replace all occurences of `<componentName/>` or `<componentName>` or `<componentName></componentName>` with contents of componentName.html. (& optional css in head)
            - note that if it is `<componentName>aaa</componentName>` then `aaa` would be passed to `$$COMPONENT_BODY` var inside componentName.html.
  - reusable component for menubar and footer
  - allow components to be used inside of eachother (in src) but hard fail on circular component usage bc infinite loop

- support for replacing $$COMPONENT_BODY inside of javascript too, since its only supported in html right now with `<!-- $$COMPONENT_BODY -->`

- integrate this into the web packing tool:

```html
<!-- should automatically place at the VERY START of the head tag -->
<style>
    html {
        visibility: hidden;
        background: #1e1e1e;
    }
</style>

<!-- should place at the VERY END of the body tag -->
<script>window.addEventListener("load", () => document.documentElement.style.visibility = "visible");</script>
```

- create a JSON schema for `sklair.json` files:
- https://json-schema.org/understanding-json-schema/reference/index.html
- note that any file paths like `input` and `output` are RELATIVE to the sklair.json file
- allow components to be FULL html files with a head and body. If it is a full file (as opposed to just a regular component body that is bare), then basically build a "cache" of things from the head that will be inserted into the source document (deduplication) and the rest of the body just gets inserted as usual