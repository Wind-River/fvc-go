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
#include <stdio.h>
#include <unistd.h>
#include <string.h>

char* bytesToHex(const unsigned char* bytes, size_t len) {
    char* ret = malloc(len*2+1);
    for(int i = 0; i < len; i++) {
        int j = i*2;
        sprintf(ret+j, "%02x", bytes[i]);
    }
    ret[len*2] = '\0';

    return ret;
}

char* charToHex(const char* s, size_t len) {
    if(len == 0) {
        len = strlen(s);
    }

    char* ret = malloc(len*2+1);
    for(int i = 0; i < len; i++) {
        int j = i*2;
        sprintf(ret+j, "%02x", s[i]);
        printf("%c -> %02x\n", s[i], s[i]);
    }
    ret[len*2] = '\0';

    return ret;
}