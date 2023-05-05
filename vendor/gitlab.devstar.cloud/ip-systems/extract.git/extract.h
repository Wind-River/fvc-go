// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include "status.h"

//extract requires pwd to be the destination directory
//if filename is null, it is expected _extract will extract the archive without problem
//if filename is not null, it is expected that the archive is not a tar, so _decompress should be tried if _extract fails
status extract(char *filepath, char *filename, char *dest);