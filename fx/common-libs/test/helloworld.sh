grpcurl \
    -plaintext \
    -format=json \
    -proto=api/helloworld.proto \
    -d @ \
    localhost:50051 \
    helloworld.Greeter.SayHello << EOM
{
  "name": "$1"
}
EOM
