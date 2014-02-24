polo
====

[![Build Status](https://travis-ci.org/agonzalezro/polo.png)](https://travis-ci.org/agonzalezro/polo)

**Disclaimer**: this project is in a really early state, but don't be scare of
using it. I am currently using it as the generator for my site:
http://agonzalezro.github.io and it works!

What's it?
----------

polo is a static blog rendering tool created with Golang.

Yes, I know that there a lot of them out there but I just want mine to learn a
little bit of Go coding.

How to use it
-------------

For now I am not providing binaries, you will need to compile it yourself, but
you can use the ``Makefile`` included:

    $ make

It will generate the file ``bin/polo``:

    $ bin/polo -help
    Usage of bin/polo:
      -config="config.json": the settings file to create your site.
      -input=".": path to your articles source files.
      -output=".": path where you want to creat the html files.

If you want try it with the examples:

    $ rm /tmp/db.sqlite;bin/polo -input examples -output /tmp
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

This is one auto explainable example:

    Title: My super title
    Date: 2010-12-03 10:20
    Tags: thats, awesome
    Slug: my-super-post

    And here is just the content.

In this case we are overriding the title and the slug.
