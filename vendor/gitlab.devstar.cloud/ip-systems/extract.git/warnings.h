// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#ifndef WR_WARN
#define WR_WARN

#include <unistd.h>
#include <stdio.h>

typedef struct warn_struct {
    int line;
    char* file;
    char* message;
} warning;

typedef struct warn_array {
    size_t nwarn;
    warning** warnings;
} warning_array;

warning_array* warning_array_init();
void warn(warning_array* wa, int line, const char* file, const char* message);
void warning_array_free(warning_array* wa);
void warning_fprint(FILE* f, warning* w);

#endif