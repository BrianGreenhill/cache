# cache

components:
- mysql database server
- memcached server
- cache application
- observability server

## Usage

```bash
vagrant up
```

`localhost:3000` to view `grafana`.

## Servers
- `10.0.0.23` - memcache server listening on `11211`
- `10.0.0.24` - mysql server listening on `3306`
- `10.0.0.25` - cache application
- `10.0.0.26` - observability stack qryn, grafana `3000` for grafana interface
