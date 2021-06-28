# uniex
Golang library with an unified interface for stock market and cryptocurrency exchanges. (Basically CCTX for Golang).

Implement a single interface in your code instead of having to deal with each exchange requests and data models.

## Supported Exchanges

> Kraken

TODO Add list of methods

> Kucoin

TODO Add list of methods

## Usage


TODO


## Development

### Testing

Since multiple endpoints are private, you need an API Key for each exchange test suite you want to run.

Current format is:

`EXCHANGENAME_API_KEY` for the API Key

`EXCHANGENAME_SECRET_KEY` for the Secret Key

Suggestion is to create .env file:

```
KRAKEN_API_KEY=<Your Key>
KRAKEN_API_SECRET=<Your Secret>
```

And export it in your terminal

`export $(grep -v '^#' .env | xargs)`


**NEVER COMMIT A FILE THAT CONTAINS YOUR API KEYS**. 
The keys will be public and I guarantee that your funds will be stolen.

**Run tests**

`make test`