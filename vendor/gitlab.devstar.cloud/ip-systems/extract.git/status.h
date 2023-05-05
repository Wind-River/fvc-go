// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#ifndef WR_STATUS
#define WR_STATUS

#include "warnings.h"

typedef struct exit_struct {
	int code;
	char* message;
	char* tag;
	warning_array* warnings;
}* status;

status report_status(int code, const char* message, const char* tag, warning_array* warnings);
status success(warning_array* wa);
void status_free(status stat);

#endif