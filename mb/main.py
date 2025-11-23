# Imports go at the top
from microbit import *
import time

while True:
    x = accelerometer.get_x()
    print("A0:", x)
    time.sleep_ms(100)
#    print("DEBUG:", time.ticks_ms() / 1000.0)
