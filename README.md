polo
====

[![Build Status](https://travis-ci.org/agonzalezro/polo.png)](https://travis-ci.org/agonzalezro/polo)

polo is a static blog rendering tool created with Golang.

I'm happily using it on my blog: http://agonzalezro.github.io, which means that
works fine :)

Yes, I know that there a lot of them out there but I just want mine to learn a
little bit of Go coding.


How to use it
-------------

### Install

Let's assume that you have already installed some other Go package, so, you
already have Go installed and `$GOPATH` in place:

	go get github.com/agonzalezro/polo

This will create a binary for you called `polo`:

    $ polo -help
    Usage of polo:
      -config="config.json": the settings file to create your site.
      -input=".": path to your articles source files.
      -output=".": path where you want to creat the html files.

### Test

If you want try it with the examples:

    $ rm /tmp/db.sqlite;go run -input examples -output /tmp
    $ cd /tmp
    $ python -m SimpleHTTPServer

**Note**: the ``rm /tmp/db.sqlite`` is because of a bug. I can not run the
database in memory yet, so, you will need to manually delete this file.

And now, you can go to http://localhost:8000 and see your generated blog.


Just markdown!
--------------

I am using markdown only. Whatever thing that is supported by [blackfriday
library](https://github.com/russross/blackfriday) is supported here. The only
difference is that I am adding some metadata to the files.

This metadata is using exactly the same format than the one used on Pelican,
but we don't care about all of it (if you do, let me know and I can think on
adding it). Supported tags:

- **title**: if it's not on the metadata info, the first line is going to be
  used to create it.
- **date**: format YYYY-MM-DD hh:mm
- **tags**: comma separated.
- **slug**: if it is not defined the first line is going to be slugified.
- **status**: if it's draft the page is not going to be rendered.
- **summary**: an introductory paragraph. It will be empty if the metadata tag
  is not defined.
- **author**: this will override the default author in the config file.

This is one auto explainable example:

    Title: My super title
    Date: 2010-12-03 10:20
    Tags: thats, awesome
    Slug: my-super-post
    Author: Federico

    And here is just the content.

In this case we are overriding the title and the slug.


Configuration
-------------

An example configuration can be found on the file `config.json`:
https://github.com/agonzalezro/polo/blob/master/config.json

This is what you can configure there:

- **author**: if it's not override with the Metadata it's the name that is
  going to be shown on the articles.
- **title**: title of the blog, for the `<title>` element and the header.
- **url**: sometimes the full url is needed.
- **show(Archive|Categories|Tags)**: if it's true the pages are going to be
  created and the links are going to be added.
- **paginationSize**: set it to -1 if you want to show all the posts.
- **favicon**: the favicon path if you have one.

### 3rd party

- **disqusSitename**: if you want comments on your blog.
- **googleAnalyticsId**: the Google Analytics ID.
- **sharethisPublisher**: it's not implemented yet.


Templating
----------

This functionality is in a kinda early stage. I think, that we will need to
split the templates in much more files, and this way we will be able to
override them quite easily, but that will require some work (PRs welcomed!) :)

You have a default theme on `templates/`, but you can create your own themes
creating the same folder struct in your page. Polo will default always to this
base theme (which is included on the binary).

Example: imagine that you have your markdown files on `$HOME/site`. You could
create a file `$HOME/site/templates/article.html` which would override our
default template for article rendering.

If you want to do changes on the default theme, you need to remember that you
need to recreate the binary data, you should do it this way:

	go-bindata -o templates/bindata.go \
            -pkg=templates \
            -ignore=bindata.go \
            templates templates/base
