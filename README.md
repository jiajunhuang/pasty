# pasty

pasty is a paste service build on gRPC.

## server

run it like this:

```bash
$ pasty -isServer -rpcAddr="0.0.0.0:9527"
```

if you omit the `-rpcAddr` part, it will listen at `127.0.0.1:9527` by default

create your token by running sqlite3 client and insert data into table `tokens` like this:

```sql
insert into tokens (id, created_at, updated_at, user_id, token) values (1,'2018-04-10 00:00:00','2018-04-10 00:00:00',1,'helloworld');
```

## client

create `~/.pasty.json`, and set:

```json
{
    "token": "helloworld",
    "rpc_addr": "127.0.0.1:9527"
}
```

replace `rpc_addr` to your server address.

and then, just use it!

```bash
jiajun@idea ~: type p
p is /home/jiajun/golang/bin/p
jiajun@idea ~: echo "hello" > test.txt
jiajun@idea ~: p < test.txt
jiajun@idea ~: p -limit 1
file id 3 created at 2018-04-10 23:19:03 +0800 CST      ==================
hello
file id 3 created at 2018-04-10 23:19:03 +0800 CST over ==================
jiajun@idea ~: p -limit 1 -hint=false
hello
jiajun@idea ~: echo "hello world" | p
jiajun@idea ~: p -limit 1 -hint=false
hello world
jiajun@idea ~: p -limit 2 -hint=false
hello world
hello
jiajun@idea ~: p -limit 2
file id 4 created at 2018-04-10 23:19:35 +0800 CST      ==================
hello world
file id 4 created at 2018-04-10 23:19:35 +0800 CST over ==================
file id 3 created at 2018-04-10 23:19:03 +0800 CST      ==================
hello
file id 3 created at 2018-04-10 23:19:03 +0800 CST over ==================
```
