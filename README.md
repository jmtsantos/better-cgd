# cgdbreak
A small app to check your Caixa Break meal card balance

## Usage

Just set the `USER` and `PASS` env variables a run the script

To run it with docker...

```
docker build . -t cgd-break
docker run -e USER=111 -e PASS=222 cgd-break
```

