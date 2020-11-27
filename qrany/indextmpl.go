package main

const indexTmpl = `<html>
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
    </style> 
  </head>
  <body>
    <div><button onclick="downloadAll()">Download all files</button></div>
    <ol>
      {{range .}}<li>
				<a class="file" href="/download?f={{.Path}}" download="{{.Name}}">{{.Name}}</a>
			</li>
      {{end -}}
		</ol>
		
    <script>
			function downloadAll() {
        let els = document.querySelectorAll(".file")
        for (let el of els) {
          download(el) 
          el.classList.add("done")
				}
			}
      
      async function download(el) {
				el.click()
				await downloadComplete(el.href)
			}

			async function downloadComplete(href) {
				let resp = await fetch("/status")
				if (!resp.ok) {
					throw href
				}
        console.log(await resp.text())
			}
    </script>
  </body>
</html>
`
