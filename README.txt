Peg your CPUs! ...or your disks! Or maybe something else?!

Pegger is a suite of tools for exploring performance and behavior of
Go programs in extreme scenarios. For example, you might want to
understand what Go's profiling tools look like when your program is
under pure CPU load, or doing lots of small I/Os to disk. You might
want to see how your cloud hardware, or a new kernel, or whatever
responds under these conditions. Pegger is all about giving you a
baseline to work with so that you can understand your tools better,
and therefore understand your own software better.

Install:

`make install`

Usage:

Try `pegger -h` or `disker -h`.

