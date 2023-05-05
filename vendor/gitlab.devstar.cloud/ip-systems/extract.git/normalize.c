// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include <stdlib.h>
#include <string.h>
#include <stdio.h>

char* normalize(const char* s) {
    char* ret = malloc(strlen(s)+1);
    int rIndex = 0;
    int sIndex;
    for(sIndex = 0; sIndex < strlen(s); sIndex++) {
        unsigned int cur = s[sIndex];
        if(cur == 0xc3) {
            unsigned int next = s[sIndex+1];
            if(next < 0x87) {
                ret[rIndex++] = 'A';
            } else if(next == 0x87) {
                ret[rIndex++] = 'C';
            } else if (next < 0x8c) {
                ret[rIndex++] = 'E';
            } else if (next < 0x90) {
                ret[rIndex++] = 'I';
            } else if (next == 0x90) {
                ret[rIndex++] = 'D';
            } else if (next == 0x91) {
                ret[rIndex++] = 'N';
            } else if (next < 0x97) {
                ret[rIndex++] = 'O';
            } else if (next == 0x97) {
                ret[rIndex++] = 'x';
            } else if (next == 0x98) {
                ret[rIndex++] = '0';
            } else if (next < 0x9d) {
                ret[rIndex++] = 'U';
            } else if (next < 0x9f) {
                ret[rIndex++] = 'Y';
            } else if (next == 0x9f) {
                ret[rIndex++] = 'S';
            } else if (next < 0xa7) {
                ret[rIndex++] = 'a';
            } else if (next == 0xa7) {
                ret[rIndex++] = 'c';
            } else if (next < 0xac) {
                ret[rIndex++] = 'e';
            } else if (next < 0xb0) {
                ret[rIndex++] = 'i';
            } else if (next == 0xb0) {
                ret[rIndex++] = 'd';
            } else if (next == 0xb1) {
                ret[rIndex++] = 'n';
            } else if (next < 0xb7) {
                ret[rIndex++] = 'o';
            } else if (next == 0xb7) {
                ret[rIndex++] = '%';
            } else if (next == 0xb8) {
                ret[rIndex++] = '0';
            } else if (next < 0xbd) {
                ret[rIndex++] = 'u';
            } else if (next <= 0xbf) {
                ret[rIndex++] = 'y';
            } else {
                printf("normalize unknown code point: %d\n", next);
            }
            sIndex += 2;
        } else {
            printf("%c -> %x\n", s[sIndex], (unsigned int)s[sIndex]);
            ret[rIndex++] = s[sIndex];
        }
    }
    ret[rIndex++] = '\0';

    return ret;
}