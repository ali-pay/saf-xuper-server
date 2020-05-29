#!/bin/bash

set -e

cd `dirname $0`
#xdev所在目录
export PATH=`pwd`:$PATH
#当前工作目录
export XDEV_ROOT=`pwd`

# install docker in precondition
if ! command -v docker &>/dev/null; then
    echo "missing docker command, please install docker first."
    exit 1
fi

# check if xdev available
if ! command -v xdev &>/dev/null; then
    echo "missing xdev command, please cd xuperchain && make"
    exit 1
fi

mkdir -p ../contract_wasm
cc=../contract_code/$1

# build single cc file
if [[ -f $cc ]]; then
    out=../contract_wasm/$(basename $1 .cc).wasm
    echo "build $cc"
    xdev build -o $out $cc
fi

# build package
if [[ -d $cc ]]; then
    echo "build $cc"
    bash -c "cd $cc && xdev build && mv -v $1.wasm ../../contract_wasm/"
fi
