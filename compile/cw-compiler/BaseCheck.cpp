//
// Created by kami on 2020/7/22.
//

#include <cstdio>
#include <cstring>
#include <sys/stat.h>
#include "BaseCheck.h"
#include "error.h"


off_t BaseCheck::file_size(char* filename)
{
    struct stat statbuf;
    stat(filename,&statbuf);
    off_t size=statbuf.st_size;
    return size;
}

int BaseCheck::LoadConfig(char *path) {
    char    szTest[4096] = {0};
    int     readed = 0;
    int     fslen = file_size(path);
    char*   writer = NULL;
    if (fslen <= 0){
        printf("CW-COMPILER(Error): Failed to get file size %s\n",path);
        return ERROR_CODE_FAILED_GET_FILE_SIZE;
    }
    FILE *fp = fopen(path, "r");
    if (NULL == fp) {
        printf("CW-COMPILER(Error): Failed to open %s\n",path);
        return ERROR_CODE_FAILED_OPEN_FILE;
    }
    while (!feof(fp)) {
        memset(szTest, 0, sizeof(szTest));
        fgets(szTest, sizeof(szTest) - 1, fp);

        if (strncmp(szTest,"@export_set:",MAX_CONFIG_ITEM_TITLE_LEN) == 0){
            //set export check list
            m_Export[export_idx] = (char*)malloc(4096);
            memset(m_Export[export_idx],0,4096);
            writer = m_Export[export_idx];
            export_idx++;
        }else if (strncmp(szTest,"@import_set:",MAX_CONFIG_ITEM_TITLE_LEN) == 0){
            //set import check list
            m_Import[import_idx] = (char*)malloc(4096);
            memset(m_Import[import_idx],0,4096);
            writer = m_Import[import_idx];
            import_idx++;
        }else if (strncmp(szTest,"@feature_st:",MAX_CONFIG_ITEM_TITLE_LEN) == 0){
            //set feature check list
            m_feature[feature_idx] = (char*)malloc(4096);
            memset(m_feature[feature_idx],0,4096);
            writer = m_feature[feature_idx];
            feature_idx++;
        }else{
            printf("Can not analyze target string: (%s)\n",szTest);
            writer = NULL;
        }
        if (writer){
            strcpy(writer,&szTest[MAX_CONFIG_ITEM_TITLE_LEN]);
            writer[strlen(writer) -1] = '\0';
        }
        readed += strlen(szTest);
        if(readed >= fslen){    //not very accurately, but that's ok
            break;
        }
    }
    fclose(fp);
    return 0;
}

int BaseCheck::Checking(bool isResolve,char* path) {

    //executing all check as `only` include policy
    char    szTest[4096] = {0};
    int     readed = 0;
    int     fslen = file_size(path);
    char    newPath[1024] = {0};

    ResetMap();

    FILE *fp = fopen(path, "r");
    if (NULL == fp) {
        printf("CW-COMPILER(Error): Failed to open %s\n",path);
        return ERROR_CODE_FAILED_OPEN_FILE;
    }
    strcpy(newPath,path);
    strcat(newPath,".n.wbt");
    FILE *fp_write = fopen(newPath, "w");
    if (NULL == fp) {
        printf("CW-COMPILER(Error): Failed to open %s\n",path);
        return ERROR_CODE_FAILED_OPEN_FILE;
    }
    while (!feof(fp)){
        memset(szTest, 0, sizeof(szTest));
        fgets(szTest, sizeof(szTest) - 1, fp);
        if (!checking_export(szTest)){
            continue;
        }else if (!checking_import(szTest)){
            printf("CW-COMPILER(Error): Unrecognized import func:(%s)",szTest);
            return ERROR_CODE_UNRECOGNIZE_IMPORT;
        }else if (!checking_feature(szTest)){
            printf("CW-COMPILER(Error): Unrecognized feature symbol:(%s)",szTest);
            return ERROR_CODE_UNRECOGNIZE_FEATURE;
        } else{
            fwrite(szTest,strlen(szTest),1,fp_write);
        }
    }

    fclose(fp);
    fclose(fp_write);

    if (!Verify()){
        remove(newPath);
    }
    return 0;
}

bool BaseCheck::lowerCheck(char* strCheck, char* prefix, char** BaseChecking, int sizeofBaseChecking,bool* Mapping){
    if (!strCheck || !prefix || !BaseChecking || sizeofBaseChecking <= 0){
        return true;//whatever, just pass it through
    }
    int len = strlen(prefix);
    if (!strCheck || strlen(strCheck) <= len){
        return true;    //passthrough
    }
    if (strncmp(prefix,strCheck,len) != 0){
        return true;    //not export section, pass it through
    }
    for (int i = 0; i < sizeofBaseChecking; ++i) {
        if (strncmp(BaseChecking[i],&strCheck[len],strlen(BaseChecking[i])) == 0){
            Mapping[i] = true;
            return true;
        }
    }
    //Aha, got you, illegal symbol, need resolve
    return false;
}

bool BaseCheck::checking_export(char* strline) {
    return lowerCheck(strline,(char*)PREFIX_EXPORT,m_Export,export_idx,m_MapOfExport);
}

bool BaseCheck::Verify(){
    bool ret = true;
    for (int i = 0; i < export_idx; ++i) {
        if (!m_MapOfExport[i]){
            //missing some export func
            printf("CW-COMPILER(Error): Missing export func (%s)\n",m_Export[i]);
            ret = false;
        }
    }
    return ret;
}

bool BaseCheck::checking_feature(char* strline) {
    return true;
}

bool BaseCheck::checking_import(char* strline) {
    return lowerCheck(strline,(char*)PREFIX_IMPORT,m_Import,import_idx,m_MapOfImport);
}
