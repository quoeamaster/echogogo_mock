# build project as plugin instead of executable or binary
rm mock-module.so
go build -buildmode=plugin -o mock-module.so
