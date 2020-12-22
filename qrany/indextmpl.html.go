package main

import "html/template"

var indextmpl = template.Must(template.New("indextmpl").Parse(`<html>
  <head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Download Files</title>
    <style>
      .done { text-decoration: line-through; }
      button {
          font-size: 1.2em;
          cursor: pointer;
      }
      .file {
          border: none;
          background: none;
      }
    </style> 
  </head>
  <body>
    <div><button onclick="downloadAll()">Download all files</button></div>
    <ol>
      {{range .}}<li><button class="file" onclick="startDownload(this)" mypath="{{.Path}}" myname="{{.Name}}">{{.Name}}</button>
      </li>
      {{end -}}
    </ol>
    
    <script>
      async function downloadAll() {
        let els = document.querySelectorAll(".file")
        for (let el of els) {
          startDownload(el, "1")
          await downloadComplete()
          el.classList.add("done")
        }
      }

      async function startDownload(btn, auto="0") {
        let path = btn.getAttribute("mypath")
        let name = btn.getAttribute("myname")
` + 
	"        console.log(`start download ${name}`)\n" + 
`        let a = document.createElement('a');
` + 
	"        a.setAttribute(\"href\", `/download?f=${path}&auto=${auto}`)\n" + 
`        a.setAttribute("download", name)
        a.style.display = 'none'
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a);
      }

      async function downloadComplete() {
        let resp = await fetch("/status")
        if (!resp.ok) {
          console.log("fetch status err"); console.log(resp)
          throw "fetch status"
        }
        let t = await resp.text()
        console.log("done " + t)
      }
    </script>
  </body>
</html>
`))
