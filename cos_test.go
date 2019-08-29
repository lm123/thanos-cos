package cos

import (
	"os"
	"strings"
	"flag"
	"context"
        "testing"

	"github.com/go-kit/kit/log"
	"github.com/improbable-eng/thanos/pkg/objstore"
	"github.com/improbable-eng/thanos/pkg/objstore/cos"
	"github.com/improbable-eng/thanos/pkg/testutil"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

//go test -v -args -srcdir=${SRCDIR} -dstdir=${DSTDIR} -oper=upload|download
var srcdir = flag.String("srcdir", "", "dir of source")
var dstdir = flag.String("dstdir", "", "dir of destination")
var oper = flag.String("oper", "upload", "operation")

func NewTestBucket(t testing.TB) (objstore.Bucket, func(), error) {
        c := cos.Config{
                Bucket:    os.Getenv("COS_BUCKET"),
                AppId:     os.Getenv("COS_APP_ID"),
                Region:    os.Getenv("COS_REGION"),
                SecretId:  os.Getenv("COS_SECRET_ID"),
                SecretKey: os.Getenv("COS_SECRET_KEY"),
        }

	if c.AppId == "" ||
           c.Region == "" ||
           c.SecretId == "" ||
           c.SecretKey == "" {
           return nil, nil, errors.New("insufficient cos configuration information")
        }

        if c.Bucket != "" {
                if os.Getenv("THANOS_ALLOW_EXISTING_BUCKET_USE") == "" {
                        return nil, nil, errors.New("COS_BUCKET is defined. Normally this tests will create temporary bucket " +
                                "and delete it after test. Unset COS_BUCKET env variable to use default logic. If you really want to run " +
                                "tests against provided (NOT USED!) bucket, set THANOS_ALLOW_EXISTING_BUCKET_USE=true. WARNING: That bucket " +
                                "needs to be manually cleared. This means that it is only useful to run one test in a time. This is due " +
                                "to safety (accidentally pointing prod bucket for test) as well as COS not being fully strong consistent.")
                }

                bc, err := yaml.Marshal(c)
                if err != nil {
                        return nil, nil, err
                }

                b, err := cos.NewBucket(log.NewNopLogger(), bc, "thanos-e2e-test")
                if err != nil {
                        t.Log(err.Error())
                }

                err = b.Iter(context.Background(), "", func(f string) error {
                        return errors.Errorf("bucket %s is not empty", c.Bucket)
                })
                if err != nil {
                        return b, nil, errors.Wrapf(err, "cos check bucket %s", c.Bucket)
                }
                t.Log("WARNING. Reusing", c.Bucket, "COS bucket for COS tests. Manual cleanup afterwards is required")
                return b, func() {}, nil
        } else {
		return nil, nil, errors.New("Bucket is NULL")
	}
}

func TestObjStore(t *testing.T) {
    t.Run("Tencent cos", func(t *testing.T) {
	t.Log("start")

        bkt, _, err := NewTestBucket(t)

	if err != nil {
	    t.Log(err.Error())
	}

	if bkt == nil {
	    t.Log("Bucket is NULL")
	    t.FailNow()
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	src := *srcdir
	dst := *dstdir
        op := *oper

	if strings.Compare(op, "upload") == 0 {
	  _, err = os.Stat(src)
	  if err != nil {
	    t.Logf("%s is not existed", src)

	    t.FailNow()
	  }

	  err = objstore.UploadDir(ctx, logger, bkt, src, dst)

	  if err != nil {
	    t.Logf("UploadDir %s to %s failed, %s", src, dst, err.Error())
	  }

	} else { // download
	  _, err = os.Stat(dst)
          if err != nil {
            t.Logf("%s is not existed", dst)

            t.FailNow()
          }

	  err = objstore.DownloadDir(ctx, logger, bkt, src, dst)

          if err != nil {
            t.Logf("DownloadDir %s to %s failed, %s", src, dst, err.Error())
          }

	}

	t.Log("end")

	testutil.Ok(t, err)
    })
}

