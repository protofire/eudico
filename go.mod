module github.com/filecoin-project/lotus

go 1.19

retract v1.14.0 // Accidentally force-pushed tag, use v1.14.1+ instead.

retract v1.20.2 // Wrongfully cherry picked PR, use v1.20.2+ instead.

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.0
	github.com/BurntSushi/toml v1.1.0
	github.com/DataDog/zstd v1.4.5
	github.com/GeertJohan/go.rice v1.0.3
	github.com/Gurpartap/async v0.0.0-20180927173644-4f7f499dd9ee
	github.com/Kubuxu/imtui v0.0.0-20210401140320-41663d68d0fa
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/alecthomas/jsonschema v0.0.0-20200530073317-71f438968921
	github.com/buger/goterm v1.0.3
	github.com/chzyer/readline v1.5.0
	github.com/consensus-shipyard/go-ipc-types v0.1.6-0.20230605095729-e7269af8a509
	github.com/containerd/cgroups v1.0.4
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/detailyang/go-fallocate v0.0.0-20180908115635-432fa640bd2e
	github.com/dgraph-io/badger/v2 v2.2007.4
	github.com/docker/go-units v0.5.0
	github.com/drand/drand v1.4.9
	github.com/drand/kyber v1.1.15
	github.com/dustin/go-humanize v1.0.0
	github.com/elastic/go-elasticsearch/v7 v7.14.0
	github.com/elastic/go-sysinfo v1.7.0
	github.com/elastic/gosigar v0.14.2
	github.com/etclabscore/go-openrpc-reflect v0.0.36
	github.com/fatih/color v1.13.0
	github.com/filecoin-project/dagstore v0.5.2
	github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200910194244-f640612a1a1f
	github.com/filecoin-project/go-address v1.1.0
	github.com/filecoin-project/go-amt-ipld/v4 v4.0.0
	github.com/filecoin-project/go-bitfield v0.2.4
	github.com/filecoin-project/go-cbor-util v0.0.1
	github.com/filecoin-project/go-commp-utils v0.1.3
	github.com/filecoin-project/go-crypto v0.0.1
	github.com/filecoin-project/go-data-transfer/v2 v2.0.0-rc4
	github.com/filecoin-project/go-fil-commcid v0.1.0
	github.com/filecoin-project/go-fil-commp-hashhash v0.1.0
	github.com/filecoin-project/go-fil-markets v1.27.0-rc1
	github.com/filecoin-project/go-jsonrpc v0.2.3
	github.com/filecoin-project/go-padreader v0.0.1
	github.com/filecoin-project/go-paramfetch v0.0.4
	github.com/filecoin-project/go-state-types v0.10.0
	github.com/filecoin-project/go-statemachine v1.0.3
	github.com/filecoin-project/go-statestore v0.2.0
	github.com/filecoin-project/go-storedcounter v0.1.0
	github.com/filecoin-project/mir v0.3.5
	github.com/filecoin-project/pubsub v1.0.0
	github.com/filecoin-project/specs-actors v0.9.15
	github.com/filecoin-project/specs-actors/v2 v2.3.6
	github.com/filecoin-project/specs-actors/v3 v3.1.2
	github.com/filecoin-project/specs-actors/v4 v4.0.2
	github.com/filecoin-project/specs-actors/v5 v5.0.6
	github.com/filecoin-project/specs-actors/v6 v6.0.2
	github.com/filecoin-project/specs-actors/v7 v7.0.1
	github.com/filecoin-project/specs-actors/v8 v8.0.1
	github.com/filecoin-project/test-vectors/schema v0.0.5
	github.com/gbrlsnchs/jwt/v3 v3.0.1
	github.com/gdamore/tcell/v2 v2.2.0
	github.com/go-openapi/spec v0.19.11
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/rpc v1.2.0
	github.com/gorilla/websocket v1.5.0
	github.com/gregdhill/go-openrpc v0.0.0-20220114144539-ae6f44720487
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/hannahhoward/go-pubsub v0.0.0-20200423002714-8d62886cc36e
	github.com/hashicorp/go-hclog v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru/v2 v2.0.2
	github.com/hashicorp/raft v1.3.10
	github.com/hashicorp/raft-boltdb v0.0.0-20171010151810-6e5ba93211ea
	github.com/icza/backscanner v0.0.0-20210726202459-ac2ffc679f94
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab
	github.com/ipfs/bbloom v0.0.4
	github.com/ipfs/go-blockservice v0.5.0
	github.com/ipfs/go-cid v0.4.0
	github.com/ipfs/go-cidutil v0.1.0
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-ds-badger2 v0.1.3
	github.com/ipfs/go-ds-leveldb v0.5.0
	github.com/ipfs/go-ds-measure v0.2.0
	github.com/ipfs/go-fs-lock v0.0.7
	github.com/ipfs/go-graphsync v0.14.3
	github.com/ipfs/go-ipfs-blockstore v1.2.0
	github.com/ipfs/go-ipfs-blocksutil v0.0.1
	github.com/ipfs/go-ipfs-chunker v0.0.5
	github.com/ipfs/go-ipfs-ds-help v1.1.0
	github.com/ipfs/go-ipfs-exchange-interface v0.2.0
	github.com/ipfs/go-ipfs-exchange-offline v0.3.0
	github.com/ipfs/go-ipfs-http-client v0.5.0
	github.com/ipfs/go-ipfs-routing v0.3.0
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipld-cbor v0.0.6
	github.com/ipfs/go-ipld-format v0.4.0
	github.com/ipfs/go-libipfs v0.7.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/ipfs/go-merkledag v0.9.0
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/go-unixfs v0.4.3
	github.com/ipfs/go-unixfsnode v1.5.2
	github.com/ipfs/interface-go-ipfs-core v0.11.1
	github.com/ipld/go-car v0.5.0
	github.com/ipld/go-car/v2 v2.7.0
	github.com/ipld/go-codec-dagpb v1.5.0
	github.com/ipld/go-ipld-prime v0.20.0
	github.com/ipld/go-ipld-selector-text-lite v0.0.1
	github.com/ipni/index-provider v0.11.0
	github.com/ipni/storetheindex v0.5.10
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/koalacxr/quantile v0.0.1
	github.com/libp2p/go-buffer-pool v0.1.0
	github.com/libp2p/go-libp2p v0.26.2
	github.com/libp2p/go-libp2p-consensus v0.0.1
	github.com/libp2p/go-libp2p-gorpc v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.21.0
	github.com/libp2p/go-libp2p-pubsub v0.9.3
	github.com/libp2p/go-libp2p-raft v0.4.0
	github.com/libp2p/go-libp2p-record v0.2.0
	github.com/libp2p/go-libp2p-routing-helpers v0.4.0
	github.com/libp2p/go-maddr-filter v0.1.0
	github.com/libp2p/go-msgio v0.3.0
	github.com/mattn/go-isatty v0.0.17
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-base32 v0.1.0
	github.com/multiformats/go-multiaddr v0.8.0
	github.com/multiformats/go-multiaddr-dns v0.3.1
	github.com/multiformats/go-multibase v0.1.1
	github.com/multiformats/go-multihash v0.2.1
	github.com/multiformats/go-varint v0.0.7
	github.com/open-rpc/meta-schema v0.0.0-20201029221707-1b72ef2ea333
	github.com/polydawn/refmt v0.89.0
	github.com/prometheus/client_golang v1.14.0
	github.com/raulk/clock v1.1.0
	github.com/raulk/go-watchdog v1.3.0
	github.com/stretchr/testify v1.8.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/urfave/cli/v2 v2.16.3
	github.com/whyrusleeping/bencher v0.0.0-20190829221104-bb6607aa8bba
	github.com/whyrusleeping/cbor-gen v0.0.0-20230126041949-52956bd4c9aa
	github.com/whyrusleeping/ledger-filecoin-go v0.9.1-0.20201010031517-c3dcc1bddce4
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xorcare/golden v0.6.1-0.20191112154924-b87f686d7542
	github.com/zyedidia/generic v1.2.1
	go.opencensus.io v0.24.0
	go.opentelemetry.io/otel v1.12.0
	go.opentelemetry.io/otel/bridge/opencensus v0.33.0
	go.opentelemetry.io/otel/exporters/jaeger v1.2.0
	go.opentelemetry.io/otel/sdk v1.11.1
	go.uber.org/atomic v1.10.0
	go.uber.org/fx v1.19.2
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.5.0
	golang.org/x/exp v0.0.0-20230129154200-a960b3787bd2
	golang.org/x/net v0.7.0
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.5.0
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9
	golang.org/x/tools v0.3.0
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2
	google.golang.org/protobuf v1.28.1
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/GeertJohan/go.incremental v1.0.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/akavel/rsrc v0.8.0 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/armon/go-metrics v0.3.9 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bep/debounce v1.2.1 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cilium/ebpf v0.4.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/daaku/go.zipexe v1.0.2 // indirect
	github.com/dave/jennifer v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/drand/kyber-bls12381 v0.2.3 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/etclabscore/go-jsonschema-walk v0.0.6 // indirect
	github.com/filecoin-project/go-amt-ipld/v2 v2.1.0 // indirect
	github.com/filecoin-project/go-amt-ipld/v3 v3.1.0 // indirect
	github.com/filecoin-project/go-commp-utils/nonffi v0.0.0-20220905160352-62059082a837 // indirect
	github.com/filecoin-project/go-ds-versioning v0.1.2 // indirect
	github.com/filecoin-project/go-hamt-ipld v0.1.5 // indirect
	github.com/filecoin-project/go-hamt-ipld/v2 v2.0.0 // indirect
	github.com/filecoin-project/go-hamt-ipld/v3 v3.1.0 // indirect
	github.com/flynn/noise v1.0.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-openapi/jsonpointer v0.19.3 // indirect
	github.com/go-openapi/jsonreference v0.19.4 // indirect
	github.com/go-openapi/swag v0.19.11 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/pprof v0.0.0-20221203041831-ce31453925ec // indirect
	github.com/hannahhoward/cbor-gen-for v0.0.0-20230214144701-5d17c9d5243c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/golang-lru v0.6.0 // indirect
	github.com/huin/goupnp v1.0.3 // indirect
	github.com/iancoleman/orderedmap v0.1.0 // indirect
	github.com/ipfs/go-bitfield v1.1.0 // indirect
	github.com/ipfs/go-block-format v0.1.1 // indirect
	github.com/ipfs/go-filestore v1.2.0 // indirect
	github.com/ipfs/go-ipfs-cmds v0.8.2 // indirect
	github.com/ipfs/go-ipfs-delay v0.0.1 // indirect
	github.com/ipfs/go-ipfs-files v0.3.0 // indirect
	github.com/ipfs/go-ipfs-posinfo v0.0.1 // indirect
	github.com/ipfs/go-ipfs-pq v0.0.3 // indirect
	github.com/ipfs/go-ipld-legacy v0.1.1 // indirect
	github.com/ipfs/go-ipns v0.3.0 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-path v0.3.1 // indirect
	github.com/ipfs/go-peertaskqueue v0.8.1 // indirect
	github.com/ipfs/go-verifcid v0.0.2 // indirect
	github.com/ipld/go-ipld-adl-hamt v0.0.0-20220616142416-9004dbd839e0 // indirect
	github.com/ipsn/go-secp256k1 v0.0.0-20180726113642-9d62b9f0bc52 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/go-random v0.0.0-20190219211222-123a90aedc0c // indirect
	github.com/jbenet/go-temp-err-catcher v0.1.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/koron/go-ssdp v0.0.3 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.2.0 // indirect
	github.com/libp2p/go-libp2p-gostream v0.6.0 // indirect
	github.com/libp2p/go-libp2p-kbucket v0.5.0 // indirect
	github.com/libp2p/go-nat v0.1.0 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
	github.com/libp2p/go-reuseport v0.2.0 // indirect
	github.com/libp2p/go-yamux/v3 v3.1.2 // indirect
	github.com/libp2p/go-yamux/v4 v4.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/magefile/mage v1.9.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multicodec v0.8.0 // indirect
	github.com/multiformats/go-multistream v0.4.1 // indirect
	github.com/nikkolasg/hexjson v0.1.0 // indirect
	github.com/nkovacs/streamquote v1.0.0 // indirect
	github.com/onsi/ginkgo/v2 v2.5.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/petar/GoLLRB v0.0.0-20210522233825-ae3b015fd3e9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.2.1 // indirect
	github.com/quic-go/qtls-go1-20 v0.1.1 // indirect
	github.com/quic-go/quic-go v0.33.0 // indirect
	github.com/quic-go/webtransport-go v0.5.2 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/tidwall/gjson v1.14.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.0.1 // indirect
	github.com/whyrusleeping/cbor v0.0.0-20171005072247-63513f603b11 // indirect
	github.com/whyrusleeping/chunker v0.0.0-20181014151217-fe64bd25879f // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/zondax/hid v0.9.1 // indirect
	github.com/zondax/ledger-go v0.12.1 // indirect
	go.opentelemetry.io/otel/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.12.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	go4.org v0.0.0-20200411211856-f5505b9728dd // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3 // indirect
	google.golang.org/grpc v1.50.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

replace github.com/filecoin-project/test-vectors => ./extern/test-vectors

// uncomment for local development. Keeping it here for now for convenience.
// replace github.com/consensus-shipyard/go-ipc-types => ../go-ipc-types

// uncomment for local development. Keeping it here for now for convenience.
// replace github.com/filecoin-project/mir => ../tmp/mir
