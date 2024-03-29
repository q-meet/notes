# badger 学习笔记

我的学习仓库地址：[github.com](https://github.com/q-meet/badger_test)

BadgerDB 是一个非常特殊的包，从我们可以对它做出的最重要的改变不是在它的 API 上，而是在数据如何存储在磁盘上的角度来看。

这就是为什么我们遵循与语义版本控制不同的版本命名模式。

当磁盘上的数据格式以不兼容的方式发生变化时，就会发布新的主要版本。
每当 API 发生变化但保持数据兼容性时，就会发布新的次要版本。请注意，API 上的更改可能是向后不兼容的 - 与语义版本控制不同。
当数据格式和 API 没有变化时，就会发布新的补丁版本。
遵循以下规则：

v1.5.0 和 v1.6.0 可以在相同文件之上使用，无需担心，因为它们的主要版本相同，因此磁盘上的数据格式是兼容的。
正如其主要版本所暗示的那样，v1.6.0 和 v2.0.0 的数据不兼容，因此使用 v1.6.0 创建的文件需要先转换为新格式，然后才能被 v2.0.0 使用。
有关使用新版本控制命名模式背后原因的详细说明，您可以阅读 [VERSIONING.md](https://github.com/dgraph-io/badger/blob/master/VERSIONING.md)

Badger 的编写考虑了以下设计目标：

用纯 Go 编写一个键值数据库。
利用最新研究成果为跨 TB 的数据集构建最快的 KV 数据库。
针对 SSD 进行优化。
Badger 的设计基于题为WiscKey：在 SSD 感知存储中将键与值分离的论文。Badger 的设计目标是在 SSD 上提供与 RocksDB 相同的性能，但是在某些情况下，它可以提供更好的性能。

1 WISCKEY论文（Badger 所基于的论文）在将值与键分开方面取得了巨大胜利，与典型的 LSM 树相比，显着减少了写入放大。

2 RocksDB是LevelDB的SSD优化版本，专为旋转磁盘设计。因此 RocksDB 的设计并不是针对 SSD。

3 SSI：可串行化快照隔离。有关更多详细信息，请参阅博客文章Badger 中的并发 ACID 事务

4 Badger 通过其 Iterator API 提供对值版本的直接访问。用户还可以通过选项指定每个密钥保留多少个版本。

## 比较

| 特征        | Badger         | RocksDB | BoltDB |
|:----------|:---------------|:--------|:-------|
| 设计        | 具有值日志的LSM树     | 仅LSM树   | B+树    |
| 高读取吞吐量    | 是的             | 不       | 是的     |
| 高写入吞吐量    | 是的             | 是的      | 不      |
| 专为SSD设计   | 是（最新研究1）       | 不具体2    | 不      |
| 可嵌入       | 是的             | 是的      | 是的     |
| 排序KV访问    | 是的             | 是的      | 是的     |
| 纯Go（无Cgo） | 是的             | 不       | 是的     |
| 交易        | 是，ACID，与SSI3并发 | 是（但非酸性） | 是的，酸性  |
| 快照        | 是的             | 是的      | 是的     |
| TTL支持     | 是的             | 是的      | 不      |
| 3D        | 访问（键值版本）是4     | 不       | 不      |
| 压缩        | 是的             | 是的      | 不      |
| 原子性       | 是的             | 是的      | 不      |
| 事务        | 是的             | 是的      | 是的     |
| 并发        | 是的             | 是的      | 是的     |
| 无锁        | 是的             | 不       | 不      |
| 无GC        | 是的             | 不       | 不      |
| 无内存分配    | 是的             | 不       | 不      |
| 无Cgo      | 是的             | 不       | 是的     |
| 无依赖       | 是的             | 不       | 不      |
