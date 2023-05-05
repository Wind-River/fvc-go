// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include "warnings.h"

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>

warning_array* warning_array_init() {
    warning_array* ret = malloc(sizeof(warning_array));
    ret->nwarn = 0;
    ret->warnings = NULL;

    return ret;
}

void add_warning(warning_array* wa, warning *w) {
    if(wa->warnings == NULL) {
        // init array
        wa->warnings = malloc(sizeof(warning*)*1);
        wa->nwarn = 1;
        wa->warnings[0] = w;
    } else {
        // append to array
        wa->warnings = realloc(wa->warnings, sizeof(warning*)*(wa->nwarn+1));
        wa->warnings[wa->nwarn] = w;
        wa->nwarn += 1;
    }
}

void warn(warning_array* wa, int line, const char* file, const char* message) {
    warning* w = malloc(sizeof(warning));
    w->line = line;

    if(file != NULL) {
        w->file = strdup(file);
    } else {
        w->file = NULL;
    }
    
    if(message != NULL) {
        w->message = strdup(message);
    } else {
        w->message = NULL;
    }

    add_warning(wa, w);
}

void warning_array_free(warning_array* wa) {
    if(wa==NULL) {
        printf("NULL\n");
        return;
    }
    for(int i = 0; i < wa->nwarn; i++) {
        warning* w = wa->warnings[i];
        if(w == NULL) continue;

        if(w->file != NULL) free(w->file);
        if(w->message != NULL) free(w->message);
        free(w);
    }
    free(wa->warnings);

    free(wa);
}

void warning_fprint(FILE* f, warning* w) {
    if(w != NULL) {
        fprintf(f, "[%d %s]: %s\n", w->line, w->file, w->message);
    }
}