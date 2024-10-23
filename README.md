#Degrees of Separation

With cinema going global these days, every one of the [A-Z]ollywoods are now connected. Use the wealth of data available at [Moviebuff](http://www.moviebuff.com) to see how. 

Write a Go program that behaves the following way:

```
$ degrees amitabh-bachchan robert-de-niro

Degrees of Separation: 3

1. Movie: The Great Gatsby
Supporting Actor: Amitabh Bachchan
Actor: Leonardo DiCaprio

2. Movie: The Wolf of Wall Street
Actor: Leonardo DiCaprio
Director: Martin Scorsese

3. Movie: Taxi Driver
Director: Martin Scorsese
Actor: Robert De Niro
```

Your solution should use the Moviebuff data available to figure out the smallest degree of separation between the two people. 
All the inputs should be Moviebuff URLs for their respective people: For Amitabh Bachchan, his page is on http://www.moviebuff.com/amitabh-bachchan and his Moviebuff URL is `amitabh-bachchan`.

Please do not attempt to scrape the Moviebuff website - All the data is available on an S3 bucket in an easy to parse JSON format here: `https://data.moviebuff.com/{moviebuff_url}`

To solve the example above, your solution would fetch at least the following:

http://data.moviebuff.com/amitabh-bachchan

http://data.moviebuff.com/the-great-gatsby

http://data.moviebuff.com/leonardo-dicaprio

http://data.moviebuff.com/the-wolf-of-wall-street

http://data.moviebuff.com/martin-scorsese

http://data.moviebuff.com/taxi-driver

##Notes
* If you receive HTTP errors when trying to fetch the data, that might be the CDN throttling you. Luckily, Go has some very elegant idioms for rate limiting :)
* There may be a discrepancy in some cases where a movie appears on an actor's list but not vice versa. This usually happens when we edit data while exporting it, so feel free to either ignore these mismatches or handle them in some way.

Write a program in any language you want (If you're here from Gophercon, use Go :D) that does this. Feel free to make your own input and output format / command line tool / GUI / Webservice / whatever you want. Feel free to hold the dataset in whatever structure you want, but try not to use external databases - as far as possible stick to your langauage without bringing in MySQL/Postgres/MongoDB/Redis/Etc.

To submit a solution, fork this repo and send a Pull Request on Github.

For any questions or clarifications, raise an issue on this repo and we'll answer your questions as fast as we can.
