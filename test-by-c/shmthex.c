#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>
#include <memory.h>
#include <time.h>
#include <errno.h>
#include <pthread.h>

#define MAX_CONTENT_LEN 10240
#define MAX_SHARE_MEM_SIZE 1024*1024*50
#define TOPIC_LEN 64

#define MAX_THREAD_NUM 2

//make
//gcc -g shmthex.c -lpthread -o testshm


pthread_t tidps[MAX_THREAD_NUM];
unsigned long long writeOffSet = 0;


static int shmid = 0;

typedef struct tagHead
{
	unsigned long long readOffSet;
	unsigned long long WriteOffSet;
} Head, *PHead;

Head h;
typedef struct tagTLV
{
	unsigned long long Tag;
	unsigned long long len;
	char Topic[TOPIC_LEN];
	char Value[MAX_CONTENT_LEN];

} TLV, *PTLV;


const unsigned int blockHeadSize = 7;
const unsigned int  maxBlockCount = 100;
const unsigned int  blockDataSize = 16;
const unsigned int blockSize = blockDataSize + blockHeadSize;

typedef struct Block
{
    /* data */
     unsigned int  completed;
     unsigned int  blockCount;
     unsigned int  blockIndex;
     unsigned int  blockLen;//this is len of  data ,
    PTLV data;
}TBlock,*PBlock;

typedef struct threadParam
{
	key_t key;
	char topic[64];
	char jsonPath[256];
	char *json;
	unsigned long long Tag;
	unsigned long long writeOffSet;
	int jsonlen;
	void *shmadd;
	int id;
} THParam;

static void *create(void *arg)
{
	THParam *param = (THParam *)arg;

	PTLV td = (PTLV)malloc(sizeof(TLV));
	memset(td, 0, sizeof(TLV));
	int lun = 0;

	param->writeOffSet = sizeof(Head);
	while (1)
	{
		lun++;
		clock_t time = clock();
		param->Tag++;
		size_t count = 500000;
		size_t i;
		for (i = 0; i < count; i++)
		{

			if (i % 10000 == 0)
			{
				printf("topic:%s index:%ld,offset:%lld Tag:%lld\n", param->topic, i, param->writeOffSet, param->Tag);
			}

			// sprintf(temp,"topic %lld ",Tag);
			int topiclen = strlen(param->topic);
			td->Tag = param->Tag;
			td->len = param->jsonlen;
			memcpy(td->Value, param->json, param->jsonlen);
			memcpy(td->Topic, param->topic, topiclen);

			int len = param->jsonlen;
			len += TOPIC_LEN;
			len += sizeof(unsigned long long) * 2;

			if (param->writeOffSet + len < MAX_SHARE_MEM_SIZE)
			{
				memcpy(param->shmadd + param->writeOffSet, td, len);
				param->writeOffSet += (MAX_CONTENT_LEN + TOPIC_LEN + sizeof(unsigned long long) * 2);
				if(param->writeOffSet >= MAX_SHARE_MEM_SIZE){
				    param->writeOffSet = 16;
				}
			}
			else
			{
				printf("to head:%lld,%s\n", param->writeOffSet, param->topic);
				memset(param->shmadd + param->writeOffSet, 0x00, MAX_SHARE_MEM_SIZE - param->writeOffSet - 1);
				param->writeOffSet = sizeof(Head);

				memcpy(param->shmadd + param->writeOffSet, td, len);
				param->writeOffSet += (MAX_CONTENT_LEN + TOPIC_LEN + sizeof(unsigned long long) * 2);
			}
			param->Tag++;
		}
		clock_t second_time = clock();
		printf("index:%d,ms:%lf\n", lun, (double)((second_time - time) / 1000));
		printf("index:%d,s:%lf\n", lun, (double)((second_time - time) / CLOCKS_PER_SEC));

		sleep(1);
	}
	free(param->json);
	free(td);

	// 共享内存区段的脱离
	int ret = shmdt(param->shmadd);
	if (ret == 0)
		printf("Successfully detach memory.\n");
	else
		printf("Memory detached failed %d\n", errno);

	// 删除该共享内存区
	ret = shmctl(param->id, IPC_RMID, 0);
	if (ret == 0)
		printf("Share memory removed.\n");
	else
		printf("Share memory remove failed.\n");

	return 0;
}
static void *createEx(void *arg)
{
	THParam *param = (THParam *)arg;

    PBlock bl = (PBlock)malloc(sizeof(TBlock));
    memset(bl, 0, sizeof(TBlock));

    PTLV td = (PTLV)malloc(sizeof(TLV));
    memset(td, 0, sizeof(TLV));

    bl->data = (PTLV)malloc(sizeof(TLV));
    memset(bl->data, 0, sizeof(TLV));

	int lun = 0;
	param->writeOffSet = sizeof(Head);
	while (1)
	{
		lun++;
		clock_t time = clock();
		param->Tag++;
		size_t count = 500000;
		size_t i;
		for (i = 0; i < count; i++)
		{

			if (i % 10000 == 0)
			{
				printf("topic:%s index:%ld,offset:%lld Tag:%lld\n", param->topic, i, param->writeOffSet, param->Tag);
			}

			// sprintf(temp,"topic %lld ",Tag);
			int topiclen = strlen(param->topic);
			td->Tag = param->Tag;
			td->len = param->jsonlen;
			memcpy(td->Value, param->json, param->jsonlen);
			memcpy(td->Topic, param->topic, topiclen);

			int len = param->jsonlen;
			len += TOPIC_LEN;
			len += sizeof(unsigned long long) * 2;
			bl->blockLen=len;
            memcpy(bl->data,td,len);
            int bloblen = bl->blockLen+sizeof(unsigned int)*4;
			if (param->writeOffSet + bloblen < MAX_SHARE_MEM_SIZE)
			{
				memcpy(param->shmadd + param->writeOffSet, bl, bloblen);
				param->writeOffSet += bloblen;
			}
			else
			{
				printf("to head:%lld,%s\n", param->writeOffSet, param->topic);
				memset(param->shmadd + param->writeOffSet, 0x00, MAX_SHARE_MEM_SIZE - param->writeOffSet - 1);
				param->writeOffSet = sizeof(Head);

				memcpy(param->shmadd + param->writeOffSet, bl, bloblen);
				param->writeOffSet += bloblen;
			}
			param->Tag++;
		}
		clock_t second_time = clock();
		printf("index:%d,ms:%lf\n", lun, (double)((second_time - time) / 1000));
		printf("index:%d,s:%lf\n", lun, (double)((second_time - time) / CLOCKS_PER_SEC));

		sleep(1);
	}
	free(param->json);
	free(td);

	// 共享内存区段的脱离
	int ret = shmdt(param->shmadd);
	if (ret == 0)
		printf("Successfully detach memory.\n");
	else
		printf("Memory detached failed %d\n", errno);

	// 删除该共享内存区
	ret = shmctl(param->id, IPC_RMID, 0);
	if (ret == 0)
		printf("Share memory removed.\n");
	else
		printf("Share memory remove failed.\n");

	return 0;
}

void *createSHM(key_t key, int *id)
{
	writeOffSet = sizeof(Head);
	h.WriteOffSet = writeOffSet;
	int ret;

	char *shmadd;

	/*
    // ftok(".", 2012);
    if (key == -1)
    {
        perror("ftok");
    }*/
	/*创建共享内存*/
	shmid = shmget(key, MAX_SHARE_MEM_SIZE, IPC_CREAT | 666);
	if (shmid < 0)
	{
		perror("shmget");
		exit(-1);
	}
	*id = shmid;
	/*映射*/
	shmadd = shmat(shmid, NULL, 0);

	struct shmid_ds shmds;
	ret = shmctl(shmid, IPC_STAT, &shmds);
	if (ret == 0)
	{
		printf("Size of memory segment is %d bytes.\n", (int)shmds.shm_segsz);
		printf("Number of attach %d\n", (int)shmds.shm_nattch);
	}
	else
	{
		printf("shmctl() call failed.\n");
	}
	bzero(shmadd, MAX_SHARE_MEM_SIZE);
	memcpy(shmadd, &h, sizeof(h));
	return shmadd;
}

char *getJsonData(const char *jsonFile, int *jsonLen)
{
	FILE *f = fopen(jsonFile, "rb");
	char *json;
	if (f)
	{
		fseek(f, 0, SEEK_END);
		fpos_t pos;

		fgetpos(f, &pos);
		int len = (int)pos.__pos;
		fseek(f, 0, SEEK_SET);
		json = (char *)malloc(len + 1);
		fread(json, sizeof(char), len, f);
		json[len] = '\0';
		*jsonLen = len;
		fclose(f);
		printf("%s\n", json);
		sleep(2);
	}
	return json;
}
void NormalShareMemory();


typedef struct SHMInfo
{
    unsigned long long max_topic_len;
    unsigned long long max_content_len;
    unsigned long long max_shm_size;
    unsigned long long count;
    key_t key[200];
} SHMI, *PSHMI;

void* createSHMDefault(key_t key, int *id)
{
    int ret;

    char *shmadd;

    /*创建共享内存*/
    int size = sizeof(SHMI);
    int shmid = shmget(key, size, IPC_CREAT | 666);
    if (shmid < 0)
    {
        perror("shmget");
        exit(-1);
    }

    struct shmid_ds shmds;
    ret = shmctl(shmid, IPC_STAT, &shmds);
    if (ret == 0)
    {
        printf("Size of memory segment is %d bytes.\n", (int)shmds.shm_segsz);
        printf("Number of attach %d\n", (int)shmds.shm_nattch);
    }
    else
    {
        printf("shmctl() call failed.\n");
        return NULL;
    }
    *id = shmid;
    /*映射*/
    shmadd = shmat(shmid, NULL, 0);
    bzero(shmadd, size);

    return shmadd;
}
#define TEST_SHMI
int main(/*int argc, char *argv[]*/)
{
#ifdef TEST_SHMI
    int id = 0;
   void* shmadd= createSHMDefault(999999,&id);
    printf("shmid:%d\n",id);

    printf("sizeof:shmi:%d\n",sizeof(SHMI));
    SHMI shmi;
    shmi.max_topic_len =TOPIC_LEN;
    shmi.max_content_len=MAX_CONTENT_LEN;
    shmi.max_shm_size=MAX_SHARE_MEM_SIZE;
    shmi.count=MAX_THREAD_NUM;
#endif//


	printf("share memory\n");
	for (int j = 0; j < MAX_THREAD_NUM; j++)
	{
		THParam *tp = (THParam *)malloc(sizeof(THParam));
		tp->key = 202107 + j;
#ifdef TEST_SHMI
        shmi.key[j]=202107+j;
#endif//
		sprintf(tp->topic, "kill_kafa_%d", j + 1);
		sprintf(tp->jsonPath, "./example.json", j + 1);
		tp->writeOffSet = sizeof(Head);

		int len = 0;
		tp->json = getJsonData(tp->jsonPath, &len);
		tp->jsonlen = len;

		int id = 0;
		tp->shmadd = createSHM(tp->key, &id);
		printf("id = %d, len:%d", id, len);
		tp->id = id;
		tp->Tag = 0;
		pthread_t tidp;
		int error;
		void *tret;
		error = pthread_create(&tidps[j], NULL, create, (void *)tp);
	}
#ifdef TEST_SHMI
    memcpy(shmadd,&shmi,sizeof(SHMI));
#endif//
	for (int k = 0; k < MAX_THREAD_NUM; k++)
	{
		void *tret;
		pthread_join(tidps[k], &tret);
		printf("pid:%ld\n", tidps[k]);
	}

	while (1)
	{
	    printf("running ....\n");
	}
	return 1;
}

void NormalShareMemory()
{
	char *json;
	int jsonlen;
	int ret = 0;
	json = getJsonData("example.json", &jsonlen);

	key_t key = 202106;
	int id = 0;
	void *shmadd = createSHM(key, &id);
	if (shmadd < 0)
	{
		perror("shmat");
		_exit(-1);
	}

	printf("copy data to shared-memory\n");

	unsigned long long Tag = 0;
	int lun = 0;
	unsigned long long writeOffSet = sizeof(Head);
	PTLV td = (PTLV)malloc(sizeof(TLV));
	memset(td, 0, sizeof(TLV));

	while (1)
	{
		lun++;
		clock_t time = clock();
		Tag++;
		size_t count = 500000;
		size_t i;
		for (i = 0; i < count; i++)
		{
			if (i % 10000 == 0)
			{
				printf("index:%ld,offset:%lld Tag:%lld\n", i, writeOffSet, Tag);
			}
			char *temp = "kill_kafa";
			// sprintf(temp,"topic %lld ",Tag);
			int topiclen = strlen(temp);
			td->Tag = Tag;
			td->len = jsonlen;
			memcpy(td->Value, json, jsonlen);
			memcpy(td->Topic, temp, topiclen);

			int len = jsonlen;
			len += TOPIC_LEN;
			len += sizeof(unsigned long long) * 2;

			if (writeOffSet + len < MAX_SHARE_MEM_SIZE)
			{
				memcpy(shmadd + writeOffSet, td, len);
				writeOffSet += len;
			}
			else
			{
				printf("to head:%lld", writeOffSet);
				memset(shmadd + writeOffSet, 0x00, MAX_SHARE_MEM_SIZE - writeOffSet - 1);
				writeOffSet = sizeof(Head);
				memcpy(shmadd + writeOffSet, td, len);
				writeOffSet += len;
			}
			Tag++;
		}
		clock_t second_time = clock();
		printf("index:%d,ms:%lf\n", lun, (double)((second_time - time) / 1000));
		printf("index:%d,s:%lf\n", lun, (double)((second_time - time) / CLOCKS_PER_SEC));

		sleep(1);
	}
	free(json);
	free(td);

	// 共享内存区段的脱离
	ret = shmdt(shmadd);
	if (ret == 0)
		printf("Successfully detach memory.\n");
	else
		printf("Memory detached failed %d\n", errno);

	// 删除该共享内存区
	ret = shmctl(shmid, IPC_RMID, 0);
	if (ret == 0)
		printf("Share memory removed.\n");
	else
		printf("Share memory remove failed.\n");
}
