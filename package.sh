#! /bin/bash

meepo_version=${GIT_TAG_NAME}
echo "build version: ${meepo_version}"

GO111MODULE=on go mod vendor

make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux darwin'
arch_all='386 amd64 arm arm64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        meepo_bin="meepo_${os}_${arch}"
        meepo_dir_name="meepo_${meepo_version}_${os}_${arch}"
        meepo_path="packages/${meepo_dir_name}"

        if [ ! -f ${meepo_bin} ]; then
            continue
        fi

        mkdir ${meepo_path}
        mv ${meepo_bin} ${meepo_path}/meepo
        cp ../LICENSE ${meepo_path}

        cd packages
        tar -zcf ${meepo_dir_name}.tar.gz ${meepo_dir_name}
        sha1sum ${meepo_dir_name}.tar.gz > ${meepo_dir_name}.tar.gz.sha1.txt
        cd ..
        rm -rf ${meepo_path}
    done
done

cd ..
