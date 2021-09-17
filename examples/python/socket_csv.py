import inspect
import socket
import json
import queue
import threading
import os


__initialized__ = False
__q = queue.Queue(150)
__s = socket.socket()


def start_csv(ip: str, client_port: int, server_port: int):
    global __s
    __s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    __s.bind((ip, client_port))

    global __initialized__
    __initialized__ = True

    threading.Thread(
        target=__send_msgs__, args=[(ip, server_port)]).start()


def append_row(file_name: str, row: list):
    caller = inspect.getframeinfo(inspect.stack()[1][0])
    _, f = os.path.split(caller.filename)

    if __initialized__:
        m = {}
        m["csv_filename"] = file_name
        m["caller"] = f"{f}:{caller.lineno}"
        m["row"] = row
        __q.put_nowait(bytes(json.dumps(m), encoding="utf-8"))
    else:
        print(
            f'WARNING!!! Not initialized!!! file_name: {file_name} data: {row}')


def __send_msgs__(addr):
    while True:
        bts = __q.get()
        __s.sendto(bts, addr)
