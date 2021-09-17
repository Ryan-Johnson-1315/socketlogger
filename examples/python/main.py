import time
import threading
import socket_logger
import socket_csv


def run_logs():
    socket_logger.start_logger("127.0.0.1", 0, 40000)
    obj = {}
    obj["name"] = "testing logging an object"
    for i in range(0, 25):
        if i % 5 == 0:
            socket_logger.log(f'UDP testing LOG from python {i}')
        elif i % 4 == 0:
            socket_logger.dbg(f'UDP testing DBG from python {i}')
        elif i % 3 == 0:
            socket_logger.success(f'UDP testing SUCCESS from python {i}')
        elif i % 2 == 0:
            socket_logger.err(f'UDP testing DBG from python {i}')
        elif i % 1 == 0:
            socket_logger.warn(f'UDP testing WARN from python {i}')
            socket_logger.warn(obj)
        time.sleep(.03)


def run_csv():
    socket_csv.start_csv("127.0.0.1", 0, 50000)

    socket_csv.append_row("python_csv.csv", ['id', 'name', 'number', 'i'])
    for i in range(0, 25):
        socket_csv.append_row("python_csv.csv", [i*-1, f'name-{i}', i*i, i])


if __name__ == '__main__':
    threading.Thread(target=run_logs).start()
    threading.Thread(target=run_csv).start()
