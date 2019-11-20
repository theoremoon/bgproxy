# bgproxy

A proxy simluates Blue-Green deploy.


## Installation

```
$ go get github.com/theoremoon/bgproxy/cmd/bgproxy
$ go get github.com/theoremoon/bgproxy/cmd/bgproxyctl
```

## Usage

Situation:

- We have a *blue* server (php) listening on `localhost:8888`
- We want to publish the service on `:80`
- We may replace the server with updated one

```
$ sudo bgproxy -addr 0.0.0.0:80 -blue http://localhost:8888/ -stop "kill \$(ps -ao cmd,pid|grep '^php -S localhost:8888'|awk '{print \$NF}')"
bgproxy:2019/11/15 00:17:21 Start
...
```

(Or more simply we can use -cmd option)
```
$ sudo bgproxy -addr 0.0.0.0:80 -blue http://localhost:8888/ -cmd "php -S localhost:8888 -t /var/www/html/"
bgproxy:2019/11/15 00:17:21 Start
```


Now we are serving the *blue* server at `:80`.




Next, we want to release new server (*green*). Let's up the server listening on `localhost:8889`. And order `bgproxy` to set new *green* server.

```
$ bgproxyctl green -addr http://localhost:8889/ -stop "kill \$(ps -ao cmd,pid|grep '^php -S localhost:8889'|awk '{print \$NF}')"
OK
```

(Also bgproxyctl provides -cmd option)
```
$ bgproxyctl green -addr http://localhost:8889/ -cmd "php -S localhost:8889 -t /home/user/public/"
OK
```

Then we serves a *green* server at `:80`. However if the *green* server has a problem, we can roll it back.

```
$ bgproxyctl rollback
```

`bgproxy` automatically stops the *green* server and switch back to serve *blue*.

Or if the *green* sever is unhealthy, `bpproxy` automatically roll back.


If the *green* server is perfect and completely healthy until 1 hour, then `bgproxy` stops the *blue* and turn *green* into new *blue*.


...In more details, you can see help and the source code.

## Author

theoremoon

## LICENSE

Apache 2.0
