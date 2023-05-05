// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include <stdbool.h>

typedef struct filename_struct {
    char* name;
    char* basename;
    const char* ext;
    bool tar;
}* filename_ptr;

const char* getBasename( filename_ptr fp);
const char* getExtension( filename_ptr fp);
bool compressedBinary( filename_ptr fp );
filename_ptr parseFilename( char *str );
void filename_free( filename_ptr fp );
