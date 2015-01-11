# Key Maps


## Development

### Dependencies

- (go) [https://github.com/go-martini/martini]()
- (go) [https://github.com/rubenv/sql-migrate]()
- (node) [https://github.com/gulpjs/gulp]()


### Setting up

1. install dependencies listed above
2. `npm install`
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
