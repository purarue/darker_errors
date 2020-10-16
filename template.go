package main

import (
	"html/template"
	"log"
)

func SetupTemplate() *template.Template {
	tmpl, err := template.New("darker_errors").Parse(
		`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><style>
html, body {
  margin: 0px;
  padding: 0px;
  border: 0px;
  width: 100%;
  min-height: 100vh
	background-color: #222;
	color: #eee;
}
main {
	display: flex;
  justify-content: center;
	align-items: center;
}
</style>
<title>{{ .Title }}</title>{{ .HeadHtml }}
</head>
<body>
  <main>
    <div id="error-container">{{ .BeforeHeading }}
      <div id="heading">{{ .Heading }}</div>{{ .AfterHeading }}
      <div id="message">{{ .Message }}</div>{{ .AfterMessage }}
    </div>
	</main>
</body>
</html>`)
	if err != nil {
		log.Fatalf("Could not create HTML template: %s\n", err)
	}
	return tmpl
}
