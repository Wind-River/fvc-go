// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include <string.h>
#include <stdlib.h>

char* join(const char *a, const char *b) {
    int aLen = strlen(a);
    int requiredLength = aLen;
    if (aLen > 0 && a[aLen-1] == '/') {
        requiredLength--;
        aLen--;
    }

    int bStart = 0;
    int bLen = 0;
    if (b != NULL) {
        bLen = strlen(b);
        requiredLength += bLen;

        if (bLen > 0 && b[0] == '/') {
            bStart = 1;
            bLen--;
            requiredLength--;
        }
    }

    char* joinedPath = malloc(sizeof(char*) * requiredLength + 1);
    joinedPath[0] = '\0';
    strncat(joinedPath, a, aLen);
    if (bLen > 0) {
        strncat(joinedPath, "/", 1);
        strncat(joinedPath, &b[bStart], bLen);
    }

    return joinedPath;
}