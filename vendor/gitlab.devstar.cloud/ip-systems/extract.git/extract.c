// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.


#include "extract.h"
#include "filename.h"
#include "logger.h"
#include "filepath.h"

#include <stdlib.h>
#include <archive.h>
#include <archive_entry.h>
#include <string.h>
#include <stdbool.h>
#include <errno.h>

status copy_data(struct archive *ar, const char* pathname, struct archive *aw, warning_array* wa)
{
	int r;
	const void *buff;
	size_t size;
	la_int64_t offset;

	for (int i = 0;;i++) {
    	r = archive_read_data_block(ar, &buff, &size, &offset);
    	if (r == ARCHIVE_EOF) {// end loop
			return report_status(ARCHIVE_OK, "", NULL, NULL);
		}
    	if (ARCHIVE_WARN < r && r < ARCHIVE_OK) {// log warning on read
			warn(wa, __LINE__, pathname, archive_error_string(ar));
			//return report_status(r, archive_error_string(ar), NULL, NULL);
		}
		if (r < ARCHIVE_WARN) {// fatal read
			return report_status(r, archive_error_string(ar), NULL, wa);
		}

    	r = archive_write_data_block(aw, buff, size, offset);
    	if (ARCHIVE_WARN < r && r < ARCHIVE_OK) {// log warning on write
			warn(wa, __LINE__, pathname, archive_error_string(aw));
			//return report_status(r, archive_error_string(aw), NULL, NULL);
    	}
		if (r < ARCHIVE_WARN) {// fatal write
			return report_status(r, archive_error_string(aw), NULL, wa);
		}
  	}
}

// Strips / from pathname if exists
void force_relative_entry(struct archive_entry *entry) {
	const char* pathname = archive_entry_pathname(entry);

	if (pathname == NULL || pathname[0] != '/') {
		return;
	}
	
	archive_entry_set_pathname(entry, pathname+1);
}

status _extract(const char *filename, const char* destination)
{
	vlog("_extract(%s)\n", filename);
	struct archive *a;
  	struct archive *ext;
  	struct archive_entry *entry;
  	int flags;
  	int r;
	warning_array* wa = warning_array_init();

	// Select which attributes we want to restore.
  	flags = ARCHIVE_EXTRACT_TIME;
  	//flags |= ARCHIVE_EXTRACT_PERM;
  	flags |= ARCHIVE_EXTRACT_ACL;
  	flags |= ARCHIVE_EXTRACT_FFLAGS;

  	a = archive_read_new();
	archive_read_support_format_all(a);
  	archive_read_support_filter_all(a);
  	ext = archive_write_disk_new();
  	archive_write_disk_set_options(ext, flags);
  	archive_write_disk_set_standard_lookup(ext);
	if ((r = archive_read_open_filename(a, filename, 10240))) {
		status stat = report_status(r, archive_error_string(a), NULL, wa);
		archive_read_free(a);
		archive_write_free(ext);
		return stat;
	}
  	for (;;) {
    	r = archive_read_next_header(a, &entry);
    	if (r == ARCHIVE_EOF) {
			break;
		}
    	if (r < ARCHIVE_OK) {
			warn(wa, __LINE__, archive_entry_pathname(entry), archive_error_string(a));
		}
		if (r < ARCHIVE_WARN) {
			status stat = report_status(r, archive_error_string(a), NULL, wa);
			archive_read_free(a);
			archive_write_free(ext);
			return stat;
		}
		if (archive_entry_filetype(entry) | AE_IFDIR) { // Fix:? Directory extracting with no execute permissions
			// If directory, set read and execute permissions
			archive_entry_set_perm(entry, 0755);
		} else {
			// If file, set only read permissions
			archive_entry_set_perm(entry, 0644);
		}
		force_relative_entry(entry);
		// update entry path with destination
		char* entry_destination = join(destination, archive_entry_pathname(entry));
		archive_entry_set_pathname(entry, entry_destination);
		// write entry
    	r = archive_write_header(ext, entry);
    	if (r < ARCHIVE_OK) {
			warn(wa, __LINE__, archive_entry_pathname(entry), archive_error_string(ext));
		}
    	else if (archive_entry_size(entry) > 0) {
      		status stat = copy_data(a, archive_entry_pathname(entry), ext, wa);
			r = stat->code;
			if (r < ARCHIVE_WARN) {
				archive_read_free(a);
				archive_write_free(ext);
				free(entry_destination);
				vlog("returning stat: %p\n", stat);
				return stat;
			}
			status_free(stat);
    	}
    	r = archive_write_finish_entry(ext);
    	if (r < ARCHIVE_OK) {
			warn(wa, __LINE__, NULL, archive_error_string(ext));
		}
		if (r < ARCHIVE_WARN) {
			status stat = report_status(r, archive_error_string(ext), NULL, wa);
			archive_read_free(a);
			archive_write_free(ext);
			free(entry_destination);
			return stat;
		}
		free(entry_destination);
	}
  	archive_read_close(a);
  	archive_read_free(a);
  	archive_write_close(ext);
	archive_write_free(ext);
  	return success(wa);
}

status _decompress(const char *filepath, const char *destname)
{
	struct archive *a;
  	struct archive_entry *entry;
  	int flags;
  	int r;
	warning_array* wa = warning_array_init();

  	a = archive_read_new();
	archive_read_support_format_raw(a);
  	archive_read_support_filter_all(a);
	if ((r = archive_read_open_filename(a, filepath, 10240))) {
		status stat = report_status(r, archive_error_string(a), NULL, wa);
		archive_read_free(a);
		return stat;
	}

	r = archive_read_next_header(a, &entry);
	if (r == ARCHIVE_EOF) {
		status stat = report_status(r, "EOF found", "archive_read_next_header", wa);
		archive_read_free(a);
		return stat;
	}
	if (r < ARCHIVE_OK) {
		fprintf(stderr, "[%d %s] %s\n", __LINE__, archive_entry_pathname(entry), archive_error_string(a));
		warn(wa, __LINE__, archive_entry_pathname(entry), archive_error_string(a));
	}
	if (r < ARCHIVE_WARN) {
		status stat = report_status(r, archive_error_string(a), "archive_read_next_header", wa);
		archive_read_free(a);
		return stat;
	}

	FILE *fp = fopen(destname, "wb+");
	if (fp == NULL) {
		extern int errno;
		status stat = report_status(errno, strerror(errno), "fopen", wa);
		archive_read_free(a);
		return stat;
	}

	char buf[1024];
	for(int size = archive_read_data(a, buf, 1024); size > 0; size = archive_read_data(a, buf, 1024)) {
		fwrite(buf, sizeof(char), size, fp);
	}

	archive_read_data_skip(a);

  	archive_read_close(a);
  	archive_read_free(a);
	fclose(fp);
  	return success(wa);
}

//extract requires pwd to be the destination directory
//if filename is null, it is expected _extract will extract the archive without problem
//if filename is not null, it is expected that the archive is not a tar, so _decompress should be tried if _extract fails
status extract(char *filepath, char *filename, char *dest)
{
	fprintf(stderr, "extract(\"%s\", \"%s\", \"%s\")\n", filepath, filename, dest);
	filename_ptr fp;
	if(filename == NULL) {
		fp = parseFilename(filepath);
	} else {
		fp = parseFilename(filename);
	}

	char* destination = join(dest, getBasename(fp));
	fprintf(stderr, "join(\"%s\", \"%s\") = %s\n", dest, getBasename(fp), destination);

	status ret = NULL;
	if(compressedBinary(fp)) {
		ret = _decompress(filepath, destination);
	} else {
		ret =  _extract(filepath, destination);
		if(ret->code != 0) {
			//on failure, check filename to see if might be compressed binary file
			//e.g. a data.gz rather than a data.tar.gz
			if(filename != NULL) {
				status de = _decompress(filepath, destination);

				status_free(ret);
				free(destination);
				filename_free(fp);
				//free(filepath);
				//free(filename);
				return de;
			} else {
				free(destination);
				filename_free(fp);
				//free(filepath);
				//free(filename);
				return ret;
			}
		}
	}

	free(destination);
	filename_free(fp);
	return ret;
}