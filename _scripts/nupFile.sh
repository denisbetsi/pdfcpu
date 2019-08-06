#!/bin/sh

# Copyright 2019 The pdfcpu Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# eg: ./nupFile.sh ~/pdf/1mb/a.pdf ~/pdf/out

if [ $# -ne 2 ]; then
    echo "usage: ./nupFile.sh inFile outDir"
    echo "nup all pages into new pages showing 4 original pages on each page"
    exit 1
fi

new=_nup

f=${1##*/}
f1=${f%.*}
out=$2

cp $1 $out/$f 

out1=$out/$f1$new.pdf
pdfcpu nup -verbose $out1 4 $out/$f &> $out/$f1.log
if [ $? -eq 1 ]; then
    echo "nup error: $1 -> $out1"
    exit $?
else
    echo "nup success: $1 -> $out1"
    pdfcpu validate -verbose -mode=relaxed $out1 >> $out/$f1.log 2>&1
    if [ $? -eq 1 ]; then
        echo "validation error: $out1"
        exit $?
    else
        echo "validation success: $out1"
    fi    
fi
	

	
