#!/usr/bin/env bash
git clone https://github.com/golang/blog.git $GOPATH/src/github.com/golang/blog
git clone https://github.com/golang/crypto.git $GOPATH/src/github.com/golang/crypto
git clone https://github.com/golang/exp.git $GOPATH/src/github.com/golang/exp
git clone https://github.com/golang/image.git $GOPATH/src/github.com/golang/image
git clone https://github.com/golang/mobile.git $GOPATH/src/github.com/golang/mobile
git clone https://github.com/golang/net.git $GOPATH/src/github.com/golang/net
git clone https://github.com/golang/sys.git $GOPATH/src/github.com/golang/sys
git clone https://github.com/golang/talks.git $GOPATH/src/github.com/golang/talks
git clone https://github.com/golang/text.git $GOPATH/src/github.com/golang/text
git clone https://github.com/golang/tools.git $GOPATH/src/github.com/golang/tools
ln -s $GOPATH/src/github.com/golang $GOPATH/src/golang.org/x