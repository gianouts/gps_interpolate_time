# gps_interpolate_time

***This is currently a work in progress and not completely working.***

Takes two gpx files as input (one with timestamps (file a), the other without (file b)), merges them together, interpolates the timestamp for points in file b using timestamps based on known points in file a and generates an outputfile. 

Input file a has lat, long, timestamp
Input file b only has lat, long

Input files are not altered.

./gps_interpolate_time -a [input_file_a] -b [input_file_b] -o [output_file]
