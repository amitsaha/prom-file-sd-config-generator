# Prometheus File SD config generator

[![Build Status](https://travis-ci.org/amitsaha/prom-file-sd-config-generator.svg?branch=master)](https://travis-ci.org/amitsaha/prom-file-sd-config-generator)

This is a utility program which will generate a file SD config (in JSON) for Prometheus
given a target HTTP URL which has a list of the different targets to scrape.

If your infrastructure has a large number of targets which you cannot specify via one
of the prometheus service discovery mechanisms, you need to resort to use the file SD
config. However, manually editing the manually SD config is a chore especially
when you have a growing number of such targets. Why not have a centralized list of
these targets and then generate the file SD config? And once you have the list centralized,
you can add any new target to it and this program will automatically generate a new file SD config
which will have the new target for prometheus to scrape.

## Usage

The program has three options (two of them optional):

```
$ prom-file-sd-config 
Usage of ./prom-file-sd-config:
  -config-path string
    	Path of the SD config JSON file (default "./file_sd_config.json")
  -scrape-interval int
    	Interval in seconds between consecutive scrapes (default 5)
  -target-source string
    	HTTP URL of the target source
```

The only required argument is the `target-source` which is the HTTP resource which acts
as the centralized repository of the targets. An example of such a page is:

```
<a href="http://127.0.0.1:9100/bar1/metrics">target1</a>
<a href="http://127.0.0.1:9100/bar2/metrics">target2</a>
<a href="http://127.0.0.1:9100/bar3">target3</a>
```

The generated file SD config will be:

```
[
  {
    "targets": ["127.0.0.1:9100"],
    "labels": {
      "__metrics_path__": "/bar1/metrics"
    }
  },
  {
    "targets": ["127.0.0.1:9100"],
    "labels": {
      "__metrics_path__": "/bar2/metrics
    }
  },
  {
    "targets": ["127.0.0.1:9100"],
    "labels": {
      "__metrics_path__": "/bar3"
    }
  }
]

```

## Development

Run tests:

1. Start a fake http server serving URLs for sraping
2. Utility scrapes it, and generates a JSON file
3. Test reads the JSON file and verifies it has the expected data


## Deployment

## LICENSE

Apache (See [LICENSE](./LICENSE))
