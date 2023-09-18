import time
import os
from hashlib import md5

def file_checksum(filepath):
    with open(filepath, 'rb') as f:
        return md5(f.read()).hexdigest()

def restart_service():
    os.system("supervisorctl restart app")

if __name__ == "__main__":
    last_checksum = file_checksum("/config.json")

    while True:
        current_checksum = file_checksum("/config.json")

        if last_checksum != current_checksum:
            print("Configuration file changed, restarting the app.")
            restart_service()
            last_checksum = current_checksum

        time.sleep(5)  # Poll every 5 seconds
