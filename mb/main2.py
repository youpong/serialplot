# Imports go at the top
from microbit import *
import time

# Code in a 'while True:' loop repeats forever
while True:
    if button_a.is_pressed():
        print("A0:{}".format(accelerometer.get_x()))
    if button_b.is_pressed():
        print("A1:{}".format(accelerometer.get_y()))
    time.sleep_ms(10)
