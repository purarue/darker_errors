package darker_errors

import (
	"html/template"
	"log"
)

func DarkTheme() *template.Template {
	tmpl, err := template.New("darker_errors").Parse(
		`<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><style>
html, body {
  margin: 0px;
  padding: 0px;
  border: 0px;
  width: 100%;
	height: 100vh;
	background-color: #121212;
	color: #eee;
	font-family: Verdana, Geneva, sans-serif;
}
#flexbox {
	height: 100%;
	width: 100%;
	display: flex;
	flex-flow: column nowrap;
	align-items: center;
}
#heading {
	font-size: 3rem;
}
#message {
	font-size: 1.5rem;
}
</style>
<title>{{ .Title }}</title>{{ .HeadHtml }}
</head>
<body>
	<div id="flexbox">
		{{ .BeforeHeading }}<div id="heading">{{ .Heading }}</div>{{ .AfterHeading }}<div id="message">{{ .Message }}</div>{{ .AfterMessage }}
	</div>
</body>
</html>`)
	if err != nil {
		log.Fatalf("Could not create HTML template: %s\n", err)
	}
	return tmpl
}
