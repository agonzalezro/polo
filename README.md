polo
====

polo is a static blog rendering tool created with Golang.

I'm happily using it on my blog: http://agonzalezro.github.io, which means that
works fine :)

Yes, I know that there a lot of them out there but I just want mine to learn a
little bit of Go coding.

How to use it
-------------

### Install

We have 2 ways of doing it, the easiest and the "little bit more complex":

#### Easy way: Download it

Just grab your binary from here http://gobuild.io/github.com/agonzalezro/polo
and copy it to wherever you want to have it. The binary by itself is all what
you need, not even a folder with templates unless you want to override the
defaults.

#### Build it yourself

Let's assume that you have already installed some other Go package, so, you
already have Go installed and `$GOPATH` in place:

	go install github.com/agonzalezro/polo

This will create a binary for you called `polo`:

    $ polo -help
    Usage: polo [options] sourcedir outdir

      -config="config.json": the settings file to create your site.
      -daemon=false: create a simple HTTP server after the blog is created to see the result
      -port=8080: port where to run the server

### Test (or if you are brave, production :)

If you want try it with the examples:

    $ polo -daemon examples /tmp/test

And now, you can go to http://localhost:8080 and see your generated blog.

Every time that you do a change on any of the files inside `examples`, or you
create/delete files the blog is going to be regenerated and you will just need
to reload the web.


Just markdown!
--------------

I am using markdown only. Whatever thing that is supported by [blackfriday
library](https://github.com/russross/blackfriday) is supported here. The only
difference is that I am adding some metadata to the files.

This metadata is using exactly the same format than the one used on Pelican or
Jekyll, but we don't support exactly the same keys. If you thing that some of
the keys that they support and we don't are needed, please, [create an
issue](https://github.com/agonzalezro/polo/issues/new) or send a pull request.

Supported tags:

- **title**: if it's not on the metadata info, the first line is going to be
  used to create it.
- **date**: format YYYY-MM-DD hh:mm
- **tags**: comma separated.
- **slug**: if it is not defined the first line is going to be slugified.
- **status**: if it's draft the page is not going to be rendered.
- **summary**: an introductory paragraph. It will be empty if the metadata tag
  is not defined.
- **author**: this will override the default author in the config file.

This is one auto explainable example (Pelican):

    Title: My super title
    Date: 2010-12-03 10:20
    Tags: thats, awesome
    Slug: my-super-post
    Author: Federico

    And here is just the content.

In this case we are overriding the title, slug & author.

If you prefer to do it with Jekyll, or you are migrating a Jekyll page:

    ---
    title: My super title
    date: 2010-12-03 10:20
    tags: thats, awesome
    slug: my-super-post
    author: Federico
    ---

    And here is just the content.

The keys are case insensitive in both cases.

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
- **shareThisPublisher**: the ShareThis publisher ID. If provided, there will
  be some social buttons on the article view.


Templating
----------

This functionality is in a kinda early stage. I think, that we will need to
split the templates in much more files, and this way we will be able to
override them quite easily, but that will require some work (PRs welcomed! :)

### Creating your own theme

You have a default theme on `templates/`. If you just installed polo this theme
is going to be part of the binary but it can be override.

In case that you want to override the theme, you don't need to provide ALL the
files. Just providing the ones that you are overriding is more than enough.

For example, imaging that you want to override the header to use another
bootstrap theme:

1. Wherever you want (but it needs to be the same place where you run polo
   from) you create the folder `templates/base`.
2. Then you edit `templates/base/header.html`, adding the following content:

````html
{{define "header"}}
  <link rel="stylesheet" href="/static/css/bootstrap.min.css">
  <link
    href="http://netdna.bootstrapcdn.com/bootstrap/3.0.0/css/bootstrap-glyphicons.css"
    rel="stylesheet">
{{end}}
````

3. Now you run polo from the folder that owns `templates/` and
4. PROFIT!

### Modifying the one that is going to be included on the binary.

If you want to do changes on the default theme, you need to remember that you
MUST recreate the binary data, you should do it this way:

    go-bindata -o templates/bindata.go -pkg=templates -ignore=bindata.go templates/...

Auto deploy
-----------

You can use your favourite CI/CD platform to generate the sites. I wrote a blog
post about how to do it with [CircleCI](http://circleci.com) here:
http://agonzalezro.github.io/how-to-automagically-generate-your-polo-blog-with-circleci.html
