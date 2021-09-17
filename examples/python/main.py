import time
import socker


if __name__ == '__main__':
    socker.start_logger("127.0.0.1", 0, 40000)
    obj = {}
    obj["name"] = "testing logging an object"
    for i in range(0, 1000):
        if i % 5 == 0:
            socker.log(f'UDP testing LOG from python {i}')
        elif i % 4 == 0:
            socker.dbg(f'UDP testing DBG from python {i}')
        elif i % 3 == 0:
            socker.success(f'UDP testing SUCCESS from python {i}')
        elif i % 2 == 0:
            socker.err(f'UDP testing DBG from python {i}')
        elif i % 1 == 0:
            socker.warn(f'UDP testing WARN from python {i}')
            socker.warn(obj)
        time.sleep(.1)
