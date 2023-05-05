/* nl_types.h - dummy POSIX localization message catalog implementation */

/*
 * Copyright (c) 2019, 2021 Wind River Systems, Inc.
 *
 * The right to copy, distribute, modify, or otherwise make use
 * of this software may be licensed only pursuant to the terms
 * of an applicable Wind River license agreement.
 */

#ifndef __INCnl_typesh
#define __INCnl_typesh

/* includes */

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

/* defines */

/* typedefs */

typedef int 	nl_catd;

/* function declarations */

static inline int catclose(nl_catd catd)
    {
    return -1;
    }

static inline char *catgets(nl_catd catd, int set_id, int msg_id, const char *s)
    {
    return (char *)s;
    }

static inline nl_catd catopen(const char *name, int oflag)
    {
    return -1;
    }

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __INCnl_typesh */
