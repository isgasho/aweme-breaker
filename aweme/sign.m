#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <dlfcn.h>
#include <objc/runtime.h>
#include <objc/message.h>

#include "aweme.h"

#import <Foundation/Foundation.h>

static void *lib = NULL;
static id instance = NULL;

int
Init(const char *path) {
    lib = dlopen(path, RTLD_LAZY);
    if (!lib) {
        fprintf(stderr, "dlopen: %s\n", dlerror());
        return -1;
    }

    instance = ((id (*)(id, SEL))objc_msgSend)((id)objc_getClass("IESAntiSpam"), sel_registerName("sharedInstance"));

    return 0;
}

void
AdjustWithServerTime(long long timestamp) {
    ((void (*)(id, SEL, long long))objc_msgSend)(instance, sel_registerName("adjustWithServerTime:"), timestamp);
}

int
Signature(
    const char *url, const char *stub, const char *cookie,
    void *gorgon, void *khronos) {

    if (!url) {
        return -1;
    }

    NSString *sUrl = [NSString stringWithCString:url encoding:[NSString defaultCStringEncoding]];

    NSMutableDictionary *header = [[NSMutableDictionary alloc] init];

    if (stub) {
        [header setValue:[NSString stringWithCString:stub encoding:[NSString defaultCStringEncoding]]
                  forKey:@"x-ss-stub"];
    }

    if (cookie) {
        [header setValue:[NSString stringWithCString:cookie encoding:[NSString defaultCStringEncoding]]
                  forKey:@"cookie"];
    }

    NSURL *nUrl = [[NSURL alloc] initWithString:sUrl];

    if (!nUrl) {
        return -2;
    }

    NSDictionary *result = (NSDictionary *)((id (*)(id, SEL, id, id))objc_msgSend)(instance, sel_registerName("sgm_encryptWithURL:msg:"), nUrl, header);

    if (gorgon) {
        NSString *nsGorgon = (NSString *)[result valueForKey:@"X-Gorgon"];
        NSUInteger nsGorgonLen = [nsGorgon lengthOfBytesUsingEncoding:[NSString defaultCStringEncoding]];
        const char *cGorgon = [nsGorgon cStringUsingEncoding:[NSString defaultCStringEncoding]];
        memcpy(gorgon, cGorgon, nsGorgonLen);
    }

    if (khronos) {
        NSString *nsKhronos = (NSString*)[result valueForKey:@"X-Khronos"];
        NSUInteger nsKhronosLen = [nsKhronos lengthOfBytesUsingEncoding:[NSString defaultCStringEncoding]];
        const char *cKhronos = [nsKhronos cStringUsingEncoding:[NSString defaultCStringEncoding]];
        memcpy(khronos, cKhronos, nsKhronosLen);
    }

    return 0;
}

void
Dispose() {
    dlclose(lib);
}
