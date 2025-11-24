# Imports go at the top
from microbit import *
import time

while True:
    x = accelerometer.get_x()
    y = accelerometer.get_y()
    print("A0:", x)
    print("A1:", y)
    time.sleep_ms(100)
