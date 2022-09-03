# xkcdsearcher
Downloads XKCD descriptions and creates a searchable offline index

## Usage
To see the program's help menu:
```bash
go run . -h
```

Sync the descriptions of *all* available XKCD comics from the web and store them locally:
```bash
go run . -u
```

To search all locally stored XKCD comics:
```bash
go run . -s "bobby tables"
```

To get the description of a particular XKCD by number (tries to retrieve from local storage first):
```bash
go run . -n 523
```

To get a random XKCD (tries to retrieve from local storage first):
```bash
go run . -r
```

## Examples
```bash
go run . -u

Latest available comic:         2664
Latest locally stored comic:    2661
Successfully downloaded 3 missing comics from web and updated local cache
```

```bash
go run . -s "bobby tables"

Nr      Title & Alternative Text
327     Exploits of a Mom
        Her daughter is named Help I'm trapped in a driver's license factory.
```

```bash
go run . -r

Returning 1433 from local storage

Nr      Title & Alternative Text
1433    Lightsaber
        A long time in the future, in a galaxy far, far, away, astronomers in the year 2008 sight an unusual gamma-ray burst originating from somewhere far across the universe.
```