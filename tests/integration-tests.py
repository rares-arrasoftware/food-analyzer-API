# This Python test runner starts the real Go server, runs integration tests, and cleans up the DB.
import os
import subprocess
import sys
import time
import requests
import platform

TEST_DB = "test.db"
BINARY_NAME = "food-analyzer-api.exe" if platform.system() == "Windows" else "food-analyzer-api"
PORT = "4000"
BASE_URL = f"http://localhost:{PORT}"


def build_server():
    print("Building Go app...")
    result = subprocess.run(["go", "build", "-o", BINARY_NAME, "../cmd/food-analyzer-api"], capture_output=True)
    if result.returncode != 0:
        print("Build failed:", result.stderr.decode())
        sys.exit(1)


def start_server():
    print("Starting server...")
    return subprocess.Popen([
        BINARY_NAME,
        "--port", PORT,
        "--db", TEST_DB,
        "--jwt-secret", "testsecret",
        "--jwt-expiry", "1"
    ])

def wait_for_server():
    for _ in range(20):
        try:
            r = requests.get(f"{BASE_URL}/auth/login")
            if r.status_code in (400, 405):
                return
        except requests.exceptions.ConnectionError:
            pass
        time.sleep(0.5)
    print("Server didn't start in time.")
    sys.exit(1)


def cleanup():
    if os.path.exists(TEST_DB):
        os.remove(TEST_DB)
    if os.path.exists(BINARY_NAME):
        os.remove(BINARY_NAME)


def run_tests():
    print("Running integration tests...")

    # Register
    r = requests.post(f"{BASE_URL}/auth/register", json={
        "email": "test@example.com",
        "password": "secret123",
        "firstName": "Test",
        "lastName": "User"
    })
    assert r.status_code == 200, f"Register failed: {r.text}"

    # Login
    r = requests.post(f"{BASE_URL}/auth/login", json={
        "email": "test@example.com",
        "password": "secret123"
    })
    assert r.status_code == 200, f"Login failed: {r.text}"
    token = r.json().get("token")
    assert token, "Missing token in login response"

    # Analyze food
    with open("../images/food.jpg", "rb") as f:
        r = requests.post(
            f"{BASE_URL}/food/analyze",
            headers={"Authorization": f"Bearer {token}"},
            files={"image": ("food.jpg", f, "image/jpeg")}
        )
    assert r.status_code == 200, f"Analyze failed: {r.text}"
    print("All tests passed.")


if __name__ == '__main__':
    cleanup()
    build_server()
    proc = start_server()
    try:
        wait_for_server()
        run_tests()
    finally:
        print("Cleaning up...")
        proc.terminate()
        proc.wait()
        cleanup()
