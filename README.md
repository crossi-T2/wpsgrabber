# wpsgrabber

wpsgrabber is a tool for watching processing reports from a [52north WPS server](https://52north.org). For completed executions, it produceds a different encoding (CSV + XML files) for further analysis.

## Usage

First define a *config.yaml* YAML configuration file:

```yaml
# Root directory from where start watching for reports
RootDir: /tmp/foo

# Directory where to write output
OutputDir: /tmp/bar

# Optional. If set, the tool will scan for reports
# with a modification time greater than ScanFrom.
# It is useful to convert reports in the past.
# Setting the zero value for Go's time, that is
# January 1, year 1, 00:00:00 UTC
# is equivalent to NOT setting the parameter
ScanFrom: 1983-05-22T14:13:00Z

# Optional. If set, the tool will override the value
# provided in the WPS Execute Response document
ProcessIdentifier: nextgeoss-sentinel2-biopar

# Optional. If set, the tool will override the value
# provided in the WPS Execute Response document
ProcessVersion: 1.4
```

Get the code, build and run it:

```bash
git clone https://github.com/crossi-T2/wpsgrabber
cd wpsgrabber
go build -o wpsgrabber cmd/wpsgrabber/*.go
./wpsgrabber -config /path/to/the/config.yaml
```

## Installation via RPM

Instructions here.

## Related Projects

* [52north](https://52north.org)
* [fsnotify](https://github.com/fsnotify/fsnotify)
