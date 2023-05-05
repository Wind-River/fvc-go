// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include "logger.h"
#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

int logger_verbose = 0;

int vlog(const char* format, ...) {
    if(logger_verbose > 0) {
        va_list arg;
        int done;

        va_start(arg, format);
        done = vfprintf(stdout, format, arg);
        va_end(arg);

        return done;
    }
    
    return 0;
}

int flog(FILE* stream, const char* format, ...) {
    va_list arg;
    int done;

    va_start(arg, format);
    done = vfprintf(stream, format, arg);
    va_end(arg);

    return done;
}

int elog(const char* format, ...) {
    va_list arg;
    int done;

    va_start(arg, format);
    done = vfprintf(stderr, format, arg);
    va_end(arg);

    return done;
}

char* jsonEscape(const char* s) {
    char* ret = malloc((strlen(s)*2)+1);

    int j = 0;
    for(int i = 0; i < strlen(s); i++) {
        char c = s[i];
        if(c == '"' || c == '\\' || ('\x00' <= c && c <= '\x1f' )) {
            ret[j++] = '\\';
        }
        ret[j++] = c;
    }
    ret[j] = '\0';

    return ret;
}