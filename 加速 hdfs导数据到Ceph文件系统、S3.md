> 导语 业务常将训练结果或者模型导出到ceph、但有时反馈到数据导出到ceph文件系统会很慢，特别在ceph集群负载很高的情况下，真的是这样吗？ 也许是使用姿势不对哦。另外为满足用户导数据到S3，我们也打通了hdfs导出到S3的数据流通。

## Hdfs数据读写

相较于传统磁盘有有扇区，在分布式文件系统中，常常将存储最小数据单元抽象为一个数据块。因为在分布式存储中可以存储超过一整块大磁盘的文件，将数据分块处理有利用并发操作数据提升存储读写效率，另外对数据的容错能力、高可用也会比较好管理。在hdfs中默认数据块大小为**128M**, 之前为64M。

### Hdfs中的Namenode和datanode

NameNode： 管理文件系统的命名空间、维护文件系统数及整棵树内所有文件和目录。一般为主备模式。

dataNode： 文件系统的存储节点， 根据需要存储并检索数据块，并定期向NameNode中发送自己存储块的列表

hdfs中一般一个NameNode带多个DataNode，数据以三副本存储。



![img](http://km.oa.com/files/photos/pictures/201811/1542889141_80_w874_h604.png)



### 读流程



![img](http://km.oa.com/files/photos/pictures/201811/1542889186_21_w830_h461.png)



1. 客户端通过rpc调用Namnode
2. 获取文件起始块位置， 对于每一块，NameNode返回该块副本的dataNode地址，并且datanode根据它们与客户端的距离(同一机架、机房、数据中心)来去最近节点
3. 通过read调用，连接最近datanode
4. 读取数据，如果中间有datanode读取出错，则读取最近副本。
5. 读取到数据末尾，在读取另外一个数据块
6. close关闭数据

### 写流程



![img](http://km.oa.com/files/photos/pictures/201811/1542889236_90_w913_h525.png)



1. 客户端通过create rpc调用namenode
2. namenode在系统空间中新建文件，此时该文件还没有相应数据块。
3. 客户端写入数据时， 选出合适存储数据副本的一组datanode, 并要求namenode分配新的数据块。
4. 数据串行写入数据流管道(主->副本)，但整个过程异步。
5. 写入确认，回包
6. 关闭文件，告知Namenode文件写入完成。

## Ceph数据读写

同hdfs一样，Ceph中也有数据块的概念，其最小的数据单元为**4M**。只有两者区别在下面会进行简单比较。

Ceph底层数据有rados统一存储对象，上层在rados对象存储基础上抽象出文件系统cephfs、块存储rbd、类s3对象存储。rados提供统一数据分布，容错、并具有管理功能。

Ceph中的组件：

OSd存储节点为实际存储数据的节点，OSD的存储引擎目前有Filestore和blustore两种，Filestore一般在操作系统原有的文件系统上存储单个数据对象，文件系统通常为xfs， 后者bluestore为ceph自研的存储引擎，用于小文件存储及性能优化。

mon节点监控osd的状态，并存储有整个集群存储规则、整个集群拓扑。

上层抽象出来不同存储形式组件： cephfs文件系统组件需要元数据服务器mds， mds多为主备，也可以多主。 块存储rbd， 对象存储rgw解析s3协议并存储数据。



![img](http://km.oa.com/files/photos/pictures/201811/1542889256_71_w693_h490.png)



### 数据读写流程



![img](http://km.oa.com/files/photos/pictures/201811/1542889270_61_w1187_h589.png)



Ceph的数据分布使用crush算法，一种分片的hash算法，只不过会加上权重和错误域两个选择因子。

数据读写时首先向mon请求集群的拓扑和crush规则，加上自己对象名和要写入数据池做hash运算，获得pg。

什么是pg，为什么会有pg，在ceph中上层的三种存储接口的抽象可以说什么基本覆盖了大部分的文件存储场景，不管大的小的对象。这就导致ceph底层rados会拥有巨量的对象数，如果通过普通的hash直接映射到存储节点osd上某个对象，则非常不好管理，例如存储减点扩缩容、down掉等，mon节点发现后，会查询需要搬迁的对象也是巨量的，往往需要很长时间并且性能不一定扛不住。但是如果中间抽象出来一层，这一层数量是固定的那么hash规则比较稳定，在中间抽象层去管理对象将会好很多。所以ceph中数据管理（迁移，分配）是以pg为单位的。

pg下会有存储节点，为保证数据高可用性，一个pg下会有不同的存储策略。像hdfs一样，ceph一般的数据高可用也是由多副本来存储，一般为2-3副本。另外ceph也提供类似raid5为纠删码存储技术。

## s3(RGW)数据流读写

在ceph上层可以抽象出类s3的对象存储系统，ceph文件系统由于元数据服务mds多为主备模式提供服务。mds的负载能力是有限的，类似有Namenode在小文件巨量文件系统下也会有同样的痛点。另外通过Ceph内核模块直接挂载在本机挂载点，虽然用户可以像本地文件系统一样方便的使用，但是在和mds通信和交互过程中、或者Ceph内核模块出bug了将会导致挂载的机器直接崩溃掉。

由于S3协议是http为基础，http为无状态的请求，因此可以横向扩展。

性能比文件系统一样略差，便利性方便也不会像使用文件系统可直接挂载使用，当然也无非是连接、put、get几个操作而已。

## 原因与解决

回到正题， 业务常在大数据平台上训练数据或者模型，想把结果导出到ceph，以供线上或者共享使用。但是有时候会发现导出数据很慢，特别是在Ceph集群负载非常高的时候。

为什么呢，上文有提到分布式文件系统中数据存储中多以块为数据单元存储，在网络传输过程和生成数据的过程往往以数据流的方式发送或者读取写入而且通常数据生成和传输速度不一致，比如tcp nagle算法，hdfs中client常常将一个数据块chunk攒为一个packet再发送出去。这就是问题所在，从hdfs读取数据并不会攒到ceph数据单元大小或者ceph数据单元大小的倍数时，数据就被发送发出或者被写入，这导致网络交互会很多，在集群负载很高的情况下(网络请求多，元数据服务mds负载大，osd IO性能下降)，ceph处理不过来，最终导致性能下降很厉害。

另外业务大多使用hadoop命令来操作hdfs的数据，hadoop命令调用需要启动jvm过程很重，这里使用直接向NameNode和DataNode发起rpc调用，会快不少。

golang语言rpc调用基础库已经由[colinmarc](https://github.com/colinmarc/hdfs)封装好了，在此感谢！但有部分接口有些问题在工具中已经修复。

在golang中read，write读写数据之后并不是向C中read、write一致，前者缓冲区大小并不一定会填满，即写入、读取部分数据调用就直接返回。这不是我们想要的，这里必须攒成Ceph数据单位的整数倍才可以解决上述的问题。

golang提供标准库bufio来实现这点，经过测试将缓冲区设为4M*4=16M速度更快些。bufio是先将数据保存在内存中，在只有达到指定大小或者调用flush时候才会从真正调用写操作。另外如果buf中数据超过预设大小，其余数据将直接调用write写入，最后不要忘了flush哦。在测试中我们发现使用bufio性能比直接cp到ceph文件系统差不多有时还快一些。

```go
err = client.Walk(source, func(p string, fi os.FileInfo, err error) error {
            ...
            err = client.CopyToCeph(p, fullDest)
            ...
    })

func (c *Client) CopyToCeph (src string, dst string) error {
    remote, err := c.Open(src)
    if err != nil {
        return err
    }
    defer remote.Close()

    local, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer local.Close()

    sizeTotal := 16 * 1024 * 1024
    size := 4 * 1024 * 1024
    buf := make([]byte, size)

    bw := bufio.NewWriterSize(local, sizeTotal)

    var written int64
    for {
        nr, er := remote.Read(buf)
        if nr > 0 {
            nw, ew := bw.Write(buf[0:nr])
            if nw > 0 {
                written += int64(nw)
            }
            if ew != nil {
                err = ew
                break
            }
            if nr != nw {
                err = io.ErrShortWrite
                break
            }
        }
        if er != nil {
            if er != io.EOF {
                err = er
            }
            break
        }
    }
    bw.Flush()
    return  err
}
```

另外直接发起rpc调用比hadoop原生命令也会快不少

```shell
# time /data/hadooptest.NameNode/bin/taf-hadoop/bin/hadoop fs -ls /
Found 11 items
drwxr-xr-x   - test test          0 2018-10-09 19:56 /benchmarks
...

real    0m1.741s
user    0m3.460s
sys        0m0.484s
# time  /tmp/rike/gohdfs/hdfs ls -l /
drwxr-xr-x test  test  0 Oct  9 19:56 benchmarks
drwxr-xr-x test  test  0 Nov 22 14:22 ceph_test
...

real    0m0.084s
user    0m0.020s
sys        0m0.008s
```

## 导出数据到S3

从hdfs文件、文件夹导出数据到S3中已经支持了，你需要指定s3需要的密钥access_key, serect_key，地址endpoint

另外我们允许用户指定参数并发上传提高数据copy效率，对大文件使用分片并发上传。

```go
type req struct {
    key   string
    path  string
    isDir bool
}

type resp struct {
    key string
    err error
}

func uploadFile(sess *session.Session, bucket string, q chan req, p chan resp, client *hdfs.Client, wg *sync.WaitGroup) {
    uploader := s3manager.NewUploader(sess)
    for r := range q {
        if r.key == "" {
            continue
        }
        var rst resp
        var err error
        rst.key = r.key
        if ! r.isDir {
            remote, err := client.Open(r.path)
            if err != nil {
                rst.err = err
            }
            _, err = uploader.Upload(&s3manager.UploadInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(r.key),
                Body:   remote,
            })
            remote.Close()
        } else {
            buffer := &aws.WriteAtBuffer{}
            remote := bytes.NewReader(buffer.Bytes())
            _, err = uploader.Upload(&s3manager.UploadInput{
                Bucket: aws.String(bucket),
                Key:    aws.String(r.key),
                Body:   remote,
            })
        }

        rst.err = err
        p <- rst
    }
    wg.Done()
}

func startClients(workerNums uint, q chan req, p chan resp, client *hdfs.Client, ak, sk, bucket, endpoint string) {
    cfg := &aws.Config{
        Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
        Endpoint:         aws.String(endpoint),
        Region:           aws.String("default"),
        DisableSSL:       aws.Bool(true),
        S3ForcePathStyle: aws.Bool(true),
        //LogLevel:         aws.LogLevel(aws.LogDebug),
    }
    sess := session.New(cfg)
    wg := &sync.WaitGroup{}
    var i uint
    for i = 0; i < workerNums; i++ {
        wg.Add(1)
        go uploadFile(sess, bucket, q, p, client, wg)
    }
    wg.Wait()
    close(p)
}

func s3put(paths []string, ak, sk, bucket, s3url string, workerNums uint, exitQuick bool) {
    sources, nn, err := normalizePaths(paths[0:1])
    if err != nil {
        fatal(err)
    }

    source := sources[0]
    u, err := url.Parse(bucket)
    if err != nil || u.Scheme != "s3" {
        fatal("bucket url error")
    }

    s3Bucket := u.Host
    s3Path := strings.TrimPrefix(u.Path, "/")

    client, err := getClient(nn)
    if err != nil {
        fatal(err)
    }

    reqQueue := make(chan req, 1000)
    resqQueue := make(chan resp, 1000)

    go func() {
        err = client.Walk(source, func(p string, fi os.FileInfo, err error) error {
            if err != nil {
                fatal(err)
            }
            var key string
            if p == source && ! fi.IsDir() {
                key = filepath.Join(s3Path, filepath.Base(source))
            } else {
                key = filepath.Join(s3Path, strings.TrimPrefix(p, source))
            }

            if key == "" {
                return nil
            }
            if key[0] == '/' {
                key = key[1:]
            }
            r := req{
                key:  key,
                path: p,
            }

            if fi.IsDir() {
                r.isDir = true
            }
            reqQueue <- r
            return nil

        })
        if err != nil {
            fatal(err)
        }
        close(reqQueue)
    }()

    go startClients(workerNums, reqQueue, resqQueue, client, *s3ak, *s3sk, s3Bucket, *s3endpoint)

    for rst := range resqQueue {
        if rst.err != nil {
            fatal(fmt.Sprintf("key: %s, error: %s\n", rst.key, rst.err))
        }
    }
}
```

## 命令使用

环境变量指定hadoop配置路径和操作用户，用户默认是本地shell登录的用户

```shell
export HADOOP_CONF_DIR="/data/hadooptest.NameNode/bin/taf-hadoop/etc/hadoop/"
export HADOOP_USER_NAME=test
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
  cfget SOURCE [DEST]
  getmerge SOURCE DEST
  cfgetmerge SOURCE DEST
  put SOURCE DEST
  s3 [-asbenq] SOURCE  
  df [-h]
```

主要介绍:

- cfget: 拉起数据到ceph目录，这里使用bufio提升copy效率
- s3： 导数据到s3
  - -a， —access_key: s3 access key
  - -s ，—secret_key： s3 secret key
  - -b，—bucket： s3 bucket，支持路径，一般为s3://test_bucket, s3://test_bucket/test_folder/-
  - -e，—endpoint: s3 endpoint 地址 [http://s3test.rike.com，支持IPhttp://1.1.1.1](http://s3test.sumeru.xn--mig%2Ciphttp-qp0u36r//1.1.1.1)
  - -n，—workers： 上传并发数，默认为10
  - -q，—exitquick： 上传过程中有错误是否立马退出，默认继续。
- cfgetmerge： 拉起数据并整合为一个文件到ceph

其他不解释了，和hadoop命令差不多。

### 优化效果

找了一个线上负载非常高的集群实际测试下

```shell
# 这里是测试数据2G
# /tmp/rike/gohdfs/hdfs ls -lh /ceph_test/testfile_big
-rw-r--r-- test  test  2.0G Nov 20 16:11 /ceph_test/testfile_big

#  本地拷数据到ceph
#  time cp /tmp/rike/testdata/testfile testfile3
real    2m56.701s
user    0m0.000s
sys        0m1.784s

#  使用我们优化过的命令拷数据到ceph
#  time /tmp/rike/gohdfs/hdfs cfget /ceph_test/testfile_big testfile_big 

real    1m29.539s
user    0m0.764s
sys        0m2.528s

# 使用hadoop原生命令拷数据到ceph
# time /data/hadooptest.NameNode/bin/taf-hadoop/bin/hadoop fs -get /ceph_test/testfile_big  testfile_ori

real    48m30.887s
user    0m14.916s
sys        0m6.316s

# 检查文件完整性
f211c32724db4f1b690ef28aa39fc1f7  testfile3
f211c32724db4f1b690ef28aa39fc1f7  testfile_big
f211c32724db4f1b690ef28aa39fc1f7  testfile_ori

#  使用我们工具拷数据到s3
# time /tmp/rike/gohdfs/hdfs s3 -a TJ1xxxxxxys -s QxxxxxxxmR -e http://s3test.rike.com  -b  s3://xixi -n 20  /ceph_test/testfile_big

real    0m16.516s
user    0m12.208s
sys        0m2.392s

# 每次测试前刷系统缓冲
echo 3 > /proc/sys/vm/drop_caches && sync
```

从上面的结果来看

- 优化过的工具copy速度和系统cp命令差不多，甚至还快一些
- hadoop原生命令相当慢
- 上传到我们对象存储的测试集群是很快的

## 番外： HDFS和Ceph对比

简单对比下Hdfs和Ceph，对hdfs不是很熟所以如果有不恰当的地方请指出。先说架构，架构决定使用场景，最后从运维角度来比较下两者

- 数据分布
  - 数据分布方式： DNS中查询方式有两种，一种是递归查询被访问DNS服务器自己查询到结果最后返回给客户端，另一种是迭代查询，被访问的DNS服务器只告诉真正该去查询DNS服务器如果自己没有的话。hdfs和ceph提供给client元数据方式和这个很像。 hdfs NameNode管理所有数据块的分布、分配它是hdfs系统最核心的部分。Ceph的客户端只会向Mon节点索要整个集群的拓扑和crush规则，自己计算出文件改存去哪里，而且本地会缓存Mon返回的结果，只有在集群拓扑或者Crush改变时才会重新计算。hdfs即使内存数据结构做的完美，小文件合并再多，它的元数据承载量也是有限的，也不会承载多少量的访问。这一架构理论上限制了hdfs不是无限扩展的，而Ceph理论上可以无限的扩展。
  - 错误域
    - 数据高可用之数据容错：数据容错大致为两种副本和纠删码，hdfs只支持副本，Ceph支持副本和纠删码，所以Ceph在同量存储下同样容错能力下存储利用率上可以比Hdfs高。
    - 数据高可用之数据地理容错：默认情况下Hdfs三副本数据，第一、二副本存于同机架不同机器上，第三副本存储在另外一个机架上。Ceph则通过crush hash规则制定错误域，自由度和支持错误域更加丰富。
  - 权重（容量和性能）： Hdfs中的DataNode会定期上报自己使用容量和已有容量，Namenode在数据均衡，分配数据块时作为参考，可以做到数据比较均衡，但不会考虑性能。Ceph有weight参数会在hash作为选择存储位置的因子，即权重越大的被选择的几率就越大，存储数据也就越多。weight是在创建Osd是由脚本根据容量自动加到crush中， 在实际使用中可以根据存储介质性能手动调整。可说两者都没有智能从存储介质的容量和性能决定因子从而影响数据分布和集群的负载均衡。
  - 数据均衡： Hdfs中NameNode在均衡和分配数据块时能够权衡datanode的实际情况，并且NameNode在dataNode数据容量相差到阀值时配置均衡器均衡数据。而Ceph由于hash规则在一开始就已经确定下来，如果osd权重配置不合理或者集群对象较少或者每个对象容量相差很大，那么数据将会很不均匀。
  - 数据冷热分层存储： hdfs支持分层存储，ceph以pool为单位，提供Cache Tier实现分层存储，但是Cache Tier这玩意真不好用。
- 数据读写：
  - 读写过程: Hdfs中client常见的读写方式是将数据写入到临时buf或者文件中，当文件攒到一定大时候才发送出去，另外，Hdfs写入数据时候是串行写入，写入队列和确认队列都是串行发送和确认。写数据时三副本也是一个个写入并一个个确认最终才认为写入完成，读数据时数据块读取也是一个接一个读取。而Ceph在处理请求时大量使用队列和线程，写入时，先写主activeosd， 之后使用posix协议并行去写剩余的副本，读取可以从多个并行副本读取。只就导致了Hdfs只能用于对于时延要求不高且是大文件场景，而ceph应付小文件和延迟要求性高比较合适。
  - posix支持：由于Hdfs设计主要使用一次写入多次读取的场景，Hdfs只支持追加写。而Ceph的文件系统cephfs基本完全兼容posix协议。这使得Cephfs使用场景更广。
- 使用场景：hadoop的分而治之的思想貌似也影响这Hdfs的架构，大容量文本、日志等文件拆成不同的数据块进行分析和读写，在上面架构部分也可看到hadoop更合适大文件存储和读写时延要求不高的场景。而Ceph的三种使用接口让它使用面很广。
- 集群管理：
  - 在上面的架构中NameNode可以智能均衡数据，在Hdfs中加入新的DataNode后，nameNode可以很轻松的把数据均衡开来，运维成本很低。
  - 但是Ceph是根据hash规则来分布数据，即使Crush算法已经比较智能，但在OSD不多情况下涉及搬迁的数据依然很多，此时对正常数据访问影响很大。另外Ceph中每个osd上pg也需要控制在一定的数量一般100-200，pg太少数据分布将不均，pg太多将消耗osd更多的资源。但在一般的场景下起初集群规模比较小可能就100个osd左右，在之后数据不断上涨到必须扩容时pg数据分配到osd上就会减少，此时需要调整pg数量才能保证数据再次相对均衡，但是pg调整有涉及到hash规则改变，并又伴随大量数据搬迁，也会影响线上数据访问。
  - 机器故障或者存储节点故障时： Hdfs Datanode会定期上报自己状态，Namenode周期检查心跳如果Datanode挂掉，那么NameNode将屏蔽该故障节点，并启动数据均衡和副本保持为设定值。Ceph在Osd down之后crush hash将重新计算，随后发起数据迁移，但是Ceph在pg状态不健康时并且请求量很大的情况下在很长一段时间内将导致部分请求访问很慢，甚至会被被block住。在集群运行很长一段时间后，磁盘批量出现损坏概率很大，这点上非常影响集群的稳定，这也许需要IO调度才能解决这样的问题。