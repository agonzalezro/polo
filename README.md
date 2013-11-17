THIS PROJECT ISN'T WORKING YET
==============================

For now, I am going to use the ``README.rst`` just to write down my notes and
perhaps point to there to somebody that can help me with them.

What's it?
----------

It's a static blog rendering tool. Yes, I know that there a lot of them out
there but I just want mine to learn a little bit of Go coding.

Thinks to take care of
----------------------

- check if 2 slugs matches, perhaps show a warning, or just Panic.

TODO (to make it functional)
----------------------------

- ~~render the rst in html format, I couldn't continue doing this because I
  couldn't found any library for reStructuredText for Go :(~~ I will use
  markdown instead, I was loosing a lot of precious time for nothing.

~~In the case that it's not possible to render RST files, I will need to migrate
my posts to markdown, which will be crappy, but better than nothing.~~

- configuration yaml file:

  + title of the blog

  + author

  + favicon

  + google analytics id

  + disqus id for the comments

  + github login for the fork icon

  + sharethis id

- pagination of the articles.

~~ - render pages too, not just articles (example ``about.html``). ~~


Nice extras (preference order)
------------------------------

#. themes (if an override template exist use it, if not, fallback to the
   default).

#. search option, even with Google could be ok.

#. allow draft articles.

#. support markdown.

Format
------

I am using markdown now. Whatever thing that is supported by [blackfriday
library](https://github.com/russross/blackfriday) is supported here. The only
difference is that I am adding some metadata to the files. For example, if you
want to define the date for the file:

```
---
date: Monday, 4th April XXYY
---

And here is just the title.
```

It's really important that you let a line between the ``---`` and the beginning
of your article.
