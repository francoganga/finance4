#!/bin/bash

air&

# trap onexit INT
# function onexit() {
#     kill $air_pid
# }

tailwind -i public/style.css -o public/main.css -w
