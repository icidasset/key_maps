# Key Maps


## Development

### Dependencies

- (go) [https://github.com/gocraft/web]()
- (go) [https://github.com/rubenv/sql-migrate]()
- (node) [https://github.com/gulpjs/gulp]()


### Setting up

1. `script/go_get`
2. `npm install -g gulp` and `npm install`
3. make postgres database `keymaps_development`
4. `script/migrate`

```bash
# start/watch server
script/server

# build & watch assets
gulp

# build assets for production
gulp build --production
```
