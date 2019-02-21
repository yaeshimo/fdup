fdup
====

Find duplicate files have contents

Usage
-----

```sh
# Prints help information
fdup -help
```

```sh
# Compare two files
fdup /path/file1 /path/file2
```

```sh
# Check recursive
fdup /path/dir
```

```sh
# Check files and directories with verbose
fdup -verbose /path/file1 /path/file2 /path/dir1/ /path/dir2/
```

Install
-------

```sh
go get github.com/kamisari/fdup
```

License
-------

MIT
