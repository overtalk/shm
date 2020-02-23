#! /bin/bash

for i in `ipcs -m | tail -n +4 | awk {'print $2'}` # 共享内存
do
	ipcrm -m $i;
	echo 删除shm = $i
done