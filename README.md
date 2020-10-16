A dark-themed HTTP error page generator.

This is originally created to be used with `nginx`, but the error pages are generic, they can be plugged into any sort of webserver. This just generates static HTML, like:

![](https://sean.fish/favicon.ico)

... like...

```
```

On all my servers, I'm typically running a reverse proxy, doing a `proxy_pass` to some other location. I don't want to have a singular `nginx` error page or depend on a `nginx` module, in case this has to be used with something else.

Since this generates individual pages also means it can be plugged/routed into any other webserver, disregarding language.

## Generate HTML files

### Install

`go get gitlab.com/seanbreckenridge/darker_errors`

---

`darker_errors` implements a small template language.

It replaces the strings:

* `STATUS_CODE` (e.g. 404)
* `STATUS_MSG` (e.g. Not Found)

With the corresponding HTTP values.

---

To override the default text for each page, you can use replacement directives:

* `ERROR_TITLE` (Text in `<title>`; default: `STATUS_CODE - STATUS_MSG`
* `ERROR_HEADING` (Large Heading; default: `<h2>STATUS_CODE</h2>`)
* `ERROR_MSG` (Message; default: `<p>STATUS_MSG<p>`)

You can also inject arbitrary HTML by setting one of the following:

* `ERROR_HEAD` (insert HTML into the `<head>` tag)
* `ERROR_BEFORE_HEADING` (before `ERROR_HEADING`)
* `ERROR_AFTER_HEADING` (after `ERROR_HEADING`, before `ERROR_MSG`)
* `ERROR_AFTER_MSG` (after `ERROR_MSG`)

All positional arguments to `darker_errors` are interpreted as replacement directives.

### Examples:

Include website name in `<title>`:

```
darker_errors 'ERROR_TITLE:MyWebsite - STATUS_MSG'
```

If you want to modify the text for just one page, you can specify that by using the HTTP status code

Include text to go back home on a 404:

```
darker_errors '404:ERROR_AFTER_MSG:<p>Click <a href="/">here</a> to go back home.<p>`
```

To left-align text:

```
darker_errors 'ERROR_HEAD:<style>p { text-align: left; }</style>'
```

Refresh the page every 2 seconds if the user encountered a `502` error:

```
darker_errors
  '502:ERROR_HEAD:<meta http-equiv="refresh" content="2">'
  '502:ERROR_MSG:<p>This page is currently being updated...<br />It should reload when it's done automatically</p>'
```

If you specify `502:ERROR_MSG`, and `ERROR_MSG`, the `502` overwrites the generic replacement.

          _ _ _     _   _     _     
  ___  __| (_) |_  | |_| |__ (_)___ 
 / _ \/ _` | | __| | __| '_ \| / __|
|  __/ (_| | | |_  | |_| | | | \__ \
 \___|\__,_|_|\__|  \__|_| |_|_|___/

It also generates a `template.html` file, in-case you want to insert a title/message with a webserver at runtime. That could be changed to whatever templating language you're using if you want, or by reading the file in, and running a search/replace with your content for `||ERROR_TITLE||`, `||ERROR_HEADER||`, `||ERROR_MSG||`.


## nginx setup

Most of the time, I use nginx like:

```
# base elixir server handles everything
  location / {
    include /etc/nginx/pheonix_params;
    proxy_pass http://localhost:8082;
  }
```

.. i.e. every request which doesn't match some other location block gets sent to another 'base' webserver. That means its handles/renders its own HTTP errors, and nginx just forwards that back.

However, whenever that upstream server is down, or restarting, the user just gets the blinding `502 Bad Gateway` page. To avoid that, I typically explicitly mark at least the 502 error code in my `nginx` configuration, with the generated page from here.

To do that, this assumes you've set a `root` in a `server` directive somewhere. That can be `/var/www/html/`, or some other folder.

As an example:

```
server {
  listen [::]:443 ssl ...
  listen 443 ssl ...
  ...
  root /var/www/html;

  error 502 = /502.html

  location / {
    proxy_pass ...
  }
}
```

In the case above, since `root` is `/var/www/html`, put the generated `502.html` file at `/var/www/html/502.html`.

If you want to use this *only* with nginx, you could put the `error_html` folder in `/var/www/html`, and then map each error code to the HTML page, like:

```
error 401 = /error_html/401.html
error 404 = /error_html/404.html
error 502 = /error_html/502.html
```

To generate the configuration for that, you can run:

```
darker_errors -nginx-conf
```
