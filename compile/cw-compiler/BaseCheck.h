//
// Created by kami on 2020/7/22.
//

#ifndef CW_COMPILER_BASECHECK_H
#define CW_COMPILER_BASECHECK_H

#include <cstdlib>

#define MAX_CHECKING_ITEM   64
#define MAX_CONFIG_ITEM_TITLE_LEN 12

//prefix of symbol
#define PREFIX_EXPORT   "  (export \""

#define PREFIX_IMPORT   "  (import \""

class BaseCheck {
private:
    int     export_idx;
    char*   m_Export[MAX_CHECKING_ITEM];
    bool    m_MapOfExport[MAX_CHECKING_ITEM];
    int     import_idx;
    char*   m_Import[MAX_CHECKING_ITEM];
    bool    m_MapOfImport[MAX_CHECKING_ITEM];
    int     feature_idx;
    char*   m_feature[MAX_CHECKING_ITEM];
    bool    m_MapOfFeature[MAX_CHECKING_ITEM];
private:
    bool    checking_export(char* strline);
    bool    checking_feature(char* strline);
    bool    checking_import(char* strline); //there is no way to resolve import error，so import error will shut compile down

public:
    BaseCheck(){
        for (int i = 0; i < MAX_CHECKING_ITEM; ++i) {
            m_Export[i] = m_Import[i] = m_feature[i] = NULL;
            m_MapOfExport[i] = m_MapOfFeature[i] = m_MapOfImport[i] = false;
        }
        export_idx = 0;
        import_idx = 0;
        feature_idx = 0;
    }
    ~BaseCheck(){
        for (int i = 0; i < MAX_CHECKING_ITEM; ++i) {
            if (m_Export[i] != NULL){
                free(m_Export[i]);
                m_Export[i] = NULL;
            }
            if (m_Import[i] != NULL){
                free(m_Import[i]);
                m_Import[i] = NULL;
            }
            if (m_feature[i] != NULL){
                free(m_feature[i]);
                m_feature[i] = NULL;
            }
        }
    }

    void ResetMap(){
        for (int i = 0; i < MAX_CHECKING_ITEM; ++i) {
            m_MapOfExport[i] = m_MapOfFeature[i] = m_MapOfImport[i] = false;
        }
    }

    //loading config file from target path
    int LoadConfig(char* path);

    //if isResolve is true, we will try to resolve this problem
    int Checking(bool isResolve, char *path);

    off_t file_size(char *filename);

    bool lowerCheck(char* strCheck, char* prefix, char** BaseChecking, int sizeofBaseChecking,bool* Mapping);

    bool Verify();
};


#endif //CW_COMPILER_BASECHECK_H
