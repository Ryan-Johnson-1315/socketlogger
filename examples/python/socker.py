from inspect import getframeinfo, currentframe
import socket
import json
import queue
import threading
import os


__initialized__ = False
__q = queue.Queue(150)
__s = socket.socket()

MESSAGE_LVL_LOG = 0
MESSAGE_LVL_WARN = 1
MESSAGE_LVL_SUCCESS = 2
MESSAGE_LVL_ERROR = 3
MESSAGE_LVL_DEBUG = 4


def start_logger(ip: str, client_port: int, server_port: int):
    global __s
    __s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    __s.bind((ip, client_port))
    print(__s)
    print((ip, client_port))
    global __initialized__
    __initialized__ = True
    threading.Thread(target=__send_msgs__, args=[(ip, server_port)]).start()


def log(msg: str):
    caller = getframeinfo(currentframe())
    _, f = os.path.split(caller.filename)
    __send__(MESSAGE_LVL_LOG, f'{f}:{caller.lineno}', msg)


def warn(msg: str):
    caller = getframeinfo(currentframe())
    _, f = os.path.split(caller.filename)
    __send__(MESSAGE_LVL_WARN, f'{f}:{caller.lineno}', msg)


def success(msg: str):
    caller = getframeinfo(currentframe())
    _, f = os.path.split(caller.filename)
    __send__(MESSAGE_LVL_SUCCESS, f'{f}:{caller.lineno}', msg)


def dbg(msg: str):
    caller = getframeinfo(currentframe())
    _, f = os.path.split(caller.filename)
    __send__(MESSAGE_LVL_DEBUG, f'{f}:{caller.lineno}', msg)


def err(msg: str):
    caller = getframeinfo(currentframe())
    _, f = os.path.split(caller.filename)
    __send__(MESSAGE_LVL_ERROR, f'{f}:{caller.lineno}', msg)


def __send__(lvl: int, caller: str, msg: str):
    m = {}
    m["level"] = lvl
    m["caller"] = caller
    m["message"] = msg

    if __initialized__:
        __q.put_nowait(bytes(json.dumps(m), encoding='utf-8'))
    else:
        print(f'{m["caller"]} -- {m["message"]}')


def __send_msgs__(addr):
    print("sending to", __s)
    while True:
        bts = __q.get()
        __s.sendto(bts, addr)
        # time.sleep(0.001)
