GO Twitter Trends
=================

In order to make this work add a `config.json` file with the format
```
{
    "APIKey": {{ your_API_key }},
    "APISecret": {{ your_API_secret_key }},
    "AccessToken": {{ your_access_token }},
    "AccessSecret": {{ your_access_secret_token }},
    "WOEID": {{ WOEID_for_the_desired_location }}
}
```

The `/` route will get a POST with your computers width and return something like this:

```
[
    {
        "w":200,
        "h":200,
        "x":0,
        "y":0,
        "name":"Real Madrid",
        "url":"http://twitter.com/search?q=%22Real+Madrid%22","tweet_volume":181083
    },
]
```

This gives all the data needed to position all the Trending Topic containters. The width and height change depending on how many people are tweeting about this.
go to the [htmlTwitterTrends](https://github.com/mauroeparis/htmlTwitterTrends) repo in order to see how it looks.
