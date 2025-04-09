# FIAP data source plugin for Grafana

This is a grafana data source that uses IEEE1888 (a.k.a. FIAP in Japan).

## Usage

### Datasource Settings

![DatasourceSettings](https://raw.githubusercontent.com/SIOS-Technology-Inc/grafana-fiap-datasource/refs/heads/main/src/img/settings.png)

| Setting         | Detail                                                                                                       |
| --------------- | ------------------------------------------------------------------------------------------------------------ |
| URL             | URI of the server to connect (with port)                                                                     |
| Server timezone | If the server only handles queries by certain timezone, its timezone <br> Default: UTC <br> Format: `+09:00` |

### Query Settings

![QuerySettings](https://raw.githubusercontent.com/SIOS-Technology-Inc/grafana-fiap-datasource/refs/heads/main/src/img/query.png)

| Setting                          | Detail                                                                                                                                                                                                                             |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Point ID                         | Equivalent to FIAP key class: `id` <br> Only one entry per line <br> Point IDs are combined into a single FIAP query, and sent as a single FETCH request.                                                                          |
| - Button                         | Delete the specified Point ID                                                                                                                                                                                                      |
| + Button                         | Insert a new Point ID field below                                                                                                                                                                                                  |
| Data range                       | Equivalent to FIAP key class: `select` <br> Choose: Period, Latest, or Oldest                                                                                                                                                      |
| Period                           | Fetch data in the range specified by **Start/End time** (Equivalent to no `select` option)                                                                                                                                         |
| Latest                           | Fetch one of the data in the range specified by **Start/End time** that is the latest (Equivalent to `select="maximum"` option)                                                                                                    |
| Oldest                           | Fetch one of the data in the range specified by **Start/End time** that is the oldest (Equivalent to `select="minimum"` option)                                                                                                    |
| Start/End time                   | Equivalent to FIAP key class: `gteq`/`lteq` <br> Format: `2006-01-02 15:04:05` <br> If the time part is omitted, `00:00:00` is completed <br> **Server timezone** taken into account in [DatasourceSettings](#datasource-settings) |
| sync with grafana start/end time | If checked, start/end time will be synchronized with the Time Range in the Grafana Dashboard <br> (**Start/End time** field is disabled)                                                                                           |

## Others
The client implementation uses：
[go-fiap-client](https://pkg.go.dev/github.com/SIOS-Technology-Inc/go-fiap-client)
