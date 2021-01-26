# bloom
bloom filter - 布隆过滤器，使用的时候需要转化成`[]byte`格式进行哈希


This implementation accepts keys for setting and testing as `[]byte`. Thus, to
add a string item, `"Love"`:

    n := uint(1000)
    filter := bloom.New(20*n, 5) // load of 20, 5 keys
    filter.Add([]byte("Love"))

Similarly, to test if `"Love"` is in bloom:

    if filter.Test([]byte("Love"))

For numeric data, I recommend that you look into the encoding/binary library. But, for example, to add a `uint32` to the filter:

    i := uint32(100)
    n1 := make([]byte, 4)
    binary.BigEndian.PutUint32(n1, i)
    filter.Add(n1)

Finally, there is a method to estimate the false positive rate of a particular
bloom filter for a set of size _n_:

    if filter.EstimateFalsePositiveRate(1000) > 0.001

Given the particular hashing scheme, it's best to be empirical about this. Note
that estimating the FP rate will clear the Bloom filter.

## Installation

```bash
go get -u github.com/ducksoso/bloom
```
