## 思路
- [ ] 交易钱包存放在一个结构WalletCollect数组中
- [ ] 查询余额直接返回，暂时的数据延迟可以忽略
- [ ] 核心交易场景下，A->B交易时需要 WalletCollect中的对应位置A和B均为被其他的线程使用，Wallet中state表示正在交易
- [ ] 分别为A和B申请一把临时的锁,用来保护 WalletCollect buf结构中A和B的Balance修改顺序正确
- [ ] 由于还在存在WalletCollect 新增时数组存在自动扩容的机制，为了避免这个并发维度的扩容导致数据不安全，内置一个大容量的长度避免自动扩容