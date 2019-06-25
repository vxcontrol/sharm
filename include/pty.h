#ifndef PTY_H
#define PTY_H

// __cplusplus gets defined when a C++ compiler processes the file
#ifdef __cplusplus
// extern "C" is needed so the C++ compiler exports the symbols without name
// manging.
extern "C" {
#endif

//External API

#ifndef DWORD
#define DWORD unsigned long
#endif

/*
* PtyOpen
* pty.open(cols, rows, shell_path)
* return handle of PTY object (self)
*/
void * PtyOpen(int cols, int rows, char * cmd);

/*
* PtyResize
* pty.resize(self, cols, rows);
*/
DWORD PtyResize(void *self, int cols, int rows);

/*
* PtyKill
* pty.kill(self);
*/
DWORD PtyKill(void *self);

/*
* PtyRead
* pty.read(self, data, size)
*/
DWORD PtyRead(void *self, unsigned char *data, DWORD size);

/*
* PtyWrite
* pty.write(self, data, amount)
*/
DWORD PtyWrite(void *self, const unsigned char *data, DWORD size);

#ifdef __cplusplus
}
#endif

#endif
