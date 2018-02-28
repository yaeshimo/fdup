go-fdup
=======
find duplicate contents

Usage:
------
display files with have same contents

compare two files
```sh
fdup -- "/path/file1" "/path/file2"
```

check recursive
```sh
fdup -- "/path/dir"
```

specify file and directory
```sh
fdup -- "/path/file" "/path/dir"
```

with verbose
```sh
fdup -verbose -- "/path/file" "/path/dir"
```

Install:
--------
```sh
go get -v -u github.com/kamisari/go-fdup/fdup
```

TODO:
-----
impl

License:
--------
MIT
