# get_stats

When testing websites I collect sample urls from the webserver logs and use them to replay them in my development

So I have a file like:

```
/a/c/b
/a/d/f/g
...
```

And I call it like this

```bash
$ get_stats --prefix http://localhost:3000 urls.txt
```

and it will log the responses giving the minimum, average and maximum response times and sizes along with a list of status codes and how often they occurred

If anything needs to be added to the end of the url then the `--suffix` flag is available

If you want it to halt the moment that an error occurs (status 5xx) then add the `--dropdead` flag
