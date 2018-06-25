# zzdm
pure golang aes encryption/decryption

[Dependency]

protobuf(https://github.com/google/protobuf) & gogoprotobuf

Install gogoprootbuf

go get github.com/gogo/protobuf/proto

go get github.com/gogo/protobuf/protoc-gen-gogo

go get github.com/gogo/protobuf/gogoproto

go get github.com/gogo/protobuf/protoc-gen-gofast

[protobuf IDL]

[code]
syntax="proto3";

package zzdm;

option optimize_for=SPEED;

message Header{

    int64 frames=1;
    
    bytes name=2;
    
    bool secret=3;
    
}

message Frame{

    bytes iv=1;
    
    bytes data=2;
    
    uint32 hash=3;
    
}

