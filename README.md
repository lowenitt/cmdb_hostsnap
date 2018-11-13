a bk-cmdb data collector demo

### init

- via go get
```
go get github.com/wusendong/cmdb_hostsnap
```

- via git clone
```
mkdir $GOPATH/src/github.com/wusendong
cd $GOPATH/src/github.com/wusendong
git clone github.com/wusendong/cmdb_hostsnap
```

### quick run 
```
go build
./cmdb_hostsnap -c cmdb_hostsnap.json
```

### build
```
make linux
```

### package

```
make package
```

### run

```
go run github.com/wusendong/cmdb_hostsnap -c cmdb_hostsnap.json
```
### reload

```
go run github.com/wusendong/cmdb_hostsnap reload
```
### stop

```
go run github.com/wusendong/cmdb_hostsnap stop
```

### development

see docs [ comming soon ]
