# wpsgrabber

wpsgrabber is a tool for watching processing reports from a [52north WPS server](https://52north.org) and producing simple CSV files for further analysis.

## Usage

First define a *config.json* JSON configuration file:

```json
{
     "RootDir": "/tmp/foo",
     "CSVOutputDir": "/tmp/csv",
     "ProcessIdentifier": "nextgeoss-sentinel2-biopar",
     "ProcessVersion": "1.4"
}
```


**Note**: The optional `ProcessIdentifier` and `ProcessVersion` overrides the ones provided by the WPS Execute Response document. 

Get the code and run:

```bash
git clone https://github.com/crossi-T2/wpsgrabber
cd wpsgrabber
go run main.go -c /path/to/the/config.json
```

## Stability

The tool is under active development and not ready for production.

## Related Projects

* [52north](https://52north.org)
* [fsnotify](https://github.com/fsnotify/fsnotify)
