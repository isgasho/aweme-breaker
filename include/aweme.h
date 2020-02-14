#ifndef _H_AWEME
#define _H_AWEME

int Init(const char *path);
int Signature(const char *url, const char *stub, const char *cookie, void *gorgon, void *khronos);
void AdjustWithServerTime(long long timestamp);
void Dispose();

#endif
