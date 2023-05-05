// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.

#line 1 "src/filename.rl"

/*
 * Parse a filepath to separate full filename, basename, extension, and whether or not it is a tar
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdbool.h>
#include <linux/limits.h>

#include "filename.h"

const char* getBasename(filename_ptr fp) {
	if(fp->basename == NULL) {
		return "";
	}

	return fp->basename;
}

const char* getExtension(filename_ptr fp) {
	if(fp->ext == NULL) {
		return "";
	}

	return fp->ext;
}

bool compressedBinary(filename_ptr fp) {
	if(fp->tar || fp->ext == NULL) return false;

	const char* const compressionExtensions[] = {".gz", ".bz2", ".xz", ".lzma"};

	for(int i = 0; i < 4; i++) {
		if(strcmp(fp->ext, compressionExtensions[i]) == 0)
			return true;
	}

	return false;
}

filename_ptr newFilename(const char* name, char* ext, bool tar) {
	// printf("newFilename(%s, %s, %d)\n", name, ext, tar);
	filename_ptr ret = malloc(sizeof(struct filename_struct));

	ret->tar = tar;
	ret->name = strdup(name);

	if(ext != NULL) {
		int baseIdx = ext-name;
		ret->basename = strndup(name, baseIdx);

		ret->ext = ret->name + (ext-name);
	} else {
		ret->basename = NULL;
		ret->ext = NULL;
	}

	// printf("newFilename(%s, %s, %d) -> filename_struct{%s, %s, %s, %d}\n",
	// 	name, ext, tar, ret->name, getBasename(ret), getExtension(ret), ret->tar
	// );

	return ret;
}

void filename_free(filename_ptr fp) {
	free(fp->name);
	if(fp->basename != NULL) free(fp->basename);
	free(fp);

	// printf("-> filename_free\n");
}


#line 79 "src/filename.c"
static const char _filename_actions[] = {
	0, 1, 0, 1, 1, 1, 2, 1, 
	3, 2, 0, 1
};

static const char _filename_key_offsets[] = {
	0, 2, 11, 20, 28, 37, 46, 54, 
	62, 71, 80
};

static const char _filename_trans_keys[] = {
	46, 47, 46, 47, 116, 48, 57, 65, 
	90, 97, 122, 46, 47, 116, 48, 57, 
	65, 90, 97, 122, 46, 47, 48, 57, 
	65, 90, 97, 122, 46, 47, 97, 48, 
	57, 65, 90, 98, 122, 46, 47, 114, 
	48, 57, 65, 90, 97, 122, 46, 47, 
	48, 57, 65, 90, 97, 122, 46, 47, 
	48, 57, 65, 90, 97, 122, 46, 47, 
	97, 48, 57, 65, 90, 98, 122, 46, 
	47, 114, 48, 57, 65, 90, 97, 122, 
	46, 47, 48, 57, 65, 90, 97, 122, 
	0
};

static const char _filename_single_lengths[] = {
	2, 3, 3, 2, 3, 3, 2, 2, 
	3, 3, 2
};

static const char _filename_range_lengths[] = {
	0, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3
};

static const char _filename_index_offsets[] = {
	0, 3, 10, 17, 23, 30, 37, 43, 
	49, 56, 63
};

static const char _filename_indicies[] = {
	1, 2, 0, 1, 2, 4, 3, 3, 
	3, 0, 1, 2, 6, 5, 5, 5, 
	0, 1, 2, 3, 3, 3, 0, 1, 
	2, 7, 3, 3, 3, 0, 1, 2, 
	8, 3, 3, 3, 0, 9, 2, 3, 
	3, 3, 0, 1, 2, 5, 5, 5, 
	0, 1, 2, 10, 5, 5, 5, 0, 
	1, 2, 11, 5, 5, 5, 0, 9, 
	2, 5, 5, 5, 0, 0
};

static const char _filename_trans_targs[] = {
	0, 1, 0, 3, 4, 7, 8, 5, 
	6, 2, 9, 10
};

static const char _filename_trans_actions[] = {
	0, 7, 5, 0, 0, 3, 3, 0, 
	1, 7, 3, 9
};

static const int filename_start = 0;
static const int filename_first_final = 3;
static const int filename_error = -1;

static const int filename_en_main = 0;


#line 78 "src/filename.rl"


filename_ptr parseFilename( char *str )
{
	char *p = str, *pe = str + strlen( str );
	int cs;

	char* slash = str;
	char* tar = NULL;
	char* ext = NULL;


	
#line 163 "src/filename.c"
	{
	cs = filename_start;
	}

#line 168 "src/filename.c"
	{
	int _klen;
	unsigned int _trans;
	const char *_acts;
	unsigned int _nacts;
	const char *_keys;

	if ( p == pe )
		goto _test_eof;
_resume:
	_keys = _filename_trans_keys + _filename_key_offsets[cs];
	_trans = _filename_index_offsets[cs];

	_klen = _filename_single_lengths[cs];
	if ( _klen > 0 ) {
		const char *_lower = _keys;
		const char *_mid;
		const char *_upper = _keys + _klen - 1;
		while (1) {
			if ( _upper < _lower )
				break;

			_mid = _lower + ((_upper-_lower) >> 1);
			if ( (*p) < *_mid )
				_upper = _mid - 1;
			else if ( (*p) > *_mid )
				_lower = _mid + 1;
			else {
				_trans += (unsigned int)(_mid - _keys);
				goto _match;
			}
		}
		_keys += _klen;
		_trans += _klen;
	}

	_klen = _filename_range_lengths[cs];
	if ( _klen > 0 ) {
		const char *_lower = _keys;
		const char *_mid;
		const char *_upper = _keys + (_klen<<1) - 2;
		while (1) {
			if ( _upper < _lower )
				break;

			_mid = _lower + (((_upper-_lower) >> 1) & ~1);
			if ( (*p) < _mid[0] )
				_upper = _mid - 2;
			else if ( (*p) > _mid[1] )
				_lower = _mid + 2;
			else {
				_trans += (unsigned int)((_mid - _keys)>>1);
				goto _match;
			}
		}
		_trans += _klen;
	}

_match:
	_trans = _filename_indicies[_trans];
	cs = _filename_trans_targs[_trans];

	if ( _filename_trans_actions[_trans] == 0 )
		goto _again;

	_acts = _filename_actions + _filename_trans_actions[_trans];
	_nacts = (unsigned int) *_acts++;
	while ( _nacts-- > 0 )
	{
		switch ( *_acts++ )
		{
	case 0:
#line 91 "src/filename.rl"
	{
			tar = p-3;
		}
	break;
	case 1:
#line 95 "src/filename.rl"
	{
			if(tar == NULL) {
				tar = p;
				// printf("tar: %s\n", p);
			}
		}
	break;
	case 2:
#line 102 "src/filename.rl"
	{
			slash = p+1;
			tar = NULL;
			ext = NULL;
		}
	break;
	case 3:
#line 108 "src/filename.rl"
	{
			ext = p;
		}
	break;
#line 269 "src/filename.c"
		}
	}

_again:
	if ( ++p != pe )
		goto _resume;
	_test_eof: {}
	}

#line 124 "src/filename.rl"


	if(tar == NULL) {
		return newFilename(slash, ext, false);
	} else {
		return newFilename(slash, tar, true);
	}
};
