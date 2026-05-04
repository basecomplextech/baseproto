import (
    "pkg3/pkg3a"
)

options (
    go_package="github.com/basecomplextech/baseproto/compiler/internal/tests/pkg2"
)

message Submessage {
    key     string      1;
    value   pkg3a.Value 2;
}
