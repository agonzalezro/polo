polo
====

[![Build Status](https://travis-ci.org/agonzalezro/polo.png)](https://travis-ci.org/agonzalezro/polo)

**Disclaimer**: this project is in a really early state. You can generate your
own blog and navigate between pages but I am sure that you are going to miss
some other functions. If you want a simple blog, you CAN use it without any
problem, but if you are looking for something else, this is not for you.

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

**This is going to change soon to be compatible with the [markdown of
pelican](http://docs.getpelican.com/en/3.1.1/getting_started.html#file-metadata).**

I am using markdown only. Whatever thing that is supported by [blackfriday
library](https://github.com/russross/blackfriday) is supported here. The only
difference is that I am adding some metadata to the files. For example, if you
want to define the date for the file:

    ---
    date: Monday, 4th April XXYY
    tags: place, Mancha, name
    ---

    And here is just the title.

It's really important that you let a line between the ``---`` and the beginning
of your article. But remember, as I said before, this is going to change soon.
