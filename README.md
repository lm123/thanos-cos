# thanos-cos
## Project Purpose
   Demo how to use Thanos Object Storage interface to upload or download files and directories from Tencent COS

## Set the configuration relating with Tencent COS HTTP API invoking
   export COS_BUCKET=bucket

   export COS_APP_ID=1300103488

   export COS_REGION=ap-beijing

   export COS_SECRET_ID=...

   export COS_SECRET_KEY=...

## Build with dep
   dep init

   dep ensure

## Use go test to run
   go test -v -args -srcdir=/root/lm/gopath/src/objstore/data/01DE4BDG1AFRM20FMSB65V3KQ7 -dstdir=thanos/test

   go test -v -args -srcdir=thanos/test -dstdir=/root/lm/gopath/src/objstore/data/download -oper=download
