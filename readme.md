==========================

Fork from [colinmarc/hdfs](https://github.com/colinmarc/hdfs)， it provides rpc call with hdfs， and simplie hdfs cmd, ye! it's very good.

In our environment, we often get/put data which some mr result or models stored at hdfs to Ceph(fs， s3). In normal, original hadoop hdfs is enough and can satisfy our needs. But if Ceph cluster load high, it export data to ceph will  spend much time, because hdfs read data will not gother to 4M or Multiple of 4M(the basic data cell of ceph),  little write and multi requect will increase burden of ceph mds or osd.

So, I use bufio to resolve this problems,  In some extreme case, I test, this can increase more than 30 times.

For convenience, we support export hdfs to s3.



## Cmd Usage

set hadoop config path and operate user(default bash login user) by system env.

```shell
export HADOOP_CONF_DIR="/data/hadooptest.NameNode/bin/taf-hadoop/etc/hadoop/"
export HADOOP_USER_NAME=mqq

```

Usage：

```shell
  ls [-lah] [FILE]...
  rm [-rf] FILE...
  mv [-nT] SOURCE... DEST
  mkdir [-p] FILE...
  touch [-amc] FILE...
  chmod [-R] OCTAL-MODE FILE...
  chown [-R] OWNER[:GROUP] FILE...
  cat SOURCE...
  head [-n LINES | -c BYTES] SOURCE...
  tail [-n LINES | -c BYTES] SOURCE...
  du [-sh] FILE...
  checksum FILE...
  get SOURCE [DEST]
  cfget SOURCE [DEST]			# added to export data to cephfs
  getmerge SOURCE DEST			
  cfgetmerge SOURCE DEST		#added to export data to cephfs
  put SOURCE DEST
  s3 [-asbenq] SOURCE  			# added to export data to s3
  df [-h]
```

Parameters:

- cfget: export data to the ceph directory, here use bufio to improve copy efficiency
- s3： export data to s3
  - -a， --access_key:   s3 access key
  - -s ，--secret_key： s3 secret key
  - -b，--bucket： s3 bucket，support s3 folder，like s3://test_bucket,  s3://test_bucket/test_folder/-
  - -e，--endpoint:   s3 endpoint url  http://s3test.some.com, support IP http://1.1.1.1
  - -n，--workers： Upload the number of concurrent, the default is 10
  - -q，--exitquick： If there is an error during the upload process, it will exit immediately and the default will continue.
- cfgetmerge： export the data and integrate it into a file to ceph

Others params passed, similar to the hadoop command.

### 优化效果

找了一个线上负载非常高的集群实际测试下

```shell
# 这里是测试数据2G
# ./gohdfs/hdfs ls -lh /ceph_test/testfile_big
-rw-r--r-- mqq  mqq  2.0G Nov 20 16:11 /ceph_test/testfile_big

#  本地拷数据到ceph, ceph挂载完毕
#  time cp ./testdata/testfile testfile3
real	2m56.701s
user	0m0.000s
sys		0m1.784s

#  使用我们优化过的命令拷数据到ceph
#  time ./gohdfs/hdfs cfget /ceph_test/testfile_big testfile_big 

real	1m29.539s
user	0m0.764s
sys		0m2.528s

# 使用hadoop原生命令拷数据到ceph
# time /data/hadooptest.NameNode/bin/taf-hadoop/bin/hadoop fs -get /ceph_test/testfile_big  testfile_ori

real	48m30.887s
user	0m14.916s
sys		0m6.316s

# 检查文件完整性
f211c32724db4f1b690ef28aa39fc1f7  testfile3
f211c32724db4f1b690ef28aa39fc1f7  testfile_big
f211c32724db4f1b690ef28aa39fc1f7  testfile_ori

#  使用我们工具拷数据到s3
# time ./gohdfs/hdfs s3 -a TJ1EAD2CNV9Zxxxx -s Qs6NSUWSf5gRYxxxxxxxxxxx -e http://s3test.some.com  -b  s3://xixi -n 20  /ceph_test/testfile_big

real	0m16.516s
user	0m12.208s
sys		0m2.392s

# 每次测试前刷系统缓冲
echo 3 > /proc/sys/vm/drop_caches && sync
```

从上面的结果来看

- 优化过的工具copy速度和系统cp命令差不多，甚至还快一些
- hadoop原生命令相当慢
- 上传到我们测试集群的对象存储集群是很快的

