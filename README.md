a bk-cmdb data collector demo

### download source code

- via go get
```
go get github.com/wusendong/cmdb_hostsnap
```

- or via git clone
```
mkdir $GOPATH/src/github.com/wusendong
cd $GOPATH/src/github.com/wusendong
git clone github.com/wusendong/cmdb_hostsnap
```

### quick run 
```
cd $GOPATH/src/github.com/wusendong/cmdb_hostsnap
go build
./cmdb_hostsnap -c cmdb_hostsnap.json
```

### usage
```
NAME:
   hostsnap - hostsnap

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     reload   reload config
     stop     stop the hostsnap process
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d             enable debug logging level [$CMDB_DEBUG]
   --config FILE, -c FILE  Load configuration form FILE [$CMDB_HOSTSNAP_CONFIG]
   --help, -h              show help
   --version, -v           print the version
```

### build
```
make linux
```

### package

```
make package
```

### development

see docs [ comming soon ]
