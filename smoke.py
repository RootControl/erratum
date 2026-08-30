#!/usr/bin/env python3
"""Drives the *TUI* through a pty, answering the terminal capability queries
Bubble Tea blocks on.

The output is every frame Bubble Tea painted, concatenated. A non-alt-screen
TUI repaints its whole viewport on each keystroke, so lines repeat and `tail`
shows you an earlier frame rather than the final state — which reads exactly
like commands being swallowed and is not. Use this only to check that the TUI
itself renders and does not crash.

There is a second thing that used to look like a swallowed first command and
was real: Bubble Tea requires an answer to its background-colour query and
then does not consume it, so the reply landed in the text input and prefixed
whatever was typed first. This script clears the input before sending
anything. If you write your own driver, do the same.

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

def send(data):
    """Write to the pty, tolerating a program that has already exited.

    QUIT and the five endings both close the program on their own, so the
    trailing ctrl-C has nothing to go to and the pty raises EIO."""
    try:
        os.write(fd, data)
    except OSError:
        return False
    return True

pump(3.0)

# Bubble Tea blocks on the background-colour query at startup and will not
# paint until it is answered — but it never consumes the answer, so those
# bytes fall through to the text input and prefix whatever is typed first:
#
#     > ]11;rgb:1c1c/1a1a/1717\x writ
#     That is not a verb this terminal knows.
#
# The reply is sent once and is required; both ST and BEL terminators leak
# the same way. So clear the input widget before the first real command.
# Ctrl-U is the whole-line kill and the input is empty afterwards.
send(b"\x15")
pump(0.4)

for c in CMDS:
    if not send(c):
        break
    pump(0.7)
send(b"\x03")
pump(0.5)

t = buf.decode("utf8", "replace")
t = re.sub(r"\x1b\][^\x07\x1b]*(\x07|\x1b\\)", "", t)
t = re.sub(r"\x1b\[[0-9;?]*[a-zA-Z]", "", t)
t = t.replace("\x1b(B", "").replace("\x1b>", "").replace("\x1b=", "")
lines = [l.rstrip() for l in t.split("\n")]
print("\n".join(l for l in lines if l.strip()))
