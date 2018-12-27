all: secp256k1 tmm

tmm:
	go install github.com/tokenme/tmm

secp256k1:
	cp -r dependencies/secp256k1/src vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1/src;
	cp -r dependencies/secp256k1/include vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1/include;

install:
	cp vendor/github.com/xluohome/phonedata/phone.dat /etc/phone.dat;
	cp uaparser.yml /etc/uaparser.yml;
	cp gse.dict /etc/gse.dict;
	cp ip2region.db /etc/ip2region.db;
	cp -f /opt/go/bin/tmm /usr/local/bin/;
	chmod a+x /usr/local/bin/tmm;
