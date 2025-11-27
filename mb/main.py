# Imports go at the top
from microbit import *
import time

while True:
    print("A0:{}".format(accelerometer.get_x()))
    print("A1:{}".format(accelerometer.get_y()))
    print("A2:{}".format(accelerometer.get_z()))
    time.sleep_ms(200)
