import time
import os
from hashlib import md5

FILEPATH = "/data/config.json"

def file_checksum(filepath):
    """Return the MD5 checksum of a file. Try until the file exists."""
    while True:
        try:
            with open(filepath, 'rb') as f:
                return md5(f.read()).hexdigest()
        except FileNotFoundError:
            print("File not found, waiting for it to appear.")
            time.sleep(5)
            continue

def restart_service():
    os.system("supervisorctl restart app")

if __name__ == "__main__":
    last_checksum = file_checksum(FILEPATH)

    while True:
        current_checksum = file_checksum(FILEPATH)

        if last_checksum != current_checksum:
            print("Watcher Service: Configuration file changed, restarting the app.")
            restart_service()
            last_checksum = current_checksum

        time.sleep(3)  # Poll every 5 seconds
