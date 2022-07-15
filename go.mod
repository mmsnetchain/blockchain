module mmschainnewaccount

go 1.18

replace github.com/prestonTao/keystore => ../keystore

replace github.com/prestonTao/keystore/kstore => ../keystore/kstore

replace github.com/prestonTao/libp2parea => ../libp2parea

replace github.com/prestonTao/libp2parea/protos/go_protos => ../libp2parea/protos/go_protos

replace github.com/prestonTao/utils => ../utils

require (
	github.com/Jeiwan/eos-b58 v0.0.0-20180918133445-43bbe264af4a
	github.com/antlabs/timer v0.0.5
	github.com/astaxie/beego v1.12.3
	github.com/btcsuite/btcd v0.22.1
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1
	github.com/btcsuite/goleveldb v1.0.0
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-xorm/xorm v0.7.9
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/hyahm/golog v0.0.0-20220504012620-eb402ee1ab9d
	github.com/json-iterator/go v1.1.12
	github.com/logoove/sqlite v1.15.3
	github.com/mr-tron/base58 v1.2.0
	github.com/prestonTao/keystore v0.0.0-20220518074257-aa1ed9263f0e
	github.com/prestonTao/keystore/kstore v0.0.0-00010101000000-000000000000
	github.com/prestonTao/libp2parea v0.0.0-00010101000000-000000000000
	github.com/prestonTao/libp2parea/protos/go_protos v0.0.0-00010101000000-000000000000
	github.com/prestonTao/utils v0.0.0-20220517073115-06a4120f1fc7
	github.com/shirou/gopsutil/v3 v3.22.4
	github.com/smartystreets/goconvey v1.7.2
	github.com/syndtr/goleveldb v1.0.0
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898
	gopkg.in/gookit/color.v1 v1.1.6
)

require (
	github.com/antlabs/stl v0.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/snappy-go v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/cupcake/rdb v0.0.0-20161107195141-43ba34106c76 // indirect
	github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/ledisdb/ledisdb v0.0.0-20200510135210-d35789ec47e6 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monnand/dhkx v0.0.0-20180522003156-9e5b033f1ac4 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.7.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/siddontang/go v0.0.0-20170517070808-cb568a3e5cc0 // indirect
	github.com/siddontang/rdb v0.0.0-20150307021120-fc89ed2e418d // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	lukechampine.com/uint128 v1.1.1 // indirect
	modernc.org/cc/v3 v3.36.0 // indirect
	modernc.org/ccgo/v3 v3.16.6 // indirect
	modernc.org/libc v1.16.8 // indirect
	modernc.org/mathutil v1.4.1 // indirect
	modernc.org/memory v1.1.1 // indirect
	modernc.org/opt v0.1.1 // indirect
	modernc.org/sqlite v1.17.3 // indirect
	modernc.org/strutil v1.1.1 // indirect
	modernc.org/token v1.0.0 // indirect
	xorm.io/builder v0.3.6 // indirect
	xorm.io/core v0.7.2-0.20190928055935-90aeac8d08eb // indirect
)
