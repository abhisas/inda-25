import random
import sys
from datetime import datetime, timedelta

random.seed(42) 
# --- CONFIGURATION ---
# Change this number to generate larger or smaller files
# 1,000,000 lines is roughly 70-80 MB. 10,000,000 is ~750 MB.
NUM_LINES = 1000000 
OUTPUT_FILE = "server_logs.txt"

# --- DATA POOLS & WEIGHTS ---
# Simulating Data Skew: The homepage ('/') gets 70% of the traffic
URLS = ['/', '/products', '/about', '/contact', '/api/checkout', '/hidden-admin']
URL_WEIGHTS = [70, 15, 5, 5, 4, 1] 

# Mostly GET requests, some POSTs
METHODS = ['GET', 'POST', 'PUT', 'DELETE']
METHOD_WEIGHTS = [85, 10, 3, 2]

# Mostly successful 200s, but enough 404s to make Task 2 interesting
STATUS_CODES = ['200', '301', '404', '500']
STATUS_WEIGHTS = [80, 5, 12, 3]

def generate_random_ip():
    """Generates a random, plausible-looking IPv4 address."""
    return f"{random.randint(1, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}.{random.randint(0, 255)}"

def generate_logs(num_lines, filename):
    print(f"Generating {num_lines:,} log lines. This might take a moment...")
    
    start_time = datetime.now() - timedelta(days=30)
    
    with open(filename, 'w') as f:
        for i in range(num_lines):
            # Advance time slightly for each log entry
            start_time += timedelta(seconds=random.randint(1, 5))
            timestamp = start_time.strftime("%Y-%m-%d:%H:%M:%S")
            
            ip = generate_random_ip()
            method = random.choices(METHODS, weights=METHOD_WEIGHTS)[0]
            url = random.choices(URLS, weights=URL_WEIGHTS)[0]
            status = random.choices(STATUS_CODES, weights=STATUS_WEIGHTS)[0]
            
            # Formatting: [TIMESTAMP] [IP_ADDRESS] [HTTP_METHOD] [URL] [HTTP_STATUS]
            log_line = f"[{timestamp}] [{ip}] [{method}] [{url}] [{status}]\n"
            f.write(log_line)
            
            # Print progress every 10%
            if (i + 1) % (num_lines // 10) == 0:
                print(f"Progress: {(i + 1) / num_lines * 100:.0f}%")

    print(f"Done! Saved to {filename}")

if __name__ == "__main__":
    generate_logs(NUM_LINES, OUTPUT_FILE)