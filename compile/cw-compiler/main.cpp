
#include <stdio.h>
#include <string.h>

#include <fcntl.h>
#include <sys/stat.h>
#include <errno.h>
#include <stdlib.h>
#include <sys/timeb.h>
#include "BaseCheck.h"

#define bool int
#define false 0
#define true 1

off_t file_size(char* filename)
{
    struct stat statbuf;
    stat(filename,&statbuf);
    off_t size=statbuf.st_size;
    return size;
}

//execute wasm2wat to convert wasm to wat
//reading every single line to check export/import/feature tables
//delete illegal symbol and write to new wat file
//after then, execute wat2wasm to build a wasm file from new wat file
int callCheckCompile(int argc, char* argv[])
{
    BaseCheck   bc;
    bc.LoadConfig(argv[3]);
    bc.Checking(true,(char*)"/Users/oker/work/go/src/github.com/cosmwasm/cosmwasm-simulate/wasm_go_poc/0.wat");
    return 0;
}

typedef int (*pfnCaller) (int argc,char* argv[]);

pfnCaller GetCaller(char* argv[]){

    if(strcmp(argv[2],"-c") == 0){
        return NULL;
    }
    return NULL;
}

void printHelp(){
    printf("Usage：Compiling target Webassembly file from original wasm file to cosmwasm-vm checked wasm file");
}

int main(int argc,char* argv[]) {
    callCheckCompile(argc,argv);
    pfnCaller pf = GetCaller(argv);
    if(pf == NULL){
        printHelp();
        return 0;
    }
    pf(argc,argv);
}
