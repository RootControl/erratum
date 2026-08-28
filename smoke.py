#!/usr/bin/env python3
"""Drives the *TUI* through a pty, answering the terminal capability queries
Bubble Tea blocks on.

The output is every frame Bubble Tea painted, concatenated. A non-alt-screen
TUI repaints its whole viewport on each keystroke, so lines repeat and `tail`
shows you an earlier frame rather than the final state — which reads exactly
like commands being swallowed and is not. Use this only to check that the TUI
itself renders and does not crash.

For anything scripted, use plain mode instead, which has no repaints:

    go build -o erratum .
    printf 'x board\nx board\nfile A1\n' | ERRATUM_PLAIN=1 ./erratum
"""
import os, pty, time, select, re, fcntl, termios, struct, sys

# The binary is expected next to this script; override with ERRATUM_BIN.
HERE = os.path.dirname(os.path.abspath(__file__))
BIN = os.environ.get("ERRATUM_BIN") or os.path.join(HERE, "erratum")

CMDS = [c.encode() + b"\r" for c in sys.argv[1:]] or [b"x board\r"]

if not os.path.exists(BIN):
    sys.exit("no binary at %s\n"
             "build it first:  go build -o erratum .\n"
             "or set ERRATUM_BIN to point elsewhere" % BIN)

pid, fd = pty.fork()
if pid == 0:
    os.environ.update({"TERM": "xterm-256color", "ERRATUM_NOALT": "1"})
    os.execv(BIN, ["erratum"])
    os._exit(1)
fcntl.ioctl(fd, termios.TIOCSWINSZ, struct.pack("HHHH", 46, 96, 0, 0))
os.set_blocking(fd, False)

buf = b""
def pump(t):
    global buf
    end = time.time() + t
    while time.time() < end:
        r, _, _ = select.select([fd], [], [], 0.05)
        if not r:
            continue
        try:
            d = os.read(fd, 65536)
        except OSError:
            return
        if not d:
            return
        buf += d
        if b"\x1b[6n" in d:
            os.write(fd, b"\x1b[1;1R")
        if b"\x1b]11;?" in d:
            os.write(fd, b"\x1b]11;rgb:1c1c/1a1a/1717\x1b\\")

pump(3.0)
for c in CMDS:
    os.write(fd, c)
    pump(0.7)
os.write(fd, b"\x03")
pump(0.5)

t = buf.decode("utf8", "replace")
t = re.sub(r"\x1b\][^\x07\x1b]*(\x07|\x1b\\)", "", t)
t = re.sub(r"\x1b\[[0-9;?]*[a-zA-Z]", "", t)
t = t.replace("\x1b(B", "").replace("\x1b>", "").replace("\x1b=", "")
lines = [l.rstrip() for l in t.split("\n")]
print("\n".join(l for l in lines if l.strip()))
